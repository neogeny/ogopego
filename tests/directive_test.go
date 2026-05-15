package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/testutil"
)

func TestWhitespaceDirectiveDoubleQuote(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /[\t ]+/
		@@grammar :: Test
		test := "test" $
	`, nil)
	testutil.AssertJSONStr(t, g, "test", `"test"`)
}
