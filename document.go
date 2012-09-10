package readability

import (
	"net/http"
	"io/ioutil"
)

var regexes = map[string]string{
	"unlikelyCandidatesRe":   `/combx|comment|community|disqus|extra|foot|header|menu|remark|rss|shoutbox|sidebar|sponsor|ad-break|agegate|pagination|pager|popup/i`,
	"okMaybeItsACandidateRe": `/and|article|body|column|main|shadow/i`,
	"positiveRe":             `/article|body|content|entry|hentry|main|page|pagination|post|text|blog|story/i`,
	"negativeRe":             `/combx|comment|com-|contact|foot|footer|footnote|masthead|media|meta|outbrain|promo|related|scroll|shoutbox|sidebar|sponsor|shopping|tags|tool|widget/i`,
	"divToPElementsRe":       `/<(a|blockquote|dl|div|img|ol|p|pre|table|ul)/i`,
	"replaceBrsRe":           `/(<br[^>]*>[ \n\r\t]*){2,}/i`,
	"replaceFontsRe":         `/<(\/?)font[^>]*>/i`,
	"trimRe":                 `/^\s+|\s+$/`,
	"normalizeRe":            `/\s{2,}/`,
	"killBreaksRe":           `/(<br\s*\/?>(\s|&nbsp;?)*){1,}/`,
	"videoRe":                `/http:\/\/(www\.)?(youtube|vimeo)\.com/i`,
}

type document struct {
	body []byte
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
	return &document{body}, nil
}
