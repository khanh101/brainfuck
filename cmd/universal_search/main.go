package main

import (
	"brainfuck_go/pkg/universal_search"
	"fmt"
)

func main() {
	input := []uint8("4")
	test := func(s []uint8) bool {
		return "2x2" == string(s)
	}

	output := universal_search.UniversalSearch(30000, input, test)
	fmt.Println(string(output))
}
