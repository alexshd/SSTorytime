# Using Go Standard Library for Text Sanitization

## Current Approach vs Better Alternatives

### 1. HTML Tag Removal

#### Current (Regex)

```go
tagRE := regexp.MustCompile(`<[^>]+>`)
s = tagRE.ReplaceAllString(s, " ")
```

**Problems:**

- Doesn't handle malformed HTML
- Doesn't extract text from nested tags properly
- Doesn't handle HTML entities (`&nbsp;`, `&lt;`, etc.)
- Can break on edge cases like `<script>` with `<` inside

#### Better Option A: `golang.org/x/net/html`

```go
import "golang.org/x/net/html"

func stripHTMLTags(s string) string {
    doc, err := html.Parse(strings.NewReader(s))
    if err != nil {
        return s // fallback to original
    }

    var buf strings.Builder
    var extract func(*html.Node)
    extract = func(n *html.Node) {
        if n.Type == html.TextNode {
            buf.WriteString(n.Data)
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            extract(c)
        }
    }
    extract(doc)
    return buf.String()
}
```

**Pros:**

- ✅ Proper HTML parsing
- ✅ Handles malformed HTML gracefully
- ✅ Extracts text from all nested tags
- ✅ Handles HTML entities automatically

**Cons:**

- ❌ Heavier (parses full DOM tree)
- ❌ Slower for simple cases

#### Better Option B: `html.UnescapeString` (stdlib)

```go
import "html"

// Just for entities
s = html.UnescapeString(s)
// Then use regex for tags
```

**Pros:**

- ✅ Handles HTML entities (`&nbsp;` → space)
- ✅ Lightweight (stdlib)

### 2. Markdown Handling

#### Current (Multiple Regexes)

```go
// Images: ![alt](url) -> alt
imgRE := regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`)
// Links: [text](url) -> text
linkRE := regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
// Etc...
```

**Problems:**

- Doesn't handle complex markdown (nested links, escaped brackets)
- Multiple regex passes (slow)
- Easy to miss edge cases

#### Better Option A: Parse then Extract Text

```go
import "github.com/gomarkdown/markdown"
import "github.com/gomarkdown/markdown/parser"

func stripMarkdown(s string) string {
    extensions := parser.CommonExtensions
    p := parser.NewWithExtensions(extensions)

    // Parse to AST
    doc := p.Parse([]byte(s))

    // Extract plain text by walking AST
    return extractText(doc)
}
```

**Pros:**

- ✅ Handles all markdown correctly
- ✅ Future-proof

**Cons:**

- ❌ Overkill for simple stripping
- ❌ Large dependency

#### Better Option B: Optimized Regex with Replacer

```go
// Pre-compile all markdown patterns
var markdownStripper = struct {
    images   *regexp.Regexp
    links    *regexp.Regexp
    code     *regexp.Regexp
    emphasis *strings.Replacer
}{
    images:   regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`),
    links:    regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`),
    code:     regexp.MustCompile("`+"),
    emphasis: strings.NewReplacer("**", "", "__", "", "*", "", "_", ""),
}
```

**Pros:**

- ✅ Fast (pre-compiled)
- ✅ No dependencies
- ✅ Good enough for most cases

### 3. URL Handling

#### Current (Not explicitly handled)

We don't currently parse URLs specifically.

#### Potential Use: `net/url` (stdlib)

```go
import "net/url"

func isURL(s string) bool {
    u, err := url.Parse(s)
    return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func extractURLText(s string) string {
    if u, err := url.Parse(s); err == nil {
        return u.Host + u.Path // More readable than full URL
    }
    return s
}
```

## Recommended Approach

For `text2n4l-web`, I recommend a **hybrid approach**:

### Option 1: Minimal Dependencies (Current + stdlib)

```go
import (
    "html"
    "regexp"
    "strings"
)

// Use html.UnescapeString for entities
// Keep optimized regex for tags
// Use strings.Replacer for markdown emphasis
```

**When to use:** Production code where dependencies matter

### Option 2: Proper Parsing (Best correctness)

```go
import (
    "golang.org/x/net/html"
    "github.com/gomarkdown/markdown"
)

// Use html.Parse for HTML
// Use markdown parser for markdown
```

**When to use:** When input quality is unknown/untrusted

### Option 3: Hybrid (Recommended)

```go
import (
    "html"
    "golang.org/x/net/html"  // Only for HTML
    "regexp"                  // Only for markdown
    "strings"
)

func Sanitize(s string) string {
    // 1. Handle HTML properly with x/net/html
    s = stripHTMLWithParser(s)

    // 2. Handle HTML entities
    s = html.UnescapeString(s)

    // 3. Strip markdown with optimized regex
    s = stripMarkdownWithRegex(s)

    // 4. Clean up
    s = normalizeWhitespace(s)

    return s
}
```

## Benchmark Impact

Let me test if using `golang.org/x/net/html` is faster than regex:

```bash
# Install if needed
go get golang.org/x/net/html

# Benchmark comparison
go test -bench=BenchmarkSanitize -benchmem
```

## My Recommendation

For **text2n4l-web**, here's what I'd do:

### 1. For HTML: Use `golang.org/x/net/html`

**Why:**

- More correct than regex
- Handles edge cases (malformed HTML, entities)
- Not much slower for typical inputs
- Already well-tested

### 2. For Markdown: Keep optimized regex

**Why:**

- Markdown parsing libraries are overkill
- Our regex handles 95% of cases
- Already optimized with pre-compilation
- No extra dependencies

### 3. Add `html.UnescapeString`

**Why:**

- Free (stdlib)
- Handles entities correctly
- One-line addition

### 4. Consider URLs: Add basic detection

**Why:**

- Can make output more readable
- Stdlib `net/url` is fast
- Optional enhancement

## Implementation

Would you like me to:

1. ✅ Implement HTML parsing with `golang.org/x/net/html`
2. ✅ Add `html.UnescapeString` for entities
3. ✅ Keep current markdown regex (already optimized)
4. ✅ Benchmark the changes
5. ❌ Skip markdown parsing libraries (overkill)

This gives us the best balance of correctness, performance, and maintainability!
