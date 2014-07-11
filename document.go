package readability

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strings"
)

type Document struct {
	doc *html.HtmlDocument
}

func ParseURL(url string) (*Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return Parse(body)
}

func Parse(content []byte) (*Document, error) {

	content = regexes["replaceFontsRe"].ReplaceAll(regexes["replaceBrsRe"].ReplaceAll(content, []byte(`</p><p>`)), []byte(`<\1span>`))

	doc, err := gokogiri.ParseHtml(content)

	if err != nil {
		return nil, err
	}

	comments, err := doc.Search("//comment()")

	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		comment.Remove()
	}

	return &Document{doc}, nil
}

func (this *Document) Content() (string, error) {
	if err := this.prepare(); err != nil {
		return "", err
	}

	candidates, err := getCandidates(this.doc, 25)
	if err != nil {
		return "", err
	}

	article := this.getArticle(candidates)

	return article.Content(), nil
}

func (this *Document) prepare() error {
	scripts, err := this.doc.Search("//script")
	if err != nil {
		return err
	}
	for _, script := range scripts {
		script.Remove()
	}
	styles, err := this.doc.Search("//style")
	if err != nil {
		return err
	}
	for _, style := range styles {
		style.Remove()
	}

	this.removeUnlikelyCandidates()

	this.misusedDivsIntoParagraphs()
	return nil
}

func (this *Document) Free() {
	this.doc.Free()
}

func (this *Document) walkElements(node xml.Node, f func(xml.Node) error) error {
	f(node)
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		err := this.walkElements(child, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Document) removeUnlikelyCandidates() error {
	this.walkElements(this.doc, func(node xml.Node) error {
		attributes := node.Attributes()
		info := ""
		if class, found := attributes["class"]; found {
			info = class.Value()
		}
		if id, found := attributes["id"]; found {
			info = fmt.Sprintf("%s%s", info, id.Value())
		}
		if regexes["unlikelyCandidatesRe"].Match([]byte(info)) && regexes["okMaybeItsACandidateRe"].Match([]byte(info)) {
			name := strings.ToLower(node.Name())
			if name != "body" && name != "html" {
				node.Remove()
			}
		}
		return nil
	})
	return nil
}

func (this *Document) misusedDivsIntoParagraphs() error {
	this.walkElements(this.doc, func(node xml.Node) error {
		if strings.ToLower(node.Name()) == "div" {
			if regexes["divToPElementsRe"].Match([]byte(node.InnerHtml())) {
				node.SetName("p")
			}
		}
		return nil
	})
	return nil
}

func extend(doc *xml.ElementNode, node xml.Node) {
	dup := node.Duplicate(1).(*xml.ElementNode)
	doc.AddChild(dup)
	if strings.ToLower(dup.Name()) != "p" || strings.ToLower(dup.Name()) != "div" {
		dup.SetName("div")
	}
}

func (this *Document) getArticle(candidates map[string]*Candidate) xml.Node {
	best := bestCandidate(candidates)
	threshold := math.Max(10.0, best.score*0.2)

	doc := this.doc.CreateElementNode("div")

	for node := best.node.Parent().FirstChild(); node != nil && node.IsValid(); node = node.NextSibling() {
		if node.String() == best.node.String() {
			extend(doc, node)
		} else if candidate, found := candidates[node.String()]; found && candidate.score >= threshold {
			extend(doc, node)
		} else if strings.ToLower(node.Name()) == "p" {
			lDensity := linkDensity(node)
			length := len(node.Content())
			if lDensity < 0.25 && length > 80 {
				extend(doc, node)
			} else if length < 80 && lDensity == 0.0 {
				match, err := regexp.Match(`\.( |$)`, []byte(node.Content()))
				if err != nil {
					panic(err)
				}
				if match {
					extend(doc, node)
				}
			}
		}

	}

	return doc
}

func (this *Document) sanitize(candidates map[string]*Candidate, article xml.Node) string {
	return ""

}

// def sanitize(node, candidates, options = {})
//   node.css("h1, h2, h3, h4, h5, h6").each do |header|
//     header.remove if class_weight(header) < 0 || get_link_density(header) > 0.33
//   end
//
//   node.css("form, object, iframe, embed").each do |elem|
//     elem.remove
//   end
//
//   if @options[:remove_empty_nodes]
//     # remove <p> tags that have no text content - this will also remove p tags that contain only images.
//     node.css("p").each do |elem|
//       elem.remove if elem.content.strip.empty?
//     end
//   end
//
//   # Conditionally clean <table>s, <ul>s, and <div>s
//   clean_conditionally(node, candidates, "table, ul, div")
//
//   # We'll sanitize all elements using a whitelist
//   base_whitelist = @options[:tags] || %w[div p]
//   # We'll add whitespace instead of block elements,
//   # so a<br>b will have a nice space between them
//   base_replace_with_whitespace = %w[br hr h1 h2 h3 h4 h5 h6 dl dd ol li ul address blockquote center]
//
//   # Use a hash for speed (don't want to make a million calls to include?)
//   whitelist = Hash.new
//   base_whitelist.each {|tag| whitelist[tag] = true }
//   replace_with_whitespace = Hash.new
//   base_replace_with_whitespace.each { |tag| replace_with_whitespace[tag] = true }
//
//   ([node] + node.css("*")).each do |el|
//     # If element is in whitelist, delete all its attributes
//     if whitelist[el.node_name]
//       el.attributes.each { |a, x| el.delete(a) unless @options[:attributes] && @options[:attributes].include?(a.to_s) }
//
//       # Otherwise, replace the element with its contents
//     else
//       # If element is root, replace the node as a text node
//       if el.parent.nil?
//         node = Nokogiri::XML::Text.new(el.text, el.document)
//         break
//       else
//         if replace_with_whitespace[el.node_name]
//           el.swap(Nokogiri::XML::Text.new(' ' << el.text << ' ', el.document))
//         else
//           el.swap(Nokogiri::XML::Text.new(el.text, el.document))
//         end
//       end
//     end
//
//   end
//
//   s = Nokogiri::XML::Node::SaveOptions
//   save_opts = s::NO_DECLARATION | s::NO_EMPTY_TAGS | s::AS_XHTML
//   html = node.serialize(:save_with => save_opts)
//
//   # Get rid of duplicate whitespace
//   return html.gsub(/[\r\n\f]+/, "\n" )
// end
