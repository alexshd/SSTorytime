package maze

import (
	"io"
	"testing"
)

// TestOpen verifies that Open() initializes a LinkedSST context correctly
func TestOpen(t *testing.T) {
	graph := NewLinkedSST()
	if graph == nil {
		t.Fatal("Open() returned nil")
	}

	if graph.nodes == nil {
		t.Error("nodes map not initialized")
	}

	if graph.arrow == nil {
		t.Error("arrow map not initialized")
	}

	if graph.inverse == nil {
		t.Error("inverse map not initialized")
	}

	// Verify forward and backward arrows are configured
	fwdArrow := graph.arrow["fwd"]
	bwdArrow := graph.arrow["bwd"]
	if graph.inverse[fwdArrow] != bwdArrow {
		t.Error("Expected fwd->bwd inverse mapping")
	}
	if graph.inverse[bwdArrow] != fwdArrow {
		t.Error("Expected bwd->fwd inverse mapping")
	}
}

// TestVertex verifies node creation and retrieval
func TestVertex(t *testing.T) {
	graph := NewLinkedSST()

	// Create a new vertex
	node1 := Vertex(graph, "A", "chapter1")
	if node1.label != "A" {
		t.Errorf("Expected name='A', got '%s'", node1.label)
	}
	if node1.chap != "chapter1" {
		t.Errorf("Expected chapter='chapter1', got '%s'", node1.chap)
	}
	if node1 == nil {
		t.Error("Expected non-nil node pointer")
	}

	// Retrieve same vertex (should be idempotent)
	node2 := Vertex(graph, "A", "chapter1")
	if node1 != node2 {
		t.Error("Vertex not idempotent: got different node pointers")
	}

	// Create different vertex
	node3 := Vertex(graph, "B", "chapter1")
	if node1 == node3 {
		t.Error("Different vertices should have different pointers")
	}
}

// TestEdge verifies link creation between nodes
func TestEdge(t *testing.T) {
	graph := NewLinkedSST()

	nodeA := Vertex(graph, "A", "ch1")
	nodeB := Vertex(graph, "B", "ch1")

	// Create edge A -[fwd]-> B
	if _, _, err := Edge(graph, nodeA, "fwd", nodeB, []string{"test context"}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}

	// Verify forward link exists in node's forward list
	if len(nodeA.forward) == 0 {
		t.Fatal("Forward link not created in node")
	}

	foundLink := false
	for _, link := range nodeA.forward {
		if link.dst == nodeB {
			foundLink = true
			if link.weight != 1.0 {
				t.Errorf("Expected weight=1.0, got %f", link.weight)
			}
			// Verify arrow is fwd
			if link.arrow.short != "fwd" {
				t.Errorf("Expected arrow 'fwd', got '%s'", link.arrow.short)
			}
			break
		}
	}
	if !foundLink {
		t.Error("Expected forward link A->B not found")
	}

	// Verify backward link exists in node B's backward list (bwd arrow)
	if len(nodeB.backward) == 0 {
		t.Fatal("Backward link not created in node")
	}

	foundBackLink := false
	for _, link := range nodeB.backward {
		if link.dst == nodeA {
			foundBackLink = true
			if link.arrow.short != "bwd" {
				t.Errorf("Expected inverse arrow 'bwd', got '%s'", link.arrow.short)
			}
			break
		}
	}
	if !foundBackLink {
		t.Error("Expected backward link with bwd arrow not found")
	}
}

