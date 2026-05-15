package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestNameInOption(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start = expr_range $
		expr_range =
			| [from: expr] '..' [to: expr]
			| expr
		expr = /[\d]+/
	`, nil)
	ogopego.AssertJSONStr(t, g, "1 .. 10", `{"from": "1", "to": "10"}`)
	ogopego.AssertJSONStr(t, g, "10", `"10"`)
}

func TestMixedReturn(t *testing.T) {
	t.Skip("optional named capture folding not yet implemented")
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := ('a' b='b') 'c' d='d'?
	`, nil)
	ogopego.AssertJSONStr(t, g, "a b c", `{"b": "b", "d": null}`)
}
