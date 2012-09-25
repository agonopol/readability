package readability

import (
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/html"
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
	return &htmlParser{body}, nil
}

func (this *htmlParser) Body() []byte {
	return []byte(this.doc.Content())
}
