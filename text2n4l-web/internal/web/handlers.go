package web

import (
	"net/http"

	"text2n4l-web/internal/analyzer"
)

// ConvertHandler handles text conversion requests (API only)
func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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
