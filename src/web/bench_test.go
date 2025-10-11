package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func BenchmarkUploadHandler(b *testing.B) {
	fileContent := "benchmark content"
	for i := 0; i < b.N; i++ {
		file, err := os.CreateTemp("", "benchfile.txt")
		if err != nil {
			b.Fatalf("Failed to create temp file: %v", err)
		}
		file.WriteString(fileContent)
		file.Seek(0, io.SeekStart)

		body := &strings.Builder{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "benchfile.txt")
		if err != nil {
			b.Fatalf("Failed to create form file: %v", err)
		}
		io.Copy(part, file)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(body.String()))
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		uploadHandler(w, req)
		file.Close()
		os.Remove(file.Name())
	}
}
