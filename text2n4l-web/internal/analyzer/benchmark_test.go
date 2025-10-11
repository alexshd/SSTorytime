package analyzer

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

// Benchmark tests for core analyzer functions

func BenchmarkSanitize(b *testing.B) {
	input := `# Markdown Header

This is **bold** text with [links](https://example.com) and <strong>HTML tags</strong>.

- List item with *emphasis*
- Another item with ` + "`" + `code` + "`" + `

> Blockquote with complex content

Final paragraph with mixed **formatting** and <em>more</em> content.`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Sanitize(input)
	}
}

func BenchmarkFractionateTextFile(b *testing.B) {
	// Load test data
	testData, err := ioutil.ReadFile(filepath.Join("testdata", "sample.txt"))
	if err != nil {
		b.Skip("Test data not available:", err)
	}
	input := string(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FractionateTextFile(input)
	}
}

func BenchmarkExtractIntentionalTokens(b *testing.B) {
	fragments := []string{
		"The quick brown fox jumps over the lazy dog",
		"A lazy dog sleeps peacefully in the warm sunshine",
		"Quick movements help the clever fox hunt effectively in the forest",
		"The forest provides shelter for many animals including foxes and dogs",
		"Effective hunting requires both speed and patience from the fox",
		"Sunshine warms the earth and helps plants grow in the forest",
		"Many animals depend on the forest ecosystem for survival",
		"The clever hunter must understand the behavior of its prey",
		"Peaceful moments in nature provide rest for weary animals",
		"Forest ecosystems support complex food webs and relationships",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ExtractIntentionalTokens(fragments)
	}
}

func BenchmarkN4LSkeletonOutput(b *testing.B) {
	// Load test data
	testData, err := ioutil.ReadFile(filepath.Join("testdata", "sample.txt"))
	if err != nil {
		b.Skip("Test data not available:", err)
	}
	content := string(testData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = N4LSkeletonOutput("test.txt", content, 50.0)
	}
}

func BenchmarkN4LSkeletonOutputLarge(b *testing.B) {
	// Create larger test content
	baseContent := `This is a comprehensive test document designed to evaluate the performance characteristics of the N4L conversion system. The document contains multiple paragraphs with varying complexity and content depth.

Each paragraph represents a different aspect of the analysis process. Some paragraphs focus on technical implementation details, while others explore conceptual frameworks and theoretical foundations.

The intentionality analysis system must process each sentence and evaluate its significance within the broader context. This involves complex natural language processing operations including tokenization, n-gram extraction, and statistical analysis.

Performance benchmarking requires consistent test data that represents realistic usage patterns. The content should be neither too simple nor overwhelmingly complex, striking a balance that reflects actual user scenarios.

Statistical significance emerges through repeated analysis of meaningful content patterns. The system identifies key concepts and relationships while filtering out noise and irrelevant information.

Quality metrics include processing speed, memory efficiency, and accuracy of content extraction. These benchmarks help identify optimization opportunities and regression risks during development.

Scalability testing ensures the system performs well with documents of varying sizes. From small snippets to comprehensive articles, the analyzer must maintain consistent performance characteristics.

Real-world applications demand robust handling of diverse content types including technical documentation, narrative text, and mixed-format documents with embedded markup and formatting.`

	// Repeat content to create larger document
	var largeContent string
	for i := 0; i < 10; i++ {
		largeContent += baseContent + "\n\n"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = N4LSkeletonOutput("large_test.txt", largeContent, 30.0)
	}
}

func BenchmarkConvertTextToN4LWithResult(b *testing.B) {
	content := `Artificial intelligence represents a transformative technology with far-reaching implications for society, economy, and human interaction.

Machine learning algorithms enable systems to improve performance through experience rather than explicit programming. This capability opens new possibilities for automation and decision support.

Natural language processing allows computers to understand and generate human language with increasing sophistication. Applications range from translation services to conversational interfaces.

Computer vision systems can analyze and interpret visual information, enabling applications in medical diagnosis, autonomous vehicles, and security systems.

Ethical considerations surrounding AI development include bias mitigation, privacy protection, and ensuring beneficial outcomes for all stakeholders.

The future of AI development depends on collaborative efforts between researchers, policymakers, and industry practitioners to establish responsible development practices.`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ConvertTextToN4LWithResult(content)
	}
}

// Memory allocation benchmarks
func BenchmarkMemoryAllocation(b *testing.B) {
	content := "This is a test sentence for memory allocation benchmarking. It contains enough content to trigger meaningful allocations during processing."

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fragments := FractionateTextFile(content)
		_ = ExtractIntentionalTokens(fragments)
	}
}

// Concurrent processing benchmark
func BenchmarkConcurrentProcessing(b *testing.B) {
	contents := []string{
		"First document with sample content for concurrent processing analysis.",
		"Second document containing different content patterns and structures.",
		"Third document with technical terminology and complex sentence structures.",
		"Fourth document exploring various topics and content complexity levels.",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			content := contents[i%len(contents)]
			_ = N4LSkeletonOutput("concurrent_test.txt", content, 50.0)
			i++
		}
	})
}
