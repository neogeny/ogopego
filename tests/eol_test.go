package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestBasicEOL(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'hello' $-> 'world'
	`, nil)
	util.AssertJSONStr(t, g, "hello\nworld", `["hello", "world"]`)
	util.AssertJSONStr(t, g, "hello  \nworld", `["hello", "world"]`)
	util.ParseFail(t, g, "hello world")
	util.ParseFail(t, g, "helloX\nworld")
}

func TestEOLAtEndOfText(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'hello' $-> $
	`, nil)
	util.AssertJSONStr(t, g, "hello\n", `"hello"`)
	util.AssertJSONStr(t, g, "hello  \n", `"hello"`)
	util.ParseFail(t, g, "hello world")
}

func TestMultipleEOLs(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'line1' $-> 'line2' $-> 'line3'
	`, nil)
	util.AssertJSONStr(t, g, "line1\nline2\nline3", `["line1", "line2", "line3"]`)
	util.AssertJSONStr(t, g, "line1  \nline2\n  line3", `["line1", "line2", "line3"]`)
}

func TestEOLWithIndentation(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'start' $-> 'indented' $-> 'end'
	`, nil)
	util.AssertJSONStr(t, g, "start\n  indented\nend", `["start", "indented", "end"]`)
}

func TestEOLInClosure(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := ('item' $->)* 'end'
	`, nil)
	util.AssertJSONStr(t, g, "item\nitem\nend", `[["item", "item"], "end"]`)
	util.AssertJSONStr(t, g, "item  \nitem\nend", `[["item", "item"], "end"]`)
	util.AssertJSONStr(t, g, "end", `[[], "end"]`)
}

func TestEOLWithMixedWhitespace(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'start' $-> 'next'
	`, nil)
	util.AssertJSONStr(t, g, "start \t \nnext", `["start", "next"]`)
	util.AssertJSONStr(t, g, "start   \nnext", `["start", "next"]`)
	util.AssertJSONStr(t, g, "start\t\nnext", `["start", "next"]`)
}

func TestEOLInTatsuEBNFEndrule(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'a' (';' | $->) 'b'
	`, nil)
	util.AssertJSONStr(t, g, "a\nb", `["a", "b"]`)
	util.AssertJSONStr(t, g, "a;b", `["a", ";", "b"]`)
}
