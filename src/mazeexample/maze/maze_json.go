package maze

import (
	"fmt"
	"io"
)

// SolveMazeJSON solves the maze and returns results as JSON.
// This is the JSON-output equivalent of SolveMaze().
func SolveMazeJSON() (*MazeResult, error) {
	return SolveMazeJSONWithOutput(io.Discard)
}

// SolveMazeJSONWithOutput solves the maze and returns results as JSON.
// The writer w receives human-readable text output (can be io.Discard to suppress).
// Returns a MazeResult struct that can be marshaled to JSON.
func SolveMazeJSONWithOutput(w io.Writer) (*MazeResult, error) {
	graph := NewPoSST()
	defer Close(graph)

	// Build the maze
	for p := range path {
		for leg := 1; leg < len(path[p]); leg++ {
			chap := "solve maze"
			context := []string{""}
			var weight float32 = 1.0
			nfrom := Vertex(graph, path[p][leg-1], chap)
			nto := Vertex(graph, path[p][leg], chap)
			if _, _, err := Edge(graph, nfrom, "fwd", nto, context, weight); err != nil {
				return nil, err
			}
		}
	}

	// Solve and collect results
	result, err := solveJSON(graph, w)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// solveJSON runs bidirectional search and collects results as structured data.
func solveJSON(graph *PoSST, w io.Writer) (*MazeResult, error) {
	const maxdepth = 16
	ldepth, rdepth := 1, 1
	var Lnum, Rnum int
	var leftPaths, rightPaths [][]Link

	startBC := "maze_a7"
	endBC := "maze_i6"

	leftPtrs := GetDBNodePtrMatchingName(graph, startBC, "")
	rightPtrs := GetDBNodePtrMatchingName(graph, endBC, "")

	if leftPtrs == nil || rightPtrs == nil {
		return nil, fmt.Errorf("no paths available from end points (start=%s, end=%s)", startBC, endBC)
	}

	cntx := []string{""}
	const limit = 10

	result := &MazeResult{
		StartNode:   startBC,
		EndNode:     endBC,
		MaxDepth:    maxdepth,
		Solutions:   []Solution{},
		Loops:       []Solution{},
		SearchSteps: []SearchStep{},
		Statistics:  Statistics{},
	}

	solutionCount := 0

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {
		leftPaths, Lnum = GetEntireNCConePathsAsLinks(graph, "fwd", leftPtrs[0], ldepth, "", cntx, limit)
		rightPaths, Rnum = GetEntireNCConePathsAsLinks(graph, "bwd", rightPtrs[0], rdepth, "", cntx, limit)

		// Collect search step
		step := SearchStep{
			Turn:          turn,
			LeftDepth:     ldepth,
			RightDepth:    rdepth,
			LeftFrontier:  getFrontierNodes(graph, leftPaths, Lnum),
			RightFrontier: getFrontierNodes(graph, rightPaths, Rnum),
		}

		solutions, loopCorrections := waveFrontsOverlap(graph, w, leftPaths, rightPaths, Lnum, Rnum, ldepth, rdepth)

		step.Solutions = len(solutions)
		step.Loops = len(loopCorrections)
		result.SearchSteps = append(result.SearchSteps, step)

		// Convert solutions to JSON format
		if len(solutions) > 0 {
			for _, sol := range solutions {
				result.Solutions = append(result.Solutions, Solution{
					ID:          solutionCount,
					Type:        "tree",
					LeftDepth:   ldepth,
					RightDepth:  rdepth,
					TotalLength: len(sol),
					Path:        linksToPathLinks(graph, sol),
				})
				solutionCount++
			}
		}

		// Convert loop corrections
		if len(loopCorrections) > 0 {
			for _, loop := range loopCorrections {
				result.Loops = append(result.Loops, Solution{
					ID:          len(result.Loops),
					Type:        "loop",
					LeftDepth:   ldepth,
					RightDepth:  rdepth,
					TotalLength: len(loop),
					Path:        linksToPathLinks(graph, loop),
				})
			}
		}

		if turn%2 == 0 {
			ldepth++
		} else {
			rdepth++
		}
	}

	// Update statistics
	result.Statistics.TotalSolutions = len(result.Solutions)
	result.Statistics.TotalLoops = len(result.Loops)
	result.Statistics.MaxLeftDepth = ldepth - 1
	result.Statistics.MaxRightDepth = rdepth - 1
	result.Statistics.TotalSearchSteps = len(result.SearchSteps)

	return result, nil
}

// getFrontierNodes extracts the frontier node names from paths.
func getFrontierNodes(graph *PoSST, paths [][]Link, count int) []string {
	nodeSet := make(map[string]bool)
	for i := 0; i < count && i < len(paths); i++ {
		if len(paths[i]) > 0 {
			lastLink := paths[i][len(paths[i])-1]
			node := GetDBNodeByNodePtr(graph, lastLink.Dst)
			nodeSet[node.S] = true
		}
	}

	nodes := make([]string, 0, len(nodeSet))
	for name := range nodeSet {
		nodes = append(nodes, name)
	}
	return nodes
}

// linksToPathLinks converts internal Link representation to JSON-friendly PathLink format.
func linksToPathLinks(graph *PoSST, links []Link) []PathLink {
	if len(links) == 0 {
		return []PathLink{}
	}

	pathLinks := make([]PathLink, 0, len(links))

	// First node is the starting point
	var prevNode string
	if len(links) > 0 {
		firstNode := GetDBNodeByNodePtr(graph, links[0].Dst)
		prevNode = firstNode.S
	}

	// Convert each link
	for _, link := range links {
		node := GetDBNodeByNodePtr(graph, link.Dst)
		arrow := GetDBArrowByPtr(graph, link.Arr)

		pathLinks = append(pathLinks, PathLink{
			From:   prevNode,
			To:     node.S,
			Arrow:  arrow.Long,
			Weight: link.Wgt,
		})

		prevNode = node.S
	}

	return pathLinks
}
