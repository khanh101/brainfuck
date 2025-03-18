package main

import (
	"brainfuck_go/pkg/brainfuck"
	"brainfuck_go/pkg/universal_search"
	"fmt"
)

func main() {
	input := []uint8("6")
	test := func(s []uint8) bool {
		return "2" == string(s)
	}

	z, output := universal_search.UniversalSearch(64, input, test)
	fmt.Println(string(brainfuck.GetCodeFromInt(z)))
	fmt.Println(string(output))
}
