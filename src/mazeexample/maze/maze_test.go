package maze

import (
	"io"
	"testing"
)

// TestOpen verifies that Open() initializes a PoSST context correctly
func TestOpen(t *testing.T) {
	ctx := NewPoSST()
	if ctx == nil {
		t.Fatal("Open() returned nil")
	}

	if ctx.nextID != 1 {
		t.Errorf("Expected nextID=1, got %d", ctx.nextID)
	}

	if ctx.name2ptr == nil {
		t.Error("name2ptr map not initialized")
	}

	if ctx.ptr2node == nil {
		t.Error("ptr2node map not initialized")
	}

	if ctx.inverse == nil {
		t.Error("inverse map not initialized")
	}

	// Verify forward and backward arrows are configured
	fwdPtr := ctx.arrowName["fwd"]
	bwdPtr := ctx.arrowName["bwd"]
	if ctx.inverse[fwdPtr] != bwdPtr {
		t.Error("Expected fwd->bwd inverse mapping")
	}
	if ctx.inverse[bwdPtr] != fwdPtr {
		t.Error("Expected bwd->fwd inverse mapping")
	}
}

// TestVertex verifies node creation and retrieval
func TestVertex(t *testing.T) {
	ctx := NewPoSST()

	// Create a new vertex
	node1 := Vertex(ctx, "A", "chapter1")
	if node1.S != "A" {
		t.Errorf("Expected name='A', got '%s'", node1.S)
	}
	if node1.Chap != "chapter1" {
		t.Errorf("Expected chapter='chapter1', got '%s'", node1.Chap)
	}
	if node1.NPtr <= 0 {
		t.Errorf("Expected positive NodePtr, got %d", node1.NPtr)
	}

	// Retrieve same vertex (should be idempotent)
	node2 := Vertex(ctx, "A", "chapter1")
	if node1.NPtr != node2.NPtr {
		t.Errorf("Vertex not idempotent: got NodePtrs %d and %d", node1.NPtr, node2.NPtr)
	}

	// Create different vertex
	node3 := Vertex(ctx, "B", "chapter1")
	if node1.NPtr == node3.NPtr {
		t.Error("Different vertices should have different NodePtrs")
	}
}

