// Package maze provides a minimal in-memory implementation of the SST (Semantic Space-Time)
// graph structure needed for the maze-solving example.
//
// This package is a lightweight extraction from the full pkg/SSTorytime implementation,
// designed to work without database dependencies. It maintains the same API semantics
// but stores all data in memory using Go maps and slices.
//
// # Core Concepts
//
// The SST graph represents a semantic network where:
//   - Nodes represent concepts, entities, or locations (e.g., maze cells)
//   - Links represent directed relationships between nodes
//   - Arrows define the type/semantics of relationships (e.g., "fwd", "bwd")
//   - Contexts provide additional metadata for grouping related links
//
// # Typical Usage
//
//	graph := maze.NewPoSST()
//	defer maze.Close(graph)
//
//	// Create nodes
//	start := maze.Vertex(graph, "maze_a1", "chapter1")
//	end := maze.Vertex(graph, "maze_b2", "chapter1")
//
//	// Create directed edge
//	maze.Edge(graph, start, "fwd", end, []string{"path"}, 1.0)
//
//	// Find paths
//	paths, count := maze.GetEntireNCConePathsAsLinks(graph, "fwd", start.NPtr, 3, "", []string{}, 100)
package maze

import (
	"fmt"
)

// Core types represent the fundamental building blocks of the SST graph.
type (
	// PoSST (Pointing Semantic Space-Time) is the main graph structure that holds
	// the entire in-memory graph. It maintains:
	//   - Node registry with bidirectional name<->pointer mappings
	//   - Adjacency lists for forward (out) and backward (in) traversal
	//   - Arrow vocabulary defining relationship types
	//   - Context registry for metadata grouping
	//
	// All graph operations require a *PoSST pointer.
	PoSST struct {
		nextID   int                // Auto-incrementing ID generator for new nodes
		name2ptr map[string]NodePtr // Fast name-based node lookup
		ptr2node map[NodePtr]Node   // Node metadata storage
		out      map[NodePtr][]Link // Outgoing edges (forward traversal)
		in       map[NodePtr][]Link // Incoming edges (backward traversal)

		// Arrow vocabulary and semantics
		arrows    []ArrowDirectory      // All defined arrow types
		arrowName map[string]ArrowPtr   // Arrow lookup by name (long or short)
		inverse   map[ArrowPtr]ArrowPtr // Bidirectional arrow mappings (e.g., fwd<->bwd)

		// Context registry
		ctxName map[string]ContextPtr // Context string to pointer mapping
	}

	// Node represents a vertex in the semantic graph. Each node has:
	//   - S: The semantic label (name) - must be unique within a PoSST context
	//   - L: Length of the label (cached for performance)
	//   - Chap: Chapter/grouping metadata (does not affect identity)
	//   - NPtr: The stable pointer used to reference this node
	//   - Seq: Sequential flag (reserved for future use)
	//
	// Nodes are immutable once created. Use the NPtr for all graph operations.
	Node struct {
		// L is the cached length of the semantic label S.
		// This optimization avoids repeated len(S) calls during graph operations
		// where node label length is frequently accessed.
		L    int     // Length of S (cached)
		S    string  // Semantic label (unique identifier)
		Seq  bool    // Sequential flag (unused in maze example)
		Chap string  // Chapter/category metadata
		NPtr NodePtr // Stable pointer to this node
	}

	// Link represents a directed edge in the graph. Each link specifies:
	//   - Arr: The arrow type (defines the semantic relationship)
	//   - Dst: The destination node pointer
	//   - Wgt: Edge weight (for weighted graph algorithms)
	//   - Ctx: Optional context pointer for grouping related links
	//
	// Links are stored in adjacency lists and represent both forward and backward
	// traversal capabilities (with inverse arrows for the reverse direction).
	Link struct {
		Arr ArrowPtr   // Arrow type defining relationship semantics
		Wgt float32    // Edge weight (default: 1.0)
		Ctx ContextPtr // Optional context for metadata grouping
		Dst NodePtr    // Destination node
	}

	// NodePtr is a stable handle to a Node in the graph.
	// It remains valid for the lifetime of the PoSST graph.
	NodePtr int

	// ArrowPtr is a handle to an arrow type definition.
	// Used to identify relationship semantics efficiently.
	ArrowPtr int

	// ContextPtr is a handle to a context label.
	// Contexts group related links for filtering or analysis.
	ContextPtr int

	// ArrowDirectory defines an arrow type with its semantic properties:
	//   - STAindex: Space-Time-Arrow index (-1 for backward, 1 for forward)
	//   - Long: Full descriptive name (e.g., "forward", "contains")
	//   - Short: Abbreviated name (e.g., "fwd", "has")
	//   - Ptr: The arrow's stable pointer handle
	//
	// Arrows can be looked up by either long or short name.
	ArrowDirectory struct {
		STAindex int      // Directionality: -1=backward, 0=neutral, 1=forward
		Long     string   // Long-form name
		Short    string   // Short-form name
		Ptr      ArrowPtr // Stable pointer to this arrow
	}
)

