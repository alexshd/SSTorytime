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
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HomeHandler)
	mux.HandleFunc("/upload", server.UploadHandler)
	mux.HandleFunc("/convert", server.ConvertHandler)
	mux.HandleFunc("/batch/", server.BatchHandler)
	// Combined editor router (handles both /combined/{batch} and /combined/{batch}/sen/{i})
	mux.HandleFunc("/combined/", server.CombinedRouter)

	// Serve static files (CSS, JS, etc.)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Setup profiling endpoints (development only)
	web.SetupProfiling(mux)

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
	fmt.Println("- Performance profiling at /debug/info")

	log.Fatal(http.ListenAndServe(port, mux))
}
