// Package maze demonstrates solving a maze using bidirectional wavefronts and LT vectors.
// Extracted and refactored from API_EXAMPLE_3.go
package maze

import (
	"fmt"
	"io"
	"os"
)

var path = [9][]string{
	{"maze_a7", "maze_b7", "maze_b6", "maze_c6", "maze_c5", "maze_b5", "maze_b4", "maze_a4", "maze_a3", "maze_b3", "maze_c3", "maze_d3", "maze_d2", "maze_e2", "maze_e3", "maze_f3", "maze_f4", "maze_e4", "maze_e5", "maze_f5", "maze_f6", "maze_g6", "maze_g5", "maze_g4", "maze_h4", "maze_h5", "maze_h6", "maze_i6"},
	{"maze_d1", "maze_d2"},
	{"maze_f1", "maze_f2", "maze_e2"},
	{"maze_f2", "maze_g2", "maze_h2", "maze_h3", "maze_g3", "maze_g2"},
	{"maze_b1", "maze_c1", "maze_c2", "maze_b2", "maze_b1"},
	{"maze_b7", "maze_b8", "maze_c8", "maze_c7", "maze_d7", "maze_d6", "maze_e6", "maze_e7", "maze_f7", "maze_f8"},
	{"maze_d7", "maze_d8", "maze_e8", "maze_e7"},
	{"maze_f7", "maze_g7", "maze_g8", "maze_h8", "maze_h7"},
	{"maze_a2", "maze_a1"},
}

// SolveMaze builds the maze graph (from predefined 'path' segments)
// and solves it using contra-colliding wavefronts (bidirectional search).
// It prints any discovered solutions and loop-corrections at each step.
// Output is written to os.Stdout.
// Returns an error if the maze cannot be built or solved.
func SolveMaze() error {
	return SolveMazeWithOutput(os.Stdout)
}

// SolveMazeWithOutput builds the maze graph and solves it using bidirectional search.
// Output is written to the provided writer, allowing for custom output destinations
// (e.g., buffers for testing, files, or stdout).
// Returns an error if the maze cannot be built or solved.
func SolveMazeWithOutput(w io.Writer) error {
	ctx := NewPoSST()

	// Add the paths to a fresh database
	for p := range path {
		for leg := 1; leg < len(path[p]); leg++ {
			chap := "solve maze"
			context := []string{""}
			var weight float32 = 1.0
			nfrom := Vertex(ctx, path[p][leg-1], chap)
			nto := Vertex(ctx, path[p][leg], chap)
			if _, _, err := Edge(ctx, nfrom, "fwd", nto, context, weight); err != nil {
				return fmt.Errorf("failed to create edge from %s to %s: %w", path[p][leg-1], path[p][leg], err)
			}
		}
	}

	if err := solve(ctx, w); err != nil {
		return fmt.Errorf("failed to solve maze: %w", err)
	}
	Close(ctx)
	return nil
}

