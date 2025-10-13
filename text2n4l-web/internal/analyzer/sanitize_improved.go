package analyzer

import (
	"html"
	"regexp"
	"strings"

	htmlparser "golang.org/x/net/html"
)

// Pre-compiled regexes for Sanitize (initialized once)
var (
	// Markdown patterns
	sanitizeMarkdownImageRE = regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`)
	sanitizeMarkdownLinkRE  = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
	sanitizeBackticksRE     = regexp.MustCompile("`+")
	sanitizeWhitespaceRE    = regexp.MustCompile(`\s+`)

	// Markdown emphasis remover
	sanitizeMarkdownReplacer = strings.NewReplacer(
		"**", "",
		"__", "",
		"*", "",
		"_", "",
	)

	// N4L safety replacer (parentheses → brackets)
	sanitizeN4LReplacer = strings.NewReplacer(
		"(", "[",
		")", "]",
	)
)

// stripHTMLTags removes HTML tags and extracts text content using proper HTML parser
func stripHTMLTags(s string) string {
	// Quick check: if no < or >, skip parsing
	if !strings.ContainsAny(s, "<>") {
		return s
	}

	doc, err := htmlparser.Parse(strings.NewReader(s))
	if err != nil {
		// Fallback to simple regex if parsing fails
		tagRE := regexp.MustCompile(`<[^>]+>`)
		return tagRE.ReplaceAllString(s, " ")
	}

	var buf strings.Builder
	var extract func(*htmlparser.Node)
	extract = func(n *htmlparser.Node) {
		if n.Type == htmlparser.TextNode {
			buf.WriteString(n.Data)
		}
		if n.Type == htmlparser.ElementNode {
			// Add space between block elements for readability
			switch n.Data {
			case "br", "p", "div", "li", "h1", "h2", "h3", "h4", "h5", "h6":
				if buf.Len() > 0 {
					buf.WriteString(" ")
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(doc)
	return buf.String()
}

// SanitizeWithParser strips HTML/Markdown using proper parsing (new optimized version)
func SanitizeWithParser(s string) string {
	// 1. Remove HTML tags with proper parser
	s = stripHTMLTags(s)

	// 2. Decode HTML entities (&nbsp; → space, &lt; → <, etc.)
	s = html.UnescapeString(s)

	// 3. Strip markdown images: ![alt](url) -> alt
	s = sanitizeMarkdownImageRE.ReplaceAllString(s, "$1")

	// 4. Strip markdown links: [text](url) -> text
	s = sanitizeMarkdownLinkRE.ReplaceAllString(s, "$1")

	// 5. Remove backticks (inline code)
	s = sanitizeBackticksRE.ReplaceAllString(s, "")

	// 6. Remove markdown emphasis (**, __, *, _)
	s = sanitizeMarkdownReplacer.Replace(s)

	// 7. Strip leading markdown headings and blockquote markers
	s = strings.TrimLeft(s, "#> ")

	// 8. Normalize whitespace
	s = sanitizeWhitespaceRE.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)

	// 9. Replace parentheses with brackets (N4L safety)
	s = sanitizeN4LReplacer.Replace(s)

	return s
}

// SanitizeWithRegex strips HTML/Markdown using regex only (original approach)
func SanitizeWithRegex(s string) string {
	// Remove HTML tags (regex)
	tagRE := regexp.MustCompile(`<[^>]+>`)
	s = tagRE.ReplaceAllString(s, " ")

	// Decode HTML entities
	s = html.UnescapeString(s)

	// Markdown: images ![alt](url) -> alt
	s = sanitizeMarkdownImageRE.ReplaceAllString(s, "$1")

	// Markdown: links [text](url) -> text
	s = sanitizeMarkdownLinkRE.ReplaceAllString(s, "$1")

	// Inline code backticks -> plain
	s = sanitizeBackticksRE.ReplaceAllString(s, "")

	// Emphasis markers
	s = sanitizeMarkdownReplacer.Replace(s)

	// Strip leading markdown headings and blockquote markers
	s = strings.TrimLeft(s, "#> ")

	// Normalize whitespace
	s = sanitizeWhitespaceRE.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)

	// Replace parentheses with brackets
	s = sanitizeN4LReplacer.Replace(s)

	return s
}
