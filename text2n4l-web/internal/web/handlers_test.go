package web

import (
	"bytes"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// NewTestServer creates a server for testing without template dependencies
func NewTestServer() *Server {
	// Create a simple template for testing
	tmpl := template.Must(template.New("index.tmpl").Parse(`
<!DOCTYPE html>
<html><head><title>Test</title></head>
<body>
<h1>Upload Text File</h1>
<h2>Quick Convert</h2>
<form>Test form</form>
</body></html>
	`))

	return &Server{
		BatchSize: 50,
		templates: tmpl,
	}
}

func TestHomeHandler(t *testing.T) {
	server := NewTestServer()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server.HomeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Upload Text File") {
		t.Error("Response should contain upload form")
	}

	if !strings.Contains(w.Body.String(), "Quick Convert") {
		t.Error("Response should contain convert form")
	}
}

func TestConvertHandler(t *testing.T) {
	server := NewTestServer()

	// Test valid conversion
	t.Run("valid_conversion", func(t *testing.T) {
		data := "text=This is a test sentence. Another sentence for testing."
		req := httptest.NewRequest("POST", "/convert", strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		server.ConvertHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "@sen") {
			t.Error("Response should contain N4L sentence markers")
		}

		if !strings.Contains(body, "(extract/quote from)") {
			t.Error("Response should contain arrow definitions")
		}

		if !strings.Contains(body, "# (begin) ************") {
			t.Error("Response should contain N4L begin marker")
		}
	})

	// Test empty text
	t.Run("empty_text", func(t *testing.T) {
		data := "text="
		req := httptest.NewRequest("POST", "/convert", strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		server.ConvertHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for empty text, got %d", w.Code)
		}
	})

	// Test wrong method
	t.Run("wrong_method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/convert", nil)
		w := httptest.NewRecorder()

		server.ConvertHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for GET method, got %d", w.Code)
		}
	})
}

func TestUploadHandler(t *testing.T) {
	server := NewTestServer()

	// Create temporary test file
	tmpFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := "This is test content for file upload. It contains multiple sentences for testing the upload functionality."
	if _, err := tmpFile.Write([]byte(testContent)); err != nil {
		t.Fatal("Failed to write to temp file:", err)
	}
	tmpFile.Close()

	// Test valid file upload
	t.Run("valid_upload", func(t *testing.T) {
		// Create multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		// Add file field
		file, err := os.Open(tmpFile.Name())
		if err != nil {
			t.Fatal("Failed to open temp file:", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatal("Failed to create form file:", err)
		}

		if _, err := io.Copy(part, file); err != nil {
			t.Fatal("Failed to copy file content:", err)
		}

		writer.Close()

		req := httptest.NewRequest("POST", "/upload", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		server.UploadHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "File uploaded") {
			t.Error("Response should indicate successful upload")
		}
	})

	// Test missing file
	t.Run("missing_file", func(t *testing.T) {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		writer.Close()

		req := httptest.NewRequest("POST", "/upload", &buf)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		server.UploadHandler(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400 for missing file, got %d", w.Code)
		}
	})

	// Test wrong method
	t.Run("wrong_method", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/upload", nil)
		w := httptest.NewRecorder()

		server.UploadHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status 405 for GET method, got %d", w.Code)
		}
	})
}

// Benchmark integration tests
func BenchmarkHomeHandler(b *testing.B) {
	server := NewTestServer()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		server.HomeHandler(w, req)
	}
}

func BenchmarkConvertHandler(b *testing.B) {
	server := NewTestServer()
	data := "text=This is a test sentence for benchmarking. It contains enough content to provide meaningful performance metrics."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/convert", strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		server.ConvertHandler(w, req)
	}
}
