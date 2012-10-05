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
	doc, err := gokogiri.ParseHtml(content)
	defer doc.Free()

	if err != nil {
		return nil, err
	}

	bodies, err := doc.Search("//body")
	if err != nil {
		return nil, err
	}

	body, err := gokogiri.ParseHtml([]byte(bodies[0].Content()))
	comments, err := body.Search("//comment()")
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		comment.Remove()
	}
	return &htmlParser{body}, nil
}

func (this *htmlParser) Body() []byte {
	return []byte(this.doc.Content())
}

func (this *htmlParser) walkElements(node xml.Node, f func(xml.Node) error) error {
	f(node)
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		err := this.walkElements(child, f)
		if err != nil {
			println("GOT ERROR")
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
	return nil
}
