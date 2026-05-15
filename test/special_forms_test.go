package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestGroup(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := ('a' 'b')*
	`, nil)
	ogopego.AssertJSONStr(t, g, "abab", `[["a", "b"], ["a", "b"]]`)
}

func TestSkipGroup(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := (?: 'a' 'b')*
	`, nil)
	ogopego.AssertJSONStr(t, g, "abab", `[null, null]`)
}

func TestVoid(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a' () 'b'
	`, nil)
	ogopego.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestEOF(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'a' $
	`, nil)
	ogopego.AssertJSONStr(t, g, "a", `"a"`)
}

func TestDot(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := /./ 'b'
	`, nil)
	ogopego.AssertJSONStr(t, g, "ab", `["a", "b"]`)
}

func TestConstant(t *testing.T) {
	const grammar = "@@grammar :: Test\nstart := `constant` ;\n"
	g := ogopego.Compile(t, grammar, nil)
	ogopego.AssertJSONStr(t, g, "", `"constant"`)
}
