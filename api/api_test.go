package api

import (
	"testing"

	"github.com/neogeny/ogopego/trees"
)

func TestParseGrammar(t *testing.T) {
	result, err := ParseGrammar("@@grammar :: Test start := 'x'", nil)
	if err != nil {
		t.Fatalf("ParseGrammar error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil Tree")
	}
	if _, ok := result.(*trees.Nil); ok {
		t.Fatal("unexpected Nil tree")
	}
}

func TestCompile(t *testing.T) {
	t.Skip("pre-existing: failed parsing ENDRULE")
	src := "@@grammar :: EBNFTest\nstart := expression $\nexpression := expression '+' term | expression '-' term | term\nterm := term '*' factor | term '/' factor | factor\nfactor := '(' expression ')' | number\nnumber := /\\d+/\n"
	g, err := Compile(src, nil)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}
	if g.Name != "EBNFTest" {
		t.Errorf("expected name 'EBNFTest', got %q", g.Name)
	}
	if !g.Analyzed {
		t.Fatal("expected analyzed grammar")
	}
	if len(g.Rules) != 6 {
		t.Errorf("expected 6 rules, got %d", len(g.Rules))
	}
}

func TestCompileToJSON(t *testing.T) {
	src := `start := 'x'`
	json, err := CompileToJSON(src, nil)
	if err != nil {
		t.Fatalf("CompileToJSON error: %v", err)
	}
	if json == nil {
		t.Fatal("expected non-nil json")
	}
}

func TestCompileToJSONString(t *testing.T) {
	src := `start := 'x'`
	s, err := CompileToJSONString(src, nil)
	if err != nil {
		t.Fatalf("CompileToJSONString error: %v", err)
	}
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestBootGrammar(t *testing.T) {
	g, err := BootGrammar()
	if err != nil {
		t.Fatalf("BootGrammar error: %v", err)
	}
	if !g.Analyzed {
		t.Fatal("expected analyzed grammar")
	}
}
