package maze

import (
	"io"
	"testing"
)

// BenchmarkOpen measures the cost of initializing a new context
func BenchmarkOpen(b *testing.B) {
	for b.Loop() {
		NewLinkedSST()
	}
}

// BenchmarkVertex measures node creation performance
func BenchmarkVertex(b *testing.B) {
	graph := NewLinkedSST()

	for b.Loop() {
		Vertex(graph, "A", "chapter1")
	}
}

// BenchmarkVertexUnique measures node creation with unique names
func BenchmarkVertexUnique(b *testing.B) {
	graph := NewLinkedSST()

	for i := 0; b.Loop(); i++ {
		name := string(rune('A' + (i % 26)))
		Vertex(graph, name, "chapter1")
	}
}

// BenchmarkEdge measures link creation performance
func BenchmarkEdge(b *testing.B) {
	graph := NewLinkedSST()
	nodeA := Vertex(graph, "A", "ch1")
	nodeB := Vertex(graph, "B", "ch1")
	context := []string{"test context"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Edge(graph, nodeA, "fwd", nodeB, context, 1.0)
	}
}

// BenchmarkGraphBuilding measures the cost of building a small graph
func BenchmarkGraphBuilding(b *testing.B) {
	for i := 0; i < b.N; i++ {
		graph := NewLinkedSST()

		// Build 10 node chain
		nodes := make([]*Node, 10)
		for j := 0; j < 10; j++ {
			nodes[j] = Vertex(graph, string(rune('A'+j)), "ch1")
		}

		// Connect them sequentially
		for j := 0; j < 9; j++ {
			Edge(graph, nodes[j], "fwd", nodes[j+1], []string{}, 1.0)
		}
	}
}

