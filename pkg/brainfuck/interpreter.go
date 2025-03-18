package brainfuck

import (
	"brainfuck_go/pkg/input_output"
	"errors"
	"math/big"
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

var commandSet = map[uint8]bool{
	'>': true,
	'<': true,
	'+': true,
	'-': true,
	'.': true,
	',': true,
	'[': true,
	']': true,
	'a': true,
	's': true,
	'm': true,
	'd': true,
	'r': true,
	'z': true,
	'w': true,
	'_': true,
}

var ztoc = map[int]uint8{
	0: '<', // no useful program ends with '<' so that's ok to use 0
	1: '>',
	2: '+',
	3: '-',
	4: '.',
	5: ',',
	6: '[',
	7: ']',
	/*
		8:  'a', // extended version of brainfuck
		9:  's',
		10: 'm',
		11: 'd',
		12: 'r',
		13: 'z',
		14: 'w',
	*/
}

func GetCodeFromInt(z *big.Int) []uint8 {
	z = (&big.Int{}).Set(z) // copy
	code := make([]uint8, 0)
	zero := big.NewInt(0)
	mod := &big.Int{}
	lenZtoC := big.NewInt(int64(len(ztoc)))
	for z.Cmp(zero) > 0 {
		mod = mod.Mod(z, lenZtoC)
		code = append(code, ztoc[int(mod.Int64())])
		z = z.Div(z, lenZtoC)
	}
	return code
}

var reduciblePair = map[string]bool{
	"<>": true,
	"><": true,
	"+-": true,
	"-+": true,
	"[]": true,
}

// EasyReducibleCode : simple test if there is a shorter equivalent code
func EasyReducibleCode(source []uint8) bool {
	if len(source) < 2 {
		return true
	}
	// every useful prog ends with printout
	if source[len(source)-1] != '.' {
		return true
	}
	// reducible pair
	for i := 0; i < len(source)-1; i++ {
		cc := string(source[i : i+1])
		if reduciblePair[cc] {
			return true
		}
	}
	// [] can be parsed
	stack := 0
	for i := 0; i < len(source); i++ {
		c := source[i]
		switch c {
		case ']':
			stack--
			if stack < 0 {
				return true
			}
		case '[':
			stack++
		}
	}

	return false
}

func parseCode(source []uint8) (tokens []token) {
	i := 0
	for i < len(source) {
		c := source[i]
		if c == '#' { // skip comment
			j := i + 1
			for j < len(source) && source[j] != '\n' {
				j++
			}
			i = j // next char
		} else if commandSet[c] {
			// look ahead for digits
			count := 0
			j := i + 1
			for j < len(source) && '0' <= source[j] && source[j] <= '9' {
				count = count*10 + int(source[j]-'0')
				j++
			}
			if count == 0 {
				count = 1
			}

			tokens = append(tokens, token{operation: c, count: count})
			i = j - 1 // next char
		}
		i++
	}
	return tokens
}

func makeJumpTable(tokens []token) (map[int]int, error) {
	jumpTable := make(map[int]int)
	biStack := make([]int, 0)
	var j int
	for i, token := range tokens {
		switch token.operation {
		case '[':
			biStack = append(biStack, i)
		case ']':
			if len(biStack) == 0 {
				return nil, nil
			}
			biStack, j = biStack[:len(biStack)-1], biStack[len(biStack)-1]
			jumpTable[i] = j
			jumpTable[j] = i
		}
	}
	return jumpTable, nil
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
