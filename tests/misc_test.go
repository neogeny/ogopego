package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/testutil"
)

func TestMapping(t *testing.T) {
	g := testutil.Compile(t, `
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = key:key value:value
		key = /\w+/
		value = /\w+/
	`, nil)
	// Named items in a sequence are auto-folded into a single map.
	testutil.AssertJSONStr(t, g, "foo bar", `{"key": "foo", "value": "bar"}`)
}
