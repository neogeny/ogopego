package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestSkipTo(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start = 'x' ab $ ;
		ab = 'a' 'b' | -> 'b' ;
	`)
	testutil.AssertJSONStr(t, g, "x yb", `["x", "b"]`)
}
