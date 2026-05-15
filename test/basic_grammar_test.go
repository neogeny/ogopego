package test

import (
	"testing"
)

func TestSimpleGrammar(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start: 'hello'
	`, nil)
	AssertJSONStr(t, g, "hello", `"hello"`)
}

func TestMultipleRules(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start: choice

		choice:
			| 'a'
			| 'b'
			| 'c'
	`, nil)
	AssertJSONStr(t, g, "a", `"a"`)
	AssertJSONStr(t, g, "b", `"b"`)
	AssertJSONStr(t, g, "c", `"c"`)
}

func TestRuleReferences(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start: 'hello' 'world'
	`, nil)
	AssertJSONStr(t, g, "helloworld", `["hello", "world"]`)
}

func TestEmptyInput(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start: 'test'?
	`, nil)
	AssertJSONStr(t, g, "", `null`)
}
