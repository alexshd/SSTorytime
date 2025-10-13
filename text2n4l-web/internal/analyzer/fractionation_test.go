package analyzer

import (
	"strings"
	"testing"
)

func TestCleanTextMatchesCLI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Abbreviations",
			input:    "Mr. Smith and Dr. Jones went to St. Mary's.",
			expected: "Mr Smith and Dr Jones went to St Mary's.#",
		},
		{
			name:     "Sentence boundaries",
			input:    "First sentence. Second sentence! Third sentence?",
			expected: "First sentence.# Second sentence!# Third sentence?#",
		},
		{
			name:     "Ellipsis",
			input:    "Wait... what?",
			expected: "Wait--- what?#",
		},
		{
			name:     "Paragraphs",
			input:    "First para.\n\nSecond para.",
			expected: "First para.#>>\n Second para.#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanText(tt.input)
			if result != tt.expected {
				t.Errorf("CleanText() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFractionateTextFileMatchesCLI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "Simple sentences",
			input: "First sentence. Second sentence. Third sentence.",
			expected: []string{
				"First sentence.",
				"Second sentence.",
				"Third sentence.",
			},
		},
		{
			name:  "Obama example",
			input: "OBAMA: My fellow citizens: I stand here today humbled by the task before us, grateful for the trust you have bestowed, mindful of the sacrifices borne by our ancestors. I thank President Bush for his service to our nation, as well as the generosity and cooperation he has shown throughout this transition.",
			// CLI separates on colon boundaries created by sentence markers
			expected: []string{
				"OBAMA:",
				"My fellow citizens:",
				"I stand here today humbled by the task before us, grateful for the trust you have bestowed, mindful of the sacrifices borne by our ancestors.",
				"I thank President Bush for his service to our nation, as well as the generosity and cooperation he has shown throughout this transition.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FractionateTextFile(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("FractionateTextFile() returned %d fragments, want %d", len(result), len(tt.expected))
				t.Logf("Got fragments: %v", result)
				t.Logf("Expected: %v", tt.expected)
				return
			}

			for i := range result {
				// Normalize whitespace for comparison
				gotNorm := strings.TrimSpace(result[i])
				wantNorm := strings.TrimSpace(tt.expected[i])
				if gotNorm != wantNorm {
					t.Errorf("Fragment %d: got %q, want %q", i, gotNorm, wantNorm)
				}
			}
		})
	}
}
