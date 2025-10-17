# Pointer-Based Refactoring - Initial Assessment

## Status: PARTIALLY IMPLEMENTED - Needs Completion

## What Was Done

1. ✅ Committed baseline Handle-based implementation
2. ✅ Captured baseline benchmarks
3. ✅ Created new branch: `pointer-based-graph`
4. ✅ Created complete pointer-based `graph.go` with new architecture
5. ✅ Documented refactoring plan

## Architecture Changes Completed in graph.go

### Type Simplifications

**Removed:**

- `NodeHandle int` - replaced with `*Node` pointers
- `ArrowHandle int` - replaced with `*Arrow` pointers
- `ContextHandle int` - replaced with `string` (simplified)
- `ArrowDirectory` - replaced with `Arrow` struct

**New Structures:**

```go
type LinkedSST struct {
    nodes   map[string]*Node   // Direct pointer lookup
    arrows  []*Arrow           // Direct arrow pointers
    arrow   map[string]*Arrow  // Name to arrow
    inverse map[*Arrow]*Arrow  // Arrow inversions
}

type Node struct {
    label    string
    forward  []*Link   // Adjacency list IN node
    backward []*Link   // Adjacency list IN node
}

type Link struct {
    arrow  *Arrow   // Direct pointer
    dst    *Node    // Direct pointer
    weight float32
}
```

### Functions Refactored in graph.go

- ✅ `NewLinkedSST()` - Simplified initialization
- ✅ `Vertex()` - Returns `*Node`, single map operation
- ✅ `Edge()` - Takes `*Node` pointers, stores in node adjacency lists
- ✅ `GetNodeByName()` - Direct map lookup, returns `*Node`
- ✅ `AdjointLinkPath()` - Works with `[]*Link`
- ✅ `GetEntireNCConePathsAsLinks()` - Takes `*Node`, returns `[][]*Link`
- ✅ `PrintLinkPath()` - Works with `[][]*Link`

### Eliminated Functions

- ❌ `GetNodeHandleMatchingName()` - Replaced by `GetNodeByName()`
- ❌ `GetDBNodeByNodeHandle()` - Not needed (direct pointers!)
- ❌ `GetDBArrowByHandle()` - Not needed (link.arrow is direct!)

## Remaining Work

###Files Needing Updates:

1. **maze/maze.go** (295 lines) - Primary maze solver

   - Update `solve()` to use `*Node` instead of handles
   - Update `waveFrontsOverlap()` to use `[]*Link`
   - Update `showNode()` to take `[]*Node`
   - Update `showNodePath()` to take `[]*Link`
   - Remove calls to `GetDBNodeByNodeHandle()`
   - Replace `GetNodeHandleMatchingName()` with `GetNodeByName()`

2. **maze/maze_json.go** - JSON export

   - Update `getFrontierNodes()` to work with `[]*Node`
   - Update `linksToPathLinks()` to work with `[]*Link`
   - Remove handle-based function calls

3. **maze/maze_test.go** (378 lines) - 13 test functions

   - Update all tests to use `*Node` pointers
   - Remove `node.NHandle` references
   - Update comparisons to use pointer equality
   - Remove `GetDBNodeByNodeHandle()` calls

4. **maze/maze_bench_test.go** (247 lines) - 11 benchmarks

   - Same updates as test file
   - Ensure benchmark setup uses new API

5. **maze/example_buffer_test.go** - Simple test
   - Should work with minimal changes

## Expected Benefits

### Code Simplification

- **Removed:** ~100 lines of handle management code
- **Simpler API:** Direct pointers instead of handle lookups
- **Better locality:** Adjacency lists in nodes, not separate maps

### Performance Improvements (Expected)

Based on architecture:

- **Node access:** Eliminate 1 map lookup → ~20-30% faster
- **Link traversal:** Direct pointer dereference → ~10-20% faster
- **Memory:** Fewer maps, better cache locality

### Memory Impact (Expected)

- **Removed:** 4 maps (nameTohandle, handleToNode, forward, backward)
- **Added:** 2 slices per node (forward, backward)
- **Net:** Likely 10-20% memory reduction for typical graphs

## How to Complete

To finish this refactoring:

```bash
# Already on pointer-based-graph branch
cd /home/alex/SHDProj/SSTorytime/src/mazeexample

# Update maze.go - largest file
# 1. Change NodeHandle to *Node throughout
# 2. Change [][]Link to [][]*Link
# 3. Update function calls to new API
# 4. Remove GetDBNodeByNodeHandle calls

# Update maze_json.go
# Similar changes as maze.go

# Update all test files
# Change assertions to use pointer equality
# Update to new function signatures

# Build and test
go build && go test ./maze/...

# Run benchmarks and compare
go test ./maze/... -bench=. -benchmem -benchtime=1000x > pointer_benchmark.txt

# Compare with baseline
# Compare baseline_benchmark.txt vs pointer_benchmark.txt
```

## Rollback Available

The Handle-based version is preserved:

- Committed on `text2n4l` branch
- Also backed up as `maze/graph_handle_based.go.backup`

To rollback:

```bash
git checkout text2n4l
```

## Conclusion

The pointer-based architecture is **significantly simpler** and should be **faster**. The core graph.go refactoring is complete and demonstrates the improved design.

Completing the refactoring requires updating ~1000 lines of code across 5 files to use the new pointer-based API. This is straightforward but time-consuming mechanical work.

**Recommendation:** Complete this refactoring - the benefits are substantial and the design is much cleaner!
