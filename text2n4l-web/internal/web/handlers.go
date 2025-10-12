package web

import (
	"net/http"

	"text2n4l-web/internal/analyzer"
)

// ConvertHandler handles text conversion requests (API only)
func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "No text provided", http.StatusBadRequest)
		return
	}

	n4lOut := analyzer.N4LSkeletonOutput("uploaded.txt", text, 100.0)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(n4lOut))
}
