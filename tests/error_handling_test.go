package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/testutil"
)

func TestInvalidInputFails(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a'
	`,
		nil,
	)
	testutil.ParseFail(t, g, "b")
}

func TestPartialMatchFails(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a' 'b'
	`, nil)
	testutil.ParseFail(t, g, "a")
}

func TestEmptyInputFailsWhenRequired(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a'
	`, nil)
	testutil.ParseFail(t, g, "")
}
