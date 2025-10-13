package analyzer

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// N4LSkeletonOutput generates the _edit_me.n4l skeleton DSL output, matching the CLI
func N4LSkeletonOutput(filename string, content string, percentage float64) string {
	MemoryInit()
	// DON'T sanitize before fractionation - it destroys newlines needed for splitting
	// Sanitize only happens on individual fragments during output

	// Fractionate content directly
	fragments := FractionateTextFile(content)
	L := len(fragments)
	if L == 0 {
		return "# No processable text found\n"
	}

	// Simulate partitioning as in CLI (single partition for now)
	// psf := [][]string{fragments}

	// Select by running and static intent
	running := make([]TextRank, 0, L)
	static := make([]TextRank, 0, L)
	for i, frag := range fragments {
		running = append(running, TextRank{
			Fragment:     frag,
			Significance: RunningIntentionality(i, frag),
			Order:        i,
			Partition:    0,
		})
		static = append(static, TextRank{
			Fragment:     frag,
			Significance: AssessStaticIntent(frag, L, ExtractIntentionalTokens(fragments), 1),
			Order:        i,
			Partition:    0,
		})
	}
	selection := mergeSelections(running, static)

	// Build ambient phrase frequencies for context arrows
	freqs := ExtractIntentionalTokens(fragments)
	// Derive top context phrases (n=2,3 preferred, then 1-grams)
	topBigrams := topN(freqs[N2GRAM], 8)
	topTrigrams := topN(freqs[N3GRAM], 8)
	topUnigrams := topN(freqs[N1GRAM], 8)
	// Compose context keywords list for "with ..."
	contextWith := pickContextForWith(topTrigrams, topBigrams, topUnigrams, 2)
	// Compose a longer list for "appears close to"
	appearsClose := mergeAndDedup(append(topTrigrams, append(topBigrams, topUnigrams...)...), 18)

	// Output in _edit_me.n4l style
	var sb strings.Builder
	base := filename
	sb.WriteString(" - Samples from " + base + "\n\n# (begin) ************\n")
	// Optional context header similar to CLI if context exists
	if len(contextWith) > 0 {
		sb.WriteString("\n :: " + strings.Join(contextWith, ", ") + " ::\n")
	}
	for _, sel := range selection {
		// Sentence line
		sb.WriteString(fmt.Sprintf("\n@sen%d   %s\n", sel.Order, Sanitize(sel.Fragment)))
		// Arrow lines: extract/quote from part and appears close to keywords
		part := fmt.Sprintf("part %d of %s", 0, fileAlias(base))
		if len(contextWith) > 0 {
			sb.WriteString("              \" (extract/quote from) " + part + " with " + strings.Join(contextWith, ", ") + "\n")
		} else {
			sb.WriteString("              \" (extract/quote from) " + part + "\n")
		}
		for _, kw := range appearsClose {
			sb.WriteString("              \" (appears close to) " + kw + "\n")
		}
	}
	sb.WriteString("\n# (end) ************\n")
	sb.WriteString(fmt.Sprintf("\n# Final fraction %.2f of requested %.2f\n", float64(len(selection)*100)/float64(L), percentage))
	sb.WriteString(fmt.Sprintf("\n# Selected %d samples of %d: ", len(selection), L))
	for _, sel := range selection {
		sb.WriteString(fmt.Sprintf("%d ", sel.Order))
	}
	sb.WriteString("\n#\n")
	return sb.String()
}

// Sanitize strips HTML/Markdown noise and replaces parentheses with brackets for N4L output
func Sanitize(s string) string {
	// Remove HTML tags
	tagRE := regexp.MustCompile(`<[^>]+>`) // simple tag stripper
	s = tagRE.ReplaceAllString(s, " ")
	// Markdown: images ![alt](url) -> alt
	imgRE := regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`) // capture alt text
	s = imgRE.ReplaceAllString(s, "$1")
	// Markdown: links [text](url) -> text
	linkRE := regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
	s = linkRE.ReplaceAllString(s, "$1")
	// Inline code backticks -> plain
	btRE := regexp.MustCompile("`+")
	s = btRE.ReplaceAllString(s, "")
	// Emphasis markers
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "__", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "_", "")
	// Strip leading markdown headings and blockquote markers
	s = strings.TrimLeft(s, "#> ")
	// Normalize whitespace
	spaceRE := regexp.MustCompile(`\s+`)
	s = spaceRE.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	// Replace parentheses with brackets to avoid N4L parse conflicts inside content
	replacer := strings.NewReplacer("(", "[", ")", "]")
	return replacer.Replace(s)
}

