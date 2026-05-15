package test

import (
	"testing"
)

func TestBasicEOL(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'hello' $-> 'world'
	`, nil)
	AssertJSONStr(t, g, "hello\nworld", `["hello", "world"]`)
	AssertJSONStr(t, g, "hello  \nworld", `["hello", "world"]`)
	ParseFail(t, g, "hello world")
	ParseFail(t, g, "helloX\nworld")
}

func TestEOLAtEndOfText(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'hello' $-> $
	`, nil)
	AssertJSONStr(t, g, "hello\n", `"hello"`)
	AssertJSONStr(t, g, "hello  \n", `"hello"`)
	ParseFail(t, g, "hello world")
}

func TestMultipleEOLs(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'line1' $-> 'line2' $-> 'line3'
	`, nil)
	AssertJSONStr(t, g, "line1\nline2\nline3", `["line1", "line2", "line3"]`)
	AssertJSONStr(t, g, "line1  \nline2\n  line3", `["line1", "line2", "line3"]`)
}

func TestEOLWithIndentation(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'start' $-> 'indented' $-> 'end'
	`, nil)
	AssertJSONStr(t, g, "start\n  indented\nend", `["start", "indented", "end"]`)
}

func TestEOLInClosure(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := ('item' $->)* 'end'
	`, nil)
	AssertJSONStr(t, g, "item\nitem\nend", `[["item", "item"], "end"]`)
	AssertJSONStr(t, g, "item  \nitem\nend", `[["item", "item"], "end"]`)
	AssertJSONStr(t, g, "end", `[[], "end"]`)
}

func TestEOLWithMixedWhitespace(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'start' $-> 'next'
	`, nil)
	AssertJSONStr(t, g, "start \t \nnext", `["start", "next"]`)
	AssertJSONStr(t, g, "start   \nnext", `["start", "next"]`)
	AssertJSONStr(t, g, "start\t\nnext", `["start", "next"]`)
}

func TestEOLInTatsuEBNFEndrule(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'a' (';' | $->) 'b'
	`, nil)
	AssertJSONStr(t, g, "a\nb", `["a", "b"]`)
	AssertJSONStr(t, g, "a;b", `["a", ";", "b"]`)
}
