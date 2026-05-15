package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestNamedCapture(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := name='hello'
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello", `{"name": "hello"}`)
}

func TestOverrideSingleton(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start: ='hello'
	`, nil)
	ogopego.AssertJSONStr(t, g, "hello", `"hello"`)
}
