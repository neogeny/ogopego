package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestGroup(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := ('a' 'b')* ;
	`)
	testutil.AssertJSONStr(t, g, "abab", `[["a", "b"], ["a", "b"]]`)
}

func TestSkipGroup(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := (?: 'a' 'b')* ;
	`)
	testutil.AssertJSONStr(t, g, "abab", `[null, null]`)
}

func TestVoid(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a' () 'b' ;
	`)
	testutil.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestEOF(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a' $ ;
	`)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
}

func TestDot(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := /./ 'b' ;
	`)
	testutil.AssertJSONStr(t, g, "ab", `["a", "b"]`)
}

func TestConstant(t *testing.T) {
	t.Skip("constant not yet implemented")
	const grammar = "@@grammar :: Test\nstart := `constant` ;\n"
	g := testutil.Compile(t, grammar)
	testutil.AssertJSONStr(t, g, "", `"constant"`)
}
