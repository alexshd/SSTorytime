package analyzer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGoldenFiles performs regression testing using golden files
func TestGoldenFiles(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple_text",
			input:    "testdata/sample.txt",
			expected: "testdata/golden/simple_text.n4l",
		},
		{
			name:     "markdown_sample",
			input:    "testdata/markdown_sample.md",
			expected: "testdata/golden/markdown_sample.n4l",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read input file
			inputData, err := os.ReadFile(tt.input)
			if err != nil {
				t.Fatalf("Failed to read input file %s: %v", tt.input, err)
			}

			// Generate N4L output
			result := N4LSkeletonOutput(tt.input, string(inputData), 0.8)

			// Check if we should update golden files (set UPDATE_GOLDEN=1)
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				// Ensure golden directory exists
				dir := filepath.Dir(tt.expected)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatalf("Failed to create golden directory: %v", err)
				}

				// Write new golden file
				if err := os.WriteFile(tt.expected, []byte(result), 0o644); err != nil {
					t.Fatalf("Failed to write golden file %s: %v", tt.expected, err)
				}
				t.Logf("Updated golden file: %s", tt.expected)
				return
			}

			// Read expected output
			expectedData, err := os.ReadFile(tt.expected)
			if err != nil {
				t.Fatalf("Failed to read golden file %s: %v", tt.expected, err)
			}

			expected := string(expectedData)

			// Compare outputs
			if !compareN4LOutput(result, expected) {
				t.Errorf("Output doesn't match golden file for %s", tt.name)
				t.Errorf("Expected:\n%s", expected)
				t.Errorf("Got:\n%s", result)

				// Write actual output to temp file for debugging
				tempFile := filepath.Join("testdata", tt.name+"_actual.n4l")
				os.WriteFile(tempFile, []byte(result), 0o644)
				t.Errorf("Actual output written to: %s", tempFile)
			}
		})
	}
}

// compareN4LOutput compares two N4L outputs, ignoring minor differences
func compareN4LOutput(actual, expected string) bool {
	// Normalize line endings
	actual = strings.ReplaceAll(actual, "\r\n", "\n")
	expected = strings.ReplaceAll(expected, "\r\n", "\n")

	// Split into lines for comparison
	actualLines := strings.Split(strings.TrimSpace(actual), "\n")
	expectedLines := strings.Split(strings.TrimSpace(expected), "\n")

	if len(actualLines) != len(expectedLines) {
		return false
	}

	for i, actualLine := range actualLines {
		expectedLine := expectedLines[i]

		// Trim whitespace for comparison
		actualLine = strings.TrimSpace(actualLine)
		expectedLine = strings.TrimSpace(expectedLine)

		// Skip comment lines that might have dynamic content (like line counts)
		if strings.HasPrefix(actualLine, "# There were approximately") &&
			strings.HasPrefix(expectedLine, "# There were approximately") {
			continue
		}

		if actualLine != expectedLine {
			return false
		}
	}

	return true
}

// BenchmarkGoldenFiles benchmarks golden file processing
func BenchmarkGoldenFiles(b *testing.B) {
	// Read sample input
	inputData, err := os.ReadFile("testdata/sample.txt")
	if err != nil {
		b.Fatal("Failed to read sample input:", err)
	}
	input := string(inputData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		N4LSkeletonOutput("sample.txt", input, 0.8)
	}
}

// TestUpdateGoldenFiles is a helper test for updating golden files
func TestUpdateGoldenFiles(t *testing.T) {
	if os.Getenv("UPDATE_GOLDEN") != "1" {
		t.Skip("Set UPDATE_GOLDEN=1 to update golden files")
	}

	// This will run the golden file tests in update mode
	TestGoldenFiles(t)
}
