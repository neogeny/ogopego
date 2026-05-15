package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestDirectLeftRecursion(t *testing.T) {
	g := ogopego.Compile(t, `
		@@left_recursion :: True
		@@whitespace :: /\s+/
		@@grammar :: Test
		start = expression $
		expression = expression '+' factor | expression '-' factor | factor
		factor = number
		number = /[0-9]+/
	`, nil)
	ogopego.AssertJSONStr(t, g, "10 - 20", `["10", "-", "20"]`)
}
