// Package analyzer provides text analysis and N4L conversion functionality.
// This is a clean refactor of the core analysis logic extracted from the monolithic SSTorytime package.
package analyzer

import (
	"math"
	"strings"
)

// Constants extracted from the original SSTorytime package
const (
	DUNBAR_30  = 45 // Approximate narrative range or sentences before new point/topic
	DUNBAR_150 = 150
	DUNBAR_15  = 15

	N_GRAM_MIN = 1
	N_GRAM_MAX = 6

	// Text classification types
	N1GRAM = 1
	N2GRAM = 2
	N3GRAM = 3
	LT128  = 4
	LT1024 = 5
	GT1024 = 6
)

// TextRank represents a ranked text fragment with significance and metadata
type TextRank struct {
	Significance float64 // Intentionality score
	Fragment     string  // The text fragment
	Order        int     // Original order in document
	Partition    int     // Coherence partition index
}

// Global variables for n-gram analysis (simplified from original)
var (
	STM_NGRAM_FREQ [N_GRAM_MAX]map[string]float64
	STM_NGRAM_LOCA [N_GRAM_MAX]map[string][]int
	STM_NGRAM_LAST [N_GRAM_MAX]map[string]int
)

// MemoryInit initializes the global n-gram tracking structures
func MemoryInit() {
	for i := N_GRAM_MIN; i < N_GRAM_MAX; i++ {
		STM_NGRAM_FREQ[i] = make(map[string]float64)
		STM_NGRAM_LOCA[i] = make(map[string][]int)
		STM_NGRAM_LAST[i] = make(map[string]int)
	}
}

// RunningIntentionality calculates the running intentionality score for a text fragment
// This is a simplified version of the original algorithm
func RunningIntentionality(t int, frag string) float64 {
	// Simplified scoring based on fragment length and position
	decayrate := float64(DUNBAR_30)
	words := strings.Fields(frag)
	if len(words) == 0 {
		return 0.0
	}

	// Basic intentionality heuristics
	score := float64(len(words)) * 0.1

	// Decay over time/position
	if t > 0 {
		score *= math.Exp(-float64(t) / decayrate)
	}

	return score
}

// AssessStaticIntent calculates static intentionality based on n-gram frequencies
func AssessStaticIntent(frag string, L int, frequencies [N_GRAM_MAX]map[string]float64, min int) float64 {
	words := strings.Fields(frag)
	if len(words) == 0 {
		return 0.0
	}

	score := 0.0
	for _, word := range words {
		// Simple frequency-based scoring
		if freq, exists := frequencies[1][strings.ToLower(word)]; exists {
			score += StaticIntentionality(L, word, freq)
		} else {
			score += 0.1 // Default score for unknown words
		}
	}

	return score
}

// StaticIntentionality calculates the intentionality of a term based on its frequency
func StaticIntentionality(L int, term string, frequency float64) float64 {
	if L == 0 || frequency == 0 {
		return 0.0
	}

	// Simple TF-IDF like calculation
	tf := frequency / float64(L)
	idf := math.Log(float64(L) / (frequency + 1))

	return tf * idf
}
