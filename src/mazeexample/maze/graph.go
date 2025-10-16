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
	"log"
)

// Core types represent the fundamental building blocks of the SST graph.
type (
	// LinkedSST (Linked Semantic Space-Time) is the main graph structure that holds
	// the entire in-memory graph. It maintains:
	//   - Node registry with bidirectional name<->handle mappings
	//   - Adjacency lists for forward (out) and backward (in) traversal
	//   - Arrow vocabulary defining relationship types
	//   - Context registry for metadata grouping
	//
	// All graph operations require a *LinkedSST handle.
	LinkedSST struct {
		nextID       int                   // Auto-incrementing ID generator for new nodes
		nameTohandle map[string]NodeHandle // Fast name-based node lookup
		handleToNode map[NodeHandle]Node   // Node metadata storage
		forward      map[NodeHandle][]Link // Outgoing edges (forward traversal)
		backward     map[NodeHandle][]Link // Incoming edges (backward traversal)

		// Arrow vocabulary and semantics
		arrows  []ArrowDirectory            // All defined arrow types
		arrow   map[string]ArrowHandle      // Arrow lookup by name (long or short)
		inverse map[ArrowHandle]ArrowHandle // Bidirectional arrow mappings (e.g., fwd<->bwd)

		// Context registry
		context map[string]ContextHandle // Context string to handle mapping
	}

	// Node represents a vertex in the semantic graph. Each node has:
	//   - label: The semantic label (name) - must be unique within a LinkedSST context
	//   - len: Length of the label (cached for performance)
	//   - Chap: Chapter/grouping metadata (does not affect identity)
	//   - NHandle: The stable handle used to reference this node
	//   - Seq: Sequential flag (reserved for future use)
	//
	// Nodes are immutable once created. Use the NHandle for all graph operations.
	Node struct {
		// len is the cached length of the semantic label.
		// This optimization avoids repeated len(label) calls during graph operations
		// where node label length is frequently accessed.
		len     int        // Length of label (cached)
		label   string     // Semantic label (unique identifier)
		Seq     bool       // Sequential flag (unused in maze example)
		Chap    string     // Chapter/category metadata
		NHandle NodeHandle // Stable handle to this node
	}

	// Link represents a directed edge in the graph. Each link specifies:
	//   - Arr: The arrow type (defines the semantic relationship)
	//   - Dst: The destination node handle
	//   - Wgt: Edge weight (for weighted graph algorithms)
	//   - Ctx: Optional context handle for grouping related links
	//
	// Links are stored in adjacency lists and represent both forward and backward
	// traversal capabilities (with inverse arrows for the reverse direction).
	Link struct {
		Arr ArrowHandle   // Arrow type defining relationship semantics
		Wgt float32       // Edge weight (default: 1.0)
		Ctx ContextHandle // Optional context for metadata grouping
		Dst NodeHandle    // Destination node
	}

	// NodeHandle is a stable handle to a Node in the graph.
	// It remains valid for the lifetime of the LinkedSST graph.
	NodeHandle int

	// ArrowHandle is a handle to an arrow type definition.
	// Used to identify relationship semantics efficiently.
	ArrowHandle int

	// ContextHandle is a handle to a context label.
	// Contexts group related links for filtering or analysis.
	ContextHandle int

	// ArrowDirectory defines an arrow type with its semantic properties:
	//   - STAindex: Space-Time-Arrow index (-1 for backward, 1 for forward)
	//   - Long: Full descriptive name (e.g., "forward", "contains")
	//   - Short: Abbreviated name (e.g., "fwd", "has")
	//   - Handle: The arrow's stable handle
	//
	// Arrows can be looked up by either long or short name.
	ArrowDirectory struct {
		STAindex int         // Directionality: -1=backward, 0=neutral, 1=forward
		Long     string      // Long-form name
		Short    string      // Short-form name
		Handle   ArrowHandle // Stable handle to this arrow
	}
)

// NewLinkedSST initializes a new in-memory semantic graph (LinkedSST).
// It seeds a minimal arrow vocabulary with two inverse relations: "fwd" and "bwd".
//
// Returns:
//   - A new *LinkedSST graph ready for operations
//
// The returned graph is independent of any external services (no database).
// All data is stored in memory and will be lost when the graph is closed or
// the program exits.
//
// Example:
//
//	graph := NewLinkedSST()
//	defer Close(graph)
//
//	// Ready to create nodes and edges
//	node := Vertex(graph, "maze_a1", "chapter1")
func NewLinkedSST() *LinkedSST {
	graph := LinkedSST{
		nextID:       1,
		nameTohandle: make(map[string]NodeHandle),
		handleToNode: make(map[NodeHandle]Node),
		forward:      make(map[NodeHandle][]Link),
		backward:     make(map[NodeHandle][]Link),
		arrows:       make([]ArrowDirectory, 0, 4),
		arrow:        make(map[string]ArrowHandle),
		inverse:      make(map[ArrowHandle]ArrowHandle),
		context:      make(map[string]ContextHandle),
	}

	// Define minimal arrows: fwd and bwd as inverses of each other
	addArrow := func(long, short string, stIndex int) ArrowHandle {
		handle := ArrowHandle(len(graph.arrows))
		graph.arrows = append(graph.arrows, ArrowDirectory{STAindex: stIndex, Long: long, Short: short, Handle: handle})
		graph.arrow[long] = handle
		graph.arrow[short] = handle
		return handle
	}

	fwd := addArrow("fwd", "fwd", 1)
	bwd := addArrow("bwd", "bwd", -1)
	graph.inverse[fwd] = bwd
	graph.inverse[bwd] = fwd

	return &graph
}

