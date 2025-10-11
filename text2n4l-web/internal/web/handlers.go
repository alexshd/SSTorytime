package web

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"text2n4l-web/internal/analyzer"
)

// Server holds the web server configuration and handlers
type Server struct {
	BatchSize int
	templates *template.Template
}

// NewServer creates a new web server instance
func NewServer() *Server {
	// Load templates
	templates := template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"splitLines": func(text string) []string {
			if text == "" {
				return []string{}
			}
			lines := strings.Split(text, "\n")
			var result []string
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					result = append(result, line)
				}
			}
			return result
		},
		"contains": func(slice []int, item int) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
	})

	// Parse template files
	templates = template.Must(templates.ParseGlob("templates/*.tmpl"))

	return &Server{
		BatchSize: 50,
		templates: templates,
	}
}

// FileUploadData represents uploaded file data
type FileUploadData struct {
	Filename string
	Content  string
	Batches  []Batch
}

// Batch represents a batch of text for processing
type Batch struct {
	Index            int
	Content          string
	StartLine        int
	EndLine          int
	AmbiguousIndices []int
}

// TemplateData holds data for rendering templates
type TemplateData struct {
	Title            string
	OriginalText     string
	N4LOutput        string
	Filename         string
	CurrentBatch     int
	TotalBatches     int
	HasFile          bool
	AmbiguousIndices []int
}

// Global storage for uploaded file (in production, use proper session storage)
var currentFile *FileUploadData

// HomeHandler serves the main page
func (s *Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title:   "N4L Text Converter",
		HasFile: currentFile != nil,
	}

	if currentFile != nil {
		data.Filename = currentFile.Filename
		data.TotalBatches = len(currentFile.Batches)
	}

	s.renderTemplate(w, "index", data)
}

// UploadHandler handles file uploads
func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}
	if err := file.Close(); err != nil {
		http.Error(w, "Failed to close file", http.StatusInternalServerError)
		fmt.Printf("error closing file: %v\n", err)
		return
	}

	// Create batches
	batches := s.createBatches(string(content))

	// Store file data
	currentFile = &FileUploadData{
		Filename: header.Filename,
		Content:  string(content),
		Batches:  batches,
	}

	// Return HTMX response
	w.Header().Set("Content-Type", "text/html")
	if _, err := fmt.Fprintf(w, `
	       <div class="alert alert-success">
		       File uploaded: %s (%d batches)
		       <button hx-get="/batch/0" hx-target="#editor-container" class="btn btn-primary btn-sm ms-2">
			       Start Processing
		       </button>
	       </div>
       `, header.Filename, len(batches)); err != nil {
		http.Error(w, "Failed to write upload response", http.StatusInternalServerError)
		fmt.Printf("error writing upload response: %v\n", err)
		return
	}
}

// BatchHandler serves a specific batch for editing
func (s *Server) BatchHandler(w http.ResponseWriter, r *http.Request) {
	if currentFile == nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}

	// Get batch index from URL
	batchStr := strings.TrimPrefix(r.URL.Path, "/batch/")
	batchIndex, err := strconv.Atoi(batchStr)
	if err != nil || batchIndex < 0 || batchIndex >= len(currentFile.Batches) {
		http.Error(w, "Invalid batch index", http.StatusBadRequest)
		return
	}

	batch := currentFile.Batches[batchIndex]

	// Convert batch content to N4L DSL skeleton (_edit_me.n4l style)
	n4lOutput := analyzer.N4LSkeletonOutput(currentFile.Filename, batch.Content, 50.0)

	data := TemplateData{
		OriginalText:     batch.Content,
		N4LOutput:        n4lOutput,
		CurrentBatch:     batchIndex,
		TotalBatches:     len(currentFile.Batches),
		AmbiguousIndices: nil, // not used for now
	}

	s.renderTemplate(w, "dual-editor", data)
}

// ConvertHandler handles text conversion requests
func (s *Server) ConvertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get text from form
	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "No text provided", http.StatusBadRequest)
		return
	}

	// Convert to N4L DSL skeleton
	filename := "uploaded.txt"
	if currentFile != nil && currentFile.Filename != "" {
		filename = currentFile.Filename
	}
	n4lOut := analyzer.N4LSkeletonOutput(filename, text, 50.0)

	// Return N4L output
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(n4lOut)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		fmt.Printf("error writing N4L output: %v\n", err)
		return
	}
}

// createBatches splits content into manageable batches
func (s *Server) createBatches(content string) []Batch {
	lines := strings.Split(content, "\n")
	var batches []Batch

	for i := 0; i < len(lines); i += s.BatchSize {
		end := i + s.BatchSize
		if end > len(lines) {
			end = len(lines)
		}

		batchContent := strings.Join(lines[i:end], "\n")
		if strings.TrimSpace(batchContent) == "" {
			continue
		}

		batch := Batch{
			Index:     len(batches),
			Content:   batchContent,
			StartLine: i + 1,
			EndLine:   end,
		}

		batches = append(batches, batch)
	}

	return batches
}

// renderTemplate renders HTML templates from files
func (s *Server) renderTemplate(w http.ResponseWriter, templateName string, data TemplateData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute template
	if err := s.templates.ExecuteTemplate(w, templateName+".tmpl", data); err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
		fmt.Printf("Template error: %v\n", err)
	}
}
