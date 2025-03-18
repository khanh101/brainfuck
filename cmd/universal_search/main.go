package main

import (
	"brainfuck_go/pkg/universal_search"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		panic("Usage: universal_search_factorization <num>")
	}

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	input := []uint8(strconv.Itoa(num))
	test := func(s []uint8) bool {
		factor, err := strconv.Atoi(string(s))
		if err != nil {
			return false
		}
		if factor < 2 || factor >= num {
			return false
		}
		return num%factor == 0
	}

	output, code := universal_search.UniversalSearch(64, input, test)
	fmt.Printf("factor: %s code %s\n", string(output), string(code))
}
