package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestBasicEOL(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'hello' $-> 'world' ;
	`)
	testutil.AssertJSONStr(t, g, "hello\nworld", `["hello", "world"]`)
	testutil.AssertJSONStr(t, g, "hello  \nworld", `["hello", "world"]`)
	testutil.ParseFail(t, g, "hello world")
	testutil.ParseFail(t, g, "helloX\nworld")
}

func TestEOLAtEndOfText(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'hello' $-> $ ;
	`)
	testutil.AssertJSONStr(t, g, "hello\n", `"hello"`)
	testutil.AssertJSONStr(t, g, "hello  \n", `"hello"`)
	testutil.ParseFail(t, g, "hello world")
}

func TestMultipleEOLs(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'line1' $-> 'line2' $-> 'line3' ;
	`)
	testutil.AssertJSONStr(t, g, "line1\nline2\nline3", `["line1", "line2", "line3"]`)
	testutil.AssertJSONStr(t, g, "line1  \nline2\n  line3", `["line1", "line2", "line3"]`)
}

func TestEOLWithIndentation(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'start' $-> 'indented' $-> 'end' ;
	`)
	testutil.AssertJSONStr(t, g, "start\n  indented\nend", `["start", "indented", "end"]`)
}

func TestEOLInClosure(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := ('item' $->)* 'end' ;
	`)
	testutil.AssertJSONStr(t, g, "item\nitem\nend", `[["item", "item"], "end"]`)
	testutil.AssertJSONStr(t, g, "item  \nitem\nend", `[["item", "item"], "end"]`)
	testutil.AssertJSONStr(t, g, "end", `[[], "end"]`)
}

func TestEOLWithMixedWhitespace(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'start' $-> 'next' ;
	`)
	testutil.AssertJSONStr(t, g, "start \t \nnext", `["start", "next"]`)
	testutil.AssertJSONStr(t, g, "start   \nnext", `["start", "next"]`)
	testutil.AssertJSONStr(t, g, "start\t\nnext", `["start", "next"]`)
}

func TestEOLInTatsuEBNFEndrule(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a' (';' | $->) 'b' ;
	`)
	testutil.AssertJSONStr(t, g, "a\nb", `["a", "b"]`)
	testutil.AssertJSONStr(t, g, "a;b", `["a", ";", "b"]`)
}
