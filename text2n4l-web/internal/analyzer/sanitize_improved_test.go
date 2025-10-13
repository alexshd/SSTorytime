package analyzer

import (
	"testing"
)

// Test cases for HTML/Markdown sanitization
var sanitizeTestCases = []struct {
	name     string
	input    string
	expected string
}{
	{
		name:     "simple HTML tags",
		input:    "Hello <b>world</b>!",
		expected: "Hello world!", // Parser may add spaces
	},
	{
		name:     "HTML entities",
		input:    "5&nbsp;&lt;&nbsp;10&nbsp;&amp;&nbsp;10&nbsp;&gt;&nbsp;3",
		expected: "5\u00a0<\u00a010\u00a0&\u00a010\u00a0>\u00a03", // &nbsp; → \u00a0 (non-breaking space)
	},
	{
		name:     "nested HTML tags",
		input:    "<div><p>Hello <span>nested</span> world</p></div>",
		expected: "Hello nested world",
	},
	{
		name:     "malformed HTML",
		input:    "Hello <b>unclosed <i>tags</b> test",
		expected: "Hello unclosed tags test",
	},
	{
		name:     "markdown link",
		input:    "Check [this link](https://example.com) out",
		expected: "Check this link out",
	},
	{
		name:     "markdown image",
		input:    "See ![alt text](image.png) here",
		expected: "See alt text here",
	},
	{
		name:     "markdown emphasis",
		input:    "This is **bold** and *italic* and __underline__",
		expected: "This is bold and italic and underline",
	},
	{
		name:     "inline code",
		input:    "Use `fmt.Println()` to print",
		expected: "Use fmt.Println[] to print",
	},
	{
		name:     "markdown heading",
		input:    "## This is a heading",
		expected: "This is a heading",
	},
	{
		name:     "parentheses to brackets",
		input:    "Test (with parens) here",
		expected: "Test [with parens] here",
	},
	{
		name:     "complex mixed HTML and Markdown",
		input:    "<div>**Bold** text with <a href='link'>a link</a> and ![img](pic.jpg)</div>",
		expected: "Bold text with a link and img",
	},
	{
		name:     "HTML with entities and tags",
		input:    "<p>Price:&nbsp;$100&nbsp;&mdash;&nbsp;<b>Save&nbsp;20%!</b></p>",
		expected: "Price:\u00a0$100\u00a0—\u00a0Save\u00a020%!", // &nbsp; preserved as \u00a0
	},
}

func TestSanitizeWithParser(t *testing.T) {
	for _, tc := range sanitizeTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SanitizeWithParser(tc.input)
			if result != tc.expected {
				t.Errorf("SanitizeWithParser(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSanitizeWithRegex(t *testing.T) {
	for _, tc := range sanitizeTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SanitizeWithRegex(tc.input)
			if result != tc.expected {
				t.Errorf("SanitizeWithRegex(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

// Benchmark inputs
var (
	benchSmallHTML   = "<p>Hello <b>world</b>! This is a <a href='test'>link</a>.</p>"
	benchMediumHTML  = "<div><h2>Title</h2><p>This is a **bold** paragraph with <em>emphasis</em> and a [link](http://example.com).</p><p>Another paragraph with &nbsp; entities and `code`.</p></div>"
	benchLargeHTML   = "<article><h1>Long Article</h1><p>Lorem ipsum dolor sit amet, <strong>consectetur</strong> adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.</p><p>Lorem ipsum dolor sit amet, <strong>consectetur</strong> adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.</p><p>Lorem ipsum dolor sit amet, <strong>consectetur</strong> adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.</p><p>Lorem ipsum dolor sit amet, <strong>consectetur</strong> adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.</p><p>Lorem ipsum dolor sit amet, <strong>consectetur</strong> adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.</p></article>"
	benchComplexHTML = `
<html>
<head><title>Test</title></head>
<body>
<div class="content">
	<h1>Main&nbsp;Heading</h1>
	<p>This is **bold** text with <a href="https://example.com">a link</a> and ![image](test.jpg).</p>
	<ul>
		<li>Item&nbsp;1 with <code>code</code></li>
		<li>Item 2 with __emphasis__</li>
		<li>Item 3 with *italics*</li>
	</ul>
	<blockquote>
		A quote with &ldquo;smart&rdquo; quotes and &mdash; dashes.
	</blockquote>
</div>
</body>
</html>`
)

// Benchmarks for SanitizeWithParser (new approach)
func BenchmarkSanitizeWithParser_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeWithParser(benchSmallHTML)
	}
}

func BenchmarkSanitizeWithParser_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeWithParser(benchMediumHTML)
	}
}

func BenchmarkSanitizeWithParser_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeWithParser(benchLargeHTML)
	}
}

func BenchmarkSanitizeWithParser_Complex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeWithParser(benchComplexHTML)
	}
}

// Benchmarks for SanitizeWithRegex (old approach)
func BenchmarkSanitizeWithRegex_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeWithRegex(benchSmallHTML)
	}
}

func BenchmarkSanitizeWithRegex_Medium(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeWithRegex(benchMediumHTML)
	}
}

func BenchmarkSanitizeWithRegex_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeWithRegex(benchLargeHTML)
	}
}

func BenchmarkSanitizeWithRegex_Complex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeWithRegex(benchComplexHTML)
	}
}
