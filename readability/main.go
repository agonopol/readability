package main

import (
	"readability"
	"fmt"
)

func main() {
	doc, err := readability.Document("http://en.wikipedia.org/wiki/The_Reluctant_Fundamentalist")
	if err != nil {
		panic(err)
	}
	body, err := doc.Content()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", string(body))
}
