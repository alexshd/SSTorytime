# N4L Text Converter Web Application# N4L Text Converter API# N4L Text Converter Web Application

A clean, modern refactor of the SSTorytime text-to-N4L converter with a web interface and streaming API.A clean, API-only Go backend for converting text to N4L (Narrative for Learning) format.A clean, modern refactor of the SSTorytime text-to-N4L converter with a web interface.

## Features## Features## Features

- **Interactive Web Interface**: Upload files or paste text for immediate conversion- **REST API**: Simple HTTP POST endpoint for text conversion- **Interactive Web Interface**: Upload files or paste text for immediate conversion

- **Streaming Output**: Real-time progressive display of N4L results

- **Dual-Pane Editor**: Side-by-side view of original text and N4L output with CodeMirror- **Profiling Support**: Built-in pprof endpoints for performance analysis- **Batch Processing**: Handle large files in manageable screen-sized batches

- **Smart Highlighting**: Ambiguous lines are automatically highlighted for review

- **Real-time Profiling**: Built-in statsviz dashboard for live performance monitoring- **Clean Architecture**: Minimal, focused Go codebase- **Dual-Pane Editor**: Side-by-side view of original text and N4L output

- **Clean Architecture**: Proper Go project structure following conventions

- **Fast Processing**: Efficient N4L conversion engine- **Smart Highlighting**: Ambiguous lines are automatically highlighted for review

## What is N4L?

- **Real-time Conversion**: Uses HTMX for responsive, fast interactions

N4L (Narrative for Learning) is a domain-specific language for representing the intentional structure of text. It analyzes text using:

## What is N4L?- **Clean Architecture**: Proper Go project structure following conventions

- **Running Intentionality**: How narrative coherence develops over time

- **Static Intentionality**: Frequency-based significance of terms and phrasesN4L (Narrative for Learning) is a domain-specific language for representing the intentional structure of text. It analyzes text using:## What is N4L?

- **Ambiguity Detection**: Identifies uncertain or low-confidence conversions

- **Running Intentionality**: How narrative coherence develops over timeN4L (Narrative for Learning) is a domain-specific language for representing the intentional structure of text. It analyzes text using:

## Quick Start

- **Static Intentionality**: Frequency-based significance of terms and phrases

`````bash

# Clone or extract to your workspace- **Ambiguity Detection**: Identifies uncertain or low-confidence conversions- **Running Intentionality**: How narrative coherence develops over time

cd text2n4l-web

- **Static Intentionality**: Frequency-based significance of terms and phrases

# Download dependencies

go mod tidy## Project Structure- **Ambiguity Detection**: Identifies uncertain or low-confidence conversions



# Build the application````## Project Structure

make build

text2n4l-web/

# Run the development server (with hot reload)

make dev-server├── cmd/web/           # Main API server entry point```

`````

├── internal/ # Private application codetext2n4l-web/

Then open:

- **Web Interface**: http://localhost:5173/│ ├── analyzer/ # N4L conversion logic├── cmd/web/ # Main application entry point

- **API Server**: http://localhost:8080

- **Profiling Dashboard**: http://localhost:8080/debug/statsviz/│ └── web/ # HTTP handlers├── internal/ # Private application code

## Project Structure├── testdata/ # Test input files│ ├── analyzer/ # N4L conversion logic

````├── go.mod             # Go module definition│   └── web/           # HTTP handlers and templates

text2n4l-web/

├── cmd/web/           # Main application entry point├── Makefile           # Build and test commands├── testdata/          # Test input files

├── internal/          # Private application code

│   ├── analyzer/      # N4L conversion logic (core algorithm)└── README.md          # This file├── web/static/        # Static web assets

│   └── web/           # HTTP handlers (API endpoints)

├── testdata/          # Test input files```├── go.mod             # Go module definition

├── docs/              # Documentation files (see below)

├── go.mod             # Go module definition├── Makefile           # Build and test commands

├── Makefile           # Build and test commands

└── README.md          # This file## Installation└── README.md          # This file

````

