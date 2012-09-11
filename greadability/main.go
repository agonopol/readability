package main

import (
	"readability"
)

func main() {
	_, err := readability.Document("")
	if err != nil {
		panic(err)
	}
}
