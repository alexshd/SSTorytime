package web

import (
	"net/http"
	_ "net/http/pprof"
)

// SetupProfiling adds pprof endpoints for performance profiling
func SetupProfiling(mux *http.ServeMux) {
	// pprof endpoints are automatically registered when importing net/http/pprof
	// Available at:
	// /debug/pprof/
	// /debug/pprof/heap
	// /debug/pprof/goroutine
	// /debug/pprof/allocs
	// /debug/pprof/block
	// /debug/pprof/mutex
	// /debug/pprof/cmdline
	// /debug/pprof/profile (30-second CPU profile)
	// /debug/pprof/symbol
	// /debug/pprof/trace

	// Add custom profile endpoint info
	mux.HandleFunc("/debug/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Profiling Endpoints</title>
    <style>
        body { font-family: system-ui, sans-serif; margin: 2rem; }
        .endpoint { margin: 1rem 0; padding: 1rem; background: #f5f5f5; border-radius: 8px; }
        code { background: #e5e5e5; padding: 2px 6px; border-radius: 4px; }
    </style>
</head>
<body>
    <h1>Performance Profiling Endpoints</h1>
    
    <div class="endpoint">
        <h3>CPU Profile</h3>
        <p>30-second CPU profile: <code><a href="/debug/pprof/profile">/debug/pprof/profile</a></code></p>
        <p>Custom duration: <code>/debug/pprof/profile?seconds=60</code></p>
    </div>
    
    <div class="endpoint">
        <h3>Memory Profiles</h3>
        <p>Heap profile: <code><a href="/debug/pprof/heap">/debug/pprof/heap</a></code></p>
        <p>Allocation profile: <code><a href="/debug/pprof/allocs">/debug/pprof/allocs</a></code></p>
    </div>
    
    <div class="endpoint">
        <h3>Goroutine Analysis</h3>
        <p>Goroutine dump: <code><a href="/debug/pprof/goroutine">/debug/pprof/goroutine</a></code></p>
        <p>Block profile: <code><a href="/debug/pprof/block">/debug/pprof/block</a></code></p>
        <p>Mutex profile: <code><a href="/debug/pprof/mutex">/debug/pprof/mutex</a></code></p>
    </div>
    
    <div class="endpoint">
        <h3>Trace Analysis</h3>
        <p>Execution trace: <code><a href="/debug/pprof/trace?seconds=5">/debug/pprof/trace?seconds=5</a></code></p>
    </div>
    
    <div class="endpoint">
        <h3>General Info</h3>
        <p>All profiles: <code><a href="/debug/pprof/">/debug/pprof/</a></code></p>
        <p>Command line: <code><a href="/debug/pprof/cmdline">/debug/pprof/cmdline</a></code></p>
    </div>
    
    <h2>Usage Examples</h2>
    <p>To analyze with go tool pprof:</p>
    <pre><code>
# CPU profile
go tool pprof http://localhost:8080/debug/pprof/profile

# Memory profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Save profile to file
curl -o cpu.prof http://localhost:8080/debug/pprof/profile
go tool pprof cpu.prof
    </code></pre>
</body>
</html>
        `))
	})
}
