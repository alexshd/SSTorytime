package main

import (
	"fmt"
	"io"
	"os"

	"text2n4l-web/internal/analyzer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s -  (read from stdin)\n", os.Args[0])
		os.Exit(1)
	}

	var content []byte
	var err error

	if os.Args[1] == "-" {
		// Read from stdin
		content, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Read from file
		content, err = os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", os.Args[1], err)
			os.Exit(1)
		}
	}

	// Convert to N4L
	result := analyzer.ConvertTextToN4LWithResult(string(content))

	if result.Error != nil {
		fmt.Fprintf(os.Stderr, "Conversion error: %v\n", result.Error)
		os.Exit(1)
	}

	// Output results
	fmt.Print(result.N4LOutput)

	// Print statistics to stderr
	fmt.Fprintf(os.Stderr, "\n--- Conversion Statistics ---\n")
	fmt.Fprintf(os.Stderr, "Total fragments: %d\n", result.TotalFragments)
	fmt.Fprintf(os.Stderr, "Selected fragments: %d\n", result.SelectedFragments)
	fmt.Fprintf(os.Stderr, "Ambiguous lines: %d\n", len(result.AmbiguousIndices))

	if len(result.AmbiguousIndices) > 0 {
		fmt.Fprintf(os.Stderr, "Ambiguous line indices: %v\n", result.AmbiguousIndices)
	}
}
