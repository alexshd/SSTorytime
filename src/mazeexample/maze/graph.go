// Package maze provides a minimal in-memory implementation of the SST (Semantic Space-Time)
// graph structure needed for the maze-solving example.
//
// This is a pointer-based implementation using native Go pointers instead of handle indirection.
// Nodes, Links, and Arrows are connected via direct pointer references for simplicity and performance.
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
//	graph := maze.NewLinkedSST()
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
//	paths, count := maze.GetEntireNCConePathsAsLinks(graph, "fwd", start, 3, "", []string{}, 100)
package maze

import (
	"fmt"
)

// Core types represent the fundamental building blocks of the SST graph.
type (
	// LinkedSST (Linked Semantic Space-Time) is the main graph structure.
	// It uses direct pointer references throughout - no handle indirection.
	LinkedSST struct {
		nodes   map[string]*Node  // Fast name-based node lookup
		arrows  []*Arrow          // All defined arrow types
		arrow   map[string]*Arrow // Arrow lookup by name (long or short)
		inverse map[*Arrow]*Arrow // Bidirectional arrow mappings (e.g., fwd<->bwd)
		context map[string]string // Context strings (simplified - just labels for now)
	}

	// Node represents a vertex in the semantic graph.
	// Adjacency lists are stored directly in the node for better locality.
	Node struct {
		label    string  // Semantic label (unique identifier)
		len      int     // Length of label (cached)
		chap     string  // Chapter/category metadata
		seq      bool    // Sequential flag (unused in maze example)
		forward  []*Link // Outgoing edges (forward traversal)
		backward []*Link // Incoming edges (backward traversal)
	}

	// Link represents a directed edge in the graph.
	// All references are direct pointers - no lookups needed.
	Link struct {
		arrow  *Arrow  // Arrow type defining relationship semantics
		dst    *Node   // Destination node (direct pointer!)
		weight float32 // Edge weight
		ctx    string  // Optional context label
	}

	// Arrow defines a relationship type with its semantic properties.
	Arrow struct {
		long    string // Full descriptive name (e.g., "forward")
		short   string // Abbreviated name (e.g., "fwd")
		stIndex int    // Directionality: -1=backward, 0=neutral, 1=forward
	}
)

// NewLinkedSST initializes a new in-memory semantic graph.
// It seeds a minimal arrow vocabulary with two inverse relations: "fwd" and "bwd".
//
// Returns a new *LinkedSST graph ready for operations.
//
// Example:
//
//	graph := NewLinkedSST()
//	defer Close(graph)
func NewLinkedSST() *LinkedSST {
	graph := &LinkedSST{
		nodes:   make(map[string]*Node),
		arrows:  make([]*Arrow, 0, 4),
		arrow:   make(map[string]*Arrow),
		inverse: make(map[*Arrow]*Arrow),
		context: make(map[string]string),
	}

	// Define minimal arrows: fwd and bwd as inverses of each other
	fwd := &Arrow{long: "fwd", short: "fwd", stIndex: 1}
	bwd := &Arrow{long: "bwd", short: "bwd", stIndex: -1}

	graph.arrows = append(graph.arrows, fwd, bwd)
	graph.arrow["fwd"] = fwd
	graph.arrow["bwd"] = bwd
	graph.inverse[fwd] = bwd
	graph.inverse[bwd] = fwd

	return graph
}

// Close releases resources associated with the in-memory graph.
// This is a no-op and exists for API compatibility.
func Close(linkedSST *LinkedSST) { /* no-op for in-memory */ }

// Vertex returns an existing node by name or creates a new one if not found.
//
// Parameters:
//   - linkedSST: The LinkedSST graph
//   - name: Unique identifier for the node
//   - chap: Chapter/category metadata
//
// Returns a *Node pointer (existing or newly created).
//
// Example:
//
//	start := Vertex(graph, "maze_a7", "chapter1")
//	// Later: same := Vertex(graph, "maze_a7", "chapter2")
//	// start == same (pointer equality)
func Vertex(linkedSST *LinkedSST, name, chap string) *Node {
	if node, ok := linkedSST.nodes[name]; ok {
		return node
	}

	node := &Node{
		label:    name,
		len:      len(name),
		chap:     chap,
		forward:  make([]*Link, 0),
		backward: make([]*Link, 0),
	}
	linkedSST.nodes[name] = node
	return node
}

