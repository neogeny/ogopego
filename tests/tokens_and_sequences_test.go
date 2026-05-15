package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestTokenSequence(t *testing.T) {
	g := util.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'hello' 'world' ;`, nil)
	util.AssertJSONStr(t, g, "hello world", `["hello", "world"]`)
}

func TestOptionalToken(t *testing.T) {
	g := util.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'a' 'b'? ;`, nil)
	util.AssertJSONStr(t, g, "a b", `["a", "b"]`)
	util.AssertJSONStr(t, g, "a", `"a"`)
}

func TestClosureTokens(t *testing.T) {
	g := util.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'a'* ;`, nil)
	util.AssertJSONStr(t, g, "aaa", `["a", "a", "a"]`)
}

func TestPositiveClosure(t *testing.T) {
	g := util.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'a'+ ;`, nil)
	util.AssertJSONStr(t, g, "aaa", `["a", "a", "a"]`)
}

func TestChoiceAlternatives(t *testing.T) {
	g := util.Compile(t, `@@whitespace :: /\s+/ @@grammar :: Test start := 'a' | 'b' | 'c' ;`, nil)
	util.AssertJSONStr(t, g, "a", `"a"`)
	util.AssertJSONStr(t, g, "b", `"b"`)
	util.AssertJSONStr(t, g, "c", `"c"`)
}

func TestMultiLineGrammar(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test

		start := 'hello' 'world'
	`, nil)
	util.AssertJSONStr(t, g, "hello world", `["hello", "world"]`)
}
