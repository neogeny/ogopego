package test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/trees"
)

func printTreeRec(t *testing.T, indent string, node trees.Tree) {
	t.Helper()
	switch n := node.(type) {
	case *trees.Text:
		t.Logf("%sText(%q)", indent, n.Value)
	case *trees.Number:
		t.Logf("%sNumber(%v)", indent, n.Value)
	case *trees.Bool:
		t.Logf("%sBool(%v)", indent, n.Value)
	case *trees.Nil:
		t.Logf("%sNil", indent)
	case *trees.Seq:
		t.Logf("%sSeq [", indent)
		for _, item := range n.Items {
			printTreeRec(t, indent+"  ", item)
		}
		t.Logf("%s]", indent)
	case *trees.Array:
		t.Logf("%sList [", indent)
		for _, item := range n.Items {
			printTreeRec(t, indent+"  ", item)
		}
		t.Logf("%s]", indent)
	case *trees.MapNode:
		t.Logf("%sMapNode {", indent)
		for k, v := range n.Entries {
			t.Logf("%s  %s:", indent, k)
			printTreeRec(t, indent+"    ", v)
		}
		t.Logf("%s}", indent)
	case *trees.Named:
		t.Logf("%sNamed(%s):", indent, n.Name)
		printTreeRec(t, indent+"  ", n.Value)
	case *trees.Node:
		t.Logf("%sNode(%s):", indent, n.TypeName)
		printTreeRec(t, indent+"  ", n.Tree)
	case *trees.Override:
		t.Logf("%sOverride:", indent)
		printTreeRec(t, indent+"  ", n.Value)
	case *trees.OverrideAsList:
		t.Logf("%sOverrideAsList:", indent)
		printTreeRec(t, indent+"  ", n.Value)
	default:
		t.Logf("%s%T", indent, n)
	}
}

func TestExploreCalcTree(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Calc
		@@left_recursion :: True
		@@whitespace :: /\s+/

		start := expression $

		expression
			:= expression '+' term
			| expression '-' term
			| term

		term
			:= term '*' factor
			| term '/' factor
			| factor

		factor
			:= NUMBER
			| '(' expression ')'

		NUMBER := /\d+/
	`, nil)
	for _, input := range []string{"42", "1 + 2", "1 + 2 * 3", "(1 + 2) * 3"} {
		t.Run(input, func(t *testing.T) {
			tree, err := api.ParseInput(g, input, nil)
			assert.NoError(t, err, "parse")
			t.Logf("input: %s", input)
			t.Logf("JSON:  %s", trees.TreeToJSONStr(tree))
			t.Logf("type:  %T", tree)
			printTreeRec(t, "", tree)
		})
	}
}
