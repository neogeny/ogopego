package test

import (
	"testing"
)

func TestGroup(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := ('a' 'b')*
	`, nil)
	AssertJSONStr(t, g, "abab", `[["a", "b"], ["a", "b"]]`)
}

func TestSkipGroup(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := (?: 'a' 'b')*
	`, nil)
	AssertJSONStr(t, g, "abab", `[null, null]`)
}

func TestVoid(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a' () 'b'
	`, nil)
	AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestEOF(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'a' $
	`, nil)
	AssertJSONStr(t, g, "a", `"a"`)
}

func TestDot(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := /./ 'b'
	`, nil)
	AssertJSONStr(t, g, "ab", `["a", "b"]`)
}

func TestConstant(t *testing.T) {
	const grammar = "@@grammar :: Test\nstart := `constant` ;\n"
	g := Compile(t, grammar, nil)
	AssertJSONStr(t, g, "", `"constant"`)
}
