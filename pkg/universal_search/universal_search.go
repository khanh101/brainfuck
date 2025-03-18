package universal_search

import (
	"brainfuck_go/pkg/brainfuck"
	"brainfuck_go/pkg/input_output"
	"fmt"
	"time"
)

func UniversalSearch(dataLength int, input []uint8, test func([]uint8) bool) ([]uint8, []uint8) {
	type task struct {
		c []uint8
		i brainfuck.Interpreter
	}

	space := make(map[string]*task)
	counter := 0

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var codeIter brainfuck.PartialCodeIterator = nil
	queue := make([]brainfuck.PartialCodeIterator, 0)
	queue = append(queue, brainfuck.NewPartialCodeIterator())

	lastCode := []uint8{}
	lastCodeSize := 0
	for {
		select {
		case <-ticker.C: // print every tick
			fmt.Printf("running tasks: %d code size %d last code %s\n", counter, lastCodeSize, string(lastCode))
		default:
			break
		}

		// new task
		codeIter, queue = queue[0], queue[1:]
		queue = append(queue, codeIter.Next()...)
		code := codeIter.Code()
		code = brainfuck.FinalizeCode(code)
		if code == nil {
			continue
		}
		i, err := brainfuck.NewInterpreter(dataLength, code, input_output.NewStringInput(input), input_output.NewStringOutput())
		if err != nil {
			continue
		}
		lastCodeSize = len(code)
		lastCode = code
		t := &task{
			c: code,
			i: i,
		}
		k := string(t.c)
		space[k] = t
		counter++ // add one more task

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
	}
}
