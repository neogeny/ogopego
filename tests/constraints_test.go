package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestPositiveLookahead(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := &'a' 'a'
	`, nil)
	util.AssertJSONStr(t, g, "a", `"a"`)
}

func TestNegativeLookahead(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := !'b' 'a'
	`, nil)
	util.AssertJSONStr(t, g, "a", `"a"`)
}

func TestCut(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a'~'b'
	`, nil)
	util.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestPatternsWithNewlines(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /[ \t]/
		@@grammar :: Test
		start = blanklines $
		blanklines = blankline blanklines?
		blankline = /(?m)^[^\n]*\n/
	`, nil)
	util.AssertJSONStr(t, g, "\n\n", `["\n", "\n"]`)
}
