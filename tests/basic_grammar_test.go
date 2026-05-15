package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestSimpleGrammar(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'hello'
	`, nil)
	util.AssertJSONStr(t, g, "hello", `"hello"`)
}

func TestMultipleRules(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := choice
		choice := 'a' | 'b' | 'c'
	`, nil)
	util.AssertJSONStr(t, g, "a", `"a"`)
	util.AssertJSONStr(t, g, "b", `"b"`)
	util.AssertJSONStr(t, g, "c", `"c"`)
}

func TestRuleReferences(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'hello' 'world'
	`, nil)
	util.AssertJSONStr(t, g, "helloworld", `["hello", "world"]`)
}

func TestEmptyInput(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'test'?
	`, nil)
	util.AssertJSONStr(t, g, "", `null`)
}