`````

## Documentation

````bash

- **[Algorithm Differences](docs/ALGORITHM_DIFFERENCES.md)** - Differences from original implementation

- **[Clean API Summary](docs/CLEAN_API_SUMMARY.md)** - API design and endpoints# Clone or extract to your Go workspace## Installation

- **[Fixes Applied](docs/FIXES_APPLIED.md)** - Bug fixes and corrections

- **[Go Packages Summary](docs/GO_PACKAGES_SUMMARY.md)** - Package organization and structurecd text2n4l-web

- **[Optimization Summary](docs/OPTIMIZATION_SUMMARY.md)** - Performance optimizations

- **[Package Options Analysis](docs/PACKAGE_OPTIONS_ANALYSIS.md)** - Third-party package evaluation```bash

- **[Performance Report](docs/PERFORMANCE_REPORT.md)** - Benchmarks and profiling results

- **[Sanitize Comparison](docs/SANITIZE_COMPARISON.md)** - HTML/Markdown sanitization approaches# Download dependencies# Clone or extract to your Go workspace

- **[Stdlib Packages Answer](docs/STDLIB_PACKAGES_ANSWER.md)** - Using Go standard library for sanitization

go mod tidycd text2n4l-web

## Usage



### Web Interface

# Run tests# Download dependencies

1. Start the development servers:

   ```bashmake test-apigo mod tidy

   # Terminal 1: Start backend API server (with hot reload)

   cd text2n4l-web

   make dev-server

   # Build application# Run tests

   # Terminal 2: Start frontend dev server

   cd ../text2n4l-editormake buildmake test

   npm run dev

`````

2. Open http://localhost:5173/ in your browser# Build application

3. Upload a text file or paste content directly## Usagemake build

4. Watch the N4L output stream in real-time as it's generated### Start the API Server# Run development server

5. Use the cancel button to stop long-running conversionsmake run

### API Endpoints`bash`

````bash# Using make

# Buffered conversion (waits for complete result)

curl -X POST http://localhost:8080/api/convert \make dev-server## Usage

  -d "text=Hello world. This is a test."

# Or directly### Web Interface

# Streaming conversion (progressive output)

curl -X POST http://localhost:8080/api/convert/stream \./bin/text2n4l-web

  -d "text=Hello world. This is a test." \

  --no-buffer```1. Start the server: `make run`

````

2. Open http://localhost:8080

### Profiling & Monitoring

The server will start on `http://localhost:8080`3. Upload a text file or paste content directly

```bash

# Real-time profiling dashboard (recommended)4. Navigate through batches using the Previous/Next buttons

open http://localhost:8080/debug/statsviz/

### API Endpoints5. Review highlighted ambiguous lines

# Traditional pprof endpoints

curl http://localhost:8080/debug/pprof/profile?seconds=30 -o cpu.prof6. Copy the N4L output for your analysis

curl http://localhost:8080/debug/pprof/heap -o heap.prof

go tool pprof cpu.prof#### Convert Text to N4L

```

### Command Line

## Testing

````bash

```bash

# Run all tests# POST /api/convert```bash

make test

curl -X POST http://localhost:8080/api/convert \# Build the application

# Run tests with coverage

go test -cover ./...  -d "text=Hello world. This is a test sentence. We analyze narrative flow."go build -o n4l-web cmd/web/main.go



# Run specific package tests

go test ./internal/analyzer -v

# Response: N4L formatted output# Run the server

# Benchmark performance

go test -bench=. ./internal/analyzer```./n4l-web

```

````

## Development

#### Profiling Endpoints

### Available Make Targets

## API Endpoints

`````bash

make build         # Build the application````bash

make dev-server    # Start API development server (with hot reload)

make test          # Run all tests# CPU Profile (30 seconds)- `GET /` - Main interface

make fmt           # Format code

make lint          # Run linterscurl http://localhost:8080/debug/pprof/profile?seconds=30 -o cpu.prof- `POST /upload` - File upload handler

make profile-cpu   # Capture CPU profile (server must be running)

make profile-heap  # Capture heap profile (server must be running)go tool pprof cpu.prof- `POST /convert` - Direct text conversion

make clean         # Clean build artifacts

make help          # Show all available targets- `GET /batch/{id}` - Batch processing interface

`````

# Heap/Memory Profile

### Architecture Notes

