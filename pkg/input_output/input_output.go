package input_output

import (
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
)

type Input interface {
	Read() uint8
}

type Output interface {
	Write(uint8)
}

type StringOutput interface {
	Output
	String() []uint8
}

type stdInput struct{}

func (i stdInput) Read() uint8 {
	// Switch stdin into raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Read a single byte
	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%c", b[0])
	return b[0]
}

func NewStdInput() Input {
	return stdInput{}
}

type stdOutput struct{}

func (o stdOutput) Write(b uint8) {
	fmt.Printf("%c", b)
}

func NewStdOutput() Output {
	return stdOutput{}
}

type stringInput struct {
	buffer   []uint8
	index    int
	boundary uint8
}

func (i *stringInput) Read() uint8 {
	if i.index >= len(i.buffer) {
		return i.boundary
	}
	b := i.buffer[i.index]
	i.index++
	return b
}

func NewStringInput(s []uint8) Input {
	return &stringInput{
		buffer:   s,
		index:    0,
		boundary: 0,
	}
}

type stringOutput struct {
	buffer []uint8
}

func (o *stringOutput) Write(b uint8) {
	o.buffer = append(o.buffer, b)
}

func (o *stringOutput) String() []uint8 {
	return o.buffer
}

func NewStringOutput() StringOutput {
	return &stringOutput{
		buffer: make([]uint8, 0),
	}
}
