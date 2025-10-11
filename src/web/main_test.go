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

func TestUploadHandler(t *testing.T) {
	// Create a sample file to upload
	fileContent := "test content"
	file, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	file.WriteString(fileContent)
	file.Seek(0, io.SeekStart)

	// Prepare multipart form
	body := &strings.Builder{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	io.Copy(part, file)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	uploadHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
