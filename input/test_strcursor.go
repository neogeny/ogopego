// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package input

import (
	"testing"

	"github.com/neogeny/ogopego/util/pyre"
)

func TestStrCursorMatchPattern(t *testing.T) {
	p, err := pyre.Compile(`\d+`)
	if err != nil {
		t.Fatal(err)
	}
	s := NewStrCursor("abc123def")
	m, ok := s.MatchPattern(p)
	if ok {
		t.Errorf("expected no match at start, got %q", m)
	}
}

func TestStrCursorMatchToken(t *testing.T) {
	s := NewStrCursor("hello world")
	if !s.MatchToken("hello") {
		t.Error("expected MatchToken to succeed")
	}
	if s.offset != 5 {
		t.Errorf("expected offset 5, got %d", s.offset)
	}
}

func TestStrCursorPeekToken(t *testing.T) {
	s := NewStrCursor("hello")
	if !s.PeekToken("hello") {
		t.Error("expected PeekToken to succeed")
	}
	if s.offset != 0 {
		t.Errorf("expected offset unchanged, got %d", s.offset)
	}
}

func TestStrCursorNextPeek(t *testing.T) {
	s := NewStrCursor("ab")
	r, ok := s.Peek()
	if !ok || r != 'a' {
		t.Errorf("expected 'a', got %c", r)
	}
	r, ok = s.Next()
	if !ok || r != 'a' {
		t.Errorf("expected 'a', got %c", r)
	}
	if s.offset != 1 {
		t.Errorf("expected offset 1, got %d", s.offset)
	}
}

func TestStrCursorAtEnd(t *testing.T) {
	s := NewStrCursor("a")
	if s.AtEnd() {
		t.Error("expected not at end")
	}
	s.Next()
	if !s.AtEnd() {
		t.Error("expected at end")
	}
}

func TestStrCursorPos(t *testing.T) {
	s := NewStrCursor("hello\nworld")
	line, col := s.Pos()
	if line != 0 || col != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", line, col)
	}
	for i := 0; i < 6; i++ {
		s.Next()
	}
	line, col = s.Pos()
	if line != 1 || col != 5 {
		t.Errorf("expected (1,5) after newline, got (%d,%d)", line, col)
	}
}

func TestStrCursorClone(t *testing.T) {
	s := NewStrCursor("hello")
	c := s.Clone()
	s.Next()
	if s.Mark() == c.Mark() {
		t.Error("expected clone to have independent offset")
	}
}
