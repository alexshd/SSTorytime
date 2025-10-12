package web

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newServerForTest() *Server {
	tmpl := template.Must(template.New("combined-editor.tmpl").Parse(`
<div id="combined-root">Combined Editor
  <div id="orig-pane"></div>
  <div id="details-pane"></div>
</div>`))
	tmpl = template.Must(tmpl.New("sentence-details.tmpl").Parse(`
<div>@sen{{.Index}}</div>`))
	return &Server{BatchSize: 50, templates: tmpl}
}

func withTestFile(content string) func() {
	s := newServerForTest()
	batches := s.createBatches(content)
	currentFile = &FileUploadData{Filename: "test.txt", Content: content, Batches: batches}
	return func() { currentFile = nil }
}

func TestCombinedBatchHandler_RendersPartial(t *testing.T) {
	teardown := withTestFile("First sentence. Second sentence.")
	defer teardown()

	s := newServerForTest()
	req := httptest.NewRequest("GET", "/combined/0", nil)
	w := httptest.NewRecorder()
	s.CombinedRouter(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Combined Editor") {
		t.Fatalf("expected Combined Editor template, got: %s", body[:min(200, len(body))])
	}
	if strings.Contains(body, "Upload Text File") {
		t.Fatalf("unexpected index content in combined editor response")
	}
}

func TestSentenceDetailsHandler_RendersDetails(t *testing.T) {
	teardown := withTestFile("First sentence. Second sentence.")
	defer teardown()

	s := newServerForTest()
	req := httptest.NewRequest("GET", "/combined/0/sen/0", nil)
	w := httptest.NewRecorder()
	s.CombinedRouter(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "@sen0") {
		t.Fatalf("expected details for @sen0, got: %s", body[:min(200, len(body))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
