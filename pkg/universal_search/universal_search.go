package universal_search

import (
	"brainfuck_go/pkg/brainfuck"
	"brainfuck_go/pkg/input_output"
	"fmt"
	"github.com/gosuri/uilive"
	"math/big"
)

func UniversalSearch(dataLength int, input []uint8, test func([]uint8) bool) ([]uint8, []uint8) {
	// zero := big.NewInt(0)
	one := big.NewInt(1)
	// two := big.NewInt(2)
	z := big.NewInt(0)
	type task struct {
		c []uint8
		i brainfuck.Interpreter
	}

	space := make(map[string]*task)
	counter := 0
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	for {
		// new task
		code := brainfuck.GetCodeFromInt(z)
		if !brainfuck.EasyReducibleCode(code) {
			i, err := brainfuck.NewInterpreter(dataLength, code, input_output.NewStringInput(input), input_output.NewStringOutput())
			if err == nil {
				t := &task{
					c: code,
					i: i,
				}
				k := string(t.c)
				space[k] = t
				counter++ // add one more task
			}
		}
		_, _ = fmt.Fprintf(writer, "number of running tasks: %d code size %d\n", counter, len(code))
		// run all tasks
		toRemoveList := make([]string, 0)
		for k, t := range space {
			// quadratic version
			halt, err := t.i.Step()
			if halt || err != nil { // halt
				toRemoveList = append(toRemoveList, k)
				if err == nil {
					v := t.i.Output().(input_output.StringOutput).String()
					if test(v) {
						return v, t.c
					}
				}
			}
		}
		for _, k := range toRemoveList {
			delete(space, k)
		}
		counter -= len(toRemoveList) // reduce counter

		// next code
		z = z.Add(z, one)
	}
}
