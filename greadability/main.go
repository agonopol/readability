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
	body := doc.Body()
	fmt.Printf("%v\n", string(body))
}
