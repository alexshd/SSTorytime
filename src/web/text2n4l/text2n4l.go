// Package text2n4l provides text-to-N4L conversion logic for both web and CLI use.
package text2n4l

import (
	SST "SSTorytime"
	"sort"
	"strings"
)

// Text2N4LResult holds the original and generated lines, and ambiguous indices
type Text2N4LResult struct {
	Original  []string
	Generated []string
	Ambiguous []int // line indices in Generated
}

// selectByRunningIntent ranks sentences by their running intentionality score
func selectByRunningIntent(psf [][][]string, L int, percentage float64) []SST.TextRank {
	const coherence_length = SST.DUNBAR_30
	var sentences []SST.TextRank
	var sentence_counter int
	for p := range psf {
		for s := 0; s < len(psf[p]); s++ {
			score := 0.0
			var sb strings.Builder
			for f := 0; f < len(psf[p][s]); f++ {
				score += SST.RunningIntentionality(sentence_counter, psf[p][s][f])
				sb.WriteString(psf[p][s][f])
				if f < len(psf[p][s])-1 {
					sb.WriteString(", ")
				}
			}
			var this SST.TextRank
			this.Fragment = sb.String()
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences, this)
			sentence_counter++
		}
	}
	return orderAndRank(sentences, percentage)
}

// selectByStaticIntent ranks sentences by their static intentionality score
func selectByStaticIntent(psf [][][]string, L int, percentage float64) []SST.TextRank {
	const coherence_length = SST.DUNBAR_30
	var sentences []SST.TextRank
	var sentence_counter int
	for p := range psf {
		for s := 0; s < len(psf[p]); s++ {
			score := 0.0
			var sb strings.Builder
			for f := 0; f < len(psf[p][s]); f++ {
				score += SST.AssessStaticIntent(psf[p][s][f], L, SST.STM_NGRAM_FREQ, 1)
				sb.WriteString(psf[p][s][f])
				if f < len(psf[p][s])-1 {
					sb.WriteString(", ")
				}
			}
			var this SST.TextRank
			this.Fragment = sb.String()
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences, this)
			sentence_counter++
		}
	}
	return orderAndRank(sentences, percentage)
}

// orderAndRank orders sentences by significance and returns the top percentage
func orderAndRank(sentences []SST.TextRank, percentage float64) []SST.TextRank {
	var selections []SST.TextRank
	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})
	limit := int((percentage / 100.0) * float64(len(sentences)))
	if limit == 0 && len(sentences) > 0 {
		limit = 1
	}
	for i := 0; i < limit; i++ {
		selections = append(selections, sentences[i])
	}
	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})
	return selections
}

// mergeSelections merges two slices of TextRank, preserving order and uniqueness
func mergeSelections(one []SST.TextRank, two []SST.TextRank) []SST.TextRank {
	var merge []SST.TextRank
	alreadySelected := make(map[int]bool)
	for i := range one {
		merge = append(merge, one[i])
		alreadySelected[one[i].Order] = true
	}
	for i := range two {
		if !alreadySelected[two[i].Order] {
			merge = append(merge, two[i])
		}
	}
	sort.Slice(merge, func(i, j int) bool {
		return merge[i].Order < merge[j].Order
	})
	return merge
}

// ConvertTextToN4LResult runs the real N4L conversion and returns generated lines and ambiguous indices.
func ConvertTextToN4LResult(inputPath string, percentage float64) (*Text2N4LResult, error) {
	SST.MemoryInit()
	psf, L := SST.FractionateTextFile(inputPath)
	ranking1 := selectByRunningIntent(psf, L, percentage)
	ranking2 := selectByStaticIntent(psf, L, percentage)
	selection := mergeSelections(ranking1, ranking2)

	const minN = 1
	const maxN = 3
	f, s, ff, ss := SST.ExtractIntentionalTokens(L, selection, minN, maxN)

	// Generate output lines in N4L DSL format
	var origLines []string
	var genLines []string
	var ambiguous []int

	// Read original file for Original field
	origContent := SST.ReadFile(inputPath)
	origLines = strings.Split(origContent, "\n")

	// Generate N4L DSL output similar to WriteOutput in text2N4L.go
	for i, sel := range selection {
		// Format as N4L DSL with @sen prefix
		line := "@sen" + strings.Join([]string{sel.Fragment}, " ")
		genLines = append(genLines, line)

		// Mark lines with ambiguous context as needing attention
		partIdx := sel.Partition
		if partIdx < len(s) && len(s[partIdx]) > 0 {
			ambiguous = append(ambiguous, i)
		}

		// Add contextual information
		if partIdx < len(f) {
			for _, ctx := range f[partIdx] {
				genLines = append(genLines, "  # INTENT "+ctx)
			}
		}
	}

	// Add global context
	for _, ctx := range ff {
		genLines = append(genLines, " # "+ctx)
	}
	for _, ctx := range ss {
		genLines = append(genLines, "  # "+ctx)
	}

	return &Text2N4LResult{Original: origLines, Generated: genLines, Ambiguous: ambiguous}, nil
}

func containsAmbiguity(s string) bool {
	// Simple heuristic: mark lines with '?' as ambiguous
	for _, c := range s {
		if c == '?' {
			return true
		}
	}
	return false
}

// ConvertTextToN4L returns only the generated lines for compatibility
func ConvertTextToN4L(inputPath string, percentage float64) ([]string, error) {
	result, err := ConvertTextToN4LResult(inputPath, percentage)
	if err != nil {
		return nil, err
	}
	return result.Generated, nil
}
