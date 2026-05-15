package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestInvalidInputFails(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a' ;
	`)
	testutil.ParseFail(t, g, "b")
}

func TestPartialMatchFails(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a' 'b' ;
	`)
	testutil.ParseFail(t, g, "a")
}

func TestEmptyInputFailsWhenRequired(t *testing.T) {
	g := testutil.Compile(t, `
		@@grammar :: Test
		start := 'a' ;
	`)
	testutil.ParseFail(t, g, "")
}
