# N4L Text Converter API# N4L Text Converter Web Application

A clean, API-only Go backend for converting text to N4L (Narrative for Learning) format.A clean, modern refactor of the SSTorytime text-to-N4L converter with a web interface.

## Features## Features

- **REST API**: Simple HTTP POST endpoint for text conversion- **Interactive Web Interface**: Upload files or paste text for immediate conversion

- **Profiling Support**: Built-in pprof endpoints for performance analysis- **Batch Processing**: Handle large files in manageable screen-sized batches

- **Clean Architecture**: Minimal, focused Go codebase- **Dual-Pane Editor**: Side-by-side view of original text and N4L output

- **Fast Processing**: Efficient N4L conversion engine- **Smart Highlighting**: Ambiguous lines are automatically highlighted for review

- **Real-time Conversion**: Uses HTMX for responsive, fast interactions

## What is N4L?- **Clean Architecture**: Proper Go project structure following conventions

N4L (Narrative for Learning) is a domain-specific language for representing the intentional structure of text. It analyzes text using:## What is N4L?

- **Running Intentionality**: How narrative coherence develops over timeN4L (Narrative for Learning) is a domain-specific language for representing the intentional structure of text. It analyzes text using:

- **Static Intentionality**: Frequency-based significance of terms and phrases

- **Ambiguity Detection**: Identifies uncertain or low-confidence conversions- **Running Intentionality**: How narrative coherence develops over time

- **Static Intentionality**: Frequency-based significance of terms and phrases

## Project Structure- **Ambiguity Detection**: Identifies uncertain or low-confidence conversions

````## Project Structure

text2n4l-web/

├── cmd/web/           # Main API server entry point```

├── internal/          # Private application codetext2n4l-web/

│   ├── analyzer/      # N4L conversion logic├── cmd/web/           # Main application entry point

│   └── web/           # HTTP handlers├── internal/          # Private application code

├── testdata/          # Test input files│   ├── analyzer/      # N4L conversion logic

├── go.mod             # Go module definition│   └── web/           # HTTP handlers and templates

├── Makefile           # Build and test commands├── testdata/          # Test input files

└── README.md          # This file├── web/static/        # Static web assets

```├── go.mod             # Go module definition

├── Makefile           # Build and test commands

## Installation└── README.md          # This file

````

````bash

# Clone or extract to your Go workspace## Installation

cd text2n4l-web

```bash

# Download dependencies# Clone or extract to your Go workspace

go mod tidycd text2n4l-web



# Run tests# Download dependencies

make test-apigo mod tidy



# Build application# Run tests

make buildmake test

````

# Build application

## Usagemake build

### Start the API Server# Run development server

make run

`bash`

# Using make

make dev-server## Usage

# Or directly### Web Interface

./bin/text2n4l-web

```1. Start the server: `make run`

2. Open http://localhost:8080

The server will start on `http://localhost:8080`3. Upload a text file or paste content directly

4. Navigate through batches using the Previous/Next buttons

### API Endpoints5. Review highlighted ambiguous lines

6. Copy the N4L output for your analysis

#### Convert Text to N4L

### Command Line

````bash

# POST /api/convert```bash

curl -X POST http://localhost:8080/api/convert \# Build the application

  -d "text=Hello world. This is a test sentence. We analyze narrative flow."go build -o n4l-web cmd/web/main.go



# Response: N4L formatted output# Run the server

```./n4l-web

````

#### Profiling Endpoints

## API Endpoints

````bash

# CPU Profile (30 seconds)- `GET /` - Main interface

curl http://localhost:8080/debug/pprof/profile?seconds=30 -o cpu.prof- `POST /upload` - File upload handler

go tool pprof cpu.prof- `POST /convert` - Direct text conversion

- `GET /batch/{id}` - Batch processing interface

# Heap/Memory Profile

curl http://localhost:8080/debug/pprof/heap -o heap.prof## Testing

go tool pprof heap.prof

```bash

# Goroutine Info# Run all tests

curl http://localhost:8080/debug/pprof/goroutinego test ./...



# Profiling Index (web interface)# Run tests with coverage

open http://localhost:8080/debug/pprof/go test -cover ./...

````

# Run specific test

## Developmentgo test ./internal/analyzer -v

````

### Available Make Targets

## Architecture Notes

```bash

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

````

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
