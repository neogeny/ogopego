package test

import (
	"testing"
)

func TestKeywordsInRuleNames(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Test
		start = whitespace
		whitespace = {'x'}+
	`), nil)
	AssertJSONStr(t, g, "x", `["x"]`)
}
