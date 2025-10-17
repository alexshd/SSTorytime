package main

import (
	"flag"
	"fmt"
	"main/mazeexample/maze"
	"os"
)

func main() {
	jsonOutput := flag.Bool("json", false, "Output results as JSON")
	flag.Parse()

	if *jsonOutput {
		// JSON output mode
		result, err := maze.SolveMazeJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		jsonBytes, err := result.ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonBytes))
	} else {
		// Text output mode (default)
		fmt.Println("=== Maze Solver Example ===")
		fmt.Printf("Solving maze from %s to %s\n", maze.StartNode, maze.EndNode)
		fmt.Println("Using bidirectional wavefront search")
		fmt.Println()

		// Solve the maze and print results to stdout
		if err := maze.SolveMaze(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println("=== Maze solving complete ===")
	}

	os.Exit(0)
}
