// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package input

import (
	"testing"

	"github.com/neogeny/ogopego/util/pyre"
)

func TestMatchPatternSuccess(t *testing.T) {
	p, err := pyre.Compile(`\d+`)
	if err != nil {
		t.Fatal(err)
	}
	s := NewStrCursor("abc123def")
	// skip past "abc"
	s.Next()
	s.Next()
	s.Next()
	m, ok := s.MatchPattern(p)
	if !ok {
		t.Fatal("expected match")
	}
	if m != "123" {
		t.Errorf("expected '123', got %q", m)
	}
	if s.offset != 6 {
		t.Errorf("expected offset 6, got %d", s.offset)
	}
}

func TestMatchPatternWithGroup(t *testing.T) {
	p, err := pyre.Compile(`(foo)\s+(bar)`)
	if err != nil {
		t.Fatal(err)
	}
	s := NewStrCursor("foo bar baz")
	m, ok := s.MatchPattern(p)
	if !ok {
		t.Fatal("expected match")
	}
	if m != "foo" {
		t.Errorf("expected 'foo' from group 1, got %q", m)
	}
}

func TestMatchEOL(t *testing.T) {
	tests := []struct {
		input string
		want  bool
		off   int
	}{
		{"\n", true, 1},
		{"\r\n", true, 2},
		{"\r", true, 1},
		{"abc", false, 0},
		{"", false, 0},
	}
	for _, tt := range tests {
		s := NewStrCursor(tt.input)
		got := s.MatchEOL()
		if got != tt.want {
			t.Errorf("MatchEOL(%q) = %v, want %v", tt.input, got, tt.want)
		}
		if s.offset != tt.off {
			t.Errorf("MatchEOL(%q) offset = %d, want %d", tt.input, s.offset, tt.off)
		}
	}
}

func TestMatchEOLWithWhitespace(t *testing.T) {
	p, err := pyre.Compile(`[ \t]+`)
	if err != nil {
		t.Fatal(err)
	}
	s := NewStrCursor("   \nnext")
	s.SetPatterns(&TokenizingPatterns{Wsp: p})
	if !s.MatchEOL() {
		t.Error("expected MatchEOL to skip whitespace and match newline")
	}
	if s.offset != 4 {
		t.Errorf("expected offset 4 (whitespace+newline), got %d", s.offset)
	}
}

func TestPosAt(t *testing.T) {
	s := NewStrCursor("hello\nworld\nfoo")
	line, col := s.PosAt(0)
	if line != 0 || col != 0 {
		t.Errorf("expected (0,0) at pos 0, got (%d,%d)", line, col)
	}
	line, col = s.PosAt(5)
	if line != 1 || col != 5 {
		t.Errorf("expected (1,5) at pos 5, got (%d,%d)", line, col)
	}
	line, col = s.PosAt(6)
	if line != 1 || col != 5 {
		t.Errorf("expected (1,5) at pos 6, got (%d,%d)", line, col)
	}
	line, col = s.PosAt(11)
	if line != 2 || col != 5 {
		t.Errorf("expected (2,5) at pos 11, got (%d,%d)", line, col)
	}
	line, col = s.PosAt(12)
	if line != 2 || col != 5 {
		t.Errorf("expected (2,5) at pos 12, got (%d,%d)", line, col)
	}
	line, col = s.PosAt(15)
	if line != 3 || col != 3 {
		t.Errorf("expected (3,3) at pos 15, got (%d,%d)", line, col)
	}
}

func TestPosAtPastEnd(t *testing.T) {
	s := NewStrCursor("hi")
	line, col := s.PosAt(100)
	if line != 1 || col != 2 {
		t.Errorf("expected (1,2) at past-end, got (%d,%d)", line, col)
	}
}

func TestLocation(t *testing.T) {
	s := NewStrCursorFromSource("test.txt", "hello\nworld", 0)
	loc := s.Location()
	if loc.Source != "test.txt" {
		t.Errorf("expected Source 'test.txt', got %q", loc.Source)
	}
	if loc.Line != 0 || loc.Col != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", loc.Line, loc.Col)
	}
	s.Next()
	s.Next()
	s.Next()
	s.Next()
	s.Next()
	s.Next() // past newline
	loc = s.Location()
	if loc.Source != "test.txt" {
		t.Errorf("expected Source 'test.txt', got %q", loc.Source)
	}
	if loc.Line != 1 || loc.Col != 5 {
		t.Errorf("expected (1,5), got (%d,%d)", loc.Line, loc.Col)
	}
}

