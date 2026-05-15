package test

import (
	"strings"
	"testing"
)

func TestPrettyGrammar(t *testing.T) {
	g := Compile(t, `
		@@grammar :: PrettyTest
		start := 'a'
		`,
		nil,
	)
	pretty := g.PrettyPrint()
	if !strings.Contains(pretty, "PrettyTest") {
		t.Errorf("expected pretty print to contain grammar name")
	}
	if !strings.Contains(pretty, "start") {
		t.Errorf("expected pretty print to contain rule name")
	}
}

func TestPrettySlashedPattern(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := ?"[a-z]+/[0-9]+" $
	`, nil)
	AssertJSONStr(t, g, "abc/123", `"abc/123"`)
}
