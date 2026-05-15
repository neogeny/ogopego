package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestGrammarDirective(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: MyGrammar
		start := 'test' ;
	`)
	if g.Name != "MyGrammar" {
		t.Errorf("expected name 'MyGrammar', got %q", g.Name)
	}
	testutil.AssertJSONStr(t, g, "test", `"test"`)
}

func TestWhitespaceDirective(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /[\t ]+/
		@@grammar :: Test
		start := 'a' 'b' ;
	`)
	testutil.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestWhitespaceNoneDirective(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: None
		@@nameguard :: False
		@@grammar :: Test
		start := 'a' 'b' ;
	`)
	testutil.AssertJSONStr(t, g, "ab", `["a", "b"]`)
}

func TestDefaultWhitespace(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a' 'b' ;
	`)
	testutil.AssertJSONStr(t, g, "a b", `["a", "b"]`)
}

func TestLeftRecursionDirective(t *testing.T) {
	g := testutil.Compile(t, `
		@@left_recursion :: False
		@@grammar :: Test
		start := 'test' ;
	`)
	testutil.AssertJSONStr(t, g, "test", `"test"`)
}

func TestParseInfoDirective(t *testing.T) {
	g := testutil.Compile(t, `
		@@parseinfo :: True
		@@grammar :: Test
		start := 'test' ;
	`)
	testutil.AssertJSONStr(t, g, "test", `"test"`)
}

func TestNameGuardDirective(t *testing.T) {
	g := testutil.Compile(t, `
		@@nameguard :: False
		@@grammar :: Test
		start := 'ab' ;
	`)
	testutil.AssertJSONStr(t, g, "ab", `"ab"`)
}

func TestCommentsDirective(t *testing.T) {
	g := testutil.Compile(t, `
		@@comments :: /#[^\n]*/
		@@grammar :: Test
		start := 'a' ;
	`)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
}
