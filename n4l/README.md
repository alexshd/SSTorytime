# N4L Package

This package contains the N4L (Notes for Learning) command-line tool and related components.

## Contents

- `N4L.go` - The main N4L parser and compiler
- `sst.go` - SSTorytime library with graph processing functions
- `sst/superduper.go` - Copy of the SSTorytime library
- `tests/` - Test files for validating N4L functionality
- `run_tests` - Script to run the test suite
- `docs/` - Documentation files

## Documentation

- **[N4L Language Guide](docs/N4L.md)** - Complete reference for the N4L language
- **[Formatter Specification](docs/N4L_FORMATTER.md)** - Code formatting rules and guidelines
- **[Syntax Highlighting](docs/N4L_SYNTAX_HIGHLIGHTING.md)** - Editor integration and highlighting guide

## N4L Command Line Tool

The N4L tool ingests files of "notes" written in a simple language and turns them into a machine representation in the form of a "Semantic Spacetime" graph (a form of knowledge graph).

### Usage

```bash
# Parse and validate a file
./N4L file.dat

# Multiple files
./N4L file1.dat file2.dat file3.dat

# Verbose output
./N4L -v file.dat

# Upload to database
./N4L -u file.dat

# Extract subgraph with specific relations
./N4L -v -s -adj="pe,he" file.dat

# Show help
./N4L -h
```

### Command Options

- `-adj string` - Comma-separated list of short link names (default "none")
- `-d` - Diagnostic mode
- `-s` - Summary (node, links...)
- `-u` - Upload mode
- `-v` - Verbose output

## Building and Testing

Use the provided Taskfile:

```bash
# Build the N4L binary
task build

# Run tests
task test

# Build and test
task all

# Clean up
task clean

# Show all available tasks
task help
```

## Testing

Run the test suite:

```bash
task test
# or
./run_tests
```

This will test various N4L files and validate that the parser works correctly.

## Documentation

See the included markdown files for detailed documentation:

- **[N4L Language Guide](docs/N4L.md)** - Complete reference for the N4L language
- **[Formatter Specification](docs/N4L_FORMATTER.md)** - Code formatting rules and guidelines
- **[Syntax Highlighting](docs/N4L_SYNTAX_HIGHLIGHTING.md)** - Editor integration and highlighting guide
