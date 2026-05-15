package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestSkipTo(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start = 'x' ab $
		ab = 'a' 'b' | -> 'b'
	`, nil)
	ogopego.AssertJSONStr(t, g, "x yb", `["x", "b"]`)
}
