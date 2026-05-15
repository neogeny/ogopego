package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestSimpleGrammar(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start: 'hello'
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello", `"hello"`)
}

func TestMultipleRules(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start: choice

		choice:
			| 'a'
			| 'b'
			| 'c'
	`, nil)
	ogopego.AssertJSONStr(t, g, "a", `"a"`)
	ogopego.AssertJSONStr(t, g, "b", `"b"`)
	ogopego.AssertJSONStr(t, g, "c", `"c"`)
}

func TestRuleReferences(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start: 'hello' 'world'
	`, nil)
	ogopego.AssertJSONStr(t, g, "helloworld", `["hello", "world"]`)
}

func TestEmptyInput(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start: 'test'?
	`, nil)
	ogopego.AssertJSONStr(t, g, "", `null`)
}
