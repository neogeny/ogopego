package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestPositiveJoin(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /\s+/
		@@nameguard :: False
		@@grammar :: Test
		start := ','%{'x' 'y'}+
	`, nil)
	util.AssertJSONStr(t, g, "x y, x y", `[["x", "y"], ",", ["x", "y"]]`)
	util.AssertJSONStr(t, g, "x y x y", `[["x", "y"]]`)
	util.ParseFail(t, g, "y x")
}