// NewPoSST initializes a new in-memory semantic graph (PoSST).
// It seeds a minimal arrow vocabulary with two inverse relations: "fwd" and "bwd".
//
// Returns:
//   - A new *PoSST graph ready for operations
//
// The returned graph is independent of any external services (no database).
// All data is stored in memory and will be lost when the graph is closed or
// the program exits.
//
// Example:
//
//	graph := NewPoSST()
//	defer Close(graph)
//
//	// Ready to create nodes and edges
//	node := Vertex(graph, "maze_a1", "chapter1")
func NewPoSST() *PoSST {
	poSST := PoSST{
		nextID:    1,
		name2ptr:  make(map[string]NodePtr),
		ptr2node:  make(map[NodePtr]Node),
		out:       make(map[NodePtr][]Link),
		in:        make(map[NodePtr][]Link),
		arrows:    make([]ArrowDirectory, 0, 4),
		arrowName: make(map[string]ArrowPtr),
		inverse:   make(map[ArrowPtr]ArrowPtr),
		ctxName:   make(map[string]ContextPtr),
	}

	// Define minimal arrows: fwd and bwd as inverses of each other
	addArrow := func(long, short string, stIndex int) ArrowPtr {
		ptr := ArrowPtr(len(poSST.arrows))
		poSST.arrows = append(poSST.arrows, ArrowDirectory{STAindex: stIndex, Long: long, Short: short, Ptr: ptr})
		poSST.arrowName[long] = ptr
		poSST.arrowName[short] = ptr
		return ptr
	}

	fwd := addArrow("fwd", "fwd", 1)
	bwd := addArrow("bwd", "bwd", -1)
	poSST.inverse[fwd] = bwd
	poSST.inverse[bwd] = fwd

	return &poSST
}

// Close releases resources associated with the in-memory graph.
// This is a no-op for the in-memory implementation and exists for API parity
// with the full database-backed SST implementation.
//
// Parameters:
//   - poSST: The PoSST graph to close
//
// Note: While this function does nothing in the current implementation,
// it's good practice to call it via defer to maintain compatibility with
// future versions that might need cleanup:
//
//	graph := NewPoSST()
//	defer Close(graph)
func Close(poSST *PoSST) { /* no-op for in-memory */ }

// Vertex returns an existing node by name or creates a new one if not found.
//
// Parameters:
//   - poSST: The PoSST graph
//   - name: Unique identifier for the node (e.g., "maze_a1", "room_entrance")
//   - chap: Chapter/category metadata (does not affect node identity)
//
// Returns:
//   - The Node (existing or newly created)
//
// Node name uniqueness is enforced within the PoSST graph. If a node with
// the given name already exists, it is returned unchanged (the chap parameter
// is ignored for existing nodes).
//
// Example:
//
//	// First call creates the node
//	start := Vertex(graph, "maze_a7", "chapter1")
//
//	// Second call returns the existing node
//	same := Vertex(graph, "maze_a7", "chapter2") // chap is ignored
//	// start.NPtr == same.NPtr (true)
func Vertex(poSST *PoSST, name, chap string) Node {
	if np, ok := poSST.name2ptr[name]; ok {
		n := poSST.ptr2node[np]
		return n
	}
	np := NodePtr(poSST.nextID)
	poSST.nextID++
	n := Node{L: len(name), S: name, Chap: chap, NPtr: np}
	poSST.name2ptr[name] = np
	poSST.ptr2node[np] = n
	return n
}

