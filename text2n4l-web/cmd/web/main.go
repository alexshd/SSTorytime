package main

import (
	"fmt"
	"log"
	"net/http"

	"text2n4l-web/internal/web"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/convert", web.ConvertHandler)
	port := ":8080"
	fmt.Printf("N4L API server starting on http://localhost%s\n", port)
	fmt.Println("POST /api/convert with text/plain or JSON {\"text\":...} to get N4L output.")
	log.Fatal(http.ListenAndServe(port, mux))
}
