package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestPositiveLookahead(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := &'a' 'a'
	`, nil)
	ogopego.AssertJSONStr(t, g, "a", `"a"`)
}

func TestNegativeLookahead(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := !'b' 'a'
	`, nil)
	ogopego.AssertJSONStr(t, g, "a", `"a"`)
}

func TestCut(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a'~'b'
	`, nil)
	ogopego.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestPatternsWithNewlines(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /[ \t]/
		@@grammar :: Test
		start = blanklines $
		blanklines = blankline blanklines?
		blankline = /(?m)^[^\n]*\n/
	`, nil)
	ogopego.AssertJSONStr(t, g, "\n\n", `["\n", "\n"]`)
}
