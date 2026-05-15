package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestTokenSequence(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		start: 'hello' 'world'
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello world", `["hello", "world"]`)
}

func TestOptionalToken(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		start := 'a' 'b'?
	`, nil)
	ogopego.AssertJSONStr(t, g, "a b", `["a", "b"]`)
	ogopego.AssertJSONStr(t, g, "a", `"a"`)
}

func TestClosureTokens(t *testing.T) {
	g := ogopego.Compile(t, `
		start := 'a'*
	`, nil)
	ogopego.AssertJSONStr(t, g, "aaa", `["a", "a", "a"]`)
}

func TestPositiveClosure(t *testing.T) {
	g := ogopego.Compile(t, `
		start := 'a'+
	`, nil)
	ogopego.AssertJSONStr(t, g, "aaa", `["a", "a", "a"]`)
}

func TestChoiceAlternatives(t *testing.T) {
	g := ogopego.Compile(t, `
		start := 'a' | 'b' | 'c'
	`, nil)
	ogopego.AssertJSONStr(t, g, "a", `"a"`)
	ogopego.AssertJSONStr(t, g, "b", `"b"`)
	ogopego.AssertJSONStr(t, g, "c", `"c"`)
}

func TestMultiLineGrammar(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'hello' 'world'
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello world", `["hello", "world"]`)
}
