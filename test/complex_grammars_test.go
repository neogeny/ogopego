package test

import (
	"testing"
)

func TestCalculatorGrammar(t *testing.T) {
	g := Compile(t, `
		@@grammar :: CALC
		@@left_recursion :: True
		@@whitespace :: /\s+/

		start := expression $

		expression
			:= expression '+' term
			| expression '-' term
			| term

		term
			:= term '*' factor
			| term '/' factor
			| factor

		factor
			:= NUMBER
			| '(' expression ')'

		NUMBER := /\d+/
	`, nil)
	result := ParseJSON(t, g, "3 + 5 * (10 - 20 )")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestJSONLikeGrammar(t *testing.T) {
	g := Compile(t, `
		@@grammar :: MiniJSON
		@@nameguard :: False
		@@whitespace :: /\s+/

		start := value $

		value := object | array | string | number | 'true' | 'false' | 'null'

		object := '{' members? '}'

		array := '[' elements? ']'

		members := pair (',' pair)*

		elements := value (',' value)*

		pair := string ':' value

		string := '"' CONTENT '"'

		CONTENT := /[^"]*/

		number := /-?\d+(\.\d+)?/
	`, nil)
	result := ParseJSON(t, g, `{"key": "value"}`)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
