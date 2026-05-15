package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/testutil"
)

func TestMultiLineInput(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'hello' 'world'
	`, nil)
	testutil.AssertJSONStr(t, g, "hello\nworld", `["hello", "world"]`)
}
