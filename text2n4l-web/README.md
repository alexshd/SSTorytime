# N4L Text Converter Web Application

A clean, modern refactor of the SSTorytime text-to-N4L converter with a web interface.

## Features

- **Interactive Web Interface**: Upload files or paste text for immediate conversion
- **Batch Processing**: Handle large files in manageable screen-sized batches
- **Dual-Pane Editor**: Side-by-side view of original text and N4L output
- **Smart Highlighting**: Ambiguous lines are automatically highlighted for review
- **Real-time Conversion**: Uses HTMX for responsive, fast interactions
- **Clean Architecture**: Proper Go project structure following conventions

## What is N4L?

N4L (Narrative for Learning) is a domain-specific language for representing the intentional structure of text. It analyzes text using:

- **Running Intentionality**: How narrative coherence develops over time
- **Static Intentionality**: Frequency-based significance of terms and phrases
- **Ambiguity Detection**: Identifies uncertain or low-confidence conversions

## Project Structure

```
text2n4l-web/
├── cmd/web/           # Main application entry point
├── internal/          # Private application code
│   ├── analyzer/      # N4L conversion logic
│   └── web/           # HTTP handlers and templates
├── testdata/          # Test input files
├── web/static/        # Static web assets
├── go.mod             # Go module definition
├── Makefile           # Build and test commands
└── README.md          # This file
```

## Installation

```bash
# Clone or extract to your Go workspace
cd text2n4l-web

# Download dependencies
go mod tidy

# Run tests
make test

# Build application
make build

# Run development server
make run
```

## Usage

### Web Interface

1. Start the server: `make run`
2. Open http://localhost:8080
3. Upload a text file or paste content directly
4. Navigate through batches using the Previous/Next buttons
5. Review highlighted ambiguous lines
6. Copy the N4L output for your analysis

### Command Line

```bash
# Build the application
go build -o n4l-web cmd/web/main.go

# Run the server
./n4l-web
```

## API Endpoints

- `GET /` - Main interface
- `POST /upload` - File upload handler
- `POST /convert` - Direct text conversion
- `GET /batch/{id}` - Batch processing interface

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/analyzer -v
```

## Architecture Notes

This is a clean refactor of the original SSTorytime monolithic codebase:

- **Extracted Core Logic**: Only the essential N4L conversion functions
- **Separated Concerns**: Web interface, analysis logic, and tests in separate packages
- **Go Conventions**: Proper module structure, internal packages, and cmd pattern
- **Simplified Dependencies**: Removed external complexity while preserving functionality

## Original Test Compatibility

The analyzer preserves the behavior validated by the original test suite. Test files from the parent project are included in `testdata/` to ensure correctness.

## Development

- **Language**: Go 1.21+
- **Web Framework**: Standard library HTTP with HTMX
- **Templates**: Go html/template
- **Styling**: Bootstrap 5
- **Architecture**: Clean architecture with internal packages

## License

Inherits license from the parent SSTorytime project.
