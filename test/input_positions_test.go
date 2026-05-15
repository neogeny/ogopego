package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestMultiLineInput(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start := 'hello' 'world'
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello\nworld", `["hello", "world"]`)
}
