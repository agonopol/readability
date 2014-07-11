package readability

import (
	"github.com/moovweb/gokogiri/xml"
)

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
