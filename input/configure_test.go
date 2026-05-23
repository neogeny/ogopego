// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import (
	"testing"
)

func TestConfigureSetsIgnoreCase(t *testing.T) {
	s := NewStrCursor("hello")
	cfg := Cfg{IgnoreCase: true}
	s.Configure(cfg)
	if !s.IgnoreCase() {
		t.Error("expected IgnoreCase to be set")
	}
}

func TestConfigureSetsNameGuard(t *testing.T) {
	s := NewStrCursor("hello")
	cfg := Cfg{NameGuard: true}
	s.Configure(cfg)
	if !s.NameGuard() {
		t.Error("expected NameGuard to be set")
	}
}

func TestConfigureSetsSource(t *testing.T) {
	s := NewStrCursor("hello")
	cfg := Cfg{Source: "test source"}
	s.Configure(cfg)
	if s.InputSource() != "test source" {
		t.Errorf("expected 'test source', got %q", s.InputSource())
	}
}

func TestConfigureSkipsEmptySource(t *testing.T) {
	s := NewStrCursor("hello")
	s.heavy.Source = "original"
	cfg := Cfg{}
	s.Configure(cfg)
	if s.InputSource() != "original" {
		t.Errorf("expected source to remain 'original', got %q", s.InputSource())
	}
}

func TestConfigureSetsPatterns(t *testing.T) {
	wsp := `[ \t]+`
	cmt := `//[^\n]*`
	eol := `\r?\n`
	ws := wsp
	s := NewStrCursor("  // comment \n  hello")
	s.Configure(Cfg{Whitespace: &ws, Comments: cmt, EolComments: eol})
	s.NextToken()
	if !s.MatchToken("hello") {
		t.Error("expected to match 'hello' after skipping whitespace and comment")
	}
}

func TestConfigureEmptyPatterns(t *testing.T) {
	ws := ``
	s := NewStrCursor("  hello")
	s.Configure(Cfg{Whitespace: &ws})
	s.NextToken()
	if s.MatchToken("hello") {
		t.Error("expected NOT to match 'hello' with empty whitespace pattern (should not skip)")
	}
}

func TestConfigureNilWhitespacePattern(t *testing.T) {
	s := NewStrCursor("  hello")
	s.Configure(Cfg{})
	s.NextToken()
	if s.MatchToken("hello") {
		t.Error("expected NOT to match 'hello' with nil whitespace (should not skip)")
	}
}

func TestConfigureIgnoresBadPattern(t *testing.T) {
	ws := `[invalid`
	wsPtr := ws
	s := NewStrCursor("hello")
	s.Configure(Cfg{Whitespace: &wsPtr})
	_ = s
}
