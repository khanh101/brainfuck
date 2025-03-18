package main

import (
	"brainfuck_go/pkg/brainfuck"
	"brainfuck_go/pkg/input_output"
	"fmt"
	"io"
	"os"
	"strconv"
)

// readFileToBytes reads the content of the specified file and returns it as a byte slice.
func readFileToBytes(filename string) ([]byte, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Read the file content into a byte slice
	data := make([]byte, fileInfo.Size())
	if _, err := io.ReadFull(file, data); err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return data, nil
}

func main() {
	if len(os.Args) < 3 {
		panic("Usage: brainfuck <data_length> <code_file>")
	}
	dataLength, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	codeFile := os.Args[2]

	code, err := readFileToBytes(codeFile)
	if err != nil {
		panic(err)
	}

	i, err := brainfuck.NewInterpreter(dataLength, code, input_output.NewStdInput(), input_output.NewStdOutput())
	if err != nil {
		panic(err)
	}
	for {
		if i.Step() {
			break
		}
	}
}