// Edge creates a directed link from 'from' to 'to' with the given arrow type.
//
// Parameters:
//   - graph: The LinkedSST graph
//   - from: Source node
//   - arrowName: Arrow type name (e.g., "fwd", "bwd")
//   - to: Destination node
//   - context: Optional context labels for this edge
//   - weight: Edge weight
//
// Returns:
//   - *Arrow: The arrow used for this edge
//   - int: Placeholder (always 0)
//   - error: Error if arrow name not found
//
// This function creates bidirectional links automatically:
//   - Forward link in from.forward
//   - Backward link in to.backward (with inverse arrow)
//
// Example:
//
//	start := Vertex(graph, "a1", "ch1")
//	end := Vertex(graph, "a2", "ch1")
//	arrow, _, err := Edge(graph, start, "fwd", end, []string{}, 1.0)
func Edge(graph *LinkedSST, from *Node, arrowName string, to *Node, context []string, weight float32) (*Arrow, int, error) {
	arrow, ok := graph.arrow[arrowName]
	if !ok {
		return nil, 0, fmt.Errorf("unknown arrow: %s", arrowName)
	}

	ctxLabel := ""
	if len(context) > 0 && context[0] != "" {
		ctxLabel = context[0]
		graph.context[ctxLabel] = ctxLabel // Just store it
	}

	// Create forward link
	link := &Link{
		arrow:  arrow,
		dst:    to,
		weight: weight,
		ctx:    ctxLabel,
	}
	from.forward = append(from.forward, link)

	// Create backward link with inverse arrow
	invArrow := graph.inverse[arrow]
	backLink := &Link{
		arrow:  invArrow,
		dst:    from,
		weight: weight,
		ctx:    ctxLabel,
	}
	to.backward = append(to.backward, backLink)

	return arrow, 0, nil
}

// GetNodeByName looks up a node by exact name match.
//
// Parameters:
//   - graph: The LinkedSST graph
//   - name: Node name to search for
//   - chap: Chapter filter (ignored in this implementation)
//
// Returns:
//   - *Node: The matching node, or nil if not found
//
// Example:
//
//	node := GetNodeByName(graph, "a7", "")
//	if node != nil {
//	    fmt.Println("Found:", node.label)
//	}
func GetNodeByName(graph *LinkedSST, name, chap string) *Node {
	return graph.nodes[name]
}

// AdjointLinkPath returns the reverse traversal of a path with inverted arrows.
//
// Parameters:
//   - graph: The LinkedSST graph
//   - path: Original path as a slice of *Link
//
// Returns:
//   - []*Link: Reversed path with inverted arrow directions
//
// Example:
//
//	// Original: A --fwd--> B --fwd--> C
//	// Adjoint: C --bwd--> B --bwd--> A
//	reversed := AdjointLinkPath(graph, originalPath)
func AdjointLinkPath(graph *LinkedSST, path []*Link) []*Link {
	var adjoint []*Link
	var prevArrow *Arrow
	if len(path) > 0 {
		prevArrow = path[len(path)-1].arrow
	}

	for j := len(path) - 1; j >= 0; j-- {
		l := path[j]
		// Create new link with inverted arrow
		invLink := &Link{
			arrow:  graph.inverse[prevArrow],
			dst:    l.dst,
			weight: l.weight,
			ctx:    l.ctx,
		}
		adjoint = append(adjoint, invLink)
		prevArrow = path[j].arrow
	}
	return adjoint
}

// GetEntireNCConePathsAsLinks performs bounded breadth-first path enumeration.
//
// Parameters:
//   - graph: The LinkedSST graph
//   - orientation: "fwd" for outgoing edges, "bwd" for incoming edges
//   - start: Starting node
//   - depth: Exact path length to enumerate (number of hops)
//   - chapter: Chapter filter (ignored)
//   - context: Context filter (ignored)
//   - limit: Maximum number of paths to return
//
// Returns:
//   - [][]*Link: Slice of paths
//   - int: Number of paths found
//
// Example:
//
//	paths, count := GetEntireNCConePathsAsLinks(
//	    graph, "fwd", startNode, 3, "", []string{}, 100)
func GetEntireNCConePathsAsLinks(graph *LinkedSST, orientation string, start *Node, depth int, chapter string, context []string, limit int) ([][]*Link, int) {
	type path struct {
		last  *Node
		links []*Link
	}
	frontier := []path{{last: start, links: nil}}
	useForward := orientation == "fwd"

	for d := 0; d < depth; d++ {
		var next []path
		for _, p := range frontier {
			var adj []*Link
			if useForward {
				adj = p.last.forward
			} else {
				adj = p.last.backward
			}

			for _, link := range adj {
				// Extend path
				newPath := make([]*Link, len(p.links)+1)
				copy(newPath, p.links)
				newPath[len(p.links)] = link
				next = append(next, path{last: link.dst, links: newPath})

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

	var results [][]*Link
	for _, p := range frontier {
		if len(p.links) == depth {
			results = append(results, p.links)
		}
	}
	return results, len(results)
}

// PrintLinkPath prints a human-readable representation of a path.
//
// Parameters:
//   - graph: The LinkedSST graph
//   - paths: Slice of paths
//   - pathIndex: Index of the path to print
//   - prefix: String to prepend to output
//   - chapter: Chapter filter (ignored)
//   - context: Context filter (ignored)
//
// Output format: prefix + node1 -(arrow)-> node2 -(arrow)-> node3 ...
//
// Example:
//
//	PrintLinkPath(graph, paths, 0, " - story 0: ", "", nil)
func PrintLinkPath(graph *LinkedSST, paths [][]*Link, pathIndex int, prefix string, chapter string, context []string) {
	if pathIndex < 0 || pathIndex >= len(paths) {
		return
	}
	path := paths[pathIndex]

	rstring := prefix
	if len(path) == 0 {
		fmt.Println(rstring)
		return
	}

	rstring += path[0].dst.label
	for i := 1; i < len(path); i++ {
		rstring += fmt.Sprintf(" -(%s)-> %s", path[i].arrow.long, path[i].dst.label)
	}
	fmt.Println(rstring)
}
