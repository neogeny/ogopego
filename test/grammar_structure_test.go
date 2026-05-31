package test

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestGrammarHasRules(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'a' | 'b' | 'c'
	`, nil)
	assert.NotZero(t, len(g.Rules), "expected at least 1 rule")
}

func TestFirstRuleIsDefault(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'a'
	`, nil)
	AssertJSONStr(t, g, "a", `"a"`)
}

func TestPrettyPrint(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'a'
	`, nil)
	s := g.PrettyPrint()
	assert.True(t, strings.Contains(s, "Test") || strings.Contains(s, "start"),
		"expected pretty print to contain 'Test' or 'start', got %q", s)
}
