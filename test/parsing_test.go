package test

import (
	"testing"
)

func TestEscapeSequences(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = 'hello\nworld' $
	`, nil)
	AssertJSONStr(t, g, "hello\\nworld", `"hello\\nworld"`)
}

func TestStart(t *testing.T) {
	const grammar = "@@grammar :: Test\n" +
		"true = 'test' @:`True` $ ;\n"
	g := Compile(t, grammar, nil)
	AssertJSONStr(t, g, "test", `"True"`)
}

func TestSkipWhitespace(t *testing.T) {
	g := Compile(t, Dedent(`
		@@whitespace :: /\s+/
		@@grammar :: Test
		statement = 'FOO' subject $
		subject = name:id
		id = /[a-z]+/
	`), nil)
	ParseJSON(t, g, "FOO something")
	ParseFail(t, g, "something")
	ParseFail(t, g, "FOO")
}

func TestNodeParseInfo(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		@@parseinfo :: True
		start = 'test' $
	`, nil)
	AssertJSONStr(t, g, "test", `"test"`)
}

func TestCutScope(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start =
			| one
			| two
			;

		one =
			| ~ !()
			| 'abc'

		two = 'something'
	`, nil)
	AssertJSONStr(t, g, "something", `"something"`)
}
