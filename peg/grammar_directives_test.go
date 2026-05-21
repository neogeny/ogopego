// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"testing"
)

func TestBootGrammarCfgFromDirectives(t *testing.T) {
	g, err := BootGrammar()
	if err != nil {
		t.Fatalf("LoadBootGrammar: %v", err)
	}

	cfg := *g.CfgFromDirectives()

	if cfg.Grammar != "TatSu" {
		t.Errorf("expected Grammar 'TatSu', got %q", cfg.Grammar)
	}
	if cfg.Name != "" {
		t.Errorf("expected empty Name, got %q", cfg.Name)
	}
	if cfg.Whitespace == nil {
		t.Fatal("expected Whitespace pattern")
	}
	if *cfg.Whitespace != `(?m)\s+` {
		t.Errorf("expected whitespace pattern, got %q", *cfg.Whitespace)
	}
	if cfg.Comments == "" {
		t.Error("expected Comments pattern")
	}
	if cfg.EolComments == "" {
		t.Error("expected EolComments pattern")
	}
	if !cfg.ParseInfo {
		t.Error("expected ParseInfo to be true")
	}
	if !cfg.NoLeftRecursion {
		t.Error("expected NoLeftRecursion to be true (from left_recursion: false)")
	}
	if cfg.IgnoreCase {
		t.Error("expected IgnoreCase to be false")
	}
	if cfg.Source != "" {
		t.Errorf("expected empty Source, got %q", cfg.Source)
	}
	if cfg.Keywords != nil {
		t.Errorf("expected Keywords to be nil, got %v", cfg.Keywords)
	}
	if cfg.Semantics == nil {
		t.Error("expected Semantics to be non-nil")
	}
}
