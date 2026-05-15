package test

import (
	"testing"
)

func TestEBNFParsing(t *testing.T) {
	g := Compile(t, `
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
