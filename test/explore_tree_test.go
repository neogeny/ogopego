package test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/asjson"
)

func printTreeRec(t *testing.T, indent string, node any) {
	t.Helper()
	switch n := node.(type) {
	case string:
		t.Logf("%s%q", indent, n)
	case float64:
		t.Logf("%s%v", indent, n)
	case bool:
		t.Logf("%s%v", indent, n)
	case nil:
		t.Logf("%snil", indent)
	case []any:
		t.Logf("%s[", indent)
		for _, item := range n {
			printTreeRec(t, indent+"  ", item)
		}
		t.Logf("%s]", indent)
	case map[string]any:
		t.Logf("%s{", indent)
		for k, v := range n {
			t.Logf("%s  %s:", indent, k)
			printTreeRec(t, indent+"    ", v)
		}
		t.Logf("%s}", indent)
	default:
		t.Logf("%s%T %+v", indent, n, n)
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
			j := asjson.AsJSON(tree)
			t.Logf("input: %s", input)
			t.Logf("JSON:  %s", asjson.AsJSONStr(tree))
			t.Logf("type:  %T", j)
			printTreeRec(t, "", j)
		})
	}
}
