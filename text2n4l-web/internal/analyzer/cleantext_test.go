package analyzer

import (
	"strings"
	"testing"
)

func TestCleanText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTML tags removed",
			input:    "Hello <b>world</b> and <div>stuff</div>",
			expected: "Hello :\n world :\n and :\n stuff :\n",
		},
		{
			name:     "Abbreviations preserved",
			input:    "Mr. Smith and Dr. Jones went to St. Mary's.",
			expected: "Mr Smith and Dr Jones went to St Mary's. #",
		},
		{
			name:     "Sentence boundaries marked",
			input:    "First sentence. Second sentence! Third sentence?",
			expected: "First sentence. #Second sentence! #Third sentence? #",
		},
		{
			name:     "Double newlines become paragraph markers",
			input:    "First paragraph.\n\nSecond paragraph.",
			expected: "First paragraph. #>> Second paragraph. #",
		},
		{
			name:     "Ellipsis converted",
			input:    "Wait... what happened?",
			expected: "Wait--- what happened? #",
		},
		{
			name:     "Em-dash converted",
			input:    "The answer—quite simply—is no.",
			expected: "The answer, quite simply, is no. #",
		},
		{
			name:     "Brackets removed",
			input:    "Some [bracketed] text here.",
			expected: "Some bracketed text here. #",
		},
		{
			name:     "Multiple newlines consolidated",
			input:    "Line one.\n\n\nLine two.",
			expected: "Line one. #>> Line two. #",
		},
		{
			name:     "Real Obama example",
			input:    "OBAMA: My fellow citizens:\n\nI stand here today humbled by the task before us.",
			expected: "OBAMA: My fellow citizens:>> I stand here today humbled by the task before us. #",
		},
		{
			name:     "Complex mixed content",
			input:    "Mr. Jones said, \"Hello!\"\n\nDr. Smith replied... \"Yes?\"",
			expected: "Mr Jones said, \"Hello! #\">> Dr Smith replied--- \"Yes? #\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanText(tt.input)
			// Normalize for comparison (spaces may vary)
			resultNorm := strings.Join(strings.Fields(result), " ")
			expectedNorm := strings.Join(strings.Fields(tt.expected), " ")

			if resultNorm != expectedNorm {
				t.Errorf("CleanText() failed\nInput:    %q\nGot:      %q\nExpected: %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanTextPreservesImportantWhitespace(t *testing.T) {
	input := "Sentence one.\n\nSentence two."
	result := CleanText(input)

	// Should contain >> marker
	if !strings.Contains(result, ">>") {
		t.Errorf("CleanText() should preserve paragraph markers (>>), got: %q", result)
	}

	// Should contain # markers
	if !strings.Contains(result, "#") {
		t.Errorf("CleanText() should add sentence markers (#), got: %q", result)
	}
}

// Benchmark the original implementation
func BenchmarkCleanText(b *testing.B) {
	input := `Text of President Barack Obama's first inaugural address on Tuesday, as
prepared for delivery and released by the Presidential Inaugural
Committee.

OBAMA: My fellow citizens:

I stand here today humbled by the task before us, grateful for the trust
you have bestowed, mindful of the sacrifices borne by our ancestors. I
thank President Bush for his service to our nation, as well as the
generosity and cooperation he has shown throughout this transition.

Forty-four Americans have now taken the presidential oath. The words
have been spoken during rising tides of prosperity and the still waters
of peace. Yet, every so often the oath is taken amidst gathering clouds
and raging storms. At these moments, America has carried on not simply
because of the skill or vision of those in high office, but because we
the people have remained faithful to the ideals of our forebears, and
true to our founding documents.

So it has been. So it must be with this generation of Americans.

That we are in the midst of crisis is now well understood. Our nation
is at war, against a far-reaching network of violence and hatred. Our
economy is badly weakened... a consequence of greed and irresponsibility
on the part of some—but also our collective failure to make hard choices
and prepare the nation for a new age.`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CleanText(input)
	}
}

func BenchmarkCleanTextSmall(b *testing.B) {
	input := "Mr. Smith said hello. Dr. Jones replied!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CleanText(input)
	}
}

func BenchmarkCleanTextLarge(b *testing.B) {
	// Create a large input by repeating the medium text
	medium := `Text of President Barack Obama's first inaugural address on Tuesday, as
prepared for delivery and released by the Presidential Inaugural
Committee.

OBAMA: My fellow citizens:

I stand here today humbled by the task before us, grateful for the trust
you have bestowed, mindful of the sacrifices borne by our ancestors.`

	input := strings.Repeat(medium, 100) // ~30KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CleanText(input)
	}
}
