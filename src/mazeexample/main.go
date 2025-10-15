// Package main provides the command-line interface for the maze solver example.
//
// This program demonstrates solving a maze using bidirectional wavefront search
// with the SST (Semantic Space-Time) graph structure. The maze is represented as
// a directed graph where nodes are maze cells and edges represent valid moves.
//
// Usage:
//
//	go run main.go
//	# or after building:
//	./mazeexample
//
// The program will:
//  1. Build a graph from predefined maze paths
//  2. Run bidirectional search from start (maze_a7) to goal (maze_i6)
//  3. Display the wavefront expansion at each depth
//  4. Show the final solution path when found
//
// For more details on the algorithm, see the maze package documentation.
package main

import (
	"fmt"
	"os"

	"main/mazeexample/maze"
)

func main() {
	fmt.Println("=== Maze Solver Example ===")
	fmt.Println("Solving maze from maze_a7 to maze_i6")
	fmt.Println("Using bidirectional wavefront search")
	fmt.Println()

	// Solve the maze and print results to stdout
	if err := maze.SolveMaze(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("=== Maze solving complete ===")
	os.Exit(0)
}
