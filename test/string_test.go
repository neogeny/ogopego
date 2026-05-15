package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
	"github.com/neogeny/ogopego/util"
)

func TestMultilineString(t *testing.T) {
	g := ogopego.Compile(t, util.Dedent(`
		@@grammar :: Test
		start := longone | shortone $
		shortone := "short"
		longone := """
			this "text"
			is a long "string"
			"""
	`), nil)
	ogopego.AssertJSONStr(t, g, "short", `"short"`)
}
