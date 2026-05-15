package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestNamedCapture(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := name='hello'
	`, nil)
	testutil.AssertJSONStr(t, g, "hello", `{"name": "hello"}`)
}

func TestOverrideSingleton(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start: ='hello'
	`, nil)
	testutil.AssertJSONStr(t, g, "hello", `"hello"`)
}
