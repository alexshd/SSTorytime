package maze

import (
	"io"
	"testing"
)

// BenchmarkOpen measures the cost of initializing a new context
func BenchmarkOpen(b *testing.B) {
	for b.Loop() {
		NewPoSST()
	}
}

// BenchmarkVertex measures node creation performance
func BenchmarkVertex(b *testing.B) {
	ctx := NewPoSST()

	for b.Loop() {
		Vertex(ctx, "A", "chapter1")
	}
}

// BenchmarkVertexUnique measures node creation with unique names
func BenchmarkVertexUnique(b *testing.B) {
	ctx := NewPoSST()

	for i := 0; b.Loop(); i++ {
		name := string(rune('A' + (i % 26)))
		Vertex(ctx, name, "chapter1")
	}
}

// BenchmarkEdge measures link creation performance
func BenchmarkEdge(b *testing.B) {
	ctx := NewPoSST()
	nodeA := Vertex(ctx, "A", "ch1")
	nodeB := Vertex(ctx, "B", "ch1")
	context := []string{"test context"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Edge(ctx, nodeA, "fwd", nodeB, context, 1.0)
	}
}

// BenchmarkGraphBuilding measures the cost of building a small graph
func BenchmarkGraphBuilding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := NewPoSST()

		// Build 10 node chain
		nodes := make([]Node, 10)
		for j := 0; j < 10; j++ {
			nodes[j] = Vertex(ctx, string(rune('A'+j)), "ch1")
		}

		// Connect them sequentially
		for j := 0; j < 9; j++ {
			Edge(ctx, nodes[j], "fwd", nodes[j+1], []string{}, 1.0)
		}
	}
}