// Close releases resources associated with the in-memory graph.
// This is a no-op for the in-memory implementation and exists for API parity
// with the full database-backed SST implementation.
//
// Parameters:
//   - linkedSST: The LinkedSST graph to close
//
// Note: While this function does nothing in the current implementation,
// it's good practice to call it via defer to maintain compatibility with
// future versions that might need cleanup:
//
//	graph := NewLinkedSST()
//	defer Close(graph)
func Close(linkedSST *LinkedSST) { /* no-op for in-memory */ }

// Vertex returns an existing node by name or creates a new one if not found.
//
// Parameters:
//   - linkedSST: The LinkedSST graph
//   - name: Unique identifier for the node (e.g., "maze_a1", "room_entrance")
//   - chap: Chapter/category metadata (does not affect node identity)
//
// Returns:
//   - The Node (existing or newly created)
//
// Node name uniqueness is enforced within the LinkedSST graph. If a node with
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
//	// start.NHandle == same.NHandle (true)
func Vertex(linkedSST *LinkedSST, name, chap string) Node {
	if nh, ok := linkedSST.nameTohandle[name]; ok {
		n := linkedSST.handleToNode[nh]
		return n
	}
	nh := NodeHandle(linkedSST.nextID)
	linkedSST.nextID++
	n := Node{len: len(name), label: name, Chap: chap, NHandle: nh}
	linkedSST.nameTohandle[name] = nh
	linkedSST.handleToNode[nh] = n
	return n
}

// TryContext registers a context label and returns a stable ContextHandle.
//
// Parameters:
//   - graph: The LinkedSST graph
//   - context: String slice containing the context label (only first element is used)
//
// Returns:
//   - ContextHandle (0 if no valid context provided)
//
// Contexts provide a way to group or tag related links with metadata.
// If the context string was previously registered, returns the existing handle.
// If the context slice is empty or the first element is empty, returns 0.
//
// Example:
//
//	ctxHandle := TryContext(graph, []string{"maze_level_1"})
//	// Later links can reference this context
//	Edge(graph, from, "fwd", to, []string{"maze_level_1"}, 1.0)
func TryContext(graph *LinkedSST, context []string) ContextHandle {
	if len(context) == 0 || context[0] == "" {
		return 0
	}
	if ch, ok := graph.context[context[0]]; ok {
		return ch
	}
	ch := ContextHandle(len(graph.context) + 1)
	graph.context[context[0]] = ch
	return ch
}

// Edge creates a directed link from 'from' to 'to' with the given arrow type.
//
// Parameters:
//   - poSST: The LinkedSST graph
//   - from: Source node
//   - arrow: Arrow type name (e.g., "fwd", "bwd") - must exist in arrow vocabulary
//   - to: Destination node
//   - context: Optional context labels for this edge
//   - weight: Edge weight (typically 1.0 for unweighted graphs)
//
// Returns:
//   - ArrowHandle: The arrow handle used for this edge
//   - int: ST-type placeholder (always 0 in this implementation)
//   - error: Error if arrow name not found
//
// This function:
//  1. Creates a forward link in the outgoing adjacency list (poSST.forward)
//  2. Creates a reverse link in the incoming adjacency list (poSST.backward) using the inverse arrow
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
//	ah, _, err := Edge(graph, start, "fwd", end, []string{}, 1.0)
//	if err != nil {
//	    return err
//	}
//
//	// Internally also creates: maze_a2 --bwd--> maze_a1
func Edge(graph *LinkedSST, from Node, arrow string, to Node, context []string, weight float32) (ArrowHandle, int, error) {
	ah, ok := graph.arrow[arrow]
	if !ok {
		return 0, 0, fmt.Errorf("unknown arrow: %s", arrow)
	}
	link := Link{Arr: ah, Wgt: weight, Ctx: TryContext(graph, context), Dst: to.NHandle}
	graph.forward[from.NHandle] = append(graph.forward[from.NHandle], link)
	// also store reverse with inverse arrow for convenience
	inv := graph.inverse[ah]
	rlink := Link{Arr: inv, Wgt: weight, Ctx: link.Ctx, Dst: from.NHandle}
	graph.backward[to.NHandle] = append(graph.backward[to.NHandle], rlink)
	return ah, 0, nil
}

