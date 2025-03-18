package universal_search

import (
	"brainfuck_go/pkg/brainfuck"
	"brainfuck_go/pkg/input_output"
	"fmt"
	"github.com/gosuri/uilive"
	"math/big"
	"sync"
)

func UniversalSearch(dataLength int, input []uint8, test func([]uint8) bool) (*big.Int, []uint8) {
	// zero := big.NewInt(0)
	one := big.NewInt(1)
	// two := big.NewInt(2)
	z := big.NewInt(0)
	type task struct {
		z *big.Int
		i brainfuck.Interpreter
	}
	type output struct {
		k string
		v []uint8
	}

	space := sync.Map{}
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
					z: (&big.Int{}).Set(z), // copy
					i: i,
				}
				k := t.z.Text(16)
				space.Store(k, t)
				counter++ // add one more task
			}
		}
		_, _ = fmt.Fprintf(writer, "number of running tasks: %d code size %d\n", counter, len(code))
		// run all tasks
		wg := &sync.WaitGroup{}

		outputList := make([]*output, 0)
		outputListMtx := &sync.Mutex{}

		space.Range(func(k1, t1 interface{}) bool {
			k := k1.(string)
			t := t1.(*task)
			wg.Add(1)
			go func(t *task) {
				defer wg.Done()
				//  linear version
				/*
					numSteps := (&big.Int{}).Exp(
						two,
						(&big.Int{}).Sub(z, t.z),
						nil,
					) // number of steps is 2^{z - t.z}
					for numSteps.Cmp(zero) > 0 {
						halt, err := t.i.Step()
						if halt || err != nil { // halt
							o := &output{
								k: k,
								v: nil,
							}
							if err == nil {
								v := t.i.Output().(input_output.StringOutput).String()
								if test(v) {
									o.v = v
								}
							}
							outputListMtx.Lock()
							outputList = append(outputList, o)
							outputListMtx.Unlock()
							break
						}
						numSteps = numSteps.Sub(numSteps, one)
					}
				*/
				// quadratic version - do fixed number of steps
				numSteps := 4096 //
				for i := 0; i < numSteps; i++ {

					halt, err := t.i.Step()
					if halt || err != nil { // halt
						o := &output{
							k: k,
							v: nil,
						}
						if err == nil {
							v := t.i.Output().(input_output.StringOutput).String()
							if test(v) {
								o.v = v
							}
						}
						outputListMtx.Lock()
						outputList = append(outputList, o)
						outputListMtx.Unlock()
						break
					}
				}
			}(t)
			return true // continue iteration
		})
		wg.Wait()

		for _, o := range outputList {
			space.Delete(o.k) // remove task
			if o.v != nil {
				z, _ = (&big.Int{}).SetString(o.k, 16)
				return z, o.v
			}
		}
		counter -= len(outputList) // reduce counter

		// next code
		z = z.Add(z, one)
	}
}
