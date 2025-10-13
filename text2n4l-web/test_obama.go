package main

import (
	"fmt"
	"io/ioutil"

	"text2n4l-web/internal/analyzer"
)

func main() {
	// Read the obama.dat file
	content, err := ioutil.ReadFile("../examples/example_data/obama.dat")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Convert using our algorithm
	output := analyzer.N4LSkeletonOutput("obama.dat", string(content), 100.0)

	// Write to temporary file
	err = ioutil.WriteFile("obama_web_test.n4l", []byte(output), 0o644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("Converted successfully. Output written to obama_web_test.n4l")

	// Show first few sentences - now without sanitize
	fragments := analyzer.FractionateTextFile(string(content))
	fmt.Printf("\nTotal fragments: %d\n", len(fragments))
	fmt.Printf("\nFirst 5 fragments:\n")
	for i := 0; i < 5 && i < len(fragments); i++ {
		fmt.Printf("@sen%d (len=%d): %s\n", i, len(fragments[i]), fragments[i])
	}
}
