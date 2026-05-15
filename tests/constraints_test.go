package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestPositiveLookahead(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := &'a' 'a' ;
	`)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
}

func TestNegativeLookahead(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := !'b' 'a' ;
	`)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
}

func TestCut(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'a'~'b' ;
	`)
	testutil.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestPatternsWithNewlines(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /[ \t]/
		@@grammar :: Test
		start = blanklines $ ;
		blanklines = blankline blanklines? ;
		blankline = /(?m)^[^\n]*\n/ ;
	`)
	testutil.AssertJSONStr(t, g, "\n\n", `["\n", "\n"]`)
}