// CleanText preprocesses text to prepare for fractionation - matches CLI algorithm
func CleanText(s string) string {
	// Strip HTML/XML tags first
	tagRE := regexp.MustCompile(`<[^>]*>`)
	s = tagRE.ReplaceAllString(s, ":\n")

	// Handle English abbreviations that end with periods
	s = strings.ReplaceAll(s, "Mr.", "Mr")
	s = strings.ReplaceAll(s, "Ms.", "Ms")
	s = strings.ReplaceAll(s, "Mrs.", "Mrs")
	s = strings.ReplaceAll(s, "Dr.", "Dr")
	s = strings.ReplaceAll(s, "St.", "St")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")

	// Mark sentence boundaries with # for later splitting
	// Match end of sentence punctuation followed by space or newline
	sentenceEndRE := regexp.MustCompile(`([?!.。]+[ \n])`)
	s = sentenceEndRE.ReplaceAllString(s, "$0#")

	// Handle ellipsis
	ellipsisRE := regexp.MustCompile(`([.][.][.])+`)
	s = ellipsisRE.ReplaceAllString(s, "---")

	// Handle em-dash
	emDashRE := regexp.MustCompile(`[—]+`)
	s = emDashRE.ReplaceAllString(s, ", ")

	// Mark paragraphs with >>
	doubleNewlineRE := regexp.MustCompile(`[\n][\n]`)
	s = doubleNewlineRE.ReplaceAllString(s, ">>\n")

	// Consolidate spurious newlines
	multiNewlineRE := regexp.MustCompile(`[\n]+`)
	s = multiNewlineRE.ReplaceAllString(s, " ")

	return s
}

// FractionateTextFile splits text into processable fragments (sentences)
// This now matches the CLI algorithm: CleanText -> SplitIntoParaSentences
func FractionateTextFile(content string) []string {
	// First clean the text using CLI algorithm
	cleanedText := CleanText(content)

	// Split into paragraphs and sentences following CLI logic
	var fragments []string

	// Split on >> for paragraphs first
	paras := strings.Split(cleanedText, ">>")

	for _, para := range paras {
		// Split on # for sentences (inserted by CleanText)
		// Use regex to split on # but preserve the character before #
		// The regex [^#]# matches any character that's not # followed by #
		// Split preserves the character before # in the left part
		sentenceRE := regexp.MustCompile(`[^#]#`)
		sentences := sentenceRE.Split(para, -1)

		for _, sentence := range sentences {
			// Clean up: remove the # markers that might remain
			sentence = strings.ReplaceAll(sentence, "#", "")
			sentence = strings.TrimSpace(sentence)

			// Filter out very short fragments
			if len(sentence) > 2 {
				fragments = append(fragments, sentence)
			}
		}
	}

	return fragments
}

// fileAlias returns the filename without extension
func fileAlias(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// topN returns the top n keys from a frequency map, ordered by descending count
func topN(m map[string]float64, n int) []string {
	type kv struct {
		K string
		V float64
	}
	list := make([]kv, 0, len(m))
	for k, v := range m {
		if k == "" {
			continue
		}
		list = append(list, kv{K: k, V: v})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].V == list[j].V {
			return list[i].K < list[j].K
		}
		return list[i].V > list[j].V
	})
	if n > len(list) {
		n = len(list)
	}
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, list[i].K)
	}
	return out
}

// pickContextForWith chooses a small set of phrases to include in the "with ..." clause
func pickContextForWith(tris, bis, uns []string, want int) []string {
	var ctx []string
	for _, s := range tris {
		if len(ctx) >= want {
			break
		}
		ctx = append(ctx, s)
	}
	for _, s := range bis {
		if len(ctx) >= want {
			break
		}
		if !contains(ctx, s) {
			ctx = append(ctx, s)
		}
	}
	for _, s := range uns {
		if len(ctx) >= want {
			break
		}
		if !contains(ctx, s) {
			ctx = append(ctx, s)
		}
	}
	return ctx
}

