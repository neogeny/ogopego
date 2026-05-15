package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestGroup(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := ('a' 'b')*
	`, nil)
	util.AssertJSONStr(t, g, "abab", `[["a", "b"], ["a", "b"]]`)
}

func TestSkipGroup(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := (?: 'a' 'b')*
	`, nil)
	util.AssertJSONStr(t, g, "abab", `[null, null]`)
}

func TestVoid(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a' () 'b'
	`, nil)
	util.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestEOF(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'a' $
	`, nil)
	util.AssertJSONStr(t, g, "a", `"a"`)
}

func TestDot(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := /./ 'b'
	`, nil)
	util.AssertJSONStr(t, g, "ab", `["a", "b"]`)
}

func TestConstant(t *testing.T) {
	t.Skip("constant not yet implemented")
	const grammar = "@@grammar :: Test\nstart := `constant` ;\n"
	g := util.Compile(t, grammar, nil)
	util.AssertJSONStr(t, g, "", `"constant"`)
}