// TryContext registers a context label and returns a stable ContextPtr.
//
// Parameters:
//   - poSST: The PoSST graph
//   - context: String slice containing the context label (only first element is used)
//
// Returns:
//   - ContextPtr handle (0 if no valid context provided)
//
// Contexts provide a way to group or tag related links with metadata.
// If the context string was previously registered, returns the existing pointer.
// If the context slice is empty or the first element is empty, returns 0.
//
// Example:
//
//	ctxPtr := TryContext(graph, []string{"maze_level_1"})
//	// Later links can reference this context
//	Edge(graph, from, "fwd", to, []string{"maze_level_1"}, 1.0)
func TryContext(poSST *PoSST, context []string) ContextPtr {
	if len(context) == 0 || context[0] == "" {
		return 0
	}
	if cp, ok := poSST.ctxName[context[0]]; ok {
		return cp
	}
	cp := ContextPtr(len(poSST.ctxName) + 1)
	poSST.ctxName[context[0]] = cp
	return cp
}

// Edge creates a directed link from 'from' to 'to' with the given arrow type.
//
// Parameters:
//   - poSST: The PoSST graph
//   - from: Source node
//   - arrow: Arrow type name (e.g., "fwd", "bwd") - must exist in arrow vocabulary
//   - to: Destination node
//   - context: Optional context labels for this edge
//   - weight: Edge weight (typically 1.0 for unweighted graphs)
//
// Returns:
//   - ArrowPtr: The arrow pointer used for this edge
//   - int: ST-type placeholder (always 0 in this implementation)
//   - error: Error if arrow name not found
//
// This function:
//  1. Creates a forward link in the outgoing adjacency list (poSST.out)
//  2. Creates a reverse link in the incoming adjacency list (poSST.in) using the inverse arrow
//
// The bidirectional storage enables efficient forward and backward graph traversal.
//
// Returns an error if the arrow name is not found in the vocabulary.
//
// Example:
//
//	start := Vertex(graph, "maze_a1", "ch1")
//	end := Vertex(graph, "maze_a2", "ch1")
//
//	// Create edge: maze_a1 --fwd--> maze_a2
//	ap, _, err := Edge(graph, start, "fwd", end, []string{}, 1.0)
//	if err != nil {
//	    return err
//	}
//
//	// Internally also creates: maze_a2 --bwd--> maze_a1
func Edge(poSST *PoSST, from Node, arrow string, to Node, context []string, weight float32) (ArrowPtr, int, error) {
	ap, ok := poSST.arrowName[arrow]
	if !ok {
		return 0, 0, fmt.Errorf("unknown arrow: %s", arrow)
	}
	link := Link{Arr: ap, Wgt: weight, Ctx: TryContext(poSST, context), Dst: to.NPtr}
	poSST.out[from.NPtr] = append(poSST.out[from.NPtr], link)
	// also store reverse with inverse arrow for convenience
	inv := poSST.inverse[ap]
	rlink := Link{Arr: inv, Wgt: weight, Ctx: link.Ctx, Dst: from.NPtr}
	poSST.in[to.NPtr] = append(poSST.in[to.NPtr], rlink)
	return ap, 0, nil
}

// GetDBNodePtrMatchingName looks up nodes by exact name match.
//
// Parameters:
//   - poSST: The PoSST graph
//   - name: Node name to search for
//   - chap: Chapter filter (ignored in this implementation)
//
// Returns:
//   - []NodePtr: Slice containing the matching node pointer, or nil if not found
//
// This in-memory version returns at most one match since node names are unique
// within a graph. The function signature matches the full database-backed
// implementation for API compatibility.
//
// Example:
//
//	ptrs := GetDBNodePtrMatchingName(graph, "maze_a7", "")
//	if ptrs != nil {
//	    node := GetDBNodeByNodePtr(graph, ptrs[0])
//	    fmt.Println("Found:", node.S)
//	}
func GetDBNodePtrMatchingName(poSST *PoSST, name, chap string) []NodePtr {
	if np, ok := poSST.name2ptr[name]; ok {
		return []NodePtr{np}
	}
	return nil
}

