package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util"
)

func TestInvalidInputFails(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'a'
	`,
		nil,
	)
	util.ParseFail(t, g, "b")
}

func TestPartialMatchFails(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'a' 'b'
	`, nil)
	util.ParseFail(t, g, "a")
}

func TestEmptyInputFailsWhenRequired(t *testing.T) {
	g := util.Compile(t, `
		@@grammar :: Test
		start := 'a'
	`, nil)
	util.ParseFail(t, g, "")
}
