package readability

import (
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