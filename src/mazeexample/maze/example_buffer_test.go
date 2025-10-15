package maze

import (
	"bytes"
	"strings"
	"testing"
)

// TestSolveMazeBufferedOutput demonstrates testing with buffered output.
// This test shows how to use bytes.Buffer to capture maze solving output
// for testing and analysis instead of printing to stdout.
func TestSolveMazeBufferedOutput(t *testing.T) {
	var buf bytes.Buffer

	// Capture output to buffer
	SolveMazeWithOutput(&buf)

	output := buf.String()

	// Verify the output contains expected content
	if !strings.Contains(output, "-- T R E E --") {
		t.Error("Expected to find solution marker in output")
	}

	if !strings.Contains(output, "maze_a7") {
		t.Error("Expected to find start node in output")
	}

	if !strings.Contains(output, "maze_i6") {
		t.Error("Expected to find end node in output")
	}

	// Verify we found the connection point
	if !strings.Contains(output, "waves impinge") {
		t.Error("Expected to find wavefront collision message")
	}

	// Count solutions
	solutionCount := strings.Count(output, "-- T R E E --")
	if solutionCount != 1 {
		t.Errorf("Expected 1 solution, found %d", solutionCount)
	}
}
