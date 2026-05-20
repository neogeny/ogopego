package fragments

import (
	"github.com/neogeny/ogopego/trees"
)

// Semantics Tree represents the concrete polymorphic node (interface/any)
type Semantics func(rule string, node trees.Tree) (trees.Tree, error)

func MyGrammarSemantics(rule string, node trees.Tree) (trees.Tree, error) {
	switch rule {
	//case "json_number":
	//	val, err := strconv.ParseFloat(node.(string), 64)
	//	if err != nil {
	//		// Tell the parser engine to fail this branch and backtrack
	//		return nil, fmt.Errorf("semantic failure: float overflow")
	//	}
	//	return val, nil
	//
	default:
		return node, nil
	}
}
