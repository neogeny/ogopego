package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestTokenSequence(t *testing.T) {
	g := testutil.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'hello' 'world' ;`)
	testutil.AssertJSONStr(t, g, "hello world", `["hello", "world"]`)
}

func TestOptionalToken(t *testing.T) {
	g := testutil.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'a' 'b'? ;`)
	testutil.AssertJSONStr(t, g, "a b", `["a", "b"]`)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
}

func TestClosureTokens(t *testing.T) {
	g := testutil.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'a'* ;`)
	testutil.AssertJSONStr(t, g, "aaa", `["a", "a", "a"]`)
}

func TestPositiveClosure(t *testing.T) {
	g := testutil.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'a'+ ;`)
	testutil.AssertJSONStr(t, g, "aaa", `["a", "a", "a"]`)
}

func TestChoiceAlternatives(t *testing.T) {
	g := testutil.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'a' | 'b' | 'c' ;`)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
	testutil.AssertJSONStr(t, g, "b", `"b"`)
	testutil.AssertJSONStr(t, g, "c", `"c"`)
}

func TestMultiLineGrammar(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test

		start := 'hello' 'world' ;
	`)
	testutil.AssertJSONStr(t, g, "hello world", `["hello", "world"]`)
}