// BenchmarkGetEntireNCConePathsAsLinks_Depth1 measures path finding at depth 1
func BenchmarkGetEntireNCConePathsAsLinks_Depth1(b *testing.B) {
	ctx := NewPoSST()
	start := Vertex(ctx, "Start", "ch1")

	// Create 10 direct connections
	for i := 0; i < 10; i++ {
		end := Vertex(ctx, string(rune('A'+i)), "ch1")
		Edge(ctx, start, "fwd", end, []string{}, 1.0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEntireNCConePathsAsLinks(ctx, "fwd", start.NPtr, 1, "ch1", []string{}, 100)
	}
}

// BenchmarkGetEntireNCConePathsAsLinks_Depth5 measures path finding at depth 5
func BenchmarkGetEntireNCConePathsAsLinks_Depth5(b *testing.B) {
	ctx := NewPoSST()

	// Build a binary tree of depth 5
	nodes := []Node{Vertex(ctx, "root", "ch1")}
	for depth := 0; depth < 5; depth++ {
		currentLevel := nodes[len(nodes)-(1<<depth):]
		for _, node := range currentLevel {
			left := Vertex(ctx, node.S+"L", "ch1")
			right := Vertex(ctx, node.S+"R", "ch1")
			Edge(ctx, node, "fwd", left, []string{}, 1.0)
			Edge(ctx, node, "fwd", right, []string{}, 1.0)
			nodes = append(nodes, left, right)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEntireNCConePathsAsLinks(ctx, "fwd", nodes[0].NPtr, 5, "ch1", []string{}, 100)
	}
}

// BenchmarkAdjointLinkPath measures path reversal performance
func BenchmarkAdjointLinkPath(b *testing.B) {
	ctx := NewPoSST()

	// Create a path of length 10
	nodes := make([]Node, 11)
	for i := range nodes {
		nodes[i] = Vertex(ctx, string(rune('A'+i)), "ch1")
	}

	fwdPtr := ctx.arrowName["fwd"]
	path := make([]Link, 10)
	for i := 0; i < 10; i++ {
		path[i] = Link{
			Dst: nodes[i+1].NPtr,
			Arr: fwdPtr,
			Wgt: 1.0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AdjointLinkPath(ctx, path)
	}
}

// BenchmarkWaveFrontsOverlap measures collision detection performance
func BenchmarkWaveFrontsOverlap(b *testing.B) {
	ctx := NewPoSST()

	// Create diamond structure with 10 parallel paths
	start := Vertex(ctx, "Start", "ch1")
	end := Vertex(ctx, "End", "ch1")

	for i := 0; i < 10; i++ {
		middle := Vertex(ctx, string(rune('A'+i)), "ch1")
		Edge(ctx, start, "fwd", middle, []string{}, 1.0)
		Edge(ctx, middle, "fwd", end, []string{}, 1.0)
	}

	// Pre-compute wavefronts
	leftPaths, leftCount := GetEntireNCConePathsAsLinks(ctx, "fwd", start.NPtr, 1, "ch1", []string{}, 100)
	rightPaths, rightCount := GetEntireNCConePathsAsLinks(ctx, "bwd", end.NPtr, 1, "ch1", []string{}, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		waveFrontsOverlap(ctx, io.Discard, leftPaths, rightPaths, leftCount, rightCount, 1, 1)
	}
}

// BenchmarkSolveMaze measures the full maze solving algorithm
func BenchmarkSolveMaze(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SolveMaze()
	}
}

// BenchmarkNodesOverlap measures node intersection detection
func BenchmarkNodesOverlap(b *testing.B) {
	ctx := NewPoSST()

	// Create two sets of paths with some overlap
	start1 := Vertex(ctx, "Start1", "ch1")
	start2 := Vertex(ctx, "Start2", "ch1")
	shared := Vertex(ctx, "Shared", "ch1")

	for i := range 5 {
		node := Vertex(ctx, string(rune('A'+i)), "ch1")
		Edge(ctx, start1, "fwd", node, []string{}, 1.0)
		Edge(ctx, start2, "fwd", node, []string{}, 1.0)
	}
	Edge(ctx, start1, "fwd", shared, []string{}, 1.0)
	Edge(ctx, start2, "fwd", shared, []string{}, 1.0)

	paths1, count1 := GetEntireNCConePathsAsLinks(ctx, "fwd", start1.NPtr, 1, "ch1", []string{}, 100)
	paths2, count2 := GetEntireNCConePathsAsLinks(ctx, "fwd", start2.NPtr, 1, "ch1", []string{}, 100)

	// Extract wavefronts
	front1 := waveFront(paths1, count1)
	front2 := waveFront(paths2, count2)

	for b.Loop() {
		nodesOverlap(ctx, io.Discard, front1, front2)
	}
}

// BenchmarkIsDAG measures cycle detection performance
func BenchmarkIsDAG(b *testing.B) {
	ctx := NewPoSST()

	// Create a path that is a DAG
	nodes := make([]Node, 10)
	for i := range nodes {
		nodes[i] = Vertex(ctx, string(rune('A'+i)), "ch1")
	}

	fwdPtr := ctx.arrowName["fwd"]
	path := make([]Link, 9)
	for i := range 9 {
		path[i] = Link{
			Dst: nodes[i+1].NPtr,
			Arr: fwdPtr,
			Wgt: 1.0,
		}
	}

	for b.Loop() {
		isDAG(path)
	}
}

// BenchmarkMemoryAllocation measures allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		ctx := NewPoSST()

		// Build the maze graph
		e1 := Vertex(ctx, "e1", "chapter")
		e2 := Vertex(ctx, "e2", "chapter")
		e3 := Vertex(ctx, "e3", "chapter")
		e4 := Vertex(ctx, "e4", "chapter")
		e5 := Vertex(ctx, "e5", "chapter")
		e6 := Vertex(ctx, "e6", "chapter")
		e7 := Vertex(ctx, "e7", "chapter")
		e8 := Vertex(ctx, "e8", "chapter")
		e9 := Vertex(ctx, "e9", "chapter")
		e10 := Vertex(ctx, "e10", "chapter")

		Edge(ctx, e1, "fwd", e2, []string{"path1"}, 0)
		Edge(ctx, e2, "fwd", e3, []string{"path1"}, 0)
		Edge(ctx, e2, "fwd", e5, []string{"path2"}, 0)
		Edge(ctx, e3, "fwd", e4, []string{"path3"}, 0)
		Edge(ctx, e4, "fwd", e6, []string{"path4"}, 0)
		Edge(ctx, e5, "fwd", e6, []string{"path5"}, 0)
		Edge(ctx, e7, "fwd", e8, []string{"path6"}, 0)
		Edge(ctx, e8, "fwd", e9, []string{"path7"}, 0)
		Edge(ctx, e8, "fwd", e10, []string{"path8"}, 0)
		Edge(ctx, e9, "fwd", e6, []string{"path9"}, 0)
	}
}
