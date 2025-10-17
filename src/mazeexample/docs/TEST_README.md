# Mazeexample Package Tests and Benchmarks

This document describes the test suite and benchmarks for the mazeexample package, which implements an in-memory semantic graph for maze solving using bidirectional wavefront search.

## Test Coverage

### Unit Tests (maze_test.go)

All tests pass successfully:

1. **TestOpen** - Verifies PoSST context initialization

   - Checks nextID starts at 1
   - Validates map initialization (name2ptr, ptr2node, inverse)
   - Confirms forward/backward arrow mapping (fwd ↔ bwd)

2. **TestVertex** - Tests node creation and retrieval

   - Validates idempotent node creation (same name returns same NodePtr)
   - Checks unique nodes get different NodePtrs
   - Verifies node metadata (name, chapter)

3. **TestEdge** - Verifies link creation between nodes

   - Tests forward link in `out` map
   - Tests backward link in `in` map with inverse arrow
   - Validates link weight and arrow type

4. **TestGetEntireNCConePathsAsLinks** - BFS path enumeration

   - Tests finding paths at exact depths
   - Validates path count and length
   - Tests both forward and backward traversal

5. **TestAdjointLinkPath** - Path reversal with arrow inversion

   - Verifies link order reversal
   - Checks arrow inversion (fwd → bwd)
   - Validates reversed destinations

6. **TestWaveFrontsOverlap** - Frontier collision detection

   - Tests diamond pattern pathfinding
   - Verifies DAG path detection
   - Confirms no false cycle detection

7. **TestSolveMaze** - Integration test

   - Runs full maze solving algorithm
   - Verifies no panics
   - Tests complete bidirectional search workflow

8. **TestGraphIsolation** - Context independence

   - Validates separate contexts don't interfere
   - Tests independent node namespaces

9. **TestMultipleEdges** - Multiple edges between same nodes

   - Verifies handling of duplicate edges
   - Tests edge accumulation in adjacency lists

10. **TestEmptyGraph** - Edge case handling

    - Tests path finding in graph with no edges
    - Validates empty result sets

11. **TestPathLimit** - Respects search limits
    - Tests maximum path count enforcement
    - Verifies limit parameter functionality

## Benchmark Results

### Performance Baseline (Intel Core i7-6500U @ 2.50GHz)

| Benchmark                         | Iterations | Time/op | Memory/op | Allocs/op |
| --------------------------------- | ---------- | ------- | --------- | --------- |
| BenchmarkOpen                     | ~320k      | 4.9 μs  | 976 B     | 11        |
| BenchmarkVertex (cached)          | ~20M       | 59 ns   | 0 B       | 0         |
| BenchmarkVertexUnique             | ~4M        | 237 ns  | 4 B       | 1         |
| BenchmarkEdge                     | ~954k      | 1.2 μs  | 335 B     | 0         |
| BenchmarkGraphBuilding (10 nodes) | ~53k       | 21 μs   | 5.8 KB    | 55        |
| BenchmarkSolveMaze (full maze)    | 100        | 53.5 ms | 256 KB    | 2255      |
| BenchmarkMemoryAllocation         | ~168k      | 6.7 μs  | 6.1 KB    | 48        |

### Key Performance Insights

1. **Node Lookups**: Cached vertex lookups are extremely fast (59 ns) due to map-based storage
2. **Graph Construction**: Building a 10-node chain takes ~21 μs, showing efficient adjacency list management
3. **Maze Solving**: Full bidirectional search completes in ~50ms, demonstrating practical performance
4. **Memory Efficiency**: Graph building allocates ~6KB for 10-node maze structure

## Running Tests

### Using Taskfile (Recommended)

```bash
# Show all available tasks
task --list

# Quick start
task build       # Build the package
task test        # Run all tests
task bench       # Run benchmarks

# See detailed help
task help
```

### Using Go directly

```bash
# Run all tests with verbose output
go test -v

# Run tests without verbose maze solving output
go test

# Run specific test
go test -run TestVertex

# Run with coverage
go test -cover
```

## Running Benchmarks

### Using Taskfile (Recommended)

```bash
# Run all benchmarks
task bench

# Run core benchmarks only (faster)
task bench-core

# Save baseline for regression detection
task bench-save

# Compare with baseline (requires benchstat)
task bench-compare

# Run with profiling
task bench-cpu    # CPU profiling
task bench-mem    # Memory profiling
```

### Using Go directly

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkVertex -benchmem

# Run with longer duration for more accurate results
go test -bench=. -benchmem -benchtime=5s

# Save benchmark results for comparison
go test -bench=. -benchmem > benchmark_results.txt
```

## Regression Detection

### Using Taskfile (Recommended)

```bash
# Save current baseline
task bench-save

# After code changes, compare results
task bench-compare
```

### Using Go directly

To detect performance regressions:

```bash
# Save current baseline
go test -bench=. -benchmem > benchmark_baseline.txt

# After code changes, compare results
go test -bench=. -benchmem > benchmark_new.txt
benchstat benchmark_baseline.txt benchmark_new.txt
```

## Additional Taskfile Commands

```bash
# Development workflow
task verify      # Full verification (build + test + bench)
task ci          # CI pipeline with coverage
task lint        # Check code formatting and vet

# Testing variants
task test-v              # Verbose test output
task test-coverage       # Generate HTML coverage report
task test-coverage-func  # Show coverage by function
task test-run TEST=TestVertex  # Run specific test

# Benchmarking variants
task bench-run BENCH=BenchmarkVertex  # Run specific benchmark
task bench-time TIME=5s               # Custom benchmark duration

# Cleanup
task clean       # Remove test artifacts
```

## Test Data

The test suite uses:

- Simple linear graphs (A → B → C)
- Diamond patterns (Start → {A,B} → End)
- Binary trees (for depth testing)
- Full maze graph (9 path segments, bidirectional search)

## Assertions Tested

- ✓ Graph initialization and cleanup
- ✓ Node uniqueness and retrieval
- ✓ Bidirectional edge creation
- ✓ BFS path enumeration with depth limits
- ✓ Path reversal and arrow inversion
- ✓ Wavefront collision detection
- ✓ DAG vs cycle classification
- ✓ Context isolation
- ✓ Memory allocation patterns

## Notes

- All tests are deterministic and repeatable
- No external dependencies required
- Tests run in parallel where safe
- Benchmarks exclude test output noise (use -run=^$ flag)
- Maze solving produces verbose output showing wavefront expansion

## Last Updated

October 15, 2025 - Initial test suite creation
