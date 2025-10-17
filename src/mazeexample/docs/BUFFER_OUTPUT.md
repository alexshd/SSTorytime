# Buffer Output Usage

## Overview

The `mazeexample` package now supports buffered output, allowing you to capture maze-solving results to a `bytes.Buffer` or any `io.Writer` instead of printing directly to stdout.

## API

### Functions

```go
// SolveMaze prints output to os.Stdout (default behavior)
func SolveMaze()

// SolveMazeWithOutput writes output to the provided writer
func SolveMazeWithOutput(w io.Writer)
```

## Use Cases

### 1. Capture Output for Testing

```go
func TestMazeOutput(t *testing.T) {
    var buf bytes.Buffer
    SolveMazeWithOutput(&buf)

    output := buf.String()

    // Verify output contains expected markers
    if !strings.Contains(output, "-- T R E E --") {
        t.Error("Expected solution marker in output")
    }
}
```

### 2. Suppress Output During Benchmarks

```go
func BenchmarkSolving(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Discard output to avoid I/O overhead in benchmarks
        SolveMazeWithOutput(io.Discard)
    }
}
```

### 3. Write Output to File

```go
func SaveMazeSolution() {
    f, err := os.Create("maze_solution.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    SolveMazeWithOutput(f)
}
```

### 4. Analyze Output Programmatically

```go
func AnalyzeSolution() {
    var buf bytes.Buffer
    SolveMazeWithOutput(&buf)

    output := buf.String()

    // Count solutions found
    solutionCount := strings.Count(output, "-- T R E E --")
    fmt.Printf("Found %d solutions\n", solutionCount)

    // Extract path length
    if strings.Contains(output, "with lengths") {
        // Parse length information
    }

    // Check if specific node was visited
    if strings.Contains(output, "maze_e3") {
        fmt.Println("Path passes through maze_e3")
    }
}
```

### 5. Combine Multiple Writers

```go
func LogAndSave() {
    // Write to both stdout AND a file simultaneously
    f, _ := os.Create("solution.log")
    defer f.Close()

    multiWriter := io.MultiWriter(os.Stdout, f)
    SolveMazeWithOutput(multiWriter)
}
```

## Implementation Details

The following functions now accept an `io.Writer` parameter:

- `solve(ctx *PoSST, w io.Writer)`
- `waveFrontsOverlap(ctx *PoSST, w io.Writer, ...)`
- `nodesOverlap(ctx *PoSST, w io.Writer, ...)`

All output formatting uses `fmt.Fprintf`, `fmt.Fprintln`, and `fmt.Fprint` to write to the provided writer.

## Backward Compatibility

The original `SolveMaze()` function remains unchanged and prints to `os.Stdout`, ensuring backward compatibility with existing code.

## Testing

See `example_buffer_test.go` for a complete working example of buffered output testing.

```bash
# Run the buffered output test
go test -v -run TestSolveMazeBufferedOutput

# Run all tests
go test -v
```
