# Benchmark Comparison: Handle-Based vs Pointer-Based Architecture

## Current Status

### ✅ Completed

1. **Baseline captured** - Handle-based implementation benchmarked
2. **Committed** - Current working state saved to `text2n4l` branch
3. **Branch created** - `pointer-based-graph` branch for refactoring
4. **graph.go refactored** - Complete pointer-based implementation done!

### 🚧 In Progress

- Updating maze.go, maze_json.go, and test files (~1000 lines remaining)

## Baseline Benchmarks (Handle-Based)

```
BenchmarkOpen-4                          1000           6.880 ns/op           0 B/op          0 allocs/op
BenchmarkVertex-4                        1000           111.3 ns/op           8 B/op          1 allocs/op
BenchmarkVertexUnique-4                  1000           146.9 ns/op           9 B/op          1 allocs/op
BenchmarkEdge-4                          1000           285.3 ns/op         202 B/op          0 allocs/op
BenchmarkGraphBuilding-4                 1000          13412 ns/op          5832 B/op         55 allocs/op
BenchmarkGetEntireNCConePathsAsLinks_Depth1-4
                                         1000           3964 ns/op          2024 B/op         19 allocs/op
BenchmarkGetEntireNCConePathsAsLinks_Depth5-4
                                         1000          17369 ns/op         13800 B/op         87 allocs/op
BenchmarkAdjointLinkPath-4               1000           738.9 ns/op         991 B/op          4 allocs/op
BenchmarkWaveFrontsOverlap-4             1000          32900 ns/op          4681 B/op        186 allocs/op
BenchmarkSolveMaze-4                     1000        2837632 ns/op        499420 B/op       3338 allocs/op
BenchmarkNodesOverlap-4                  1000           2045 ns/op          280 B/op          9 allocs/op
BenchmarkIsDAG-4                         1000           2801 ns/op          328 B/op          3 allocs/op
BenchmarkMemoryAllocation-4              1000          13142 ns/op          6096 B/op         48 allocs/op
```

## Architectural Comparison

### Handle-Based (Current Baseline)

```go
type LinkedSST struct {
    nextID       int
    nameTohandle map[string]NodeHandle  // String → Int
    handleToNode map[NodeHandle]Node    // Int → Node
    forward      map[NodeHandle][]Link  // Int → Links
    backward     map[NodeHandle][]Link  // Int → Links
}

// Accessing a node requires TWO map lookups:
handle := graph.nameTohandle["a7"]    // Lookup 1
node := graph.handleToNode[handle]    // Lookup 2
links := graph.forward[handle]        // Lookup 3 (for edges)
```

**Cost per operation:**

- Node access: 2 map lookups
- Link traversal: 1 additional map lookup
- Total: 3 map lookups to go from name → node → edges

### Pointer-Based (New Architecture)

```go
type LinkedSST struct {
    nodes  map[string]*Node  // String → *Node (ONE lookup!)
}

type Node struct {
    label    string
    forward  []*Link   // Embedded in node
    backward []*Link   // Embedded in node
}

// Accessing a node requires ONE map lookup:
node := graph.nodes["a7"]       // Lookup 1 - done!
links := node.forward           // Direct field access
```

**Cost per operation:**

- Node access: 1 map lookup
- Link traversal: 0 lookups (direct field access!)
- Total: 1 map lookup + field access

## Expected Performance Improvements

Based on architecture changes:

| Operation      | Handle-Based | Pointer-Based       | Expected Speedup   |
| -------------- | ------------ | ------------------- | ------------------ |
| Node lookup    | 2 maps       | 1 map               | **~40-50%** faster |
| Link access    | 3 maps       | 1 map + field       | **~60-70%** faster |
| Path traversal | Many lookups | Direct pointers     | **~30-40%** faster |
| Memory usage   | 5 maps       | 1 map + node fields | **~15-20%** less   |

### Key Metrics to Watch

**Should improve significantly:**

- `BenchmarkVertex` - Node creation (simpler)
- `BenchmarkGetEntireNCConePathsAsLinks_*` - Path finding (fewer lookups)
- `BenchmarkWaveFrontsOverlap` - Many node accesses
- `BenchmarkSolveMaze` - Overall maze solving

**Should stay about the same:**

- `BenchmarkIsDAG` - Pure algorithm, not graph-dependent
- `BenchmarkOpen` - Initialization overhead

**Memory allocations:**

- Should see reduction in `allocs/op` due to fewer map operations
- Should see reduction in `B/op` due to consolidated storage

## How to Complete and Benchmark

Once remaining files are updated:

```bash
cd /home/alex/SHDProj/SSTorytime/src/mazeexample

# Build and test
go build && go test ./maze/... -v

# Run benchmarks (same parameters as baseline)
go test ./maze/... -bench=. -benchmem -benchtime=1000x > pointer_benchmark.txt

# Compare
diff baseline_benchmark.txt pointer_benchmark.txt

# Or use benchcmp if available:
benchstat baseline_benchmark.txt pointer_benchmark.txt
```

## Next Steps

1. **Update maze.go** (~295 lines)

   - Change all `NodeHandle` → `*Node`
   - Change `[][]Link` → `[][]*Link`
   - Remove `GetDBNodeByNodeHandle()` calls

2. **Update maze_json.go** (~100 lines)

   - Similar changes to maze.go

3. **Update test files** (~600 lines total)

   - maze_test.go
   - maze_bench_test.go
   - example_buffer_test.go

4. **Run full benchmark comparison**

5. **Document results**

## Prediction

Based on the architectural improvements, I expect:

- **20-40% performance improvement** on graph-heavy operations
- **10-20% memory reduction**
- **Significantly simpler code** (fewer maps, no handle management)

The pointer-based design is **how you would naturally write this in Go** without trying to emulate a database-backed system!

---

**Files:**

- Baseline benchmarks: `baseline_benchmark.txt`
- Plan: `POINTER_REFACTORING_PLAN.md`
- Status: `POINTER_REFACTORING_STATUS.md`
- Backup: `maze/graph_handle_based.go.backup`
