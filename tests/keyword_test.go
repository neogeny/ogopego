package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestKeywordsInRuleNames(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start = whitespace
		whitespace = {'x'}+
	`, nil)
	util.AssertJSONStr(t, g, "x", `["x"]`)
}
