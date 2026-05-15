package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestInvalidInputFails(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'a'
	`,
		nil,
	)
	ogopego.ParseFail(t, g, "b")
}

func TestPartialMatchFails(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'a' 'b'
	`, nil)
	ogopego.ParseFail(t, g, "a")
}

func TestEmptyInputFailsWhenRequired(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: Test
		start := 'a'
	`, nil)
	ogopego.ParseFail(t, g, "")
}