// GetNodeHandleMatchingName looks up nodes by exact name match.
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
// GetNodeHandleMatchingName looks up nodes by exact name match.
//
// Parameters:
//   - poSST: The LinkedSST graph
//   - name: Node name to search for
//   - chap: Chapter filter (ignored in this implementation)
//
// Returns:
//   - []NodeHandle: Slice containing the matching node handle, or nil if not found
//
// This in-memory version returns at most one match since node names are unique
// within a graph. The function signature matches the full database-backed
// implementation for API compatibility.
//
// Example:
//
//	handles := GetNodeHandleMatchingName(graph, "a7", "")
//	if handles != nil {
//	    node := GetDBNodeByNodeHandle(graph, handles[0])
//	    fmt.Println("Found:", node.label)
//	}
func GetNodeHandleMatchingName(poSST *LinkedSST, name, chap string) []NodeHandle {
	nh, ok := poSST.nameTohandle[name]
	if !ok {
		log.Default().Println("GetNodeHandleMatchingName: not found:", name)
		return nil
	}

	return []NodeHandle{nh}
}

// GetDBNodeByNodeHandle resolves a NodeHandle to its Node metadata.
//
// Parameters:
//   - poSST: The LinkedSST graph
//   - dbNhandle: Node handle to resolve
//
// Returns:
//   - Node: The node metadata (zero value if handle not found)
//
// This is the primary way to access node information after obtaining a NodeHandle
// from graph traversal or lookup operations.
//
// Example:
//
//	handles := GetNodeHandleMatchingName(graph, "a7", "")
//	if handles != nil {
//	    node := GetDBNodeByNodeHandle(graph, handles[0])
//	    fmt.Printf("Node: %s (length=%d)\n", node.label, node.len)
//	}
func GetDBNodeByNodeHandle(poSST *LinkedSST, dbNhandle NodeHandle) Node {
	return poSST.handleToNode[dbNhandle]
}

// GetDBArrowByHandle returns arrow metadata for a given ArrowHandle.
//
// Parameters:
//   - poSST: The LinkedSST graph
//   - arrowhandle: Arrow handle to resolve
//
// Returns:
//   - ArrowDirectory: Arrow metadata including long/short names and index
//
// If the arrow handle is invalid, returns a placeholder with "unknown" as the name.
//
// Example:
//
//	arr := GetDBArrowByHandle(graph, link.Arr)
//	fmt.Printf("Arrow: %s (short: %s, index: %d)\n", arr.Long, arr.Short, arr.STAindex)
func GetDBArrowByHandle(poSST *LinkedSST, arrowhandle ArrowHandle) ArrowDirectory {
	if int(arrowhandle) >= 0 && int(arrowhandle) < len(poSST.arrows) {
		return poSST.arrows[arrowhandle]
	}
	return ArrowDirectory{Long: "unknown", Short: "?", Handle: arrowhandle}
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
func AdjointLinkPath(poSST *LinkedSST, LL []Link) []Link {
	var adjoint []Link
	var prev ArrowHandle
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
//   - poSST: The LinkedSST graph
//   - orientation: "fwd" for outgoing edges, "bwd" for incoming edges
//   - start: Starting node handle
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
//	// Find all 3-hop forward paths from a7
//	paths, count := GetEntireNCConePathsAsLinks(
//	    graph,
//	    "fwd",        // forward direction
//	    startHandle,  // starting node
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
func GetEntireNCConePathsAsLinks(poSST *LinkedSST, orientation string, start NodeHandle, depth int, chapter string, context []string, limit int) ([][]Link, int) {
	// paths are represented as []Link where first element's Dst is the first hop node
	var results [][]Link
	type path struct {
		last  NodeHandle
		links []Link
	}
	frontier := []path{{last: start, links: nil}}

	useOut := orientation == "fwd"

	for d := 0; d < depth; d++ {
		var next []path
		for _, p := range frontier {
			var adj []Link
			if useOut {
				adj = poSST.forward[p.last]
			} else {
				adj = poSST.backward[p.last]
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
func PrintLinkPath(poSST *LinkedSST, cone [][]Link, p int, prefix string, chapter string, context []string) {
	if p < 0 || p >= len(cone) {
		return
	}
	rstring := prefix
	if len(cone[p]) == 0 {
		fmt.Println(rstring)
		return
	}
	start := GetDBNodeByNodeHandle(poSST, cone[p][0].Dst)
	rstring += start.label
	for i := 1; i < len(cone[p]); i++ {
		arr := GetDBArrowByHandle(poSST, cone[p][i].Arr)
		node := GetDBNodeByNodeHandle(poSST, cone[p][i].Dst)
		rstring += fmt.Sprintf(" -(%s)-> %s", arr.Long, node.label)
	}
	fmt.Println(rstring)
}
