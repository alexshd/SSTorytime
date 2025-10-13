package analyzer

import (
	"bufio"
	"fmt"
	"net/http"
)

// StreamN4LOutput generates N4L output with streaming, sending data as it's processed
func StreamN4LOutput(w *bufio.Writer, flusher http.Flusher, filename string, content string, percentage float64) error {
	MemoryInit()

	// Send header immediately
	w.WriteString(" - Samples from " + filename + "\n\n# (begin) ************\n")
	w.Flush()
	flusher.Flush()

	// Fractionate content
	fragments := FractionateTextFile(content)
	L := len(fragments)
	if L == 0 {
		w.WriteString("# No processable text found\n")
		w.Flush()
		flusher.Flush()
		return nil
	}

	// Process fragments and build rankings
	running := make([]TextRank, 0, L)
	static := make([]TextRank, 0, L)

	// Extract tokens once for all fragments
	freqs := ExtractIntentionalTokens(fragments)

	// Process each fragment
	for i, frag := range fragments {
		running = append(running, TextRank{
			Fragment:     frag,
			Significance: RunningIntentionality(i, frag),
			Order:        i,
			Partition:    0,
		})
		static = append(static, TextRank{
			Fragment:     frag,
			Significance: AssessStaticIntent(frag, L, freqs, 1),
			Order:        i,
			Partition:    0,
		})

		// Stream progress indicator every 100 fragments
		if (i+1)%100 == 0 {
			w.WriteString(fmt.Sprintf("# Processing... %d/%d fragments\n", i+1, L))
			w.Flush()
			flusher.Flush()
		}
	}

	// Merge selections
	selection := mergeSelections(running, static)

	// Build context
	topBigrams := topN(freqs[N2GRAM], 8)
	topTrigrams := topN(freqs[N3GRAM], 8)
	topUnigrams := topN(freqs[N1GRAM], 8)
	contextWith := pickContextForWith(topTrigrams, topBigrams, topUnigrams, 2)
	appearsClose := mergeAndDedup(append(topTrigrams, append(topBigrams, topUnigrams...)...), 18)

	// Stream context header
	if len(contextWith) > 0 {
		w.WriteString("\n :: " + Join(contextWith, ", ") + " ::\n")
		w.Flush()
		flusher.Flush()
	}

	// Stream each sentence as it's processed
	base := filename
	part := fmt.Sprintf("part %d of %s", 0, fileAlias(base))

	for _, sel := range selection {
		// Write sentence
		w.WriteString(fmt.Sprintf("\n@sen%d   %s\n", sel.Order, Sanitize(sel.Fragment)))

		// Write arrows
		if len(contextWith) > 0 {
			w.WriteString("              \" (extract/quote from) " + part + " with " + Join(contextWith, ", ") + "\n")
		} else {
			w.WriteString("              \" (extract/quote from) " + part + "\n")
		}
		for _, kw := range appearsClose {
			w.WriteString("              \" (appears close to) " + kw + "\n")
		}

		// Flush after EVERY sentence for visible streaming
		w.Flush()
		flusher.Flush()
	}

	// Stream footer
	w.WriteString("\n# (end) ************\n")
	w.WriteString(fmt.Sprintf("\n# Final fraction %.2f of requested %.2f\n", float64(len(selection)*100)/float64(L), percentage))
	w.WriteString(fmt.Sprintf("\n# Selected %d samples of %d: ", len(selection), L))
	for _, sel := range selection {
		w.WriteString(fmt.Sprintf("%d ", sel.Order))
	}
	w.WriteString("\n#\n")
	w.Flush()
	flusher.Flush()

	return nil
}

// Join is a helper to join strings (in case strings.Join isn't available in scope)
func Join(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
