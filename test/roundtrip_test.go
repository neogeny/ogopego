package test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/peg"
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
	assert.NoError(t, err, "compile")

	r1, err := api.ParseInputToJSONString(g1, "3 + 5 * (10 - 20 )", nil)
	assert.NoError(t, err, "g1 parse")

	jsonStr := peg.ModelToJSONStr(g1)

	g2, err := peg.LoadGrammarFromJSON([]byte(jsonStr))
	assert.NoError(t, err, "parse JSON:\nJSON:\n%s", jsonStr)
	assert.NoError(t, g2.Initialize(), "init")

	assert.Equal(t, len(g1.Rules), len(g2.Rules), "rule count")
	for i := range g1.Rules {
		assert.Equal(t, g1.Rules[i].Name, g2.Rules[i].Name, "rule %d name", i)
	}
	assert.Equal(t, g1.Name, g2.Name, "name")

	r2, err := api.ParseInputToJSONString(g2, "3 + 5 * (10 - 20 )", nil)
	assert.NoError(t, err, "g2 parse")
	assert.Equal(t, r1, r2, "parse mismatch")
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
	assert.NoError(t, err, "compile")

	r1, err := api.ParseInputToJSONString(g1, "3 + 5", nil)
	assert.NoError(t, err, "g1 parse")

	pretty := g1.PrettyPrint()

	g2, err := api.Compile(pretty, nil)
	assert.NoError(t, err, "recompile:\npretty:\n%s", pretty)

	r2, err := api.ParseInputToJSONString(g2, "3 + 5", nil)
	assert.NoError(t, err, "g2 parse")

	assert.Equal(t, r1, r2, "mismatch")
	assert.Equal(t, len(g1.Rules), len(g2.Rules), "rule count")
}
