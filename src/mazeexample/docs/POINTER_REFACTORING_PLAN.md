# Pointer-Based Refactoring Plan

## Goal

Replace handle-based indirection with native Go pointers for simpler, faster, more idiomatic code.

## Current Architecture (Handle-Based)

```go
type NodeHandle int
type LinkedSST struct {
    nameTohandle map[string]NodeHandle  // name → int
    handleToNode map[NodeHandle]Node    // int → Node
    forward      map[NodeHandle][]Link  // int → []Link
    backward     map[NodeHandle][]Link  // int → []Link
}
type Link struct {
    Dst NodeHandle  // int reference
    Arr ArrowHandle // int reference
}
// Usage: node = graph.handleToNode[nodeHandle]  // Extra lookup!
```

## New Architecture (Pointer-Based)

```go
type LinkedSST struct {
    nodes   map[string]*Node           // name → *Node (direct!)
    arrows  []*Arrow                   // Direct arrow pointers
    arrow   map[string]*Arrow          // name → *Arrow
    inverse map[*Arrow]*Arrow          // *Arrow → *Arrow
    context map[string]*Context        // context strings
}

type Node struct {
    label   string
    chap    string
    forward []*Link    // Adjacency list IN the node!
    backward []*Link   // Adjacency list IN the node!
}

type Link struct {
    arrow  *Arrow     // Direct pointer!
    dst    *Node      // Direct pointer to destination!
    weight float32
    ctx    *Context
}

type Arrow struct {
    long     string
    short    string
    stIndex  int
}

// Usage: link.dst  // Direct access - NO lookup needed!
```

## Key Improvements

### 1. Eliminate Handle Lookups

**Before:**

```go
dstNode := graph.handleToNode[link.Dst]  // Map lookup
```

**After:**

```go
dstNode := link.dst  // Direct pointer - instant!
```

### 2. Adjacency Lists Move to Nodes

**Before:**

```go
links := graph.forward[nodeHandle]  // Separate map
```

**After:**

```go
links := node.forward  // Part of the node itself!
```

### 3. Simpler Node Creation

**Before:**

```go
handle := NodeHandle(graph.nextID)
graph.nextID++
node := Node{...}
graph.nameTohandle[name] = handle
graph.handleToNode[handle] = node
```

**After:**

```go
node := &Node{...}
graph.nodes[name] = node  // One operation!
```

### 4. Type Simplification

**Remove entirely:**

- `NodeHandle` type
- `ArrowHandle` type
- `ContextHandle` type
- `handleToNode` map
- `nameTohandle` map
- `nextID` counter
- `forward` map (moves to Node)
- `backward` map (moves to Node)

## API Changes

### Function Signatures

| Function                        | Before                             | After                           |
| ------------------------------- | ---------------------------------- | ------------------------------- |
| `Vertex()`                      | Returns `Node` (value)             | Returns `*Node` (pointer)       |
| `Edge()`                        | Takes `Node`, `Node`               | Takes `*Node`, `*Node`          |
| `GetEntireNCConePathsAsLinks()` | Takes `NodeHandle`                 | Takes `*Node`                   |
| `GetNodeHandleMatchingName()`   | Returns `[]NodeHandle`             | Returns `*Node` or nil          |
| `GetDBNodeByNodeHandle()`       | Takes `NodeHandle`, returns `Node` | REMOVED (not needed!)           |
| `GetDBArrowByHandle()`          | Takes `ArrowHandle`                | REMOVED (link.arrow is direct!) |
| `AdjointLinkPath()`             | No change                          | No change                       |
| `PrintLinkPath()`               | Takes `[][]Link`, `int`            | Takes `[][]*Link`, `int`        |

### Test Updates

All tests need to update:

- `node.NHandle` → just use `node` pointer
- `graph.handleToNode[...]` → removed
- `GetDBNodeByNodeHandle(graph, handle)` → just use the pointer
- Comparisons: `if handle1 == handle2` → `if node1 == node2`

## Performance Expectations

### Memory Impact

- **Removed:** 2 maps (nameTohandle, handleToNode) + 2 maps (forward, backward)
- **Added:** Slices in each Node (forward, backward)
- **Net:** Likely slight reduction + better locality

### Speed Impact

- **Saved:** 1-2 map lookups per node access
- **Expected:** 10-30% faster for graph operations
- **Benchmarks will tell us!**

## Migration Steps

1. ✅ Commit current state (Handle-based)
2. ✅ Run baseline benchmarks
3. Create new branch: `pointer-based-graph`
4. Refactor `graph.go`:
   - Update type definitions
   - Update LinkedSST struct
   - Update NewLinkedSST()
   - Update Vertex()
   - Update Edge()
   - Update traversal functions
5. Update `maze.go` - change all function calls
6. Update `maze_json.go` - change all function calls
7. Update all test files
8. Run tests - fix any issues
9. Run benchmarks - compare with baseline
10. Document results

## Rollback Plan

If pointer-based version has issues:

```bash
git checkout text2n4l  # Back to Handle-based version
```

The Handle-based version is preserved and working!