// BenchmarkGetEntireNCConePathsAsLinks_Depth1 measures path finding at depth 1
func BenchmarkGetEntireNCConePathsAsLinks_Depth1(b *testing.B) {
	graph := NewLinkedSST()
	start := Vertex(graph, "Start", "ch1")

	// Create 10 direct connections
	for i := 0; i < 10; i++ {
		end := Vertex(graph, string(rune('A'+i)), "ch1")
		Edge(graph, start, "fwd", end, []string{}, 1.0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEntireNCConePathsAsLinks(graph, "fwd", start, 1, "ch1", []string{}, 100)
	}
}

// BenchmarkGetEntireNCConePathsAsLinks_Depth5 measures path finding at depth 5
func BenchmarkGetEntireNCConePathsAsLinks_Depth5(b *testing.B) {
	graph := NewLinkedSST()

	// Build a binary tree of depth 5
	nodes := []*Node{Vertex(graph, "root", "ch1")}
	for depth := 0; depth < 5; depth++ {
		currentLevel := nodes[len(nodes)-(1<<depth):]
		for _, node := range currentLevel {
			left := Vertex(graph, node.label+"L", "ch1")
			right := Vertex(graph, node.label+"R", "ch1")
			Edge(graph, node, "fwd", left, []string{}, 1.0)
			Edge(graph, node, "fwd", right, []string{}, 1.0)
			nodes = append(nodes, left, right)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEntireNCConePathsAsLinks(graph, "fwd", nodes[0], 5, "ch1", []string{}, 100)
	}
}

// BenchmarkAdjointLinkPath measures path reversal performance
func BenchmarkAdjointLinkPath(b *testing.B) {
	graph := NewLinkedSST()

	// Create a path of length 10
	nodes := make([]*Node, 11)
	for i := range nodes {
		nodes[i] = Vertex(graph, string(rune('A'+i)), "ch1")
	}

	fwdArrow := graph.arrow["fwd"]
	path := make([]*Link, 10)
	for i := 0; i < 10; i++ {
		path[i] = &Link{
			dst:    nodes[i+1],
			arrow:  fwdArrow,
			weight: 1.0,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AdjointLinkPath(graph, path)
	}
}

// BenchmarkWaveFrontsOverlap measures collision detection performance
func BenchmarkWaveFrontsOverlap(b *testing.B) {
	graph := NewLinkedSST()

	// Create diamond structure with 10 parallel paths
	start := Vertex(graph, "Start", "ch1")
	end := Vertex(graph, "End", "ch1")

	for i := 0; i < 10; i++ {
		middle := Vertex(graph, string(rune('A'+i)), "ch1")
		Edge(graph, start, "fwd", middle, []string{}, 1.0)
		Edge(graph, middle, "fwd", end, []string{}, 1.0)
	}

	// Pre-compute wavefronts
	leftPaths, leftCount := GetEntireNCConePathsAsLinks(graph, "fwd", start, 1, "ch1", []string{}, 100)
	rightPaths, rightCount := GetEntireNCConePathsAsLinks(graph, "bwd", end, 1, "ch1", []string{}, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		waveFrontsOverlap(graph, io.Discard, leftPaths, rightPaths, leftCount, rightCount, 1, 1)
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
	graph := NewLinkedSST()

	// Create two sets of paths with some overlap
	start1 := Vertex(graph, "Start1", "ch1")
	start2 := Vertex(graph, "Start2", "ch1")
	shared := Vertex(graph, "Shared", "ch1")

	for i := range 5 {
		node := Vertex(graph, string(rune('A'+i)), "ch1")
		Edge(graph, start1, "fwd", node, []string{}, 1.0)
		Edge(graph, start2, "fwd", node, []string{}, 1.0)
	}
	Edge(graph, start1, "fwd", shared, []string{}, 1.0)
	Edge(graph, start2, "fwd", shared, []string{}, 1.0)

	paths1, count1 := GetEntireNCConePathsAsLinks(graph, "fwd", start1, 1, "ch1", []string{}, 100)
	paths2, count2 := GetEntireNCConePathsAsLinks(graph, "fwd", start2, 1, "ch1", []string{}, 100)

	// Extract wavefronts
	front1 := waveFront(paths1, count1)
	front2 := waveFront(paths2, count2)

	for b.Loop() {
		nodesOverlap(graph, io.Discard, front1, front2)
	}
}

// BenchmarkIsDAG measures cycle detection performance
func BenchmarkIsDAG(b *testing.B) {
	graph := NewLinkedSST()

	// Create a path that is a DAG
	nodes := make([]*Node, 10)
	for i := range nodes {
		nodes[i] = Vertex(graph, string(rune('A'+i)), "ch1")
	}

	fwdArrow := graph.arrow["fwd"]
	path := make([]*Link, 9)
	for i := range 9 {
		path[i] = &Link{
			dst:    nodes[i+1],
			arrow:  fwdArrow,
			weight: 1.0,
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
		graph := NewLinkedSST()

		// Build the maze graph
		e1 := Vertex(graph, "e1", "chapter")
		e2 := Vertex(graph, "e2", "chapter")
		e3 := Vertex(graph, "e3", "chapter")
		e4 := Vertex(graph, "e4", "chapter")
		e5 := Vertex(graph, "e5", "chapter")
		e6 := Vertex(graph, "e6", "chapter")
		e7 := Vertex(graph, "e7", "chapter")
		e8 := Vertex(graph, "e8", "chapter")
		e9 := Vertex(graph, "e9", "chapter")
		e10 := Vertex(graph, "e10", "chapter")

		Edge(graph, e1, "fwd", e2, []string{"path1"}, 0)
		Edge(graph, e2, "fwd", e3, []string{"path1"}, 0)
		Edge(graph, e2, "fwd", e5, []string{"path2"}, 0)
		Edge(graph, e3, "fwd", e4, []string{"path3"}, 0)
		Edge(graph, e4, "fwd", e6, []string{"path4"}, 0)
		Edge(graph, e5, "fwd", e6, []string{"path5"}, 0)
		Edge(graph, e7, "fwd", e8, []string{"path6"}, 0)
		Edge(graph, e8, "fwd", e9, []string{"path7"}, 0)
		Edge(graph, e8, "fwd", e10, []string{"path8"}, 0)
		Edge(graph, e9, "fwd", e6, []string{"path9"}, 0)
	}
}
