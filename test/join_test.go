package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestPositiveJoin(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		@@nameguard :: False
		@@grammar :: Test
		start := ','%{'x' 'y'}+
	`, nil)
	ogopego.AssertJSONStr(t, g, "x y, x y", `[["x", "y"], ",", ["x", "y"]]`)
	ogopego.AssertJSONStr(t, g, "x y x y", `[["x", "y"]]`)
	ogopego.ParseFail(t, g, "y x")
}
