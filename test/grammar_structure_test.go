package ogopego_test

import (
	"strings"
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestGrammarHasRules(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'a' | 'b' | 'c'
	`, nil)
	if len(g.Rules) < 1 {
		t.Fatal("expected at least 1 rule")
	}
}

func TestFirstRuleIsDefault(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'a'
	`, nil)
	ogopego.AssertJSONStr(t, g, "a", `"a"`)
}

func TestPrettyPrint(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'a'
	`, nil)
	s := g.PrettyPrint()
	if !strings.Contains(s, "Test") && !strings.Contains(s, "start") {
		t.Errorf("expected pretty print to contain 'Test' or 'start', got %q", s)
	}
}