// mergeAndDedup merges slices and returns the first k unique strings
func mergeAndDedup(in []string, k int) []string {
	seen := make(map[string]bool)
	var out []string
	for _, s := range in {
		if s == "" {
			continue
		}
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
			if len(out) >= k {
				break
			}
		}
	}
	return out
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// ExtractIntentionalTokens tokenizes text and tracks n-gram frequencies
func ExtractIntentionalTokens(fragments []string) [N_GRAM_MAX]map[string]float64 {
	var frequencies [N_GRAM_MAX]map[string]float64

	// Initialize frequency maps
	for i := N_GRAM_MIN; i < N_GRAM_MAX; i++ {
		frequencies[i] = make(map[string]float64)
	}

	for _, fragment := range fragments {
		words := strings.Fields(strings.ToLower(fragment))

		// Count unigrams
		for _, word := range words {
			// Clean word of punctuation
			word = cleanWord(word)
			if word != "" {
				frequencies[N1GRAM][word]++
			}
		}

		// Count bigrams
		for i := 0; i < len(words)-1; i++ {
			word1 := cleanWord(words[i])
			word2 := cleanWord(words[i+1])
			if word1 != "" && word2 != "" {
				bigram := word1 + " " + word2
				frequencies[N2GRAM][bigram]++
			}
		}

		// Count trigrams
		for i := 0; i < len(words)-2; i++ {
			word1 := cleanWord(words[i])
			word2 := cleanWord(words[i+1])
			word3 := cleanWord(words[i+2])
			if word1 != "" && word2 != "" && word3 != "" {
				trigram := word1 + " " + word2 + " " + word3
				frequencies[N3GRAM][trigram]++
			}
		}
	}

	return frequencies
}

// cleanWord removes punctuation and normalizes words
func cleanWord(word string) string {
	// Remove common punctuation
	word = strings.Trim(word, ".,!?;:()[]{}\"'")
	return strings.ToLower(word)
}

// ConvertTextToN4L performs the main conversion from text to N4L format
func ConvertTextToN4L(content string) (string, []int, error) {
	if content == "" {
		return "", nil, nil
	}

	// Initialize memory structures
	MemoryInit()

	// Fractionate text into sentences
	fragments := FractionateTextFile(content)
	if len(fragments) == 0 {
		return "# No processable text found\n", nil, nil
	}

	// Extract n-gram frequencies
	frequencies := ExtractIntentionalTokens(fragments)
	L := len(fragments)

	// Calculate intentionality scores for each fragment
	rankings := make([]TextRank, len(fragments))
	for i, fragment := range fragments {
		runningScore := RunningIntentionality(i, fragment)
		staticScore := AssessStaticIntent(fragment, L, frequencies, 1)

		rankings[i] = TextRank{
			Significance: runningScore + staticScore,
			Fragment:     fragment,
			Order:        i,
			Partition:    i / DUNBAR_30, // Simple partitioning
		}
	}

	// Select most significant fragments using both methods
	runningSelection := selectByRunningIntent(rankings, L)
	staticSelection := selectByStaticIntent(rankings, L)

	// Merge selections
	merged := mergeSelections(runningSelection, staticSelection)

	// Convert to N4L format
	n4lOutput, ambiguousIndices := formatAsN4L(merged, fragments)

	return n4lOutput, ambiguousIndices, nil
}

// selectByRunningIntent selects fragments based on running intentionality
func selectByRunningIntent(rankings []TextRank, totalLength int) []TextRank {
	var selected []TextRank

	// Select top fragments based on running intentionality
	threshold := 0.5 // Adjustable threshold

	for _, rank := range rankings {
		if rank.Significance > threshold {
			selected = append(selected, rank)
		}
	}

	// Ensure we have at least some fragments
	if len(selected) == 0 && len(rankings) > 0 {
		// Take the top fragment if nothing meets threshold
		best := rankings[0]
		for _, rank := range rankings {
			if rank.Significance > best.Significance {
				best = rank
			}
		}
		selected = append(selected, best)
	}

	return selected
}

