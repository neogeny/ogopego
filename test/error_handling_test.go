package test

import (
	"testing"
)

func TestInvalidInputFails(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'a'
	`,
		nil,
	)
	ParseFail(t, g, "b")
}

func TestPartialMatchFails(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'a' 'b'
	`, nil)
	ParseFail(t, g, "a")
}

func TestEmptyInputFailsWhenRequired(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Test
		start := 'a'
	`, nil)
	ParseFail(t, g, "")
}
