package test

import (
	"testing"

	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/peg"
)

func TestRoundtripJSON(t *testing.T) {
	grammar := `
		@@grammar :: Calc
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
	`

	g1, err := api.Compile(grammar, nil)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	r1, err := api.ParseInputToJSONString(g1, "3 + 5 * (10 - 20 )", nil)
	if err != nil {
		t.Fatalf("g1 parse: %v", err)
	}

	jsonStr := peg.ModelToJSONStr(g1)

	g2, err := peg.ParseGrammar([]byte(jsonStr))
	if err != nil {
		t.Fatalf("parse JSON: %v\nJSON:\n%s", err, jsonStr)
	}
	if err := g2.Initialize(); err != nil {
		t.Fatalf("init: %v", err)
	}

	if len(g1.Rules) != len(g2.Rules) {
		t.Errorf("rule count: %d vs %d", len(g1.Rules), len(g2.Rules))
	}
	for i := range g1.Rules {
		if g1.Rules[i].Name != g2.Rules[i].Name {
			t.Errorf("rule %d name: %q vs %q", i, g1.Rules[i].Name, g2.Rules[i].Name)
		}
	}
	if g1.Name != g2.Name {
		t.Errorf("name: %q vs %q", g1.Name, g2.Name)
	}

	r2, err := api.ParseInputToJSONString(g2, "3 + 5 * (10 - 20 )", nil)
	if err != nil {
		t.Fatalf("g2 parse: %v", err)
	}
	if r1 != r2 {
		t.Errorf("parse mismatch:\ng1: %s\ng2: %s", r1, r2)
	}
}

func TestRoundtripPrettyPrint(t *testing.T) {
	grammar := `
		@@grammar :: Calc
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
	`

	g1, err := api.Compile(grammar, nil)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	r1, err := api.ParseInputToJSONString(g1, "3 + 5", nil)
	if err != nil {
		t.Fatalf("g1 parse: %v", err)
	}

	pretty := g1.PrettyPrint()

	g2, err := api.Compile(pretty, nil)
	if err != nil {
		t.Fatalf("recompile: %v\npretty:\n%s", err, pretty)
	}

	r2, err := api.ParseInputToJSONString(g2, "3 + 5", nil)
	if err != nil {
		t.Fatalf("g2 parse: %v", err)
	}

	if r1 != r2 {
		t.Errorf("mismatch:\ng1: %s\ng2: %s", r1, r2)
	}
	if len(g1.Rules) != len(g2.Rules) {
		t.Errorf("rule count: %d vs %d", len(g1.Rules), len(g2.Rules))
	}
}
