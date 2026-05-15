package ogopego

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/config"
	"github.com/neogeny/ogopego/peg"
)

type Cfg = config.Cfg

func Compile(t testing.TB, grammar string, cfg *config.Cfg) *peg.Grammar {
	t.Helper()
	g, err := api.Compile(grammar, &Cfg{
		Trace:    true,
		Colorize: true,
	})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	return g
}

func ParseJSON(t testing.TB, g *peg.Grammar, text string) any {
	t.Helper()
	result, err := api.ParseInputToJSON(g, text, nil)
	if err != nil {
		t.Fatalf("parse %q: %v", text, err)
	}
	return result
}

func ParseFail(t testing.TB, g *peg.Grammar, text string) {
	t.Helper()
	_, err := api.ParseInputToJSON(g, text, nil)
	if err == nil {
		t.Fatalf("expected parse error for %q", text)
	}
}

func AssertJSON(t testing.TB, g *peg.Grammar, text string, want any) {
	t.Helper()
	got := ParseJSON(t, g, text)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("input %q\ngot:  %#v\nwant: %#v", text, got, want)
	}
}

func AssertJSONStr(t testing.TB, g *peg.Grammar, text string, wantJSON string) {
	t.Helper()
	var want any
	if err := json.Unmarshal([]byte(wantJSON), &want); err != nil {
		t.Fatalf("invalid want JSON %q: %v", wantJSON, err)
	}
	AssertJSON(t, g, text, want)
}
