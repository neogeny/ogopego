package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestWhitespaceDirectiveDoubleQuote(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /[\t ]+/
		@@grammar :: Test
		test := "test" $
	`, nil)
	ogopego.AssertJSONStr(t, g, "test", `"test"`)
}
