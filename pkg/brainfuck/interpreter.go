package brainfuck

import (
	"brainfuck_go/pkg/input_output"
	"errors"
)

type Interpreter interface {
	Input() input_output.Input
	Output() input_output.Output
	Step() (bool, error)
}

type token struct {
	operation uint8
	count     int
}

type interpreter struct {
	dataPtr int
	codePtr int
	data    []uint8
	code    []token

	input  input_output.Input
	output input_output.Output

	jumpTable map[int]int
}

func NewInterpreter(dataLength int, code []uint8, input input_output.Input, output input_output.Output) (Interpreter, error) {
	tokens := parseCode(code)
	jumpTable, err := makeJumpTable(tokens)
	if err != nil {
		return nil, err
	}
	return &interpreter{
		dataPtr:   0,
		codePtr:   0,
		data:      make([]uint8, dataLength),
		code:      tokens,
		input:     input,
		output:    output,
		jumpTable: jumpTable,
	}, nil
}

func mod(a int, b int) int {
	if a < 0 {
		return (a % b) + b
	} else {
		return a % b
	}
}

// (a + b) % d
func addMod(a int, b int, d int) int {
	return mod(a+b, d)
}

func (i *interpreter) Step() (bool, error) {
	// fmt.Println(string(i.data))
	if i.codePtr >= len(i.code) {
		return true, nil // halt
	}
	t := i.code[i.codePtr]
	nextDataPtr := addMod(i.dataPtr, 1, len(i.data))
	switch t.operation {
	case '>':
		i.dataPtr = addMod(i.dataPtr, t.count, len(i.data))
	case '<':
		i.dataPtr = addMod(i.dataPtr, -t.count, len(i.data))
	case '+':
		i.data[i.dataPtr] = i.data[i.dataPtr] + uint8(t.count)
	case '-':
		i.data[i.dataPtr] = i.data[i.dataPtr] - uint8(t.count)
	case '.':
		i.output.Write(i.data[i.dataPtr])
	case ',':
		i.data[i.dataPtr] = i.input.Read()
	case '[':
		if i.data[i.dataPtr] == 0 {
			i.codePtr = i.jumpTable[i.codePtr]
		}
	case ']':
		if i.data[i.dataPtr] != 0 {
			i.codePtr = i.jumpTable[i.codePtr]
		}
	case 'a':
		i.data[i.dataPtr] = i.data[i.dataPtr] + i.data[nextDataPtr]
	case 's':
		i.data[i.dataPtr] = i.data[i.dataPtr] - i.data[nextDataPtr]
	case 'm':
		i.data[i.dataPtr] = i.data[i.dataPtr] * i.data[nextDataPtr]
	case 'd':
		if i.data[nextDataPtr] == 0 {
			return true, errors.New("division by zero")
		}
		i.data[i.dataPtr] = i.data[i.dataPtr] / i.data[nextDataPtr]
	case 'r':
		if i.data[nextDataPtr] == 0 {
			return true, errors.New("division by zero")
		}
		i.data[i.dataPtr] = i.data[i.dataPtr] % i.data[nextDataPtr]
	case 'z':
		i.data[i.dataPtr] = 0
	case 'w':
		i.data[i.dataPtr], i.data[nextDataPtr] = i.data[nextDataPtr], i.data[i.dataPtr]
	case '_':
		// noop
	}
	i.codePtr++
	return false, nil
}

func (i *interpreter) Input() input_output.Input {
	return i.input
}

func (i *interpreter) Output() input_output.Output {
	return i.output
}
