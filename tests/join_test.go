package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestPositiveJoin(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@nameguard :: False
		@@grammar :: Test
		start := ','%{'x' 'y'}+
	`, nil)
	testutil.AssertJSONStr(t, g, "x y, x y", `[["x", "y"], ",", ["x", "y"]]`)
	testutil.AssertJSONStr(t, g, "x y x y", `[["x", "y"]]`)
	testutil.ParseFail(t, g, "y x")
}
