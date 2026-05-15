package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestChildren(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Calc
		start = expression $
		expression = term
		term = 'x'
	`, nil)
	util.AssertJSONStr(t, g, "x", `"x"`)
}

func TestNodeKWArgs(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start = 'value'
	`, nil)
	util.AssertJSONStr(t, g, "value", `"value"`)
}
