package ogopego_test

import (
	"strings"
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestPrettyGrammar(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: PrettyTest
		start := 'a' ;
	`)
	pretty := g.PrettyPrint()
	if !strings.Contains(pretty, "PrettyTest") {
		t.Errorf("expected pretty print to contain grammar name")
	}
	if !strings.Contains(pretty, "start") {
		t.Errorf("expected pretty print to contain rule name")
	}
}

func TestPrettySlashedPattern(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := ?"[a-z]+/[0-9]+" $ ;
	`)
	testutil.AssertJSONStr(t, g, "abc/123", `"abc/123"`)
}
