package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestSimplePattern(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := /\d+/
	`, nil)
	util.AssertJSONStr(t, g, "123", `"123"`)
}

func TestPatternWithLetters(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := /[a-z]+/
	`, nil)
	util.AssertJSONStr(t, g, "hello", `"hello"`)
}

func TestPatternWithAnchors(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := /^start/
	`, nil)
	util.AssertJSONStr(t, g, "start", `"start"`)
}

func TestPatternCaseInsensitive(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := /(?i)hello/
	`, nil)
	util.AssertJSONStr(t, g, "HELLO", `"HELLO"`)
}

func TestPatternCharacterClasses(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := /[A-Za-z_]\w*/
	`, nil)
	util.AssertJSONStr(t, g, "hello_world", `"hello_world"`)
}