// GetDBNodeByNodePtr resolves a NodePtr to its Node metadata.
//
// Parameters:
//   - poSST: The PoSST graph
//   - dbNptr: Node pointer to resolve
//
// Returns:
//   - Node: The node metadata (zero value if pointer not found)
//
// This is the primary way to access node information after obtaining a NodePtr
// from graph traversal or lookup operations.
//
// Example:
//
//	ptrs := GetDBNodePtrMatchingName(graph, "maze_a7", "")
//	if ptrs != nil {
//	    node := GetDBNodeByNodePtr(graph, ptrs[0])
//	    fmt.Printf("Node: %s (length=%d)\n", node.S, node.L)
//	}
func GetDBNodeByNodePtr(poSST *PoSST, dbNptr NodePtr) Node {
	return poSST.ptr2node[dbNptr]
}

// GetDBArrowByPtr returns arrow metadata for a given ArrowPtr.
//
// Parameters:
//   - poSST: The PoSST graph
//   - arrowptr: Arrow pointer to resolve
//
// Returns:
//   - ArrowDirectory: Arrow metadata including long/short names and index
//
// If the arrow pointer is invalid, returns a placeholder with "unknown" as the name.
//
// Example:
//
//	arr := GetDBArrowByPtr(graph, link.Arr)
//	fmt.Printf("Arrow: %s (short: %s, index: %d)\n", arr.Long, arr.Short, arr.STAindex)
func GetDBArrowByPtr(poSST *PoSST, arrowptr ArrowPtr) ArrowDirectory {
	if int(arrowptr) >= 0 && int(arrowptr) < len(poSST.arrows) {
		return poSST.arrows[arrowptr]
	}
	return ArrowDirectory{Long: "unknown", Short: "?", Ptr: arrowptr}
}

// AdjointLinkPath returns the reverse traversal of a path with inverted arrows.
//
// Parameters:
//   - poSST: The PoSST graph
//   - LL: Original path as a slice of Links
//
// Returns:
//   - []Link: Reversed path with inverted arrow directions
//
// This function reverses a path for backward traversal. Each link's arrow
// is replaced with its inverse (e.g., "fwd" becomes "bwd"), and the sequence
// is reversed so that traversing the adjoint path visits nodes in opposite order.
//
// The adjoint operation is essential for bidirectional search algorithms where
// paths found from the goal need to be reversed and joined with paths from the start.
//
// Example:
//
//	// Original path: A --fwd--> B --fwd--> C
//	originalPath := []Link{
//	    {Arr: fwdPtr, Dst: BPtr},
//	    {Arr: fwdPtr, Dst: CPtr},
//	}
//
//	// Adjoint path: C --bwd--> B --bwd--> A
//	reversed := AdjointLinkPath(graph, originalPath)
//	// Now reversed can be traversed from C back to A
func AdjointLinkPath(poSST *PoSST, LL []Link) []Link {
	var adjoint []Link
	var prev ArrowPtr
	if len(LL) > 0 {
		prev = LL[len(LL)-1].Arr
	}
	for j := len(LL) - 1; j >= 0; j-- {
		l := LL[j]
		// flip using inverse of previous arrow, consistent with original semantics
		l.Arr = poSST.inverse[prev]
		adjoint = append(adjoint, l)
		prev = LL[j].Arr
	}
	return adjoint
}

