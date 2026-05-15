package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/testutil"
)

func TestPositiveLookahead(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := &'a' 'a'
	`, nil)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
}

func TestNegativeLookahead(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := !'b' 'a'
	`, nil)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
}

func TestCut(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a'~'b'
	`, nil)
	testutil.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestPatternsWithNewlines(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /[ \t]/
		@@grammar :: Test
		start = blanklines $
		blanklines = blankline blanklines?
		blankline = /(?m)^[^\n]*\n/
	`, nil)
	testutil.AssertJSONStr(t, g, "\n\n", `["\n", "\n"]`)
}
