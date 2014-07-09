package readability

import (
	"github.com/moovweb/gokogiri/xml"
	"math"
	"regexp"
	"strings"
)

func weight(elem xml.Node) float64 {
	weight := 0.0
	class := elem.Attr("class")
	if class != "" {
		if regexes["negativeRe"].Match([]byte(class)) {
			weight -= 25.0
		} else if regexes["positiveRe"].Match([]byte(class)) {
			weight += 25.0
		}
	}

	id := elem.Attr("id")
	if id != "" {
		if regexes["negativeRe"].Match([]byte(id)) {
			weight -= 25.0
		} else if regexes["positiveRe"].Match([]byte(id)) {
			weight += 25.0
		}
	}
	return weight
}

func score(node xml.Node) float64 {
	switch strings.ToLower(node.Name()) {
	case "div":
		return weight(node) + 5.0
	case "blockquote":
		return weight(node) + 3.0
	case "form":
		return weight(node) - 3.0
	case "th":
		return weight(node) - 5.0
	default:
		return weight(node)
	}
}

func linkDensity(node xml.Node) float64 {
	links, err := node.Search("a")
	if err != nil {
		return 0.0
	}

	llength := 0.0
	for _, link := range links {
		llength += float64(len(link.Content()))
	}
	tlength := float64(len(node.Content()))
	return llength / tlength
}

func getCandidates(html *htmlParser, minLen int) (map[xml.Node]float64, error) {

	candidates := make(map[xml.Node]float64)

	paragraphs, err := html.paragraphs()
	if err != nil {
		return nil, err
	}

	for _, elem := range paragraphs {

		text := elem.InnerHtml()
		if len(text) < minLen {
			continue
		}

		sc := 1.0
		sc += float64(len(strings.Split(text, ",")))
		sc += math.Min(float64(len(text)/100.0), 3.0)

		parent := elem.Parent()
		grandParent := parent.Parent()

		if _, found := candidates[parent]; !found {
			candidates[parent] = score(parent)
		}
		candidates[parent] = candidates[parent] + sc

		if grandParent != nil && grandParent.IsValid() {
			if _, found := candidates[grandParent]; !found {
				candidates[grandParent] = score(grandParent)
			}
			candidates[grandParent] = candidates[grandParent] + (sc / 2.0)
		}

		for candidate, score := range candidates {
			candidates[candidate] = score * (1 - linkDensity(candidate))
		}

	}

	return candidates, nil
}

func bestCandidate(candidates map[xml.Node]float64) (xml.Node, float64) {
	var best xml.Node
	var champ = math.Inf(-1)
	for node, score := range candidates {
		if score > champ {
			best = node
			champ = score
		}
	}

	return best, champ
}

func extend(doc *xml.XmlDocument, node xml.Node) {
	dup := node.DuplicateTo(doc, 1)
	if strings.ToLower(dup.Name()) != "p" || strings.ToLower(dup.Name()) != "div" {
		dup.SetName("div")
	}
}

func getArticle(candidates map[xml.Node]float64) *xml.XmlDocument {
	best, score := bestCandidate(candidates)
	threshold := math.Max(10.0, score*0.2)

	doc := xml.CreateEmptyDocument(xml.DefaultEncodingBytes, xml.DefaultEncodingBytes)

	for node := best.Parent().FirstChild(); node != nil && node.IsValid(); node = node.NextSibling() {
		if node == best {
			extend(doc, node)
		} else if candidates[node] >= threshold {
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
