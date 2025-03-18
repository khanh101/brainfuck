package main

import (
	"brainfuck_go/pkg/universal_search"
	"fmt"
)

func main() {
	input := []uint8("6")
	test := func(s []uint8) bool {
		return "2" == string(s)
	}

	output, code := universal_search.UniversalSearch(64, input, test)
	fmt.Println(string(code))
	fmt.Println(string(output))
}
