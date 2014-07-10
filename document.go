package readability

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
	"io/ioutil"
	"net/http"
	"math"
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
