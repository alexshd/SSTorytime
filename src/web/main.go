package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	SST "SSTorytime"
)

// In-memory session state (single user demo)
var (
	sessionMu      sync.Mutex
	sessionData    *Text2N4LResult
	batchProcessor *BatchProcessor
)

// Text2N4LResult holds the original and generated lines, and ambiguous indices
type Text2N4LResult struct {
	Original  []string
	Generated []string
	Ambiguous []int // line indices in Generated
}

// selectByRunningIntent ranks sentences by their running intentionality score
func selectByRunningIntent(psf [][][]string, L int, percentage float64) []SST.TextRank {
	const coherence_length = SST.DUNBAR_30
	var sentences []SST.TextRank
	var sentence_counter int
	for p := range psf {
		for s := 0; s < len(psf[p]); s++ {
			score := 0.0
			var sb strings.Builder
			for f := 0; f < len(psf[p][s]); f++ {
				score += SST.RunningIntentionality(sentence_counter, psf[p][s][f])
				sb.WriteString(psf[p][s][f])
				if f < len(psf[p][s])-1 {
					sb.WriteString(", ")
				}
			}
			var this SST.TextRank
			this.Fragment = sb.String()
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences, this)
			sentence_counter++
		}
	}
	return orderAndRank(sentences, percentage)
}

// selectByStaticIntent ranks sentences by their static intentionality score
func selectByStaticIntent(psf [][][]string, L int, percentage float64) []SST.TextRank {
	const coherence_length = SST.DUNBAR_30
	var sentences []SST.TextRank
	var sentence_counter int
	for p := range psf {
		for s := 0; s < len(psf[p]); s++ {
			score := 0.0
			var sb strings.Builder
			for f := 0; f < len(psf[p][s]); f++ {
				score += SST.AssessStaticIntent(psf[p][s][f], L, SST.STM_NGRAM_FREQ, 1)
				sb.WriteString(psf[p][s][f])
				if f < len(psf[p][s])-1 {
					sb.WriteString(", ")
				}
			}
			var this SST.TextRank
			this.Fragment = sb.String()
			this.Significance = score
			this.Order = sentence_counter
			this.Partition = sentence_counter / coherence_length
			sentences = append(sentences, this)
			sentence_counter++
		}
	}
	return orderAndRank(sentences, percentage)
}

// orderAndRank orders sentences by significance and returns the top percentage
func orderAndRank(sentences []SST.TextRank, percentage float64) []SST.TextRank {
	var selections []SST.TextRank
	sort.Slice(sentences, func(i, j int) bool {
		return sentences[i].Significance > sentences[j].Significance
	})
	limit := int((percentage / 100.0) * float64(len(sentences)))
	if limit == 0 && len(sentences) > 0 {
		limit = 1
	}
	for i := 0; i < limit; i++ {
		selections = append(selections, sentences[i])
	}
	sort.Slice(selections, func(i, j int) bool {
		return selections[i].Order < selections[j].Order
	})
	return selections
}

// mergeSelections merges two slices of TextRank, preserving order and uniqueness
func mergeSelections(one []SST.TextRank, two []SST.TextRank) []SST.TextRank {
	var merge []SST.TextRank
	alreadySelected := make(map[int]bool)
	for i := range one {
		merge = append(merge, one[i])
		alreadySelected[one[i].Order] = true
	}
	for i := range two {
		if !alreadySelected[two[i].Order] {
			merge = append(merge, two[i])
		}
	}
	sort.Slice(merge, func(i, j int) bool {
		return merge[i].Order < merge[j].Order
	})
	return merge
}

