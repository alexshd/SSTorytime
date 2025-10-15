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
go build

# Run tests
go test

# Run benchmarks
go test -bench=. -benchmem
```

## Project Structure

```
mazeexample/
├── graph.go              # In-memory SST graph implementation
├── maze.go               # Maze solving algorithm
├── maze_test.go          # Unit tests
├── maze_bench_test.go    # Benchmarks
├── Taskfile.yml          # Task automation
├── TEST_README.md        # Detailed testing documentation
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

```go
package main

import "mazeexample"

func main() {
    // Run the maze solving example
    mazeexample.SolveMaze()

    // Or create your own graph
    ctx := mazeexample.Open(false)
    defer mazeexample.Close(ctx)

    // Add nodes
    a := mazeexample.Vertex(ctx, "A", "chapter1")
    b := mazeexample.Vertex(ctx, "B", "chapter1")

    // Connect them
    mazeexample.Edge(ctx, a, "fwd", b, []string{"path1"}, 1.0)
}
```

## Dependencies

- Go 1.16+ (for any Go features used)
- No external dependencies for core functionality
- Optional: `task` for task automation
- Optional: `benchstat` for benchmark comparison

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
