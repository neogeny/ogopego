package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/testutil"
)

func TestKeywordsInRuleNames(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start = whitespace
		whitespace = {'x'}+
	`, nil)
	testutil.AssertJSONStr(t, g, "x", `["x"]`)
}
