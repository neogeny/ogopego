package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestNameInOption(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start = expr_range $
		expr_range =
			| [from: expr] '..' [to: expr]
			| expr
		expr = /[\d]+/
	`, nil)
	util.AssertJSONStr(t, g, "1 .. 10", `{"from": "1", "to": "10"}`)
	util.AssertJSONStr(t, g, "10", `"10"`)
}

func TestMixedReturn(t *testing.T) {
	t.Skip("optional named capture folding not yet implemented")
	g := util.Compile(t, `
		@@grammar :: Test
		start := ('a' b='b') 'c' d='d'?
	`, nil)
	util.AssertJSONStr(t, g, "a b c", `{"b": "b", "d": null}`)
}
