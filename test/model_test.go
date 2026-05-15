package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestChildren(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Calc
		start = expression $
		expression = term
		term = 'x'
	`, nil)
	ogopego.AssertJSONStr(t, g, "x", `"x"`)
}

func TestNodeKWArgs(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start = 'value'
	`, nil)
	ogopego.AssertJSONStr(t, g, "value", `"value"`)
}
