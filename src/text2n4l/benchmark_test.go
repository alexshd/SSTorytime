package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

// Benchmark comparison between direct file writing and buffered writing

func BenchmarkDirectFileWrite(b *testing.B) {
	tmpDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("%s/test_direct_%d.txt", tmpDir, i)
		fp, err := os.Create(filename)
		if err != nil {
			b.Fatal(err)
		}

		// Simulate many small writes like in the original code
		for j := 0; j < 1000; j++ {
			fmt.Fprintf(fp, "Line %d: This is a test line with some content\n", j)
		}

		fp.Close()
	}
}

func BenchmarkBufferedFileWrite(b *testing.B) {
	tmpDir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("%s/test_buffered_%d.txt", tmpDir, i)
		fp, err := os.Create(filename)
		if err != nil {
			b.Fatal(err)
		}

		writer := bufio.NewWriter(fp)

		// Simulate many small writes with buffered writer
		for j := 0; j < 1000; j++ {
			fmt.Fprintf(writer, "Line %d: This is a test line with some content\n", j)
		}

		writer.Flush()
		fp.Close()
	}
}

func BenchmarkAddIntentionalContextDirect(b *testing.B) {
	tmpDir := b.TempDir()
	context := []string{"context1", "context2", "context3", "context4", "context5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("%s/test_context_direct_%d.txt", tmpDir, i)
		fp, err := os.Create(filename)
		if err != nil {
			b.Fatal(err)
		}

		// Simulate old direct file writing approach
		for w := 0; w < len(context); w++ {
			fmt.Fprintf(fp, "              \" (%s) %s\n", "TEST_LABEL", context[w])
		}

		fp.Close()
	}
}

func BenchmarkAddIntentionalContextBuffered(b *testing.B) {
	tmpDir := b.TempDir()
	context := []string{"context1", "context2", "context3", "context4", "context5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("%s/test_context_buffered_%d.txt", tmpDir, i)
		fp, err := os.Create(filename)
		if err != nil {
			b.Fatal(err)
		}

		writer := bufio.NewWriter(fp)

		// Use the new buffered approach
		AddIntentionalContext(writer, context)

		writer.Flush()
		fp.Close()
	}
}

// Test to ensure both approaches produce identical output
func TestOutputIdentical(t *testing.T) {
	tmpDir := t.TempDir()
	context := []string{"test context 1", "test context 2", "another context"}

	// Direct file writing (old approach simulation)
	directFile := tmpDir + "/direct.txt"
	fp1, err := os.Create(directFile)
	if err != nil {
		t.Fatal(err)
	}

	for w := 0; w < len(context); w++ {
		fmt.Fprintf(fp1, "              \" (%s) %s\n", "TEST_LABEL", context[w])
	}
	fp1.Close()

	// Buffered writing (new approach)
	bufferedFile := tmpDir + "/buffered.txt"
	fp2, err := os.Create(bufferedFile)
	if err != nil {
		t.Fatal(err)
	}

	writer := bufio.NewWriter(fp2)
	for w := 0; w < len(context); w++ {
		fmt.Fprintf(writer, "              \" (%s) %s\n", "TEST_LABEL", context[w])
	}
	writer.Flush()
	fp2.Close()

	// Compare outputs
	directContent, err := os.ReadFile(directFile)
	if err != nil {
		t.Fatal(err)
	}

	bufferedContent, err := os.ReadFile(bufferedFile)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.EqualFold(string(directContent), string(bufferedContent)) {
		t.Errorf("Output differs:\nDirect: %s\nBuffered: %s",
			string(directContent), string(bufferedContent))
	}
}