// solve runs the alternating expansion of left and right wavefronts
// until a maximum depth is reached. It uses GetEntireNCConePathsAsLinks
// to enumerate all tendrils and then splices overlapping pairs.
// Output is written to the provided writer.
// Returns an error if the start or end nodes cannot be found.
func solve(ctx *PoSST, w io.Writer) error {
	const maxdepth = 16
	ldepth, rdepth := 1, 1
	var Lnum, Rnum int
	var count int
	var leftPaths, rightPaths [][]Link

	startBC := "maze_a7"
	endBC := "maze_i6"

	leftPtrs := GetDBNodePtrMatchingName(ctx, startBC, "")
	rightPtrs := GetDBNodePtrMatchingName(ctx, endBC, "")

	if leftPtrs == nil || rightPtrs == nil {
		return fmt.Errorf("no paths available from end points (start=%s, end=%s)", startBC, endBC)
	}

	cntx := []string{""}
	const limit = 10

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {
		leftPaths, Lnum = GetEntireNCConePathsAsLinks(ctx, "fwd", leftPtrs[0], ldepth, "", cntx, limit)
		rightPaths, Rnum = GetEntireNCConePathsAsLinks(ctx, "bwd", rightPtrs[0], rdepth, "", cntx, limit)

		solutions, loopCorrections := waveFrontsOverlap(ctx, w, leftPaths, rightPaths, Lnum, Rnum, ldepth, rdepth)

		if len(solutions) > 0 {
			fmt.Fprintln(w, "-- T R E E ----------------------------------")
			fmt.Fprintf(w, "Path solution %d from %s to %s with lengths %d -%d\n", count, startBC, endBC, ldepth, rdepth)
			for s := 0; s < len(solutions); s++ {
				prefix := fmt.Sprintf(" - story %d: ", s)
				PrintLinkPath(ctx, solutions, s, prefix, "", nil)
			}
			count++
			fmt.Fprintln(w, "-------------------------------------------")
		}

		if len(loopCorrections) > 0 {
			fmt.Fprintln(w, "++ L O O P S +++++++++++++++++++++++++++++++")
			fmt.Fprintf(w, "Path solution %d from %s to %s with lengths %d -%d\n", count, startBC, endBC, ldepth, rdepth)
			for s := range loopCorrections {
				prefix := fmt.Sprintf(" - story %d: ", s)
				PrintLinkPath(ctx, loopCorrections, s, prefix, "", nil)
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
func waveFrontsOverlap(ctx *PoSST, w io.Writer, leftPaths, rightPaths [][]Link, Lnum, Rnum, ldepth, rdepth int) ([][]Link, [][]Link) {
	var solutions [][]Link
	var loops [][]Link

	leftfront := waveFront(leftPaths, Lnum)
	rightfront := waveFront(rightPaths, Rnum)

	fmt.Fprintf(w, "\n  Left front radius %d : %s\n", ldepth, showNode(ctx, leftfront))
	fmt.Fprintf(w, "  Right front radius %d : %s\n", rdepth, showNode(ctx, rightfront))

	incidence := nodesOverlap(ctx, w, leftfront, rightfront)

	for lp := range incidence {
		rp := incidence[lp]
		var LRsplice []Link
		LRsplice = leftJoin(LRsplice, leftPaths[lp])
		adjoint := AdjointLinkPath(ctx, rightPaths[rp])
		LRsplice = rightComplementJoin(LRsplice, adjoint)
		fmt.Fprintf(w, "...SPLICE PATHS L%d with R%d.....\n", lp, rp)
		fmt.Fprintln(w, "Left tendril", showNodePath(ctx, leftPaths[lp]))
		fmt.Fprintln(w, "Right tendril", showNodePath(ctx, rightPaths[rp]))
		fmt.Fprintln(w, "Right adjoint:", showNodePath(ctx, adjoint))
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
func waveFront(path [][]Link, num int) []NodePtr {
	var front []NodePtr
	for l := 0; l < num; l++ {
		front = append(front, path[l][len(path[l])-1].Dst)
	}
	return front
}

// nodesOverlap finds overlaps between the left and right frontier node sets,
// returning an index map: left-path-index -> right-path-index.
// It also logs the names of the touching nodes for visibility.
// Output is written to the provided writer.
func nodesOverlap(ctx *PoSST, w io.Writer, left, right []NodePtr) map[int]int {
	LRsplice := make(map[int]int)
	var list string
	for l := 0; l < len(left); l++ {
		for r := 0; r < len(right); r++ {
			if left[l] == right[r] {
				node := GetDBNodeByNodePtr(ctx, left[l])
				list += node.S + ", "
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
func leftJoin(LRsplice, seq []Link) []Link {
	for i := 0; i < len(seq); i++ {
		LRsplice = append(LRsplice, seq[i])
	}
	return LRsplice
}

// rightComplementJoin appends the right-side path (already adjointed/reversed)
// except its first element, to avoid duplicating the meet node.
func rightComplementJoin(LRsplice, adjoint []Link) []Link {
	for j := 1; j < len(adjoint); j++ {
		LRsplice = append(LRsplice, adjoint[j])
	}
	return LRsplice
}

// isDAG checks the spliced path for repeated destination nodes, which would
// imply a loop. Returns true if the path is acyclic.
func isDAG(seq []Link) bool {
	freq := make(map[NodePtr]int)
	for i := range seq {
		freq[seq[i].Dst]++
	}
	for n := range freq {
		if freq[n] > 1 {
			return false
		}
	}
	return true
}

// showNode renders a comma-separated list of node labels for a frontier set.
func showNode(ctx *PoSST, nptr []NodePtr) string {
	var ret string
	for n := range nptr {
		node := GetDBNodeByNodePtr(ctx, nptr[n])
		ret += node.S + ","
	}
	return ret
}

// showNodePath renders a human-readable representation of a link path with
// arrows and node labels.
func showNodePath(ctx *PoSST, lnk []Link) string {
	var ret string
	for n := range lnk {
		node := GetDBNodeByNodePtr(ctx, lnk[n].Dst)
		arrs := GetDBArrowByPtr(ctx, lnk[n].Arr).Long
		ret += fmt.Sprintf("(%s) -> %s ", arrs, node.S)
	}
	return ret
}
