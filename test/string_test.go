package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestMultilineString(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := longone | shortone $
		shortone := "short"
		longone := """
			this "text"
			is a long "string"
			"""
	`, nil)
	ogopego.AssertJSONStr(t, g, "short", `"short"`)
}
