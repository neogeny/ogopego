// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	_ "embed"
	"testing"

	"github.com/neogeny/ogopego"
)

func TestImportCalcJSON(t *testing.T) {
	data := ogopego.CalcJSON
	g, err := LoadGrammarFromJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if g.Name != "CALC" {
		t.Errorf("expected name CALC, got %q", g.Name)
	}
	if len(g.Rules) == 0 {
		t.Fatal("expected at least one rule")
	}
	startRule := g.Rules[0]
	if startRule.Name != "start" {
		t.Errorf("expected first rule 'start', got %q", startRule.Name)
	}
	seq, ok := startRule.Exp.(*Sequence)
	if !ok {
		t.Fatalf("expected start rule exp to be Sequence, got %T", startRule.Exp)
	}
	if len(seq.Sequence) != 2 {
		t.Fatalf("expected 2 items in start sequence, got %d", len(seq.Sequence))
	}
	call, ok := seq.Sequence[0].(*Call)
	if !ok {
		t.Fatalf("expected first item to be Call, got %T", seq.Sequence[0])
	}
	if call.Name != "expression" {
		t.Errorf("expected call to 'expression', got %q", call.Name)
	}
	eof, ok := seq.Sequence[1].(*EOF)
	if !ok {
		t.Fatalf("expected second item to be EOF, got %T", seq.Sequence[1])
	}
	_ = eof
}

func TestImportTatsuJSON(t *testing.T) {
	data := ogopego.TatSuGrammarJSON
	g, err := LoadGrammarFromJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if g.Name == "" {
		t.Error("expected non-empty grammar name")
	}
	if len(g.Rules) < 10 {
		t.Errorf("expected at least 10 rules in TatSu grammar, got %d", len(g.Rules))
	}
	// Check for key rules
	names := make(map[string]bool)
	for _, r := range g.Rules {
		names[r.Name] = true
	}
	for _, want := range []string{"start", "grammar", "rule", "expre", "choice", "sequence"} {
		if !names[want] {
			t.Errorf("missing rule %q", want)
		}
	}
}

func TestImportInvalidJSON(t *testing.T) {
	_, err := LoadGrammarFromJSON([]byte(`{"__class__": "NotAGrammar"}`))
	if err == nil {
		t.Fatal("expected error for non-Grammar root")
	}
}

func TestImportMalformedJSON(t *testing.T) {
	_, err := LoadGrammarFromJSON([]byte(`{not valid json`))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}
