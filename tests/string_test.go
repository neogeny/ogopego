package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestMultilineString(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := longone | shortone $
		shortone := "short"
		longone := """
			this "text"
			is a long "string"
			"""
	`, nil)
	testutil.AssertJSONStr(t, g, "short", `"short"`)
}
