// Package maze demonstrates solving a maze using bidirectional wavefronts and LT vectors.
// Extracted and refactored from API_EXAMPLE_3.go
package maze

import (
	"fmt"
	"io"
	"os"
)

// Cell type: 0 = wall (blocked), 1 = path (walkable)
type Cell byte

const (
	Wall Cell = 0
	Path Cell = 1
)

// MazeGrid defines a 9x9 maze using binary notation (0=wall, 1=path)
// Rows are A-I (top to bottom), columns are 1-9 (left to right)
// Human-readable: you can see the maze structure visually!
var MazeGrid = [9][9]Cell{
	// 		1  2  3  4  5  6  7  8  9
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, // A
	{1, 1, 1, 1, 1, 1, 1, 1, 0}, // B
	{0, 1, 1, 0, 0, 0, 1, 1, 0}, // C
	{0, 1, 1, 1, 1, 1, 1, 1, 0}, // D
	{0, 1, 1, 1, 1, 0, 1, 1, 0}, // E
	{0, 1, 1, 1, 1, 1, 0, 1, 1}, // F
	{0, 1, 0, 0, 1, 1, 1, 0, 0}, // G
	{0, 1, 1, 1, 1, 1, 1, 0, 0}, // H
	{0, 0, 0, 0, 0, 0, 0, 0, 0}, // I
}

var (
	StartNode = "f9"
	EndNode   = "b1"
)

// SolveMaze builds the maze graph from MazeGrid and solves it using
// contra-colliding wavefronts (bidirectional search).
// It prints any discovered solutions and loop-corrections at each step.
// Output is written to os.Stdout.
// Returns an error if the maze cannot be built or solved.
func SolveMaze() error {
	return SolveMazeWithOutput(os.Stdout)
}

// SolveMazeWithOutput builds the maze graph from MazeGrid and solves it using bidirectional search.
// Output is written to the provided writer, allowing for custom output destinations
// (e.g., buffers for testing, files, or stdout).
// Returns an error if the maze cannot be built or solved.
func SolveMazeWithOutput(w io.Writer) error {
	graph := NewLinkedSST()

	// Build graph from maze grid structure
	if err := buildGraphFromGrid(graph); err != nil {
		return fmt.Errorf("failed to build maze graph: %w", err)
	}

	if err := solve(graph, w); err != nil {
		return fmt.Errorf("failed to solve maze: %w", err)
	}
	Close(graph)
	return nil
}

// buildGraphFromGrid creates graph nodes and edges from the MazeGrid.
// For each walkable cell, it creates a node and connects it to adjacent walkable cells.
func buildGraphFromGrid(graph *LinkedSST) error {
	chap := "solve maze"
	context := []string{""}
	var weight float32 = 1.0

	// Iterate through the grid and create edges for adjacent walkable cells
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if MazeGrid[row][col] == Wall {
				continue // Skip walls
			}

			currentNode := cellToNodeID(row, col)
			current := Vertex(graph, currentNode, chap)

			// Check right neighbor (col+1)
			if col+1 < 9 && MazeGrid[row][col+1] == Path {
				rightNode := cellToNodeID(row, col+1)
				right := Vertex(graph, rightNode, chap)
				if _, _, err := Edge(graph, current, "fwd", right, context, weight); err != nil {
					return fmt.Errorf("failed to create edge from %s to %s: %w", currentNode, rightNode, err)
				}
			}

			// Check down neighbor (row+1)
			if row+1 < 9 && MazeGrid[row+1][col] == Path {
				downNode := cellToNodeID(row+1, col)
				down := Vertex(graph, downNode, chap)
				if _, _, err := Edge(graph, current, "fwd", down, context, weight); err != nil {
					return fmt.Errorf("failed to create edge from %s to %s: %w", currentNode, downNode, err)
				}
			}
		}
	}

	return nil
}

// cellToNodeID converts grid coordinates to node ID format: "a1", "b7", etc.
// row: 0-8 (A-I), col: 0-8 (1-9)
func cellToNodeID(row, col int) string {
	return fmt.Sprintf("%c%d", 'a'+row, col+1)
}

// solve runs the alternating expansion of left and right wavefronts
// until a maximum depth is reached. It uses GetEntireNCConePathsAsLinks
// to enumerate all tendrils and then splices overlapping pairs.
// Output is written to the provided writer.
// Returns an error if the start or end nodes cannot be found.
func solve(graph *LinkedSST, w io.Writer) error {
	const maxdepth = 16
	ldepth, rdepth := 1, 1
	var Lnum, Rnum int
	var count int
	var leftPaths, rightPaths [][]*Link

	leftNode := GetNodeByName(graph, StartNode, "")
	rightNode := GetNodeByName(graph, EndNode, "")

	if leftNode == nil || rightNode == nil {
		return fmt.Errorf("no paths available from end points (start=%s, end=%s)", StartNode, EndNode)
	}

	cntx := []string{""}
	const limit = 10

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {
		leftPaths, Lnum = GetEntireNCConePathsAsLinks(graph, "fwd", leftNode, ldepth, "", cntx, limit)
		rightPaths, Rnum = GetEntireNCConePathsAsLinks(graph, "bwd", rightNode, rdepth, "", cntx, limit)

		solutions, loopCorrections := waveFrontsOverlap(graph, w, leftPaths, rightPaths, Lnum, Rnum, ldepth, rdepth)

		if len(solutions) > 0 {
			fmt.Fprintln(w, "-- T R E E ----------------------------------")
			fmt.Fprintf(w, "Path solution %d from %s to %s with lengths %d -%d\n", count, StartNode, EndNode, ldepth, rdepth)
			for s := 0; s < len(solutions); s++ {
				prefix := fmt.Sprintf(" - story %d: ", s)
				PrintLinkPath(graph, solutions, s, prefix, "", nil)
			}
			count++
			fmt.Fprintln(w, "-------------------------------------------")
		}

		if len(loopCorrections) > 0 {
			fmt.Fprintln(w, "++ L O O P S +++++++++++++++++++++++++++++++")
			fmt.Fprintf(w, "Path solution %d from %s to %s with lengths %d -%d\n", count, StartNode, EndNode, ldepth, rdepth)
			for s := range loopCorrections {
				prefix := fmt.Sprintf(" - story %d: ", s)
				PrintLinkPath(graph, loopCorrections, s, prefix, "", nil)
			}
			count++
			fmt.Fprintln(w, "+++++++++++++++++++++++++++++++++++++++++++")
		}

		if turn%2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}
	return nil
}

