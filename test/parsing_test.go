package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestEscapeSequences(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = 'hello\nworld' $
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello\\nworld", `"hello\\nworld"`)
}

func TestStart(t *testing.T) {
	t.Skip("constant evaluation (backtick syntax) not yet implemented")
	const grammar = "@@grammar :: Test\n" +
		"true = 'test' @:`True` $ ;\n"
	g := ogopego.Compile(t, grammar, nil)
	ogopego.AssertJSONStr(t, g, "test", `"True"`)
}

func TestSkipWhitespace(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		statement = 'FOO' subject $
		subject = name:id
		id = /[a-z]+/
	`, nil)
	ogopego.ParseJSON(t, g, "FOO something")
	ogopego.ParseFail(t, g, "something")
	ogopego.ParseFail(t, g, "FOO")
}

func TestNodeParseInfo(t *testing.T) {
	t.Skip("parseinfo directive not yet implemented")
	g := ogopego.Compile(t, `
		@@grammar :: Test
		@@parseinfo :: True
		start = 'test' $
	`, nil)
	ogopego.AssertJSONStr(t, g, "test", `"test"`)
}

func TestCutScope(t *testing.T) {
	g := ogopego.Compile(t, `
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
	ogopego.AssertJSONStr(t, g, "something", `"something"`)
}
