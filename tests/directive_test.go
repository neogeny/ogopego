package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestWhitespaceDirectiveDoubleQuote(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /[\t ]+/
		@@grammar :: Test
		test := "test" $
	`, nil)
	util.AssertJSONStr(t, g, "test", `"test"`)
}
