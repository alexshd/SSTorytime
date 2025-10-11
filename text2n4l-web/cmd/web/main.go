package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"text2n4l-web/internal/web"
)

func main() {
	// Create server instance
	server := web.NewServer()

	// Setup routes
	http.HandleFunc("/", server.HomeHandler)
	http.HandleFunc("/upload", server.UploadHandler)
	http.HandleFunc("/convert", server.ConvertHandler)
	http.HandleFunc("/batch/", server.BatchHandler)

	// Serve static files (CSS, JS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Touch reload.txt for browser auto-refresh
	reloadFile := "tmp/reload.txt"
	f, err := os.OpenFile(reloadFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err == nil {
		f.WriteString(time.Now().Format(time.RFC3339Nano))
		f.Close()
	}

	// Start server
	port := ":8080"
	fmt.Printf("N4L Text Converter server starting on http://localhost%s\n", port)
	fmt.Println("Features:")
	fmt.Println("- Interactive text to N4L conversion")
	fmt.Println("- File upload with batch processing")
	fmt.Println("- Real-time ambiguous line highlighting")
	fmt.Println("- Dual-pane editor view")

	log.Fatal(http.ListenAndServe(port, nil))
}
