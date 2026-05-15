package test

import (
	"testing"
)

func TestNameInOption(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Test
		start = expr_range $
		expr_range =
			| [from: expr] '..' [to: expr]
			| expr
		expr = /[\d]+/
	`), nil)
	AssertJSONStr(t, g, "1 .. 10", `{"from": "1", "to": "10"}`)
	AssertJSONStr(t, g, "10", `"10"`)
}

func TestMixedReturn(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := ('a' b='b') 'c' d='d'?
	`, nil)
	AssertJSONStr(t, g, "a b c", `{"b": "b", "d": null}`)
}