// TestEdge verifies link creation between nodes
func TestEdge(t *testing.T) {
	ctx := NewPoSST()

	nodeA := Vertex(ctx, "A", "ch1")
	nodeB := Vertex(ctx, "B", "ch1")

	// Create edge A -[fwd]-> B
	if _, _, err := Edge(ctx, nodeA, "fwd", nodeB, []string{"test context"}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}

	// Verify forward link exists in out map
	outLinks, exists := ctx.out[nodeA.NPtr]
	if !exists || len(outLinks) == 0 {
		t.Fatal("Forward link not created in out map")
	}

	foundLink := false
	for _, link := range outLinks {
		if link.Dst == nodeB.NPtr {
			foundLink = true
			if link.Wgt != 1.0 {
				t.Errorf("Expected weight=1.0, got %f", link.Wgt)
			}
			// Verify arrow is fwd
			arrInfo := GetDBArrowByPtr(ctx, link.Arr)
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
	inLinks, exists := ctx.in[nodeB.NPtr]
	if !exists || len(inLinks) == 0 {
		t.Fatal("Backward link not created in in map")
	}

	foundBackLink := false
	for _, link := range inLinks {
		if link.Dst == nodeA.NPtr {
			foundBackLink = true
			arrInfo := GetDBArrowByPtr(ctx, link.Arr)
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
	ctx := NewPoSST()

	// Create simple linear path: A -> B -> C
	nodeA := Vertex(ctx, "A", "ch1")
	nodeB := Vertex(ctx, "B", "ch1")
	nodeC := Vertex(ctx, "C", "ch1")

	if _, _, err := Edge(ctx, nodeA, "fwd", nodeB, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(ctx, nodeB, "fwd", nodeC, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}

	// Find paths from A with depth 1
	paths1, count1 := GetEntireNCConePathsAsLinks(ctx, "fwd", nodeA.NPtr, 1, "ch1", []string{}, 100)
	if count1 != 1 {
		t.Errorf("Expected 1 path at depth 1, got %d", count1)
	}
	if len(paths1) > 0 && len(paths1[0]) != 1 {
		t.Errorf("Expected path length 1, got %d", len(paths1[0]))
	}

	// Find paths from A with depth 2
	paths2, count2 := GetEntireNCConePathsAsLinks(ctx, "fwd", nodeA.NPtr, 2, "ch1", []string{}, 100)
	if count2 != 1 {
		t.Errorf("Expected 1 path at depth 2, got %d", count2)
	}
	if len(paths2) > 0 && len(paths2[0]) != 2 {
		t.Errorf("Expected path length 2, got %d", len(paths2[0]))
	}
}

// TestAdjointLinkPath verifies path reversal with arrow inversion
func TestAdjointLinkPath(t *testing.T) {
	ctx := NewPoSST()

	// Create nodes A, B, C
	Vertex(ctx, "A", "ch1") // nodeA - not directly used but establishes context
	nodeB := Vertex(ctx, "B", "ch1")
	nodeC := Vertex(ctx, "C", "ch1")

	// Get arrow pointer for fwd
	fwdPtr := ctx.arrowName["fwd"]

	// Create forward path A -> B -> C
	// In the Link representation, we only store the destination
	link1 := Link{Dst: nodeB.NPtr, Arr: fwdPtr, Wgt: 1.0} // represents step to B
	link2 := Link{Dst: nodeC.NPtr, Arr: fwdPtr, Wgt: 1.0} // represents step to C
	path := []Link{link1, link2}

	// Reverse the path
	reversed := AdjointLinkPath(ctx, path)

	if len(reversed) != 2 {
		t.Fatalf("Expected reversed path length 2, got %d", len(reversed))
	}

	// The reversed path should have links in reverse order
	// Original: [B, C]  -> Reversed: [C, B]
	if reversed[0].Dst != nodeC.NPtr {
		t.Errorf("First reversed link should point to C, got %d", reversed[0].Dst)
	}
	if reversed[1].Dst != nodeB.NPtr {
		t.Errorf("Second reversed link should point to B, got %d", reversed[1].Dst)
	}

	// Check arrows are inverted to bwd
	bwdPtr := ctx.arrowName["bwd"]
	if reversed[0].Arr != bwdPtr || reversed[1].Arr != bwdPtr {
		t.Error("Arrows not properly inverted to bwd")
	}
}

// TestWaveFrontsOverlap verifies frontier collision detection
func TestWaveFrontsOverlap(t *testing.T) {
	ctx := NewPoSST()

	// Create diamond pattern: Start -> A, B -> End
	start := Vertex(ctx, "Start", "ch1")
	nodeA := Vertex(ctx, "A", "ch1")
	nodeB := Vertex(ctx, "B", "ch1")
	end := Vertex(ctx, "End", "ch1")

	if _, _, err := Edge(ctx, start, "fwd", nodeA, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(ctx, start, "fwd", nodeB, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(ctx, nodeA, "fwd", end, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(ctx, nodeB, "fwd", end, []string{}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}

	// Get forward paths from start
	leftPaths, leftCount := GetEntireNCConePathsAsLinks(ctx, "fwd", start.NPtr, 1, "ch1", []string{}, 100)

	// Get backward paths from end
	rightPaths, rightCount := GetEntireNCConePathsAsLinks(ctx, "bwd", end.NPtr, 1, "ch1", []string{}, 100)

	// Check overlap - should find connections at A and B
	// Use io.Discard to suppress output during test
	dags, cycles := waveFrontsOverlap(ctx, io.Discard, leftPaths, rightPaths, leftCount, rightCount, 1, 1)

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
	ctx1 := NewPoSST()
	ctx2 := NewPoSST()

	node1 := Vertex(ctx1, "A", "ch1")
	node2 := Vertex(ctx2, "A", "ch1")

	// Same name but different contexts should have independent NodePtrs
	if ctx1.name2ptr == nil || ctx2.name2ptr == nil {
		t.Fatal("Context maps not initialized")
	}

	// Each context should have exactly one node
	if len(ctx1.name2ptr) != 1 {
		t.Errorf("ctx1 should have 1 node, got %d", len(ctx1.name2ptr))
	}
	if len(ctx2.name2ptr) != 1 {
		t.Errorf("ctx2 should have 1 node, got %d", len(ctx2.name2ptr))
	}

	// Verify they are distinct in their respective contexts
	if node1.NPtr == node2.NPtr && ctx1 == ctx2 {
		t.Error("Contexts are not properly isolated")
	}
}

// TestMultipleEdges verifies handling of multiple edges between same nodes
func TestMultipleEdges(t *testing.T) {
	ctx := NewPoSST()

	nodeA := Vertex(ctx, "A", "ch1")
	nodeB := Vertex(ctx, "B", "ch1")

	// Add multiple edges with same arrow type
	if _, _, err := Edge(ctx, nodeA, "fwd", nodeB, []string{"context1"}, 1.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}
	if _, _, err := Edge(ctx, nodeA, "fwd", nodeB, []string{"context2"}, 2.0); err != nil {
		t.Fatalf("Edge creation failed: %v", err)
	}

	outLinks := ctx.out[nodeA.NPtr]
	if len(outLinks) != 2 {
		t.Errorf("Expected 2 outgoing links, got %d", len(outLinks))
	}
}

// TestEmptyGraph verifies behavior with no edges
func TestEmptyGraph(t *testing.T) {
	ctx := NewPoSST()

	nodeA := Vertex(ctx, "A", "ch1")

	// Try to find paths in empty graph
	paths, count := GetEntireNCConePathsAsLinks(ctx, "fwd", nodeA.NPtr, 1, "ch1", []string{}, 100)

	if count != 0 {
		t.Errorf("Expected no paths in empty graph, got %d", count)
	}
	if len(paths) != 0 {
		t.Errorf("Expected empty paths slice, got %d paths", len(paths))
	}
}

// TestPathLimit verifies depth and count limits work
func TestPathLimit(t *testing.T) {
	ctx := NewPoSST()

	// Create a node with multiple outgoing edges
	start := Vertex(ctx, "Start", "ch1")
	for i := 0; i < 10; i++ {
		end := Vertex(ctx, string(rune('A'+i)), "ch1")
		if _, _, err := Edge(ctx, start, "fwd", end, []string{}, 1.0); err != nil {
			t.Fatalf("Edge creation failed: %v", err)
		}
	}

	// Request limited number of paths
	paths, count := GetEntireNCConePathsAsLinks(ctx, "fwd", start.NPtr, 1, "ch1", []string{}, 5)

	if count > 5 {
		t.Errorf("Expected at most 5 paths due to limit, got %d", count)
	}
	if len(paths) > 5 {
		t.Errorf("Expected at most 5 paths in result slice, got %d", len(paths))
	}
}

// TestEdgeUnknownArrow tests that Edge returns an error for unknown arrow names
func TestEdgeUnknownArrow(t *testing.T) {
	ctx := NewPoSST()
	defer Close(ctx)

	start := Vertex(ctx, "node1", "ch1")
	end := Vertex(ctx, "node2", "ch1")

	// Try to create an edge with an unknown arrow name
	_, _, err := Edge(ctx, start, "unknown_arrow", end, []string{}, 1.0)
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
