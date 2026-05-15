package test

import (
	"testing"
)

func TestMultilineString(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Test
		start := longone | shortone $
		shortone := "short"
		longone := """
			this "text"
			is a long "string"
			"""
	`), nil)
	AssertJSONStr(t, g, "short", `"short"`)
}
