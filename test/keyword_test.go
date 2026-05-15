package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestKeywordsInRuleNames(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start = whitespace
		whitespace = {'x'}+
	`, nil)
	ogopego.AssertJSONStr(t, g, "x", `["x"]`)
}
