package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestEscapeSequences(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = 'hello\nworld' $
	`, nil)
	util.AssertJSONStr(t, g, "hello\\nworld", `"hello\\nworld"`)
}

func TestStart(t *testing.T) {
	t.Skip("constant evaluation (backtick syntax) not yet implemented")
	const grammar = "@@grammar :: Test\n" +
		"true = 'test' @:`True` $ ;\n"
	g := util.Compile(t, grammar, nil)
	util.AssertJSONStr(t, g, "test", `"True"`)
}

func TestSkipWhitespace(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		statement = 'FOO' subject $
		subject = name:id
		id = /[a-z]+/
	`, nil)
	util.ParseJSON(t, g, "FOO something")
	util.ParseFail(t, g, "something")
	util.ParseFail(t, g, "FOO")
}

func TestNodeParseInfo(t *testing.T) {
	t.Skip("parseinfo directive not yet implemented")
	g := util.Compile(t, `
		@@grammar :: Test
		@@parseinfo :: True
		start = 'test' $
	`, nil)
	util.AssertJSONStr(t, g, "test", `"test"`)
}

func TestCutScope(t *testing.T) {
	g := util.Compile(t, `
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
	util.AssertJSONStr(t, g, "something", `"something"`)
}
