package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestMultiLineInput(t *testing.T) {
	g := util.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'hello' 'world'
	`, nil)
	util.AssertJSONStr(t, g, "hello\nworld", `["hello", "world"]`)
}
