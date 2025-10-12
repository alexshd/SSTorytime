package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"text2n4l-web/internal/web"
)

func main() {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/convert", web.ConvertHandler)

	// Profiling endpoints (accessible at /debug/pprof/)
	// For CPU profiling: curl http://localhost:8080/debug/pprof/profile?seconds=30 -o cpu.prof
	// For heap profiling: curl http://localhost:8080/debug/pprof/heap -o heap.prof
	// For goroutine info: curl http://localhost:8080/debug/pprof/goroutine
	mux.Handle("/debug/", http.DefaultServeMux)

	port := ":8080"
	fmt.Printf("N4L API server starting on http://localhost%s\n", port)
	fmt.Println("API Endpoints:")
	fmt.Println("  POST /api/convert - Convert text to N4L format")
	fmt.Println("Profiling Endpoints:")
	fmt.Println("  GET /debug/pprof/ - Profiling index")
	fmt.Println("  GET /debug/pprof/profile - CPU profile")
	fmt.Println("  GET /debug/pprof/heap - Heap profile")
	fmt.Println("  GET /debug/pprof/goroutine - Goroutine info")
	log.Fatal(http.ListenAndServe(port, mux))
}
