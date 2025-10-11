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
}

// NewServer creates a new web server instance
func NewServer() *Server {
	return &Server{
		BatchSize: 20, // Default batch size
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

// renderTemplate renders HTML templates
func (s *Server) renderTemplate(w http.ResponseWriter, templateName string, data TemplateData) {
	templates := map[string]string{
		"index":       getIndexTemplate(),
		"dual-editor": getDualEditorTemplate(),
	}

	// Create template with helper functions
	tmpl := template.New(templateName).Funcs(template.FuncMap{
		"sub":        func(a, b int) int { return a - b },
		"add":        func(a, b int) int { return a + b },
		"splitLines": func(s string) []string { return strings.Split(s, "\n") },
		"contains": func(slice []int, item int) bool {
			for _, v := range slice {
				if v == item {
					return true
				}
			}
			return false
		},
	})

	// Parse template
	tmpl, err := tmpl.Parse(templates[templateName])
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute template
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}

func getIndexTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.Title}}</title>
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
	<script src="https://unpkg.com/htmx.org@1.8.4"></script>
	<style>
		.ambiguous { background-color: #fff3cd; }
		.editor-container { min-height: 400px; }
	</style>
	<script>
	// --- Auto-refresh: poll reload.txt every 2s ---
	let lastReload = null;
	setInterval(function() {
		fetch('/static/../tmp/reload.txt', {cache: 'no-store'})
			.then(r => r.text())
			.then(txt => {
				if (lastReload && txt.trim() !== lastReload) {
					window.location.reload();
				}
				lastReload = txt.trim();
			});
	}, 2000);

	// --- Session persistence using localStorage ---
	function saveSession(data) {
		localStorage.setItem('n4l_session', JSON.stringify(data));
	}
	function loadSession() {
		let s = localStorage.getItem('n4l_session');
		if (!s) return null;
		try { return JSON.parse(s); } catch { return null; }
	}
	function clearSession() {
		localStorage.removeItem('n4l_session');
	}

	// On file upload, save session info
	document.addEventListener('htmx:afterOnLoad', function(evt) {
		if (evt.detail && evt.detail.target && evt.detail.target.id === 'upload-result') {
			// Parse filename and batch count from upload-result
			let html = evt.detail.target.innerHTML;
			let m = html.match(/File uploaded: ([^<]+) \((\d+) batches?\)/);
			if (m) {
				saveSession({ filename: m[1], batch: 0 });
			}
		}
	});

	// On batch navigation, update session
	document.addEventListener('htmx:afterOnLoad', function(evt) {
		if (evt.detail && evt.detail.target && evt.detail.target.id === 'editor-container') {
			// Try to extract batch index from header
			let html = evt.detail.target.innerHTML;
			let m = html.match(/Batch (\d+) of (\d+)/);
			let session = loadSession();
			if (m && session) {
				session.batch = parseInt(m[1], 10);
				saveSession(session);
			}
		}
	});

	// On page load, restore session if present
	window.addEventListener('DOMContentLoaded', function() {
		let session = loadSession();
		if (session && session.filename) {
			// Simulate clicking the batch button to restore progress
			setTimeout(function() {
				let url = '/batch/' + (session.batch || 0);
				let container = document.getElementById('editor-container');
				if (container) {
					container.innerHTML = '<div class="alert alert-info">Restoring session for <b>' + session.filename + '</b> (batch ' + (session.batch || 0) + ')...</div>';
				}
				fetch(url).then(r => r.text()).then(html => {
					if (container) container.innerHTML = html;
				});
			}, 300);
		}
	});
	</script>
</head>
<body>
    <div class="container mt-4">
        <h1>{{.Title}}</h1>
        
        <div class="row mb-4">
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>Upload Text File</h5>
                    </div>
                    <div class="card-body">
                        <form hx-post="/upload" hx-target="#upload-result" enctype="multipart/form-data">
                            <div class="mb-3">
                                <input type="file" class="form-control" name="file" accept=".txt,.md,.n4l" required>
                            </div>
                            <button type="submit" class="btn btn-primary">Upload & Process</button>
                        </form>
                        <div id="upload-result" class="mt-3"></div>
                    </div>
                </div>
            </div>
            
            <div class="col-md-6">
                <div class="card">
                    <div class="card-header">
                        <h5>Quick Con</h5>
                    </div>
                    <div class="card-body">
                        <form hx-post="/convert" hx-target="#convert-result">
                            <div class="mb-3">
                                <textarea class="form-control" name="text" rows="5" placeholder="Enter text to convert..."></textarea>
                            </div>
                            <button type="submit" class="btn btn-success">Convert</button>
                        </form>
                        <div id="convert-result" class="mt-3"></div>
                    </div>
                </div>
            </div>
        </div>
        
        {{if .HasFile}}
        <div class="row">
            <div class="col-12">
                <div class="card">
                    <div class="card-header">
                        <h5>File: {{.Filename}} ({{.TotalBatches}} batches)</h5>
                    </div>
                    <div class="card-body">
                        <button hx-get="/batch/0" hx-target="#editor-container" class="btn btn-primary">
                            Start Batch Processing
                        </button>
                    </div>
                </div>
            </div>
        </div>
        {{end}}
        
        <div id="editor-container" class="mt-4"></div>
    </div>
</body>
</html>`
}

func getDualEditorTemplate() string {
	return `<div class="card">
	<div class="card-header d-flex justify-content-between align-items-center">
		<h5>Batch {{.CurrentBatch}} of {{.TotalBatches}}</h5>
		<div class="btn-group">
			{{if gt .CurrentBatch 0}}
			<button hx-get="/batch/{{sub .CurrentBatch 1}}" hx-target="#editor-container" class="btn btn-outline-primary btn-sm">
				← Previous
			</button>
			{{end}}
			{{if lt .CurrentBatch (sub .TotalBatches 1)}}
			<button hx-get="/batch/{{add .CurrentBatch 1}}" hx-target="#editor-container" class="btn btn-outline-primary btn-sm">
				Next →
			</button>
			{{end}}
			<button id="download-n4l" class="btn btn-success btn-sm ms-2" title="Download N4L output">
				⬇️ Save N4L
			</button>
		</div>
	</div>
	<div class="card-body">
		<style>
			.pane-scroll { overflow-y: auto; }
		</style>
		<div class="row g-3">
			<div class="col-md-6">
				<h6>Original Text</h6>
				<div id="left-pane" class="border p-3 pane-scroll">
					{{range $i, $line := splitLines .OriginalText}}
					<div class="{{if contains $.AmbiguousIndices $i}}ambiguous{{end}}">{{$line}}</div>
					{{end}}
				</div>
			</div>
			<div class="col-md-6">
				<h6>N4L Output</h6>
				<div id="right-pane" class="border p-3 pane-scroll">
					<pre id="n4l-output" style="margin:0;">{{.N4LOutput}}</pre>
				</div>
			</div>
		</div>
		{{if .AmbiguousIndices}}
		<div class="alert alert-warning mt-3">
			<strong>Ambiguous lines detected:</strong> {{len .AmbiguousIndices}} lines marked for review
		</div>
		{{end}}
		<script>
		// Download N4L output as file
		document.getElementById('download-n4l').onclick = function() {
			var n4l = document.getElementById('n4l-output').innerText;
			var blob = new Blob([n4l], {type: 'text/plain'});
			var a = document.createElement('a');
			a.href = URL.createObjectURL(blob);
			a.download = 'output.n4l';
			document.body.appendChild(a);
			a.click();
			setTimeout(function() { document.body.removeChild(a); }, 100);
		};

		// Resize panes to fill viewport height
		(function() {
			const left = document.getElementById('left-pane');
			const right = document.getElementById('right-pane');
			function resizePanes() {
				if (!left || !right) return;
				const rect = left.getBoundingClientRect();
				const available = window.innerHeight - rect.top - 24; // padding at bottom
				const h = Math.max(200, available);
				left.style.height = h + 'px';
				right.style.height = h + 'px';
			}
			window.addEventListener('resize', resizePanes);
			setTimeout(resizePanes, 0);
		})();

		// Sync scroll between panes
		(function() {
			const left = document.getElementById('left-pane');
			const right = document.getElementById('right-pane');
			let syncingLeft = false, syncingRight = false;
			function sync(from, to, dir) {
				if (!from || !to) return;
				const maxFrom = from.scrollHeight - from.clientHeight;
				const maxTo = to.scrollHeight - to.clientHeight;
				const ratio = maxFrom > 0 ? (from.scrollTop / maxFrom) : 0;
				to.scrollTop = ratio * maxTo;
			}
			left.addEventListener('scroll', function() {
				if (syncingLeft) { syncingLeft = false; return; }
				syncingRight = true;
				sync(left, right, 'L');
			});
			right.addEventListener('scroll', function() {
				if (syncingRight) { syncingRight = false; return; }
				syncingLeft = true;
				sync(right, left, 'R');
			});
		})();
		</script>
	</div>
</div>`
}
