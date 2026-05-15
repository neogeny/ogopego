package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestNamedCapture(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := name='hello'
	`, nil)
	util.AssertJSONStr(t, g, "hello", `{"name": "hello"}`)
}

func TestOverrideSingleton(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start: ='hello'
	`, nil)
	util.AssertJSONStr(t, g, "hello", `"hello"`)
}
