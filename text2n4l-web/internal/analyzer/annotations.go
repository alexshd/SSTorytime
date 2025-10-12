package analyzer

import (
	"regexp"
	"strings"
)

// Annotation represents parsed information for a single sentence from N4L output
type Annotation struct {
	Index          int      // @sen index in original ordering
	Text           string   // sentence text
	ExtractFrom    string   // extract/quote from ...
	AppearsCloseTo []string // list of context keywords
	NeedsAttention bool     // heuristic marker for highlighting
}

var (
	senLineRe      = regexp.MustCompile(`^@sen(\d+)\s+(.*)$`)
	appearsCloseRe = regexp.MustCompile(`" \(appears close to\)\s+(.*)$`)
	extractFromRe  = regexp.MustCompile(`" \(extract/quote from\)\s+(.*)$`)
)

// ParseN4LAnnotations parses the N4L skeleton output produced by N4LSkeletonOutput
// and returns structured annotations keyed by sentence index.
func ParseN4LAnnotations(n4l string) []Annotation {
	lines := strings.Split(n4l, "\n")
	var out []Annotation
	var cur *Annotation
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		if m := senLineRe.FindStringSubmatch(strings.TrimSpace(line)); m != nil {
			// Start a new annotation
			if cur != nil {
				cur.NeedsAttention = attentionHeuristic(cur)
				out = append(out, *cur)
			}
			idx := atoiSafe(m[1])
			txt := strings.TrimSpace(m[2])
			cur = &Annotation{Index: idx, Text: txt}
			continue
		}
		if cur == nil {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if m := extractFromRe.FindStringSubmatch(trimmed); m != nil {
			cur.ExtractFrom = strings.TrimSpace(m[1])
			continue
		}
		if m := appearsCloseRe.FindStringSubmatch(trimmed); m != nil {
			kw := strings.TrimSpace(m[1])
			if kw != "" {
				cur.AppearsCloseTo = append(cur.AppearsCloseTo, kw)
			}
			continue
		}
	}
	if cur != nil {
		cur.NeedsAttention = attentionHeuristic(cur)
		out = append(out, *cur)
	}
	return out
}

func attentionHeuristic(a *Annotation) bool {
	// Simple initial heuristic until SSTconfig-driven rules are wired:
	// - Many context keywords
	// - Certain diagnostic keywords present
	if len(a.AppearsCloseTo) > 12 {
		return true
	}
	flags := []string{"arrow definition", "html", "intentionality", "n-gram", "text"}
	for _, kw := range a.AppearsCloseTo {
		low := strings.ToLower(kw)
		for _, f := range flags {
			if strings.Contains(low, f) {
				return true
			}
		}
	}
	return false
}

func atoiSafe(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}
