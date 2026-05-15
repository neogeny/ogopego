package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestBasicEOL(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'hello' $-> 'world'
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello\nworld", `["hello", "world"]`)
	ogopego.AssertJSONStr(t, g, "hello  \nworld", `["hello", "world"]`)
	ogopego.ParseFail(t, g, "hello world")
	ogopego.ParseFail(t, g, "helloX\nworld")
}

func TestEOLAtEndOfText(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'hello' $-> $
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello\n", `"hello"`)
	ogopego.AssertJSONStr(t, g, "hello  \n", `"hello"`)
	ogopego.ParseFail(t, g, "hello world")
}

func TestMultipleEOLs(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'line1' $-> 'line2' $-> 'line3'
	`, nil)
	ogopego.AssertJSONStr(t, g, "line1\nline2\nline3", `["line1", "line2", "line3"]`)
	ogopego.AssertJSONStr(t, g, "line1  \nline2\n  line3", `["line1", "line2", "line3"]`)
}

func TestEOLWithIndentation(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'start' $-> 'indented' $-> 'end'
	`, nil)
	ogopego.AssertJSONStr(t, g, "start\n  indented\nend", `["start", "indented", "end"]`)
}

func TestEOLInClosure(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := ('item' $->)* 'end'
	`, nil)
	ogopego.AssertJSONStr(t, g, "item\nitem\nend", `[["item", "item"], "end"]`)
	ogopego.AssertJSONStr(t, g, "item  \nitem\nend", `[["item", "item"], "end"]`)
	ogopego.AssertJSONStr(t, g, "end", `[[], "end"]`)
}

func TestEOLWithMixedWhitespace(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'start' $-> 'next'
	`, nil)
	ogopego.AssertJSONStr(t, g, "start \t \nnext", `["start", "next"]`)
	ogopego.AssertJSONStr(t, g, "start   \nnext", `["start", "next"]`)
	ogopego.AssertJSONStr(t, g, "start\t\nnext", `["start", "next"]`)
}

func TestEOLInTatsuEBNFEndrule(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'a' (';' | $->) 'b'
	`, nil)
	ogopego.AssertJSONStr(t, g, "a\nb", `["a", "b"]`)
	ogopego.AssertJSONStr(t, g, "a;b", `["a", ";", "b"]`)
}
