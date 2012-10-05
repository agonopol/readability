package main

import (
	"readability"
	// "fmt"
)

func main() {
	doc, err := readability.Document("http://en.wikipedia.org/wiki/The_Reluctant_Fundamentalist")
	if err != nil {
		panic(err)
	}
	body, _ := doc.Content()
	// fmt.Printf("%v\n", string(body))
	if body == nil {
		panic(nil)
	}
}
