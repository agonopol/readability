package readability

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
	"strings"
)

type htmlParser struct {
	doc *html.HtmlDocument
}

func HTMLParser(content []byte) (*htmlParser, error) {
	content = regexes["replaceFontsRe"].ReplaceAll(regexes["replaceBrsRe"].ReplaceAll(content, []byte("</p><p>")), []byte("\\1span>"))

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

	return &htmlParser{doc}, nil
}

func (this *htmlParser) free() {
	this.doc.Free()
}

func (this *htmlParser) Body() []byte {
	return []byte(this.doc.Content())
}

func (this *htmlParser) walkElements(node xml.Node, f func(xml.Node) error) error {
	f(node)
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		err := this.walkElements(child, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *htmlParser) removeUnlikelyCandidates() error {
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

func (this *htmlParser) transformMisusedDivsIntoParagraphs() error {
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

func (this *htmlParser) scoreParagraphs() (map[*xml.Node]int64, error) {

	return nil, nil
}

func (this *htmlParser) prepareCandidates() error {
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

	this.transformMisusedDivsIntoParagraphs()
	return nil
}

func (this *htmlParser) paragraphs() ([]xml.Node, error) {
	p, err := this.doc.Search("p,td")
	if err != nil {
		return nil, err
	}
	return p, nil
}
