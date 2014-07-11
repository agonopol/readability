package readability

import (
	_ "fmt"
	"github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
	"math"
	"strings"
)

type Candidate struct {
	node  xml.Node
	score float64
}

func newCadidate(elem xml.Node) *Candidate {
	this := &Candidate{elem, 0}
	switch strings.ToLower(elem.Name()) {

	case "div":
		this.score = this.weight() + 5.0
		break
	case "blockquote":
		this.score = this.weight() + 3.0
		break
	case "form":
		this.score = this.weight() - 3.0
		break
	case "th":
		this.score = this.weight() - 5.0
		break
	default:
		this.score = this.weight()
		break
	}
	return this
}
func (this *Candidate) weight() float64 {
	weight := 0.0
	class := this.node.Attr("class")
	if class != "" {
		if regexes["negativeRe"].Match([]byte(class)) {
			weight -= 25.0
		} else if regexes["positiveRe"].Match([]byte(class)) {
			weight += 25.0
		}
	}

	id := this.node.Attr("id")
	if id != "" {
		if regexes["negativeRe"].Match([]byte(id)) {
			weight -= 25.0
		} else if regexes["positiveRe"].Match([]byte(id)) {
			weight += 25.0
		}
	}
	return weight
}

func getCandidates(doc *html.HtmlDocument, minLen int) (map[string]*Candidate, error) {

	candidates := make(map[string]*Candidate)

	paragraphs, err := doc.Search(`//p|//td`)
	if err != nil {
		return nil, err
	}

	for _, elem := range paragraphs {
		text := elem.Content()

		if len(text) < minLen {
			continue
		}

		sc := 1.0
		sc += float64(len(strings.Split(text, ",")))
		sc += math.Min(float64(len(text)/100.0), 3.0)

		parent := elem.Parent()
		grandParent := parent.Parent()

		if _, found := candidates[parent.String()]; !found {
			candidates[parent.String()] = newCadidate(parent)
		}
		candidates[parent.String()].score += sc

		if grandParent != nil && grandParent.IsValid() {
			if _, found := candidates[grandParent.String()]; !found {
				candidates[grandParent.String()] = newCadidate(grandParent)
			}
			candidates[grandParent.String()].score += (sc / 2.0)
		}

		for _, candidate := range candidates {
			candidate.score = (candidate.score * (1 - linkDensity(candidate.node)))
		}

	}

	return candidates, nil
}

func bestCandidate(candidates map[string]*Candidate) *Candidate {
	var best *Candidate
	champ := math.Inf(-1)
	for _, candidate := range candidates {
		if candidate.score > champ {
			best = candidate
			champ = candidate.score
		}
	}
	return best
}
