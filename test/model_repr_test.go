package test

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/asjson"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/tool"
)

func TestModelReprOutput(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Calc
		@@whitespace :: /\s+/

		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`, nil)
	code := tool.GenerateGrammarModel(*g, "calc")
	assert.NotZero(t, code, "empty output")
}

func TestPairFromTree(t *testing.T) {
	g := Compile(t, `
		@@grammar :: Calc
		@@whitespace :: /\s+/

		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`, nil)
	cfg := &config.Cfg{}
	tree, err := api.ParseInput(g, "(abc:123)", cfg)
	assert.NoError(t, err, "parse")

	j := asjson.AsJSON(tree).(map[string]any)
	assert.Equal(t, "Pair", j["__class__"])
	assert.Equal(t, 3, len(j), "expected 3 entries (__class__, key, val)")
	assert.Equal(t, "abc", j["key"], "key")
	assert.Equal(t, "123", j["val"], "val")
}

func TestModelReprTypedRef(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Calc
		start::Start = child:pair $
		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
		NUMBER::Num = /\d+/
	`), nil)
	code := tool.GenerateGrammarModel(*g, "calc")
	assert.True(t, strings.Contains(code, "Child *Pair"), "expected Start.Child *Pair field")
	assert.True(t, strings.Contains(code, "PairFromTree("), "expected PairFromTree call in StartFromTree")
	assert.True(t, strings.Contains(code, "Value any"), "expected Num.Value any field")
}

func TestTypedRefFromTree(t *testing.T) {
	g := Compile(t, Dedent(`
		@@grammar :: Calc
		start::Start = child:pair $
		pair::Pair = '(' key:/[a-z]+/ ':' val:/[0-9]+/ ')'
	`), nil)
	cfg := &config.Cfg{}
	tree, err := api.ParseInput(g, "(abc:123)", cfg)
	assert.NoError(t, err, "parse")

	j := asjson.AsJSON(tree).(map[string]any)
	assert.Equal(t, "Start", j["__class__"])
	child := j["child"].(map[string]any)
	assert.Equal(t, "Pair", child["__class__"])
	assert.Equal(t, "abc", child["key"], "key")
	assert.Equal(t, "123", child["val"], "val")
}
