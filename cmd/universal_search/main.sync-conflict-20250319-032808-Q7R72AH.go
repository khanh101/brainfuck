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

	i := brainf

	output := universal_search.UniversalSearch(64, input, test)
	fmt.Println(string(output))
}
