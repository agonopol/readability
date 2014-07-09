package readability

import (
	"strings"
	"github.com/moovweb/gokogiri/xml"
	"math"
)

func weight(elem xml.Node) float64 {
      weight := 0.0
	  class :=  elem.Attr("class")
	  if class != "" {
		  if regexes["negativeRe"].Match([]byte(class)) {
		  	weight -= 25.0
		  } else if regexes["positiveRe"].Match([]byte(class)) {
          	weight += 25.0
		  }
	  }
        
	  id :=  elem.Attr("id")
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

func candidates(html *htmlParser, minLen int) (map[xml.Node]float64, error) {
	
	candidates := make(map[xml.Node] float64)
	
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
		sc += math.Min(float64(len(text) / 100.0), 3.0)
		
		
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

func article(candidates map[xml.Node]float64) {
	best, score := bestCandidate(candidates)
	threshold = math.Max(10.0, score * 0.2)
	node := best.Parent().FirstChild() 
	
	while (node != nil && node.) {
		
		node = node.NextSibling()
	}
}

// def get_article(candidates, best_candidate)
//   # Now that we have the top candidate, look through its siblings for content that might also be related.
//   # Things like preambles, content split by ads that we removed, etc.
//
//   sibling_score_threshold = [10, best_candidate[:content_score] * 0.2].max
//   output = Nokogiri::XML::Node.new('div', @html)
//   best_candidate[:elem].parent.children.each do |sibling|
//     append = false
//     append = true if sibling == best_candidate[:elem]
//     append = true if candidates[sibling] && candidates[sibling][:content_score] >= sibling_score_threshold
//
//     if sibling.name.downcase == "p"
//       link_density = get_link_density(sibling)
//       node_content = sibling.text
//       node_length = node_content.length
//
//       append = if node_length > 80 && link_density < 0.25
//         true
//       elsif node_length < 80 && link_density == 0 && node_content =~ /\.( |$)/
//         true
//       end
//     end
//
//     if append
//       sibling_dup = sibling.dup # otherwise the state of the document in processing will change, thus creating side effects
//       sibling_dup.name = "div" unless %w[div p].include?(sibling.name.downcase)
//       output << sibling_dup
//     end
//   end
//
//   output
// end
