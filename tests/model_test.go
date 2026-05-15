package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestChildren(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Calc
		start = expression $
		expression = term
		term = 'x'
	`, nil)
	testutil.AssertJSONStr(t, g, "x", `"x"`)
}

func TestNodeKWArgs(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start = 'value'
	`, nil)
	testutil.AssertJSONStr(t, g, "value", `"value"`)
}
