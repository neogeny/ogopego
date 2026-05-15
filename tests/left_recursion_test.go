package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/testutil"
)

func TestDirectLeftRecursion(t *testing.T) {
	g := testutil.Compile(t, `
		@@left_recursion :: True
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = expression $ ;
		expression = expression '+' factor | expression '-' factor | factor ;
		factor = number ;
		number = /[0-9]+/ ;
	`)
	testutil.AssertJSONStr(t, g, "10 - 20", `["10", "-", "20"]`)
}
