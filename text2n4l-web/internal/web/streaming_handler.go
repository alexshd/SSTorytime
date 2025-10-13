package web

import (
	"bufio"
	"net/http"

	"text2n4l-web/internal/analyzer"
)

// StreamingConvertHandler handles text conversion with streaming output
func StreamingConvertHandler(w http.ResponseWriter, r *http.Request) {
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

	// Set headers for streaming
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Get flusher for immediate sending
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Create buffered writer for efficient streaming
	bw := bufio.NewWriter(w)

	// Stream the N4L output as it's generated
	err := analyzer.StreamN4LOutput(bw, flusher, "uploaded.txt", text, 100.0)
	if err != nil {
		// Error already written to stream or client disconnected
		return
	}

	// Final flush
	bw.Flush()
	flusher.Flush()
}
