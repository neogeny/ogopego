package test

import (
	"testing"

	"github.com/neogeny/ogopego/tool"
)

func TestModelReprBasic(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Calc
		@@whitespace :: /\s+/

		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`, nil)
	code := tool.ModelRepr(*g, "calc")
	t.Logf("Generated code:\n%s", code)
}
