// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"os"
	"testing"
)

func TestLoadBootGrammar(t *testing.T) {
	data, err := os.ReadFile("../grammar/tatsu.json")
	if err != nil {
		t.Fatal(err)
	}
	g, err := LoadBootGrammar(data)
	if err != nil {
		t.Fatalf("LoadBootGrammar error: %v", err)
	}
	if g.Name == "" {
		t.Error("expected non-empty grammar name")
	}
	if !g.Analyzed {
		t.Error("expected grammar to be analyzed after Initialize")
	}
	if len(g.Rules) < 10 {
		t.Errorf("expected at least 10 rules, got %d", len(g.Rules))
	}
	required := []string{"start", "grammar", "rule", "expre", "choice", "sequence"}
	for _, name := range required {
		rule, err := g.GetRule(name)
		if err != nil {
			t.Errorf("missing required rule %q: %v", name, err)
			continue
		}
		if rule == nil {
			t.Errorf("GetRule(%q) returned nil", name)
		}
	}
}