// TestGetEntireNCConePathsAsLinks verifies BFS path enumeration
func TestGetEntireNCConePathsAsLinks(t *testing.T) {
	graph := NewLinkedSST()

	// Create simple linear path: A -> B -> C
	nodeA := Vertex(graph, "A", "ch1")
	nodeB := Vertex(graph, "B", "ch1")
	nodeC := Vertex(graph, "C", "ch1")

	if _, _, err := Edge(graph, nodeA, "fwd", nodeB, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(graph, nodeB, "fwd", nodeC, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}

	// Find paths from A with depth 1
	paths1, count1 := GetEntireNCConePathsAsLinks(graph, "fwd", nodeA, 1, "ch1", []string{}, 100)
	if count1 != 1 {
		t.Errorf("Expected 1 path at depth 1, got %d", count1)
	}
	if len(paths1) > 0 && len(paths1[0]) != 1 {
		t.Errorf("Expected path length 1, got %d", len(paths1[0]))
	}

	// Find paths from A with depth 2
	paths2, count2 := GetEntireNCConePathsAsLinks(graph, "fwd", nodeA, 2, "ch1", []string{}, 100)
	if count2 != 1 {
		t.Errorf("Expected 1 path at depth 2, got %d", count2)
	}
	if len(paths2) > 0 && len(paths2[0]) != 2 {
		t.Errorf("Expected path length 2, got %d", len(paths2[0]))
	}
}

// TestAdjointLinkPath verifies path reversal with arrow inversion
func TestAdjointLinkPath(t *testing.T) {
	graph := NewLinkedSST()

	// Create nodes A, B, C
	Vertex(graph, "A", "ch1") // nodeA - not directly used but establishes context
	nodeB := Vertex(graph, "B", "ch1")
	nodeC := Vertex(graph, "C", "ch1")

	// Get arrow handle for fwd
	fwdArrow := graph.arrow["fwd"]

	// Create forward path A -> B -> C
	// In the Link representation, we only store the destination
	link1 := &Link{dst: nodeB, arrow: fwdArrow, weight: 1.0} // represents step to B
	link2 := &Link{dst: nodeC, arrow: fwdArrow, weight: 1.0} // represents step to C
	path := []*Link{link1, link2}

	// Reverse the path
	reversed := AdjointLinkPath(graph, path)

	if len(reversed) != 2 {
		t.Fatalf("Expected reversed path length 2, got %d", len(reversed))
	}

	// The reversed path should have links in reverse order
	// Original: [B, C]  -> Reversed: [C, B]
	if reversed[0].dst != nodeC {
		t.Errorf("First reversed link should point to C, got %v", reversed[0].dst)
	}
	if reversed[1].dst != nodeB {
		t.Errorf("Second reversed link should point to B, got %v", reversed[1].dst)
	}

	// Check arrows are inverted to bwd
	bwdArrow := graph.arrow["bwd"]
	if reversed[0].arrow != bwdArrow || reversed[1].arrow != bwdArrow {
		t.Error("Arrows not properly inverted to bwd")
	}
}

// TestWaveFrontsOverlap verifies collision detection at wavefront tips
func TestWaveFrontsOverlap(t *testing.T) {
	graph := NewLinkedSST()

	// Create diamond pattern: Start -> A, B -> End
	start := Vertex(graph, "Start", "ch1")
	nodeA := Vertex(graph, "A", "ch1")
	nodeB := Vertex(graph, "B", "ch1")
	end := Vertex(graph, "End", "ch1")

	if _, _, err := Edge(graph, start, "fwd", nodeA, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(graph, start, "fwd", nodeB, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(graph, nodeA, "fwd", end, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(graph, nodeB, "fwd", end, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}

	// Get forward paths from start
	leftPaths, leftCount := GetEntireNCConePathsAsLinks(graph, "fwd", start, 1, "ch1", []string{}, 100)

	// Get backward paths from end
	rightPaths, rightCount := GetEntireNCConePathsAsLinks(graph, "bwd", end, 1, "ch1", []string{}, 100)

	// Check overlap - should find connections at A and B
	// Use io.Discard to suppress output during test
	dags, cycles := waveFrontsOverlap(graph, io.Discard, leftPaths, rightPaths, leftCount, rightCount, 1, 1)

	if len(dags) == 0 {
		t.Error("Expected to find DAG paths through diamond")
	}

	if len(cycles) != 0 {
		t.Errorf("Expected no cycles in diamond pattern, got %d", len(cycles))
	}
}

// TestSolveMaze is an integration test for the full maze solving algorithm
func TestSolveMaze(t *testing.T) {
	// Run the maze solver - should not return an error
	if err := SolveMaze(); err != nil {
		t.Errorf("SolveMaze returned error: %v", err)
	}
}

// TestGraphIsolation verifies each context is independent
func TestGraphIsolation(t *testing.T) {
	graph1 := NewLinkedSST()
	graph2 := NewLinkedSST()

	node1 := Vertex(graph1, "A", "ch1")
	node2 := Vertex(graph2, "A", "ch1")

	// Same name but different contexts should have independent nodes
	if node1 == node2 {
		t.Error("Different graph contexts should produce different node pointers")
	}
}

// TestSolveMazeJSON verifies the JSON output variant
func TestSolveMazeJSON(t *testing.T) {
	result, err := SolveMazeJSON()
	if err != nil {
		t.Fatalf("SolveMazeJSON returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Verify result structure
	if result.StartNode != StartNode {
		t.Errorf("Expected start=%s, got %s", StartNode, result.StartNode)
	}
	if result.EndNode != EndNode {
		t.Errorf("Expected end=%s, got %s", EndNode, result.EndNode)
	}
}

// TestMultipleEdges verifies handling of multiple edges between same nodes
func TestMultipleEdges(t *testing.T) {
	graph := NewLinkedSST()

	nodeA := Vertex(graph, "A", "ch1")
	nodeB := Vertex(graph, "B", "ch1")

	// Add multiple edges with same arrow type
	if _, _, err := Edge(graph, nodeA, "fwd", nodeB, []string{"context1"}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(graph, nodeA, "fwd", nodeB, []string{"context2"}, 2.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}

	// Check that nodeA has multiple outgoing links
	if len(nodeA.forward) != 2 {
		t.Errorf("Expected 2 outgoing links, got %d", len(nodeA.forward))
	}
}

// TestEmptyGraph verifies behavior with no edges
func TestEmptyGraph(t *testing.T) {
	graph := NewLinkedSST()

	nodeA := Vertex(graph, "A", "ch1")

	// Try to find paths in empty graph
	paths, count := GetEntireNCConePathsAsLinks(graph, "fwd", nodeA, 1, "ch1", []string{}, 100)

	if count != 0 {
		t.Errorf("Expected no paths in empty graph, got %d", count)
	}
	if len(paths) != 0 {
		t.Errorf("Expected empty paths slice, got %d paths", len(paths))
	}
}

// TestPathLimit verifies depth and count limits work
func TestPathLimit(t *testing.T) {
	graph := NewLinkedSST()

	// Create a node with multiple outgoing edges
	start := Vertex(graph, "Start", "ch1")
	for i := 0; i < 10; i++ {
		end := Vertex(graph, string(rune('A'+i)), "ch1")
		if _, _, err := Edge(graph, start, "fwd", end, []string{}, 1.0); err != nil {
			t.Fatalf("Edge creation failed: %v", err)
		}
	}

	// Request limited number of paths
	paths, count := GetEntireNCConePathsAsLinks(graph, "fwd", start, 1, "ch1", []string{}, 5)

	if count > 5 {
		t.Errorf("Expected at most 5 paths due to limit, got %d", count)
	}
	if len(paths) > 5 {
		t.Errorf("Expected at most 5 paths in result slice, got %d", len(paths))
	}
}

// TestEdgeUnknownArrow tests that Edge returns an error for unknown arrow names
func TestEdgeUnknownArrow(t *testing.T) {
	graph := NewLinkedSST()
	defer Close(graph)

	start := Vertex(graph, "node1", "ch1")
	end := Vertex(graph, "node2", "ch1")

	// Try to create an edge with an unknown arrow name
	_, _, err := Edge(graph, start, "unknown_arrow", end, []string{}, 1.0)
	if err == nil {
		t.Error("Expected error for unknown arrow, got nil")
	}

	expectedSubstr := "unknown arrow"
	if err != nil && !contains(err.Error(), expectedSubstr) {
		t.Errorf("Error message should contain %q, got: %v", expectedSubstr, err)
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
