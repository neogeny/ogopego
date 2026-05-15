package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestSkipTo(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start = 'x' ab $
		ab = 'a' 'b' | -> 'b'
	`, nil)
	util.AssertJSONStr(t, g, "x yb", `["x", "b"]`)
}
