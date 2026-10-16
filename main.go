package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run . <input-file> <output-file>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %q: %v\n", inputPath, err)
		os.Exit(1)
	}

	result := Process(string(data))
	if err := os.WriteFile(outputPath, []byte(result), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %q: %v\n", outputPath, err)
		os.Exit(1)
	}
}
