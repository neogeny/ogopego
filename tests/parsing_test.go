package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestEscapeSequences(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = 'hello\nworld' $
	`, nil)
	testutil.AssertJSONStr(t, g, "hello\\nworld", `"hello\\nworld"`)
}

func TestStart(t *testing.T) {
	t.Skip("constant evaluation (backtick syntax) not yet implemented")
	const grammar = "@@grammar :: Test\n" +
		"true = 'test' @:`True` $ ;\n"
	g := testutil.Compile(t, grammar, nil)
	testutil.AssertJSONStr(t, g, "test", `"True"`)
}

func TestSkipWhitespace(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		statement = 'FOO' subject $
		subject = name:id
		id = /[a-z]+/
	`, nil)
	testutil.ParseJSON(t, g, "FOO something")
	testutil.ParseFail(t, g, "something")
	testutil.ParseFail(t, g, "FOO")
}

func TestNodeParseInfo(t *testing.T) {
	t.Skip("parseinfo directive not yet implemented")
	g := testutil.Compile(t, `
		@@grammar :: Test
		@@parseinfo :: True
		start = 'test' $
	`, nil)
	testutil.AssertJSONStr(t, g, "test", `"test"`)
}

func TestCutScope(t *testing.T) {
	g := testutil.Compile(t, `
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
	testutil.AssertJSONStr(t, g, "something", `"something"`)
}
