package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestSimplePattern(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := /\d+/
	`, nil)
	testutil.AssertJSONStr(t, g, "123", `"123"`)
}

func TestPatternWithLetters(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := /[a-z]+/
	`, nil)
	testutil.AssertJSONStr(t, g, "hello", `"hello"`)
}

func TestPatternWithAnchors(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := /^start/
	`, nil)
	testutil.AssertJSONStr(t, g, "start", `"start"`)
}

func TestPatternCaseInsensitive(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := /(?i)hello/
	`, nil)
	testutil.AssertJSONStr(t, g, "HELLO", `"HELLO"`)
}

func TestPatternCharacterClasses(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := /[A-Za-z_]\w*/
	`, nil)
	testutil.AssertJSONStr(t, g, "hello_world", `"hello_world"`)
}
