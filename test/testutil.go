package test

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/peg"
	"github.com/neogeny/ogopego/pkg/util"
)

type Cfg config.Cfg

func Dedent(s string) string {
	return util.Dedent(s)
}

func Compile(t testing.TB, grammar string, cfg *config.Cfg) *peg.Grammar {
	t.Helper()
	g, err := api.Compile(grammar, cfg)
	assert.NoError(t, err, "compile")
	return g
}

func ParseJSON(t testing.TB, g *peg.Grammar, text string) any {
	t.Helper()
	result, err := api.ParseInputToJSON(g, text, nil)
	assert.NoError(t, err, "parse %q", text)
	return result
}

func ParseFail(t testing.TB, g *peg.Grammar, text string) {
	t.Helper()
	_, err := api.ParseInputToJSON(g, text, nil)
	assert.Error(t, err, "expected parse error for %q", text)
}

func AssertJSON(t testing.TB, g *peg.Grammar, text string, want any) {
	t.Helper()
	got := ParseJSON(t, g, text)
	assert.Equal(t, want, got, "input %q", text)
}

func AssertJSONStr(t testing.TB, g *peg.Grammar, text string, wantJSON string) {
	t.Helper()
	var want any
	if err := json.Unmarshal([]byte(wantJSON), &want); err != nil {
		t.Fatalf("invalid want JSON %q: %v", wantJSON, err)
	}
	AssertJSON(t, g, text, want)
}