// ConvertTextToN4LResult runs the real N4L conversion and returns generated lines and ambiguous indices.
func ConvertTextToN4LResult(inputPath string, percentage float64) (*Text2N4LResult, error) {
	SST.MemoryInit()
	psf, L := SST.FractionateTextFile(inputPath)
	ranking1 := selectByRunningIntent(psf, L, percentage)
	ranking2 := selectByStaticIntent(psf, L, percentage)
	selection := mergeSelections(ranking1, ranking2)

	const minN = 1
	const maxN = 3
	f, s, ff, ss := SST.ExtractIntentionalTokens(L, selection, minN, maxN)

	// Generate output lines in N4L DSL format
	var origLines []string
	var genLines []string
	var ambiguous []int

	// Read original file for Original field
	origContent := SST.ReadFile(inputPath)
	origLines = strings.Split(origContent, "\n")

	// Generate N4L DSL output similar to WriteOutput in text2N4L.go
	for i, sel := range selection {
		// Format as N4L DSL with @sen prefix
		line := "@sen" + strings.Join([]string{sel.Fragment}, " ")
		genLines = append(genLines, line)

		// Mark lines with ambiguous context as needing attention
		partIdx := sel.Partition
		if partIdx < len(s) && len(s[partIdx]) > 0 {
			ambiguous = append(ambiguous, i)
		}

		// Add contextual information
		if partIdx < len(f) {
			for _, ctx := range f[partIdx] {
				genLines = append(genLines, "  # INTENT "+ctx)
			}
		}
	}

	// Add global context
	for _, ctx := range ff {
		genLines = append(genLines, " # "+ctx)
	}
	for _, ctx := range ss {
		genLines = append(genLines, "  # "+ctx)
	}

	return &Text2N4LResult{Original: origLines, Generated: genLines, Ambiguous: ambiguous}, nil
}

func main() {
	StartProfiler() // Start pprof profiler
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/batch/prev", batchPrevHandler)
	http.HandleFunc("/batch/next", batchNextHandler)
	http.HandleFunc("/batch/update", batchUpdateHandler)
	http.HandleFunc("/next", nextAmbiguityHandler)
	http.HandleFunc("/clarify", clarifyHandler)
	http.HandleFunc("/download", downloadHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := ":" + port
	server := &http.Server{Addr: addr}

	log.Printf("[INFO] Server starting at http://localhost:%s", port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("[INFO] Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("[ERROR] Server forced to shutdown: %v", err)
		}
		log.Println("[INFO] Server exited gracefully.")
		os.Exit(0)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[ERROR] ListenAndServe: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] Request for: %s", r.URL.Path)
	log.Printf("[DEBUG] Trying to serve: ./static/index.html")

	// Check if file exists
	if _, err := os.Stat("./static/index.html"); os.IsNotExist(err) {
		log.Printf("[ERROR] File not found: ./static/index.html")
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "./static/index.html")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form: " + err.Error()))
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get file: " + err.Error()))
		return
	}
	defer file.Close()

	tmpDir := "./web/tmp"
	os.MkdirAll(tmpDir, 0o755)
	tmpFilePath := tmpDir + "/" + handler.Filename
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create temp file: " + err.Error()))
		return
	}
	defer tmpFile.Close()

	_, err = tmpFile.ReadFrom(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to save file: " + err.Error()))
		return
	}

	// Use in-process Go function for text2n4l with real N4L DSL output
	result, err := ConvertTextToN4LResult(tmpFilePath, 50.0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("text2n4l error: " + err.Error()))
		return
	}

	sessionMu.Lock()
	sessionData = result
	batchProcessor = NewBatchProcessorWithAmbiguous(result.Original, result.Generated, result.Ambiguous)
	sessionMu.Unlock()

	renderDualEditor(w)
}

func nextAmbiguityHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Return next ambiguous section for clarification
	w.Write([]byte("Next ambiguity endpoint (to be implemented)"))
}

func clarifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
	lineIdxStr := r.FormValue("line")
	clar := r.FormValue("clarification")
	idx, err := strconv.Atoi(lineIdxStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid line index"))
		return
	}
	sessionMu.Lock()
	if sessionData != nil && idx >= 0 && idx < len(sessionData.Generated) {
		sessionData.Generated[idx] = clar
	}
	sessionMu.Unlock()
	// Re-render the split view with updated data (if needed)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Serve the final clarified file
	w.Write([]byte("Download endpoint (to be implemented)"))
}
