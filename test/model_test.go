package test

import (
	"testing"
)

func TestChildren(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Calc
		start = expression $
		expression = term
		term = 'x'
	`), nil)
	AssertJSONStr(t, g, "x", `"x"`)
}

func TestNodeKWArgs(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start = 'value'
	`, nil)
	AssertJSONStr(t, g, "value", `"value"`)
}
