# Error Handling Modernization

This document describes the changes made to modernize error handling in the maze solver example to follow idiomatic Go practices.

## Summary

The codebase has been updated to use Go's standard error handling patterns instead of panics. The `errors` package is now used throughout, and all public functions return errors where appropriate. Additionally, the unused boolean parameter was removed from the `Open()` function to follow Go best practices (avoiding boolean parameters that make function calls unclear).

## Changes Made

### 1. Package Imports

- Added `errors` package import to `maze/graph.go`
- Already had `fmt` for error formatting with `fmt.Errorf`

### 2. Function Signature Changes

#### graph.go

- **Open()**: Removed unused boolean parameter

  - Before: `Open(loadArrows bool) *PoSST`
  - After: `Open() *PoSST`
  - Reason: Boolean parameters are considered a code smell - they make function calls unclear at the call site. The parameter was unused (reserved for database-backed implementation compatibility).

- **Edge()**: Changed from `(ArrowPtr, int)` to `(ArrowPtr, int, error)`
  - Previously: `panic("unknown arrow: " + arrow)`
  - Now: `return 0, 0, fmt.Errorf("unknown arrow: %s", arrow)`
  - Returns error instead of panicking when arrow name is not found

#### maze.go

- **SolveMaze()**: Changed from `void` to `error`
  - Now returns error from `SolveMazeWithOutput()`
- **SolveMazeWithOutput(w io.Writer)**: Changed from `void` to `error`
  - Checks error from `Edge()` calls with context wrapping
  - Checks error from `solve()` with context wrapping
  - Uses `fmt.Errorf` with `%w` verb for error wrapping
- **solve(ctx \*PoSST, w io.Writer)**: Changed from `void` to `error`
  - Returns error when start/end nodes cannot be found
  - Previously: `fmt.Fprintln(w, "No paths available from end points"); return`
  - Now: `return fmt.Errorf("no paths available from end points (start=%s, end=%s)", startBC, endBC)`

#### main.go

- **main()**: Now checks error from `maze.SolveMaze()`
  - Prints error to stderr: `fmt.Fprintf(os.Stderr, "Error: %v\n", err)`
  - Exits with code 1 on error: `os.Exit(1)`
  - Exits with code 0 on success (unchanged)

### 3. Error Wrapping

Error wrapping with `%w` is used to preserve error chains:

```go
// In SolveMazeWithOutput
if _, _, err := Edge(ctx, nfrom, "fwd", nto, context, weight); err != nil {
    return fmt.Errorf("failed to create edge from %s to %s: %w",
        path[p][leg-1], path[p][leg], err)
}

if err := solve(ctx, w); err != nil {
    return fmt.Errorf("failed to solve maze: %w", err)
}
```

This allows callers to use `errors.Is()` and `errors.As()` for error inspection.

### 4. Test Updates

#### maze_test.go

- **TestSolveMaze**: Removed panic recovery, now checks for errors

  - Before: Used `defer/recover` to catch panics
  - After: `if err := SolveMaze(); err != nil { t.Errorf(...) }`

- **TestEdgeUnknownArrow**: New test added to verify error handling

  - Tests that `Edge()` returns an error for unknown arrow names
  - Verifies error message contains "unknown arrow"

- **All Edge() calls**: Updated to check returned errors
  - Pattern: `if _, _, err := Edge(...); err != nil { t.Fatalf("Edge creation failed: %v", err) }`

## Benefits

1. **Idiomatic Go**: Follows Go best practices and conventions
2. **Better Error Context**: Errors include contextual information about what failed
3. **Error Chains**: Using `%w` allows error unwrapping and inspection
4. **Graceful Failures**: Functions can return errors instead of crashing
5. **Testable**: Error conditions can be tested without recovering from panics
6. **Debuggable**: Error messages include helpful context (node names, operation details)

## Example Usage

### Success Case

```go
if err := maze.SolveMaze(); err != nil {
    log.Fatal(err)  // Will not execute if maze solves successfully
}
```

### Error Handling

```go
err := maze.SolveMazeWithOutput(buffer)
if err != nil {
    // Handle error appropriately
    if errors.Is(err, someSpecificError) {
        // Handle specific error type
    }
    log.Printf("Maze solving failed: %v", err)
}
```

## Verification

All tests pass:

```bash
go test -C maze -v          # 13 tests pass
go test -C maze -bench=.    # 11 benchmarks pass
go build -o mazeexample     # Binary builds successfully
./mazeexample               # Runs and solves maze correctly
```

## Migration Notes

If you have existing code using this package:

1. Update all `Edge()` calls to check the third return value (error)
2. Update `SolveMaze()` and `SolveMazeWithOutput()` calls to handle returned errors
3. Remove any panic recovery code that was catching panics from `Edge()`
4. Consider adding error handling logic appropriate to your application

## Go Version

This implementation uses standard Go error handling patterns compatible with Go 1.13+, which introduced error wrapping with `%w`.
