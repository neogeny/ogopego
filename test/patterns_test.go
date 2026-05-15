package test

import (
	"testing"
)

func TestSimplePattern(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := /\d+/
	`, nil)
	AssertJSONStr(t, g, "123", `"123"`)
}

func TestPatternWithLetters(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := /[a-z]+/
	`, nil)
	AssertJSONStr(t, g, "hello", `"hello"`)
}

func TestPatternWithAnchors(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := /^start/
	`, nil)
	AssertJSONStr(t, g, "start", `"start"`)
}

func TestPatternCaseInsensitive(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := /(?i)hello/
	`, nil)
	AssertJSONStr(t, g, "HELLO", `"HELLO"`)
}

func TestPatternCharacterClasses(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := /[A-Za-z_]\w*/
	`, nil)
	AssertJSONStr(t, g, "hello_world", `"hello_world"`)
}
