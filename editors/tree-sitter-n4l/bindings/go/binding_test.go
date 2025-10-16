package tree_sitter_n4l_test

import (
	"testing"

	tree_sitter "github.com/smacker/go-tree-sitter"
	"github.com/tree-sitter/tree-sitter-n4l"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_n4l.Language())
	if language == nil {
		t.Errorf("Error loading N4l grammar")
	}
}
