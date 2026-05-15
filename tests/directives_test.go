package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestGrammarDirective(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: MyGrammar
		start := 'test'
	`, nil)
	if g.Name != "MyGrammar" {
		t.Errorf("expected name 'MyGrammar', got %q", g.Name)
	}
	util.AssertJSONStr(t, g, "test", `"test"`)
}

func TestWhitespaceDirective(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /[\t ]+/
		@@grammar :: Test
		start := 'a' 'b'
	`, nil)
	util.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestWhitespaceNoneDirective(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: None
		@@nameguard :: False
		@@grammar :: Test
		start := 'a' 'b'
	`, nil)
	util.AssertJSONStr(t, g, "ab", `["a", "b"]`)
}

func TestDefaultWhitespace(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'a' 'b'
	`, nil)
	util.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestLeftRecursionDirective(t *testing.T) {
	g := util.Compile(t, `
		@@left_recursion :: False
		@@grammar :: Test
		start := 'test'
	`, nil)
	util.AssertJSONStr(t, g, "test", `"test"`)
}

func TestParseInfoDirective(t *testing.T) {
	g := util.Compile(t, `
		@@parseinfo :: True
		@@grammar :: Test
		start := 'test'
	`, nil)
	util.AssertJSONStr(t, g, "test", `"test"`)
}

func TestNameGuardDirective(t *testing.T) {
	g := util.Compile(t, `
		@@nameguard :: False
		@@grammar :: Test
		start := 'ab'
	`, nil)
	util.AssertJSONStr(t, g, "ab", `"ab"`)
}

func TestCommentsDirective(t *testing.T) {
	g := util.Compile(t, `
		@@comments :: /#[^\n]*/
		@@grammar :: Test
		start := 'a'
	`, nil)
	util.AssertJSONStr(t, g, "a", `"a"`)
}
