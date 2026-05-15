package test

import (
	"testing"
)

func TestMultiLineInput(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'hello' 'world'
	`, nil)
	AssertJSONStr(t, g, "hello\nworld", `["hello", "world"]`)
}
