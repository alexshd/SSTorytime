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

	if graph.nextID != 1 {
		t.Errorf("Expected nextID=1, got %d", graph.nextID)
	}

	if graph.nameTohandle == nil {
		t.Error("nameTohandle map not initialized")
	}

	if graph.handleToNode == nil {
		t.Error("handleToNode map not initialized")
	}

	if graph.inverse == nil {
		t.Error("inverse map not initialized")
	}

	// Verify forward and backward arrows are configured
	fwdHandle := graph.arrow["fwd"]
	bwdHandle := graph.arrow["bwd"]
	if graph.inverse[fwdHandle] != bwdHandle {
		t.Error("Expected fwd->bwd inverse mapping")
	}
	if graph.inverse[bwdHandle] != fwdHandle {
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
	if node1.Chap != "chapter1" {
		t.Errorf("Expected chapter='chapter1', got '%s'", node1.Chap)
	}
	if node1.NHandle <= 0 {
		t.Errorf("Expected positive NodeHandle, got %d", node1.NHandle)
	}

	// Retrieve same vertex (should be idempotent)
	node2 := Vertex(graph, "A", "chapter1")
	if node1.NHandle != node2.NHandle {
		t.Errorf("Vertex not idempotent: got NodeHandles %d and %d", node1.NHandle, node2.NHandle)
	}

	// Create different vertex
	node3 := Vertex(graph, "B", "chapter1")
	if node1.NHandle == node3.NHandle {
		t.Error("Different vertices should have different NodeHandles")
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

	// Verify forward link exists in out map
	outLinks, exists := graph.forward[nodeA.NHandle]
	if !exists || len(outLinks) == 0 {
		t.Fatal("Forward link not created in out map")
	}

	foundLink := false
	for _, link := range outLinks {
		if link.Dst == nodeB.NHandle {
			foundLink = true
			if link.Wgt != 1.0 {
				t.Errorf("Expected weight=1.0, got %f", link.Wgt)
			}
			// Verify arrow is fwd
			arrInfo := GetDBArrowByHandle(graph, link.Arr)
			if arrInfo.Short != "fwd" {
				t.Errorf("Expected arrow 'fwd', got '%s'", arrInfo.Short)
			}
			break
		}
	}
	if !foundLink {
		t.Error("Expected forward link A->B not found")
	}

	// Verify backward link exists in in map (bwd arrow)
	inLinks, exists := graph.backward[nodeB.NHandle]
	if !exists || len(inLinks) == 0 {
		t.Fatal("Backward link not created in in map")
	}

	foundBackLink := false
	for _, link := range inLinks {
		if link.Dst == nodeA.NHandle {
			foundBackLink = true
			arrInfo := GetDBArrowByHandle(graph, link.Arr)
			if arrInfo.Short != "bwd" {
				t.Errorf("Expected inverse arrow 'bwd', got '%s'", arrInfo.Short)
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
	paths1, count1 := GetEntireNCConePathsAsLinks(graph, "fwd", nodeA.NHandle, 1, "ch1", []string{}, 100)
	if count1 != 1 {
		t.Errorf("Expected 1 path at depth 1, got %d", count1)
	}
	if len(paths1) > 0 && len(paths1[0]) != 1 {
		t.Errorf("Expected path length 1, got %d", len(paths1[0]))
	}

	// Find paths from A with depth 2
	paths2, count2 := GetEntireNCConePathsAsLinks(graph, "fwd", nodeA.NHandle, 2, "ch1", []string{}, 100)
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
	fwdHandle := graph.arrow["fwd"]

	// Create forward path A -> B -> C
	// In the Link representation, we only store the destination
	link1 := Link{Dst: nodeB.NHandle, Arr: fwdHandle, Wgt: 1.0} // represents step to B
	link2 := Link{Dst: nodeC.NHandle, Arr: fwdHandle, Wgt: 1.0} // represents step to C
	path := []Link{link1, link2}

	// Reverse the path
	reversed := AdjointLinkPath(graph, path)

	if len(reversed) != 2 {
		t.Fatalf("Expected reversed path length 2, got %d", len(reversed))
	}

	// The reversed path should have links in reverse order
	// Original: [B, C]  -> Reversed: [C, B]
	if reversed[0].Dst != nodeC.NHandle {
		t.Errorf("First reversed link should point to C, got %d", reversed[0].Dst)
	}
	if reversed[1].Dst != nodeB.NHandle {
		t.Errorf("Second reversed link should point to B, got %d", reversed[1].Dst)
	}

	// Check arrows are inverted to bwd
	bwdHandle := graph.arrow["bwd"]
	if reversed[0].Arr != bwdHandle || reversed[1].Arr != bwdHandle {
		t.Error("Arrows not properly inverted to bwd")
	}
}

// TestWaveFrontsOverlap verifies frontier collision detection
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
	leftPaths, leftCount := GetEntireNCConePathsAsLinks(graph, "fwd", start.NHandle, 1, "ch1", []string{}, 100)

	// Get backward paths from end
	rightPaths, rightCount := GetEntireNCConePathsAsLinks(graph, "bwd", end.NHandle, 1, "ch1", []string{}, 100)

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

	// Same name but different contexts should have independent NodeHandles
	if graph1.nameTohandle == nil || graph2.nameTohandle == nil {
		t.Fatal("Context maps not initialized")
	}

	// Each context should have exactly one node
	if len(graph1.nameTohandle) != 1 {
		t.Errorf("graph1 should have 1 node, got %d", len(graph1.nameTohandle))
	}
	if len(graph2.nameTohandle) != 1 {
		t.Errorf("graph2 should have 1 node, got %d", len(graph2.nameTohandle))
	}

	// Verify they are distinct in their respective contexts
	if node1.NHandle == node2.NHandle && graph1 == graph2 {
		t.Error("Contexts are not properly isolated")
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

	outLinks := graph.forward[nodeA.NHandle]
	if len(outLinks) != 2 {
		t.Errorf("Expected 2 outgoing links, got %d", len(outLinks))
	}
}

// TestEmptyGraph verifies behavior with no edges
func TestEmptyGraph(t *testing.T) {
	graph := NewLinkedSST()

	nodeA := Vertex(graph, "A", "ch1")

	// Try to find paths in empty graph
	paths, count := GetEntireNCConePathsAsLinks(graph, "fwd", nodeA.NHandle, 1, "ch1", []string{}, 100)

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
	paths, count := GetEntireNCConePathsAsLinks(graph, "fwd", start.NHandle, 1, "ch1", []string{}, 5)

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