// selectByStaticIntent selects fragments based on static intentionality
func selectByStaticIntent(rankings []TextRank, totalLength int) []TextRank {
	var selected []TextRank

	// Select fragments with above-average significance
	totalSignificance := 0.0
	for _, rank := range rankings {
		totalSignificance += rank.Significance
	}

	if len(rankings) > 0 {
		avgSignificance := totalSignificance / float64(len(rankings))

		for _, rank := range rankings {
			if rank.Significance >= avgSignificance {
				selected = append(selected, rank)
			}
		}
	}

	return selected
}

// mergeSelections combines running and static selections
func mergeSelections(running, static []TextRank) []TextRank {
	// Create a map to track selected fragments by order
	selected := make(map[int]TextRank)

	// Add running selections
	for _, rank := range running {
		selected[rank.Order] = rank
	}

	// Add static selections (may override with higher scores)
	for _, rank := range static {
		if existing, exists := selected[rank.Order]; exists {
			// Keep the one with higher significance
			if rank.Significance > existing.Significance {
				selected[rank.Order] = rank
			}
		} else {
			selected[rank.Order] = rank
		}
	}

	// Convert back to slice and sort by order
	var result []TextRank
	for _, rank := range selected {
		result = append(result, rank)
	}

	// Sort by original order
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Order > result[j].Order {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// formatAsN4L converts selected rankings to N4L format
func formatAsN4L(rankings []TextRank, allFragments []string) (string, []int) {
	if len(rankings) == 0 {
		return "# No significant content found\n", nil
	}

	var n4l strings.Builder
	var ambiguousIndices []int

	// Write header
	n4l.WriteString("# N4L Narrative Analysis\n\n")

	// Group by partitions
	partitions := make(map[int][]TextRank)
	for _, rank := range rankings {
		partitions[rank.Partition] = append(partitions[rank.Partition], rank)
	}

	// Write each partition
	for partition := 0; partition < len(partitions); partition++ {
		ranks, exists := partitions[partition]
		if !exists {
			continue
		}

		n4l.WriteString(fmt.Sprintf("## Context %d\n\n", partition+1))

		for _, rank := range ranks {
			// Mark ambiguous lines (low significance)
			if rank.Significance < 1.0 {
				ambiguousIndices = append(ambiguousIndices, rank.Order)
				n4l.WriteString("? ")
			} else {
				n4l.WriteString("+ ")
			}

			n4l.WriteString(rank.Fragment)
			n4l.WriteString(fmt.Sprintf(" [sig: %.2f]\n", rank.Significance))
		}

		n4l.WriteString("\n")
	}

	// Add metadata
	n4l.WriteString("---\n")
	n4l.WriteString(fmt.Sprintf("Total fragments: %d\n", len(allFragments)))
	n4l.WriteString(fmt.Sprintf("Selected: %d\n", len(rankings)))
	n4l.WriteString(fmt.Sprintf("Ambiguous: %d\n", len(ambiguousIndices)))

	return n4l.String(), ambiguousIndices
}

// ReadFile is a simple file reader function
func ReadFile(filename string) (string, error) {
	// This would normally read from file system
	// For web app, content is provided directly
	return "", fmt.Errorf("ReadFile not implemented in web context")
}

// ConvertTextToN4LResult represents the result of N4L conversion
type ConvertTextToN4LResult struct {
	N4LOutput         string
	AmbiguousIndices  []int
	Error             error
	TotalFragments    int
	SelectedFragments int
}

// ConvertTextToN4LWithResult returns a structured result
func ConvertTextToN4LWithResult(content string) ConvertTextToN4LResult {
	n4lOutput, ambiguousIndices, err := ConvertTextToN4L(content)

	fragments := FractionateTextFile(content)

	return ConvertTextToN4LResult{
		N4LOutput:         n4lOutput,
		AmbiguousIndices:  ambiguousIndices,
		Error:             err,
		TotalFragments:    len(fragments),
		SelectedFragments: strings.Count(n4lOutput, "+ ") + strings.Count(n4lOutput, "? "),
	}
}
