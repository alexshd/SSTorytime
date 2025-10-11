package main

import (
	"html/template"
	"net/http"
	"strings"
)

func renderDualEditor(w http.ResponseWriter) {
	sessionMu.Lock()
	if batchProcessor == nil {
		sessionMu.Unlock()
		w.Write([]byte("No file processed yet"))
		return
	}

	origBatch, genBatch, ambiguousIndices := batchProcessor.GetCurrentBatch()
	hints := batchProcessor.GetHints()
	sessionMu.Unlock()

	tmpl := template.Must(template.New("dualeditor").Parse(`
	<div class="hints">{{.Hints}}</div>
	<div class="dual-editor">
		<div class="editor-panel">
			<h3>Original (Batch {{.CurrentBatch}}/{{.TotalBatches}})</h3>
			<textarea class="editor-textarea" readonly>{{.OriginalText}}</textarea>
		</div>
		<div class="editor-panel">
			<h3>Generated N4L</h3>
			<pre id="generated-text" class="editor-textarea" contenteditable="true" name="generated" style="white-space: pre-wrap;">{{.GeneratedHTML}}</pre>
			<div style="font-size: 0.9em; color: #888;">Highlighted lines need your attention.</div>
		</div>
	</div>
	<div class="controls">
		<div class="batch-nav">
			<button hx-post="/batch/prev" hx-target="#main-content" {{if eq .CurrentBatch 1}}disabled{{end}}>← Previous</button>
			<span>Batch {{.CurrentBatch}} of {{.TotalBatches}}</span>
			<button hx-post="/batch/next" hx-target="#main-content" {{if eq .CurrentBatch .TotalBatches}}disabled{{end}}>Next →</button>
		</div>
		<div style="margin-top: 1em;">
			<button hx-post="/batch/update" hx-target="#main-content" hx-include="#generated-text">Save Changes</button>
			<button hx-get="/download" style="margin-left: 1em;">Download Final</button>
		</div>
	</div>
	`))

	// Highlight lines using real ambiguous indices instead of '???' demo
	var generatedHTML strings.Builder
	for i, line := range genBatch {
		isAmbiguous := false
		for _, ambIdx := range ambiguousIndices {
			if ambIdx == i {
				isAmbiguous = true
				break
			}
		}
		if isAmbiguous {
			generatedHTML.WriteString(`<span class="highlight-line">`)
			generatedHTML.WriteString(template.HTMLEscapeString(line))
			generatedHTML.WriteString("</span>\n")
		} else {
			generatedHTML.WriteString(template.HTMLEscapeString(line) + "\n")
		}
	}
	data := struct {
		Hints         string
		OriginalText  string
		GeneratedHTML template.HTML
		CurrentBatch  int
		TotalBatches  int
	}{
		Hints:         hints,
		OriginalText:  strings.Join(origBatch, "\n"),
		GeneratedHTML: template.HTML(generatedHTML.String()),
		CurrentBatch:  batchProcessor.CurrentBatch + 1, // 1-indexed for display
		TotalBatches:  batchProcessor.TotalBatches,
	}

	tmpl.Execute(w, data)
}

func batchPrevHandler(w http.ResponseWriter, r *http.Request) {
	sessionMu.Lock()
	if batchProcessor != nil {
		batchProcessor.PrevBatch()
	}
	sessionMu.Unlock()
	renderDualEditor(w)
}

func batchNextHandler(w http.ResponseWriter, r *http.Request) {
	sessionMu.Lock()
	if batchProcessor != nil {
		batchProcessor.NextBatch()
	}
	sessionMu.Unlock()
	renderDualEditor(w)
}

func batchUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	generatedText := r.FormValue("generated")
	lines := strings.Split(generatedText, "\n")

	sessionMu.Lock()
	if batchProcessor != nil {
		batchProcessor.UpdateBatch(lines)
	}
	sessionMu.Unlock()

	renderDualEditor(w)
}