// waveFrontsOverlap examines the current frontier tendrils from left and right,
// detects collisions at common nodes, and splices the respective path segments
// into full solutions. Non-DAG splices are returned as loop-corrections.
// Output is written to the provided writer.
func waveFrontsOverlap(graph *LinkedSST, w io.Writer, leftPaths, rightPaths [][]*Link, Lnum, Rnum, ldepth, rdepth int) ([][]*Link, [][]*Link) {
	var solutions [][]*Link
	var loops [][]*Link

	leftfront := waveFront(leftPaths, Lnum)
	rightfront := waveFront(rightPaths, Rnum)

	fmt.Fprintf(w, "\n  Left front radius %d : %s\n", ldepth, showNode(graph, leftfront))
	fmt.Fprintf(w, "  Right front radius %d : %s\n", rdepth, showNode(graph, rightfront))

	incidence := nodesOverlap(graph, w, leftfront, rightfront)

	for lp := range incidence {
		rp := incidence[lp]
		var LRsplice []*Link
		LRsplice = leftJoin(LRsplice, leftPaths[lp])
		adjoint := AdjointLinkPath(graph, rightPaths[rp])
		LRsplice = rightComplementJoin(LRsplice, adjoint)
		fmt.Fprintf(w, "...SPLICE PATHS L%d with R%d.....\n", lp, rp)
		fmt.Fprintln(w, "Left tendril", showNodePath(graph, leftPaths[lp]))
		fmt.Fprintln(w, "Right tendril", showNodePath(graph, rightPaths[rp]))
		fmt.Fprintln(w, "Right adjoint:", showNodePath(graph, adjoint))
		fmt.Fprint(w, ".....................\n\n")
		if isDAG(LRsplice) {
			solutions = append(solutions, LRsplice)
		} else {
			loops = append(loops, LRsplice)
		}
	}
	fmt.Fprintf(w, "  (found %d touching solutions)\n", len(incidence))
	return solutions, loops
}

// waveFront returns the tip nodes of each path in a set of link paths,
// representing the current search frontier.
func waveFront(path [][]*Link, num int) []*Node {
	var front []*Node
	for l := 0; l < num; l++ {
		front = append(front, path[l][len(path[l])-1].dst)
	}
	return front
}

// nodesOverlap finds overlaps between the left and right frontier node sets,
// returning an index map: left-path-index -> right-path-index.
// It also logs the names of the touching nodes for visibility.
// Output is written to the provided writer.
func nodesOverlap(graph *LinkedSST, w io.Writer, left, right []*Node) map[int]int {
	LRsplice := make(map[int]int)
	var list string
	for l := 0; l < len(left); l++ {
		for r := 0; r < len(right); r++ {
			if left[l] == right[r] {
				list += left[l].label + ", "
				LRsplice[l] = r
			}
		}
	}
	if len(list) > 0 {
		fmt.Fprintf(w, "  (i.e. waves impinge%dtimes at: %s)\n\n", len(LRsplice), list)
	}
	return LRsplice
}

// leftJoin appends the left-side path sequence into the splice buffer.
func leftJoin(LRsplice, seq []*Link) []*Link {
	for i := 0; i < len(seq); i++ {
		LRsplice = append(LRsplice, seq[i])
	}
	return LRsplice
}

// rightComplementJoin appends the right-side path (already adjointed/reversed)
// except its first element, to avoid duplicating the meet node.
func rightComplementJoin(LRsplice, adjoint []*Link) []*Link {
	for j := 1; j < len(adjoint); j++ {
		LRsplice = append(LRsplice, adjoint[j])
	}
	return LRsplice
}

// isDAG checks the spliced path for repeated destination nodes, which would
// imply a loop. Returns true if the path is acyclic.
func isDAG(seq []*Link) bool {
	freq := make(map[*Node]int)
	for i := range seq {
		freq[seq[i].dst]++
	}
	for n := range freq {
		if freq[n] > 1 {
			return false
		}
	}
	return true
}

// showNode renders a comma-separated list of node labels for a frontier set.
func showNode(graph *LinkedSST, nodes []*Node) string {
	var s string
	for n := range nodes {
		s += nodes[n].label + ","
	}
	return s
}

// showNodePath renders a human-readable representation of a link path with
// arrows and node labels.
func showNodePath(graph *LinkedSST, lnk []*Link) string {
	var ret string
	for n := range lnk {
		ret += fmt.Sprintf("(%s) -> %s ", lnk[n].arrow.long, lnk[n].dst.label)
	}
	return ret
}
