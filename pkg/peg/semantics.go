package peg

import "github.com/neogeny/ogopego/pkg/trees"

// EBNFGrammarSemantics are the library's semantics for parsing grammars
// in the syntax defined by tatsu.ebnf and tatsu.json.
type EBNFGrammarSemantics struct{}

func (EBNFGrammarSemantics) Apply(
	node any,
	ruleName string,
	params []string) (any, bool) {
	switch ruleName {
	case "true":
		return true, true
	case "false":
		return false, true
	case "null":
		return nil, true
	case "meta":
		text := textValue(node)
		switch text {
		case "name":
			return &trees.Node{TypeName: "NameMeta"}, true
		case "int":
			return &trees.Node{TypeName: "IntMeta"}, true
		case "uint":
			return &trees.Node{TypeName: "UIntMeta"}, true
		case "float":
			return &trees.Node{TypeName: "FloatMeta"}, true
		case "bool":
			return &trees.Node{TypeName: "BoolMeta"}, true
		}
	}
	// do nothing by default
	// return false to tell caller to apply default processing
	return node, false
}
