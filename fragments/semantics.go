// Tree represents the concrete polymorphic node (interface/any)
type Semantics func(rule string, node Tree) (Tree, error)

func MyGrammarSemantics(rule string, node Tree) (Tree, error) {
	switch rule {
	case "json_number":
		val, err := strconv.ParseFloat(node.(string), 64)
		if err != nil {
			// Tell the parser engine to fail this branch and backtrack
			return nil, fmt.Errorf("semantic failure: float overflow")
		}
		return val, nil

	default:
		return node, nil
	}
}
