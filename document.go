package readability

import (
	"io/ioutil"
	"net/http"
)

type document struct {
	parser *htmlParser
}

func Document(url string) (*document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	parser, err := HTMLParser(body)
	if err != nil {
		return nil, err
	}

	return &document{parser}, nil
}

func (this *document) Content() (string, error) {
	if err := this.parser.prepareCandidates(); err != nil {
		return "", err
	}

	cadnidates, err := getCandidates(this.parser, 25)

	if err != nil {
		return "", err
	}

	article := getArticle(cadnidates)

	return article.Content(), nil
}

func (this *document) Free() {
	this.parser.free()
}
