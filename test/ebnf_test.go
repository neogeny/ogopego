package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestEBNFParsing(t *testing.T) {
	g := ogopego.Compile(t, `
		@@grammar :: EBNF

		start := expression $

		expression := expression '+' term | expression '-' term | term

		term := term '*' factor | term '/' factor | factor

		factor := '(' expression ')' | number

		number := /\d+/
	`, nil)
	if g.Name != "EBNF" {
		t.Errorf("expected name 'EBNF', got %q", g.Name)
	}
}
