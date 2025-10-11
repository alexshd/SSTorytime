package main

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name          string
		args          []string
		expectedPerc  float64
		expectedFile  string
		expectedError bool
	}{
		{
			name:          "default percentage with file",
			args:          []string{"cmd", testFile},
			expectedPerc:  50.0,
			expectedFile:  testFile,
			expectedError: false,
		},
		{
			name:          "custom percentage with file",
			args:          []string{"cmd", "-percentage", "75.5", testFile},
			expectedPerc:  75.5,
			expectedFile:  testFile,
			expectedError: false,
		},
		{
			name:          "zero percentage",
			args:          []string{"cmd", "-percentage", "0", testFile},
			expectedPerc:  0.0,
			expectedFile:  testFile,
			expectedError: false,
		},
		{
			name:          "100 percentage",
			args:          []string{"cmd", "-percentage", "100", testFile},
			expectedPerc:  100.0,
			expectedFile:  testFile,
			expectedError: false,
		},
		{
			name:          "negative percentage",
			args:          []string{"cmd", "-percentage", "-10", testFile},
			expectedError: true,
		},
		{
			name:          "percentage over 100",
			args:          []string{"cmd", "-percentage", "150", testFile},
			expectedError: true,
		},
		{
			name:          "no filename",
			args:          []string{"cmd", "-percentage", "50"},
			expectedError: true,
		},
		{
			name:          "multiple filenames",
			args:          []string{"cmd", testFile, "another.txt"},
			expectedError: true,
		},
		{
			name:          "non-existent file",
			args:          []string{"cmd", "nonexistent.txt"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset os.Args for each test
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()
			os.Args = tt.args

			// Reset flag state
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			config, err := parseArgs()

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if config.percentage != tt.expectedPerc {
				t.Errorf("Expected percentage %.2f, got %.2f", tt.expectedPerc, config.percentage)
			}

			if config.filename != tt.expectedFile {
				t.Errorf("Expected filename %s, got %s", tt.expectedFile, config.filename)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	config := &Config{
		percentage: 75.0,
		filename:   "test.txt",
	}

	if config.percentage != 75.0 {
		t.Errorf("Expected percentage 75.0, got %.2f", config.percentage)
	}

	if config.filename != "test.txt" {
		t.Errorf("Expected filename test.txt, got %s", config.filename)
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello (world)", "hello [world]"},
		{"(test)", "[test]"},
		{"no parentheses", "no parentheses"},
		{"()", "[]"},
		{"((nested))", "[[nested]]"},
		{"mixed (content) here (too)", "mixed [content] here [too]"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Sanitize(tt.input)
			if result != tt.expected {
				t.Errorf("Sanitize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSpliceSet(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: "",
		},
		{
			name:     "single element",
			input:    []string{"hello"},
			expected: "hello",
		},
		{
			name:     "multiple elements",
			input:    []string{"hello", "world", "test"},
			expected: "hello, world, test",
		},
		{
			name:     "two elements",
			input:    []string{"first", "second"},
			expected: "first, second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SpliceSet(tt.input)
			if result != tt.expected {
				t.Errorf("SpliceSet(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPartName(t *testing.T) {
	tests := []struct {
		partition int
		file      string
		context   string
		expected  string
	}{
		{1, "test.txt", "some context", "part 1 of test.txt with some context"},
		{0, "document", "", "part 0 of document with "},
		{5, "analysis.doc", "multiple, words", "part 5 of analysis.doc with multiple, words"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := PartName(tt.partition, tt.file, tt.context)
			if result != tt.expected {
				t.Errorf("PartName(%d, %q, %q) = %q, want %q",
					tt.partition, tt.file, tt.context, result, tt.expected)
			}
		})
	}
}

func TestAddIntentionalContextWithBufio(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_output.txt")

	// Create file and buffered writer
	fp, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer fp.Close()

	writer := bufio.NewWriter(fp)
	defer func() {
		if err := writer.Flush(); err != nil {
			t.Errorf("Error flushing buffer: %v", err)
		}
	}()

	// Test data
	context := []string{"test context 1", "test context 2", "another context"}

	// Call function
	AddIntentionalContext(writer, context)

	// Flush to ensure data is written
	if err := writer.Flush(); err != nil {
		t.Fatalf("Failed to flush writer: %v", err)
	}

	// Read back and verify
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	contentStr := string(content)

	// Check if all context items are present
	for _, ctx := range context {
		if !strings.Contains(contentStr, ctx) {
			t.Errorf("Expected context %q not found in output", ctx)
		}
	}

	// Check format - should contain quotes and proper spacing
	expectedLines := 3
	actualLines := strings.Count(contentStr, "\n")
	if actualLines != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, actualLines)
	}

	// Check that it contains the proper formatting
	if !strings.Contains(contentStr, "\"") {
		t.Error("Expected output to contain quotes")
	}
}
