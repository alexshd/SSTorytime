package maze

import (
	"encoding/json"
)

// MazeResult represents the complete solution output in JSON format.
type MazeResult struct {
	StartNode   string       `json:"start_node"`
	EndNode     string       `json:"end_node"`
	MaxDepth    int          `json:"max_depth"`
	Solutions   []Solution   `json:"solutions"`
	Loops       []Solution   `json:"loops"`
	SearchSteps []SearchStep `json:"search_steps"`
	Statistics  Statistics   `json:"statistics"`
}

// Solution represents a complete path from start to end.
type Solution struct {
	ID          int        `json:"id"`
	Type        string     `json:"type"` // "tree" or "loop"
	LeftDepth   int        `json:"left_depth"`
	RightDepth  int        `json:"right_depth"`
	TotalLength int        `json:"total_length"`
	Path        []PathLink `json:"path"`
}

// PathLink represents a single edge in a path.
type PathLink struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Arrow  string  `json:"arrow"`
	Weight float32 `json:"weight"`
}

// SearchStep represents one iteration of the bidirectional search.
type SearchStep struct {
	Turn          int      `json:"turn"`
	LeftDepth     int      `json:"left_depth"`
	RightDepth    int      `json:"right_depth"`
	LeftFrontier  []string `json:"left_frontier"`
	RightFrontier []string `json:"right_frontier"`
	Solutions     int      `json:"solutions_found"`
	Loops         int      `json:"loops_found"`
}

// Statistics provides summary information about the search.
type Statistics struct {
	TotalSolutions   int `json:"total_solutions"`
	TotalLoops       int `json:"total_loops"`
	MaxLeftDepth     int `json:"max_left_depth"`
	MaxRightDepth    int `json:"max_right_depth"`
	TotalSearchSteps int `json:"total_search_steps"`
}

// ToJSON converts MazeResult to JSON bytes.
func (mr *MazeResult) ToJSON() ([]byte, error) {
	return json.MarshalIndent(mr, "", "  ")
}

// ToJSONCompact converts MazeResult to compact JSON bytes.
func (mr *MazeResult) ToJSONCompact() ([]byte, error) {
	return json.Marshal(mr)
}
