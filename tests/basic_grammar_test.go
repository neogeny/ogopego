package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestSimpleGrammar(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'hello' ;
	`)
	testutil.AssertJSONStr(t, g, "hello", `"hello"`)
}

func TestMultipleRules(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := choice ;
		choice := 'a' | 'b' | 'c' ;
	`)
	testutil.AssertJSONStr(t, g, "a", `"a"`)
	testutil.AssertJSONStr(t, g, "b", `"b"`)
	testutil.AssertJSONStr(t, g, "c", `"c"`)
}

func TestRuleReferences(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'hello' 'world' ;
	`)
	testutil.AssertJSONStr(t, g, "helloworld", `["hello", "world"]`)
}

func TestEmptyInput(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'test'? ;
	`)
	testutil.AssertJSONStr(t, g, "", `null`)
}
