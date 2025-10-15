# Mazeexample Package

Self-contained maze-solving implementation using bidirectional wavefront search with an in-memory semantic graph.

## Overview

This package demonstrates pathfinding through a maze using contra-colliding wavefronts. The algorithm expands search frontiers from both start and goal simultaneously until they meet, then splices the paths together.

## Features

- **In-memory semantic graph** - No external dependencies, no database
- **Bidirectional search** - Expands from both endpoints simultaneously
- **Wavefront collision detection** - Identifies when search frontiers meet
- **DAG vs cycle classification** - Distinguishes valid paths from cycles
- **Fully tested** - Comprehensive unit tests and benchmarks

## Quick Start

### Using Taskfile (Recommended)

```bash
# Show all available tasks
task --list

# Build and test
task build
task test

# Run benchmarks
task bench-core

# Full verification
task verify
```

### Using Go directly

```bash
# Build
go build -o mazeexample

# Run with text output (default)
./mazeexample

# Run with JSON output
./mazeexample --json

# Generate JSON for visualization
./mazeexample --json > results.json

# Run tests
go test

# Run benchmarks
go test -bench=. -benchmem
```

### Visualization

Generate and view results interactively:

```bash
# Easy mode - automated script
./visualize.sh

# Manual mode
./mazeexample --json > results.json
# Then open viewer.html in your web browser
```

See [JSON_VIEWER_README.md](JSON_VIEWER_README.md) for details on the web-based visualization.

## Project Structure

```
mazeexample/
├── main.go               # CLI with text/JSON output modes
├── graph.go              # In-memory SST graph implementation
├── maze.go               # Maze solving algorithm (text output)
├── maze_json.go          # JSON output implementation
├── json_output.go        # JSON types and formatting
├── maze_test.go          # Unit tests
├── maze_bench_test.go    # Benchmarks
├── viewer.html           # Interactive web visualization
├── visualize.sh          # Helper script for visualization
├── Taskfile.yml          # Task automation
├── TEST_README.md        # Detailed testing documentation
├── JSON_VIEWER_README.md # Visualization guide
└── README.md             # This file
```

## Algorithm

The maze solver uses a bidirectional wavefront collision algorithm:

1. **Initialize** - Create graph with fwd/bwd inverse arrows
2. **Build maze** - Define nodes and edges representing maze structure
3. **Expand wavefronts** - Simultaneously grow search frontiers from start and goal
4. **Detect collision** - Find nodes where wavefronts overlap
5. **Splice paths** - Combine left and right paths through collision points
6. **Classify results** - Separate DAG paths from cycles

## Core Types

- `PoSST` - In-memory semantic graph context
- `Node` - Graph vertex with name and metadata
- `Link` - Directed edge with arrow type and weight
- `NodePtr`, `ArrowPtr`, `ContextPtr` - Type-safe integer identifiers

## Core Functions

- `Open()` - Initialize new graph context
- `Vertex()` - Create or retrieve node
- `Edge()` - Add bidirectional link
- `GetEntireNCConePathsAsLinks()` - BFS path enumeration
- `SolveMaze()` - Complete maze solving example

## Performance

Baseline on Intel Core i7-6500U @ 2.50GHz:

- Node creation: ~59 ns (cached), ~237 ns (unique)
- Edge creation: ~1.2 μs
- Full maze solve: ~53.5 ms
- Memory: ~6KB for 10-node graph

See [TEST_README.md](TEST_README.md) for detailed benchmark results.

## Testing

Run the comprehensive test suite:

```bash
task test        # All tests
task test-v      # Verbose output
task test-coverage  # Generate HTML coverage report
```

All 11 tests cover:

- Graph initialization
- Node/edge operations
- Path finding and reversal
- Wavefront collision detection
- Integration testing

## Benchmarking

Measure performance:

```bash
task bench           # All benchmarks
task bench-core      # Core operations only
task bench-save      # Save baseline
task bench-compare   # Compare with baseline
```

## Example Usage

### Command Line

```bash
# Text output (default) - shows search progress
./mazeexample

# JSON output - for programmatic use or visualization
./mazeexample --json

# Save JSON to file
./mazeexample --json > results.json

# Pretty print statistics
./mazeexample --json | jq '.statistics'

# Extract solution paths
./mazeexample --json | jq '.solutions[].path'
```

### Programmatic Usage

```go
package main

import (
    "mazeexample"
    "fmt"
    "os"
)

func main() {
    // Text output to stdout
    if err := mazeexample.SolveMaze(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Text output to custom writer
    var buf bytes.Buffer
    if err := mazeexample.SolveMazeWithOutput(&buf); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // JSON output
    result, err := mazeexample.SolveMazeJSON()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("Found %d solutions\n", result.Statistics.TotalSolutions)

    // Create your own graph
    poSST, err := mazeexample.NewPoSST()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
    defer mazeexample.Close(poSST)

    // Add nodes
    a, _ := mazeexample.Vertex(poSST, "A", "chapter1")
    b, _ := mazeexample.Vertex(poSST, "B", "chapter1")

    // Connect them
    mazeexample.Edge(poSST, a, "fwd", b, []string{"path1"}, 1.0)
}
```

## Output Formats

### Text Output (Default)

Human-readable progress and results:

```
Solving maze from maze_a7 to maze_i6...
Turn 1: L:1 R:1 (0 solutions found so far)
...
Found solution at turn 25
Solution found: maze_a7 → ... → maze_i6
```

### JSON Output (--json flag)

Structured data for UIs and automation:

```json
{
  "start_node": "maze_a7",
  "end_node": "maze_i6",
  "statistics": {
    "total_solutions": 1,
    "total_loops": 0,
    "max_left_depth": 15,
    "max_right_depth": 14,
    "total_search_steps": 29
  },
  "solutions": [...],
  "search_steps": [...]
}
```

See [JSON_VIEWER_README.md](JSON_VIEWER_README.md) for complete JSON schema.

## Dependencies

- Go 1.21+ (for errors package and modern idioms)
- No external dependencies for core functionality
- Optional: `task` for task automation
- Optional: `benchstat` for benchmark comparison
- Optional: `jq` for JSON pretty-printing
- Optional: Modern web browser for visualization

## Development

```bash
# Format code
task fmt

# Run linter
task lint

# Full CI pipeline
task ci

# Watch for changes
task watch-test
```

## Documentation

- [TEST_README.md](TEST_README.md) - Comprehensive testing guide
- [Taskfile.yml](Taskfile.yml) - Task automation definitions
- Code comments - GoDoc-style documentation on all functions

## License

See project root LICENSE file.

## Related

Extracted from the SSTorytime project's API example demonstrating maze solving using Leads-To (LT) vectors in semantic graphs.
