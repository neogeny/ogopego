package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestMultilineString(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := longone | shortone $
		shortone := "short"
		longone := """
			this "text"
			is a long "string"
			"""
	`, nil)
	util.AssertJSONStr(t, g, "short", `"short"`)
}
