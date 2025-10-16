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
	graph := NewLinkedSST()
	defer Close(graph)

	// Build the maze from grid
	if err := buildGraphFromGrid(graph); err != nil {
		return nil, fmt.Errorf("failed to build maze graph: %w", err)
	}

	// Solve and collect results
	result, err := solveJSON(graph, w)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// solveJSON runs bidirectional search and collects results as structured data.
func solveJSON(graph *LinkedSST, w io.Writer) (*MazeResult, error) {
	const maxdepth = 16
	ldepth, rdepth := 1, 1
	var Lnum, Rnum int
	var leftPaths, rightPaths [][]*Link

	leftNode := GetNodeByName(graph, StartNode, "")
	rightNode := GetNodeByName(graph, EndNode, "")

	if leftNode == nil || rightNode == nil {
		return nil, fmt.Errorf("no paths available from end points (start=%s, end=%s)", StartNode, EndNode)
	}

	cntx := []string{""}
	const limit = 10

	result := &MazeResult{
		StartNode:   StartNode,
		EndNode:     EndNode,
		MaxDepth:    maxdepth,
		Solutions:   []Solution{},
		Loops:       []Solution{},
		SearchSteps: []SearchStep{},
		Statistics:  Statistics{},
	}

	solutionCount := 0

	for turn := 0; ldepth < maxdepth && rdepth < maxdepth; turn++ {
		leftPaths, Lnum = GetEntireNCConePathsAsLinks(graph, "fwd", leftNode, ldepth, "", cntx, limit)
		rightPaths, Rnum = GetEntireNCConePathsAsLinks(graph, "bwd", rightNode, rdepth, "", cntx, limit)

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
func getFrontierNodes(graph *LinkedSST, paths [][]*Link, count int) []string {
	nodeSet := make(map[string]bool)
	for i := 0; i < count && i < len(paths); i++ {
		if len(paths[i]) > 0 {
			lastLink := paths[i][len(paths[i])-1]
			nodeSet[lastLink.dst.label] = true
		}
	}

	nodes := make([]string, 0, len(nodeSet))
	for name := range nodeSet {
		nodes = append(nodes, name)
	}
	return nodes
}

// linksToPathLinks converts internal Link representation to JSON-friendly PathLink format.
func linksToPathLinks(graph *LinkedSST, links []*Link) []PathLink {
	if len(links) == 0 {
		return []PathLink{}
	}

	pathLinks := make([]PathLink, 0, len(links))

	// First node is the starting point
	var prevNode string
	if len(links) > 0 {
		prevNode = links[0].dst.label
	}

	// Convert each link
	for _, link := range links {
		pathLinks = append(pathLinks, PathLink{
			From:   prevNode,
			To:     link.dst.label,
			Arrow:  link.arrow.long,
			Weight: link.weight,
		})

		prevNode = link.dst.label
	}

	return pathLinks
}
