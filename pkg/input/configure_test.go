// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestConfigureSetsIgnoreCase(t *testing.T) {
	s := NewStrCursor("hello")
	cfg := Cfg{IgnoreCase: true}
	s.Configure(cfg)
	assert.True(t, s.IgnoreCase(), "expected IgnoreCase to be set")
}

func TestConfigureSetsNameGuard(t *testing.T) {
	s := NewStrCursor("hello")
	cfg := Cfg{NameGuard: new(true)}
	s.Configure(cfg)
	assert.True(t, s.NameGuard(), "expected NameGuard to be set")
}

func TestConfigureSetsSource(t *testing.T) {
	s := NewStrCursor("hello")
	cfg := Cfg{Source: "test source"}
	s.Configure(cfg)
	assert.Equal(t, "test source", s.InputSource())
}

func TestConfigureSkipsEmptySource(t *testing.T) {
	s := NewStrCursor("hello")
	s.heavy.Source = "original"
	cfg := Cfg{}
	s.Configure(cfg)
	assert.Equal(t, "original", s.InputSource())
}

func TestConfigureSetsPatterns(t *testing.T) {
	wsp := `[ \t]+`
	cmt := `//[^\n]*`
	eol := `\r?\n`
	s := NewStrCursor("  // comment \n  hello")
	s.Configure(Cfg{Whitespace: new(wsp), Comments: cmt, EolComments: eol})
	s.NextToken()
	assert.True(t, s.MatchToken("hello"), "expected to match 'hello' after skipping whitespace and comment")
}

func TestConfigureEmptyPatterns(t *testing.T) {
	s := NewStrCursor("  hello")
	s.Configure(Cfg{Whitespace: new(``)})
	s.NextToken()
	assert.False(t, s.MatchToken("hello"), "expected NOT to match 'hello' with empty whitespace pattern (should not skip)")
}

func TestConfigureNilWhitespacePattern(t *testing.T) {
	s := NewStrCursor("  hello")
	s.Configure(Cfg{Whitespace: new("")})
	s.NextToken()
	assert.False(t, s.MatchToken("hello"), "expected NOT to match 'hello' with nil whitespace (should not skip)")
}

func TestConfigureNameGuardWithNameChars(t *testing.T) {
	s := NewStrCursor("if_else")
	// NameChars implies NameGuard (handled in config.Override)
	s.Configure(Cfg{NameChars: "_"})
	assert.True(t, s.NameGuard(), "expected NameGuard to be enabled when NameChars is set")
}

func TestConfigureIgnoresBadPattern(t *testing.T) {
	ws := `[invalid`
	s := NewStrCursor("hello")
	s.Configure(Cfg{Whitespace: new(ws)})
	_ = s
}