func TestLocationAt(t *testing.T) {
	s := NewStrCursorFromSource("src", "abc\ndef", 0)
	loc := s.LocationAt(4)
	if loc.Source != "src" {
		t.Errorf("expected 'src', got %q", loc.Source)
	}
	if loc.Line != 1 || loc.Col != 3 {
		t.Errorf("expected (1,3) at pos 4, got (%d,%d)", loc.Line, loc.Col)
	}
}

func TestInputSource(t *testing.T) {
	s := NewStrCursorFromSource("myfile.ebnf", "grammar", 0)
	if src := s.InputSource(); src != "myfile.ebnf" {
		t.Errorf("expected 'myfile.ebnf', got %q", src)
	}
}

func TestAsStr(t *testing.T) {
	s := NewStrCursor("some text")
	if str := s.AsStr(); str != "some text" {
		t.Errorf("expected 'some text', got %q", str)
	}
}

func TestAsRef(t *testing.T) {
	s := NewStrCursor("ref text")
	if ref := s.AsRef(); ref != "ref text" {
		t.Errorf("expected 'ref text', got %q", ref)
	}
}

func TestIgnoreCase(t *testing.T) {
	s := NewStrCursor("HELLO")
	if s.IgnoreCase() {
		t.Error("expected IgnoreCase false by default")
	}
}

func TestNameGuard(t *testing.T) {
	s := NewStrCursor("hello")
	if !s.NameGuard() {
		t.Error("expected NameGuard true by default")
	}
}

func TestLookahead(t *testing.T) {
	s := NewStrCursor("hello world\nnext line")
	la := s.Lookahead(0)
	if la != "hello world" {
		t.Errorf("expected 'hello world', got %q", la)
	}
	la = s.Lookahead(12)
	if la != "next line" {
		t.Errorf("expected 'next line', got %q", la)
	}
	la = s.Lookahead(100)
	if la != "" {
		t.Errorf("expected empty, got %q", la)
	}
}

func TestLookaheadNoNewline(t *testing.T) {
	s := NewStrCursor("single line")
	la := s.Lookahead(0)
	if la != "single line" {
		t.Errorf("expected 'single line', got %q", la)
	}
}

func TestNewStrCursorFromSourceStart(t *testing.T) {
	s := NewStrCursorFromSource("src", "hello world", 6)
	if s.offset != 6 {
		t.Errorf("expected offset 6, got %d", s.offset)
	}
	if s.InputSource() != "src" {
		t.Errorf("expected 'src', got %q", s.InputSource())
	}
}

func TestNewStrCursorFromSourceClamp(t *testing.T) {
	s := NewStrCursorFromSource("src", "hi", 100)
	if s.offset != 2 {
		t.Errorf("expected offset clamped to 2, got %d", s.offset)
	}
}

func TestSetPatternsAndNextToken(t *testing.T) {
	wsp, _ := pyre.Compile(`[ \t]+`)
	cmt, _ := pyre.Compile(`#[^\n]*`)
	eol, _ := pyre.Compile(`\r?\n`)
	s := NewStrCursor("  # comment \n  hello")
	s.SetPatterns(&TokenizingPatterns{Wsp: wsp, Cmt: cmt, Eol: eol})
	s.NextToken()
	if !s.MatchToken("hello") {
		t.Error("expected to match 'hello' after skipping whitespace and comment")
	}
}

func TestSetPatternsNoPatterns(t *testing.T) {
	s := NewStrCursor("hello")
	s.SetPatterns(nil)
	s.NextToken()
	if !s.MatchToken("hello") {
		t.Error("expected to match 'hello' with nil patterns")
	}
}

func TestReset(t *testing.T) {
	s := NewStrCursor("hello")
	s.offset = 3
	s.Reset(1)
	if s.offset != 1 {
		t.Errorf("expected offset 1 after reset, got %d", s.offset)
	}
}

func TestCloneCursor(t *testing.T) {
	s := NewStrCursor("hello world")
	s.offset = 6
	c := s.Clone()
	if c.Mark() != 6 {
		t.Errorf("expected cloned offset 6, got %d", c.Mark())
	}
	s.offset = 0
	if c.Mark() != 6 {
		t.Errorf("expected clone independent after original move, got %d", c.Mark())
	}
}

func TestPosEmpty(t *testing.T) {
	s := NewStrCursor("")
	line, col := s.Pos()
	if line != 0 || col != 0 {
		t.Errorf("expected (0,0) for empty string, got (%d,%d)", line, col)
	}
}
