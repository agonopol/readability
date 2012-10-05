package readability

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

var regexes map[string]*regexp.Regexp

func init() {
	regexes = make(map[string]*regexp.Regexp)
	regex, err := regexp.Compile(`/combx|comment|community|disqus|extra|foot|header|menu|remark|rss|shoutbox|sidebar|sponsor|ad-break|agegate|pagination|pager|popup/i`)
	if err != nil {
		panic(err)
	}
	regexes["unlikelyCandidatesRe"] = regex
	regex, err = regexp.Compile(`/and|article|body|column|main|shadow/i`)
	if err != nil {
		panic(err)
	}
	regexes["okMaybeItsACandidateRe"] = regex
	regex, err = regexp.Compile(`/article|body|content|entry|hentry|main|page|pagination|post|text|blog|story/i`)
	if err != nil {
		panic(err)
	}
	regexes["positiveRe"] = regex
	regex, err = regexp.Compile(`/combx|comment|com-|contact|foot|footer|footnote|masthead|media|meta|outbrain|promo|related|scroll|shoutbox|sidebar|sponsor|shopping|tags|tool|widget/i`)
	if err != nil {
		panic(err)
	}
	regexes["negativeRe"] = regex
	regex, err = regexp.Compile(`/<(a|blockquote|dl|div|img|ol|p|pre|table|ul)/i`)
	if err != nil {
		panic(err)
	}
	regexes["divToPElementsRe"] = regex
	regex, err = regexp.Compile(`/(<br[^>]*>[ \n\r\t]*){2,}/i`)
	if err != nil {
		panic(err)
	}
	regexes["replaceBrsRe"] = regex
	regex, err = regexp.Compile(`/<(\/?)font[^>]*>/i`)
	if err != nil {
		panic(err)
	}
	regexes["replaceFontsRe"] = regex
	regex, err = regexp.Compile(`/^\s+|\s+$/`)
	if err != nil {
		panic(err)
	}
	regexes["trimRe"] = regex
	regex, err = regexp.Compile(`/\s{2,}/`)
	if err != nil {
		panic(err)
	}
	regexes["normalizeRe"] = regex
	regex, err = regexp.Compile(`/(<br\s*\/?>(\s|&nbsp;?)*){1,}/`)
	if err != nil {
		panic(err)
	}
	regexes["killBreaksRe"] = regex
	regex, err = regexp.Compile(`/http:\/\/(www\.)?(youtube|vimeo)\.com/i`)
	if err != nil {
		panic(err)
	}
	regexes["videoRe"] = regex
}

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
	body = regexes["replaceFontsRe"].ReplaceAll(regexes["replaceBrsRe"].ReplaceAll(body, []byte("</p><p>")), []byte("<\\span>"))
	parser, err := HTMLParser(body)
	if err != nil {
		return nil, err
	}
	return &document{parser}, nil
}

func (this *document) Content() ([]byte, error) {
	if err := this.parser.prepareCandidates(); err != nil {
		return nil, err
	}
	if err := this.parser.removeUnlikelyCandidates(); err != nil {
		return nil, err
	}
	if err := this.parser.transformMisusedDivsIntoParagraphs(); err != nil {
		return nil, err
	}
	return this.parser.Body(), nil
}