// GetEntireNCConePathsAsLinks performs bounded breadth-first path enumeration.
//
// Parameters:
//   - poSST: The PoSST graph
//   - orientation: "fwd" for outgoing edges, "bwd" for incoming edges
//   - start: Starting node pointer
//   - depth: Exact path length to enumerate (number of hops)
//   - chapter: Chapter filter (ignored in this implementation)
//   - context: Context filter (ignored in this implementation)
//   - limit: Maximum number of paths to return (0 = unlimited)
//
// Returns:
//   - [][]Link: Slice of paths, where each path is a slice of Links
//   - int: Number of paths found
//
// This function explores the graph from a starting node using breadth-first
// expansion. It returns ALL paths of exactly the specified depth. Each path
// is represented as a []Link where:
//   - First link's Dst is the first hop from start
//   - Last link's Dst is the final node at depth 'depth'
//
// The limit parameter caps the total number of paths at each expansion level,
// useful for preventing combinatorial explosion in dense graphs.
//
// Orientation controls traversal direction:
//   - "fwd": Follow outgoing edges (normal forward search)
//   - "bwd": Follow incoming edges (backward/reverse search)
//
// Example:
//
//	// Find all 3-hop forward paths from maze_a7
//	paths, count := GetEntireNCConePathsAsLinks(
//	    graph,
//	    "fwd",        // forward direction
//	    startPtr,     // starting node
//	    3,            // exactly 3 hops
//	    "",           // no chapter filter
//	    []string{},   // no context filter
//	    100,          // max 100 paths
//	)
//
//	fmt.Printf("Found %d paths of length 3\n", count)
//	for i, path := range paths {
//	    fmt.Printf("Path %d has %d links\n", i, len(path))
//	}
func GetEntireNCConePathsAsLinks(poSST *PoSST, orientation string, start NodePtr, depth int, chapter string, context []string, limit int) ([][]Link, int) {
	// paths are represented as []Link where first element's Dst is the first hop node
	var results [][]Link
	type path struct {
		last  NodePtr
		links []Link
	}
	frontier := []path{{last: start, links: nil}}

	useOut := orientation == "fwd"

	for d := 0; d < depth; d++ {
		var next []path
		for _, p := range frontier {
			var adj []Link
			if useOut {
				adj = poSST.out[p.last]
			} else {
				adj = poSST.in[p.last]
			}
			for _, l := range adj {
				// extend path
				nl := make([]Link, len(p.links)+1)
				copy(nl, p.links)
				nl[len(p.links)] = l
				next = append(next, path{last: l.Dst, links: nl})
				if limit > 0 && len(next) >= limit {
					break
				}
			}
			if limit > 0 && len(next) >= limit {
				break
			}
		}
		frontier = next
		if len(frontier) == 0 {
			break
		}
	}

	for _, p := range frontier {
		if len(p.links) == depth {
			results = append(results, p.links)
		}
	}
	return results, len(results)
}

// PrintLinkPath prints a human-readable representation of a path from a set.
//
// Parameters:
//   - poSST: The PoSST graph
//   - cone: Slice of paths (typically from GetEntireNCConePathsAsLinks)
//   - p: Index of the path to print within cone
//   - prefix: String to prepend to the output line (e.g., " - story 0: ")
//   - chapter: Chapter filter (ignored in this implementation)
//   - context: Context filter (ignored in this implementation)
//
// Output format:
//
//	prefix + node1 -(arrow1)-> node2 -(arrow2)-> node3 ...
//
// The function prints directly to stdout. For each link in the path, it shows:
//   - The destination node name
//   - The arrow type used (long name)
//
// If the path index is out of bounds or the path is empty, the function
// returns early (printing just the prefix for empty paths).
//
// Example output:
//
//   - story 0: maze_a7 -(fwd)-> maze_b7 -(fwd)-> maze_c7 -(fwd)-> maze_d7
//
// Example usage:
//
//	paths, count := GetEntireNCConePathsAsLinks(graph, "fwd", start, 3, "", []string{}, 100)
//	for i := 0; i < count; i++ {
//	    prefix := fmt.Sprintf("Path %d: ", i)
//	    PrintLinkPath(graph, paths, i, prefix, "", nil)
//	}
func PrintLinkPath(poSST *PoSST, cone [][]Link, p int, prefix string, chapter string, context []string) {
	if p < 0 || p >= len(cone) {
		return
	}
	rstring := prefix
	if len(cone[p]) == 0 {
		fmt.Println(rstring)
		return
	}
	start := GetDBNodeByNodePtr(poSST, cone[p][0].Dst)
	rstring += start.S
	for i := 1; i < len(cone[p]); i++ {
		arr := GetDBArrowByPtr(poSST, cone[p][i].Arr)
		node := GetDBNodeByNodePtr(poSST, cone[p][i].Dst)
		rstring += fmt.Sprintf(" -(%s)-> %s", arr.Long, node.S)
	}
	fmt.Println(rstring)
}
