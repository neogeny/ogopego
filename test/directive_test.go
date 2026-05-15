package test

import (
	"testing"
)

func TestWhitespaceDirectiveDoubleQuote(t *testing.T) {
	g := Compile(t, `
		@@whitespace :: /[\t ]+/
		@@grammar :: Test
		test := "test" $
	`, nil)
	AssertJSONStr(t, g, "test", `"test"`)
}
