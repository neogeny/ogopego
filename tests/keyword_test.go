package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestKeywordsInRuleNames(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start = whitespace ;
		whitespace = {'x'}+ ;
	`)
	testutil.AssertJSONStr(t, g, "x", `["x"]`)
}
