package test

import (
	"testing"
)

func TestTokenSequence(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /\s+/
		start: 'hello' 'world'
	`, nil)
	AssertJSONStr(t, g, "hello world", `["hello", "world"]`)
}

func TestOptionalToken(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /\s+/
		start := 'a' 'b'?
	`, nil)
	AssertJSONStr(t, g, "a b", `["a", "b"]`)
	AssertJSONStr(t, g, "a", `"a"`)
}

func TestClosureTokens(t *testing.T) {
	g := Compile(t, `
		start := 'a'*
	`, nil)
	AssertJSONStr(t, g, "aaa", `["a", "a", "a"]`)
}

func TestPositiveClosure(t *testing.T) {
	g := Compile(t, `
		start := 'a'+
	`, nil)
	AssertJSONStr(t, g, "aaa", `["a", "a", "a"]`)
}

func TestChoiceAlternatives(t *testing.T) {
	g := Compile(t, `
		start := 'a' | 'b' | 'c'
	`, nil)
	AssertJSONStr(t, g, "a", `"a"`)
	AssertJSONStr(t, g, "b", `"b"`)
	AssertJSONStr(t, g, "c", `"c"`)
}

func TestMultiLineGrammar(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'hello' 'world'
	`, nil)
	AssertJSONStr(t, g, "hello world", `["hello", "world"]`)
}
