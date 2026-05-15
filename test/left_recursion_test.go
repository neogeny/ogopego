package test

import (
	"testing"
)

func TestDirectLeftRecursion(t *testing.T) {
	g := Compile(t, Dedent(`
		@@left_recursion :: True
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = expression $
		expression = expression '+' factor | expression '-' factor | factor
		factor = number
		number = /[0-9]+/
	`), nil)
	AssertJSONStr(t, g, "10 - 20", `["10", "-", "20"]`)
}
