package test

import (
	"testing"
)

func TestPositiveLookahead(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := &'a' 'a'
	`, nil)
	AssertJSONStr(t, g, "a", `"a"`)
}

func TestNegativeLookahead(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := !'b' 'a'
	`, nil)
	AssertJSONStr(t, g, "a", `"a"`)
}

func TestCut(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a'~'b'
	`, nil)
	AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestPatternsWithNewlines(t *testing.T) {
	g := Compile(t, Dedent(`
		@@whitespace :: /[ \t]/
		@@grammar :: Test
		start = blanklines $
		blanklines = blankline blanklines?
		blankline = /(?m)^[^\n]*\n/
	`), nil)
	AssertJSONStr(t, g, "\n\n", `["\n", "\n"]`)
}
