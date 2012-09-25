package main

import (
	"readability"
	"fmt"
)

func main() {
	doc, err := readability.Document("http://localhost:6060/pkg/strconv")
	if err != nil {
		panic(err)
	}
	body, _ := doc.Content()
	fmt.Printf("%v\n", string(body))
}
