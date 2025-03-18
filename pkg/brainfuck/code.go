package brainfuck

type PartialCodeIterator interface {
	Code() []uint8
	Next() []PartialCodeIterator
}

func NewPartialCodeIterator() PartialCodeIterator {
	return &codeIterator{code: nil}
}

type codeIterator struct {
	code []uint8
}

func (i *codeIterator) Code() []uint8 {
	return i.code
}

var reduciblePair = map[string]bool{
	"<>": true,
	"><": true,
	"+-": true,
	"-+": true,
	"[]": true,
}

func (i *codeIterator) Next() []PartialCodeIterator {
	nextCodeList := make([]PartialCodeIterator, 0)
	for next := range reducedCommandSet {
		if len(i.code) >= 1 {
			prev := i.code[len(i.code)-1]
			pair := string([]uint8{prev, next})
			if reduciblePair[pair] {
				continue // skip irreducible pair
			}
		}
		nextCode := append([]uint8{}, i.code...)
		nextCode = append(nextCode, next)
		nextCodeList = append(nextCodeList, &codeIterator{code: nextCode})
	}
	return nextCodeList
}

// FinalizeCode :
func FinalizeCode(code []uint8) []uint8 {
	// code too short
	if len(code) < 2 {
		return nil
	}
	// every useful prog ends with printout
	code = append([]uint8{}, code...)
	code = append(code, '.')

	// reducible pair
	for i := 0; i < len(code)-1; i++ {
		cc := string(code[i : i+1])
		if reduciblePair[cc] {
			return nil
		}
	}
	// [] can be parsed
	stack := 0
	for i := 0; i < len(code); i++ {
		c := code[i]
		switch c {
		case ']':
			stack--
			if stack < 0 {
				return nil
			}
		case '[':
			stack++
		}
	}
	if stack != 0 {
		return nil
	}

	return code
}

var reducedCommandSet = map[uint8]bool{
	'>': true,
	'<': true,
	'+': true,
	'-': true,
	'.': true,
	',': true,
	'[': true,
	']': true,
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
