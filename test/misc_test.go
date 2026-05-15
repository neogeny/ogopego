package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestMapping(t *testing.T) {
	g := ogopego.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = key:key value:value
		key = /\w+/
		value = /\w+/
	`, nil)
	// Named items in a sequence are auto-folded into a single map.
	ogopego.AssertJSONStr(t, g, "foo bar", `{"key": "foo", "value": "bar"}`)
}