curl http://localhost:8080/debug/pprof/heap -o heap.prof## Testing

This is a clean refactor of the original SSTorytime monolithic codebase:

go tool pprof heap.prof

- **Extracted Core Logic**: Only the essential N4L conversion functions

- **Separated Concerns**: Web interface, analysis logic, and tests in separate packages```bash

- **Go Conventions**: Proper module structure, internal packages, and cmd pattern

- **Modern Stack**: Vite + CodeMirror frontend, streaming HTTP backend# Goroutine Info# Run all tests

- **Real-time Monitoring**: Integrated statsviz for live profiling visualization

curl http://localhost:8080/debug/pprof/goroutinego test ./...

## Technology Stack

**Backend:**

- **Language**: Go 1.21+# Profiling Index (web interface)# Run tests with coverage

- **Web Framework**: Standard library HTTP with streaming support

- **Profiling**: statsviz for real-time monitoring, pprof for detailed analysisopen http://localhost:8080/debug/pprof/go test -cover ./...

- **Architecture**: Clean architecture with internal packages

```

**Frontend:**

- **Build Tool**: Vite# Run specific test

- **Editor**: CodeMirror 6

- **Styling**: Tailwind CSS## Developmentgo test ./internal/analyzer -v

- **HTTP**: Fetch API with ReadableStream for streaming

```

## Original Test Compatibility

### Available Make Targets

The analyzer preserves the behavior validated by the original test suite. Test files from the parent project are included in `testdata/` to ensure correctness.

## Architecture Notes

## License

```bash

Inherits license from the parent SSTorytime project.

make build         # Build the applicationThis is a clean refactor of the original SSTorytime monolithic codebase:

make dev-server    # Start API development server

make test-api      # Run API tests- **Extracted Core Logic**: Only the essential N4L conversion functions

make test          # Run all tests- **Separated Concerns**: Web interface, analysis logic, and tests in separate packages

make fmt           # Format code- **Go Conventions**: Proper module structure, internal packages, and cmd pattern

make lint          # Run linters- **Simplified Dependencies**: Removed external complexity while preserving functionality

make profile-cpu   # Capture CPU profile (server must be running)

make profile-heap  # Capture heap profile (server must be running)## Original Test Compatibility

make clean         # Clean build artifacts

make help          # Show all available targetsThe analyzer preserves the behavior validated by the original test suite. Test files from the parent project are included in `testdata/` to ensure correctness.

```

## Development

### Running Tests

- **Language**: Go 1.21+

```bash- **Web Framework**: Standard library HTTP with HTMX

# Run API tests- **Templates**: Go html/template

make test-api- **Styling**: Bootstrap 5

- **Architecture**: Clean architecture with internal packages

# Run all tests

make test## License



# Run with verbose outputInherits license from the parent SSTorytime project.

go test -v ./internal/web
```

### Profiling

Start the server in one terminal:

```bash
make dev-server
```

In another terminal, capture profiles:

```bash
# CPU profiling
make profile-cpu

# Memory profiling
make profile-heap

# Analyze with pprof
go tool pprof cpu.prof
```

## API Response Format

The API returns N4L formatted text with semantic arrows indicating relationships between text fragments:

```
# File: uploaded.txt

{Topic}-1-(Running index: 1, Static: 0.5)[1] sentence 0 -> Hello world
{Topic}-2-(Running index: 2, Static: 0.7)[2] sentence 1 -> This is a test sentence
{Topic}-3-(Running index: 3, Static: 0.6)[3] sentence 2 -> We analyze narrative flow

# Arrows (semantic links)
[1] +CONTAINS> [2] [0.85]
[2] +CONTAINS> [3] [0.75]
```

## Architecture

- **cmd/web/main.go**: API server entry point with routing
- **internal/web/handlers.go**: HTTP request handlers
- **internal/analyzer/**: N4L conversion engine
  - `converter.go`: Main conversion logic
  - `memory.go`: Memory and context tracking
  - `annotations.go`: Semantic arrow generation

## Performance

The API is designed for high throughput:

- Efficient text processing
- Minimal memory allocations
- Built-in profiling for optimization

Use the profiling endpoints to identify bottlenecks in your workload.

## License

See LICENSE file in the root directory.
