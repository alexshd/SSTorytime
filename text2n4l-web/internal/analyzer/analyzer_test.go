package analyzer

import (
	"os"
	"strings"
	"testing"
)

func TestConvertTextToN4L(t *testing.T) {
	// Test with simple text
	input := "This is a test sentence. This is another test sentence. A third sentence for testing."

	result := ConvertTextToN4LWithResult(input)

	if result.Error != nil {
		t.Fatalf("Conversion failed: %v", result.Error)
	}

	if result.N4LOutput == "" {
		t.Error("N4L output is empty")
	}

	if result.TotalFragments == 0 {
		t.Error("No fragments detected")
	}

	// Check that output contains N4L markers
	if !strings.Contains(result.N4LOutput, "+") && !strings.Contains(result.N4LOutput, "?") {
		t.Error("N4L output doesn't contain expected markers (+ or ?)")
	}

	// Check that header is present
	if !strings.Contains(result.N4LOutput, "# N4L Narrative Analysis") {
		t.Error("N4L output missing header")
	}
}

func TestConvertEmptyText(t *testing.T) {
	result := ConvertTextToN4LWithResult("")

	if result.Error != nil {
		t.Fatalf("Empty text conversion failed: %v", result.Error)
	}

	if result.TotalFragments != 0 {
		t.Error("Empty text should have 0 fragments")
	}
}

func TestConvertTestFile(t *testing.T) {
	// Try to read and convert a test file
	content, err := os.ReadFile("../../testdata/pass_1.in")
	if err != nil {
		t.Skip("Test file not available:", err)
		return
	}

	result := ConvertTextToN4LWithResult(string(content))

	if result.Error != nil {
		t.Fatalf("Test file conversion failed: %v", result.Error)
	}

	if result.N4LOutput == "" {
		t.Error("N4L output is empty for test file")
	}

	if result.TotalFragments == 0 {
		t.Error("No fragments detected in test file")
	}

	t.Logf("Converted %d fragments, selected %d", result.TotalFragments, result.SelectedFragments)
	if len(result.AmbiguousIndices) > 0 {
		t.Logf("Found %d ambiguous lines", len(result.AmbiguousIndices))
	}
}

func TestFractionateTextFile(t *testing.T) {
	input := "First sentence. Second sentence! Third sentence? Fourth line\nFifth line."

	fragments := FractionateTextFile(input)

	if len(fragments) == 0 {
		t.Error("No fragments generated")
	}

	// Should have multiple fragments
	if len(fragments) < 3 {
		t.Errorf("Expected at least 3 fragments, got %d", len(fragments))
	}

	// Check that fragments are not empty
	for i, frag := range fragments {
		if strings.TrimSpace(frag) == "" {
			t.Errorf("Fragment %d is empty", i)
		}
	}
}

func TestExtractIntentionalTokens(t *testing.T) {
	fragments := []string{
		"This is a test sentence.",
		"This is another test sentence.",
		"A completely different sentence.",
	}

	frequencies := ExtractIntentionalTokens(fragments)

	// Check that unigrams were counted
	if len(frequencies[N1GRAM]) == 0 {
		t.Error("No unigrams extracted")
	}

	// Word "this" should appear twice
	if freq, exists := frequencies[N1GRAM]["this"]; !exists || freq != 2.0 {
		t.Errorf("Expected 'this' to have frequency 2.0, got %f", freq)
	}

	// Check bigrams
	if len(frequencies[N2GRAM]) == 0 {
		t.Error("No bigrams extracted")
	}
}
