// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/util/pyre"
)

func ptr[T any](v T) *T {
	return &v
}

func TestRuneCursorMatchPattern(t *testing.T) {
	s := NewRuneCursor("abc123def")
	m, ok := s.MatchPattern(`\d+`)
	assert.False(t, ok, "expected no match at start, got %q", m)
}

func TestRuneCursorMatchToken(t *testing.T) {
	s := NewRuneCursor("hello world")
	assert.True(t, s.MatchToken("hello"), "expected MatchToken to succeed")
	assert.Equal(t, 5, s.Mark())
}

func TestRuneCursorPeekToken(t *testing.T) {
	s := NewRuneCursor("hello")
	assert.True(t, s.PeekToken("hello"), "expected PeekToken to succeed")
	assert.Equal(t, 0, s.Mark(), "expected mark unchanged")
}

func TestRuneCursorNextPeek(t *testing.T) {
	s := NewRuneCursor("ab")
	r, ok := s.Peek()
	assert.True(t, ok)
	assert.Equal(t, 'a', r)
	r, ok = s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'a', r)
	assert.Equal(t, 1, s.Mark())
}

func TestRuneCursorAtEnd(t *testing.T) {
	s := NewRuneCursor("a")
	assert.False(t, s.AtEnd(), "expected not at end")
	s.Next()
	assert.True(t, s.AtEnd(), "expected at end")
}

func TestRuneCursorPos(t *testing.T) {
	s := NewRuneCursor("hello\nworld")
	line, col := s.Pos()
	assert.Equal(t, 1, line)
	assert.Equal(t, 1, col)
	for range 6 {
		s.Next()
	}
	line, col = s.Pos()
	assert.Equal(t, 2, line)
	assert.Equal(t, 1, col)
}

func TestRuneCursorClone(t *testing.T) {
	s := NewRuneCursor("hello")
	c := s.Clone()
	s.Next()
	assert.True(t, s.Mark() != c.Mark(), "expected clone to have independent mark")
}

func TestRuneMatchPatternSuccess(t *testing.T) {
	s := NewRuneCursor("abc123def")
	s.Next()
	s.Next()
	s.Next()
	m, ok := s.MatchPattern(`\d+`)
	assert.True(t, ok, "expected match")
	assert.Equal(t, "123", m)
	assert.Equal(t, 6, s.Mark())
}

func TestRuneMatchPatternWithGroup(t *testing.T) {
	s := NewRuneCursor("foo bar baz")
	m, ok := s.MatchPattern(`(foo)\s+(bar)`)
	assert.True(t, ok, "expected match")
	assert.Equal(t, "foo", m, "expected 'foo' from group 1")
}

func TestRuneMatchEOL(t *testing.T) {
	tests := []struct {
		input string
		want  bool
		mark  int
	}{
		{"\n", true, 1},
		{"\r\n", true, 2},
		{"\r", true, 1},
		{"abc", false, 0},
		{"", false, 0},
	}
	for _, tt := range tests {
		s := NewRuneCursor(tt.input)
		got := s.MatchEOL()
		assert.Equal(t, tt.want, got, "MatchEOL(%q)", tt.input)
		assert.Equal(t, tt.mark, s.Mark(), "MatchEOL(%q) mark", tt.input)
	}
}

func TestRuneMatchEOLWithWhitespace(t *testing.T) {
	p, err := pyre.Compile(`[ \t]+`)
	assert.NoError(t, err)
	s := NewRuneCursor("   \nnext")
	s.SetPatterns(&TokenizingPatterns{Wsp: p})
	assert.True(t, s.MatchEOL(), "expected MatchEOL to skip whitespace and match newline")
	assert.Equal(t, 4, s.Mark(), "expected mark 4 (whitespace+newline)")
}

func TestRunePosAt(t *testing.T) {
	s := NewRuneCursor("hello\nworld\nfoo")

	testPos := func(pos int, eline int, ecol int) string {
		line, col := s.PosAt(pos)
		if line != eline || col != ecol {
			return fmt.Sprintf(
				"Expected (%d, %d) at pos %d, got (%d, %d) ",
				eline, ecol, pos, line, col,
			)
		}
		return ""
	}

	assert.Equal(t, "", testPos(0, 1, 1))
	assert.Equal(t, "", testPos(5, 1, 5))
	assert.Equal(t, "", testPos(3, 1, 3))
	assert.Equal(t, "", testPos(6, 2, 1))
	assert.Equal(t, "", testPos(11, 2, 5))
	assert.Equal(t, "", testPos(12, 3, 1))
	assert.Equal(t, "", testPos(15, 3, 3))
}

func TestRunePosAtPastEnd(t *testing.T) {
	s := NewRuneCursor("hi")
	line, col := s.PosAt(100)
	assert.Equal(t, 1, line)
	assert.Equal(t, 2, col)
}

func TestRuneLocation(t *testing.T) {
	s := NewRuneCursorFromSource("test.txt", "hello\nworld", 0)
	loc := s.Location()
	assert.Equal(t, "test.txt", loc.Source)
	assert.Equal(t, 1, loc.Line)
	assert.Equal(t, 1, loc.Col)
	s.Next()
	s.Next()
	s.Next()
	s.Next()
	s.Next()
	s.Next()
	loc = s.Location()
	assert.Equal(t, "test.txt", loc.Source)
	assert.Equal(t, 2, loc.Line)
	assert.Equal(t, 1, loc.Col)
}

func TestRuneLocationAt(t *testing.T) {
	s := NewRuneCursorFromSource("src", "abc\ndef", 0)
	loc := s.LocationAt(4)
	assert.Equal(t, "src", loc.Source)
	assert.Equal(t, 2, loc.Line)
	assert.Equal(t, 1, loc.Col)
}

func TestRuneInputSource(t *testing.T) {
	s := NewRuneCursorFromSource("myfile.ebnf", "grammar", 0)
	assert.Equal(t, "myfile.ebnf", s.InputSource())
}

func TestRuneAsStr(t *testing.T) {
	s := NewRuneCursor("some text")
	assert.Equal(t, "some text", s.AsStr())
}

func TestRuneAsRef(t *testing.T) {
	s := NewRuneCursor("ref text")
	assert.Equal(t, "ref text", s.AsRef())
}

func TestRuneIgnoreCase(t *testing.T) {
	s := NewRuneCursor("HELLO")
	assert.False(t, s.IgnoreCase(), "expected IgnoreCase false by default")
}

func TestRuneNameGuard(t *testing.T) {
	s := NewRuneCursor("hello")
	assert.False(t, s.NameGuard(), "expected NameGuard false by default")
}

func TestRuneLookahead(t *testing.T) {
	s := NewRuneCursor("hello world\nnext line")
	la := s.Lookahead(0)
	assert.Equal(t, "hello world", la)
	la = s.Lookahead(12)
	assert.Equal(t, "next line", la)
	la = s.Lookahead(100)
	assert.Equal(t, "", la)
}

func TestRuneLookaheadNoNewline(t *testing.T) {
	s := NewRuneCursor("single line")
	la := s.Lookahead(0)
	assert.Equal(t, "single line", la)
}

func TestNewRuneCursorFromSourceStart(t *testing.T) {
	s := NewRuneCursorFromSource("src", "hello world", 6)
	assert.Equal(t, 6, s.Mark())
	assert.Equal(t, "src", s.InputSource())
}

func TestNewRuneCursorFromSourceClamp(t *testing.T) {
	s := NewRuneCursorFromSource("src", "hi", 100)
	assert.Equal(t, 2, s.Mark())
}

func TestRuneSetPatternsAndNextToken(t *testing.T) {
	wsp, _ := pyre.Compile(`[ \t]+`)
	cmt, _ := pyre.Compile(`#[^\n]*`)
	eol, _ := pyre.Compile(`\r?\n`)
	s := NewRuneCursor("  # comment \n  hello")
	s.SetPatterns(&TokenizingPatterns{Wsp: wsp, Cmt: cmt, Eol: eol})
	s.NextToken()
	assert.True(t, s.MatchToken("hello"), "expected to match 'hello' after skipping whitespace and comment")
}

func TestRuneSetPatternsNoPatterns(t *testing.T) {
	s := NewRuneCursor("hello")
	s.SetPatterns(nil)
	s.NextToken()
	assert.True(t, s.MatchToken("hello"), "expected to match 'hello' with nil patterns")
}

func TestRuneReset(t *testing.T) {
	s := NewRuneCursor("hello world")
	s.Next()
	s.Next()
	s.Next()
	s.Reset(1)
	assert.Equal(t, 1, s.Mark())
}

func TestRuneCloneCursor(t *testing.T) {
	s := NewRuneCursor("hello world")
	for range 6 {
		s.Next()
	}
	c := s.Clone()
	assert.Equal(t, 6, c.Mark())
	for range 3 {
		s.Next()
	}
	assert.Equal(t, 9, s.Mark())
	assert.Equal(t, 6, c.Mark(), "expected clone independent after original move")
}

func TestRunePosEmpty(t *testing.T) {
	s := NewRuneCursor("")
	line, col := s.Pos()
	assert.Equal(t, 1, line)
	assert.Equal(t, 1, col)
}

func TestRuneCursorGetPattern(t *testing.T) {
	s := NewRuneCursor("")
	p1 := s.GetPattern(`\d+`)
	assert.NotZero(t, p1, "expected non-nil pattern")
	p2 := s.GetPattern(`\d+`)
	assert.True(t, p2 == p1, "expected cached pattern to be same instance")
}

func TestRuneCursorGetPatternInvalid(t *testing.T) {
	s := NewRuneCursor("")
	p := s.GetPattern(`[invalid`)
	assert.Zero(t, p, "expected nil for invalid pattern")
}

func TestRuneCursorNextTokenSkipWhitespace(t *testing.T) {
	wsp, _ := pyre.Compile(`[ \t]+`)
	s := NewRuneCursor("   hello")
	s.SetPatterns(&TokenizingPatterns{Wsp: wsp})
	s.NextToken()
	assert.Equal(t, 3, s.Mark())
	assert.True(t, s.MatchToken("hello"))
}

func TestRuneCursorNextTokenSkipComment(t *testing.T) {
	cmt, _ := pyre.Compile(`//[^\n]*`)
	eol, _ := pyre.Compile(`\r?\n`)
	s := NewRuneCursor("// comment\nhello")
	s.SetPatterns(&TokenizingPatterns{Cmt: cmt, Eol: eol})
	s.NextToken()
	assert.True(t, s.MatchToken("hello"), "expected to match 'hello' after comment")
}

func TestRuneCursorMatchTokenNameGuard(t *testing.T) {
	s := NewRuneCursor("if_else")
	s.SetPatterns(&TokenizingPatterns{
		NonDefault: true,
		Wsp:        nil,
	})
	s.Configure(Cfg{NameGuard: ptr(true)})
	assert.False(t, s.MatchToken("if"), "expected name guard to prevent partial match")
	assert.Equal(t, 0, s.Mark(), "expected mark unchanged after failed match")
}

func TestRuneCursorMatchTokenNameGuardDisabled(t *testing.T) {
	s := NewRuneCursor("if_else")
	assert.True(t, s.MatchToken("if"), "expected MatchToken to succeed without name guard")
	assert.Equal(t, 2, s.Mark())
}

func TestRuneCursorPeekTokenIgnoreCase(t *testing.T) {
	s := NewRuneCursor("Hello")
	s.SetIgnoreCase(true)
	assert.True(t, s.PeekToken("hello"), "expected case-insensitive PeekToken")
}

func TestRuneCursorMatchTokenIgnoreCase(t *testing.T) {
	s := NewRuneCursor("Hello World")
	s.SetIgnoreCase(true)
	assert.True(t, s.MatchToken("hello"), "expected case-insensitive MatchToken")
	assert.Equal(t, 5, s.Mark())
}

func TestRuneCursorMatchEOLComment(t *testing.T) {
	cmt, _ := pyre.Compile(`//[^\n]*`)
	wsp, _ := pyre.Compile(`[ \t]+`)
	s := NewRuneCursor("// comment")
	s.SetPatterns(&TokenizingPatterns{Cmt: cmt, Wsp: wsp})
	assert.False(t, s.MatchEOL(), "expected no EOL match on comment-only line")
}

func TestRuneCursorLineCount(t *testing.T) {
	s := NewRuneCursor("hello\nworld\nfoo")
	assert.Equal(t, 3, s.LineCount())
}

func TestRuneCursorLineAt(t *testing.T) {
	s := NewRuneCursor("hello\nworld\nfoo")
	assert.Equal(t, "hello\n", s.LineAt(0))
	assert.Equal(t, "world\n", s.LineAt(1))
	assert.Equal(t, "foo", s.LineAt(2))
	assert.Equal(t, "", s.LineAt(3))
}

func TestRuneCursorLinesAt(t *testing.T) {
	s := NewRuneCursor("hello\nworld\nfoo")
	lines := s.LinesAt(0, 2)
	assert.Equal(t, 2, len(lines))
	assert.Equal(t, "hello\n", lines[0])
	assert.Equal(t, "world\n", lines[1])
}

func TestRuneCursorLinesAtNil(t *testing.T) {
	s := NewRuneCursor("hello\nworld")
	assert.Zero(t, s.LinesAt(2, 1))
	assert.Zero(t, s.LinesAt(-1, 1))
}

func TestRuneCursorNextRune(t *testing.T) {
	s := NewRuneCursor("abc")
	r, ok := s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'a', r)
	r, ok = s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'b', r)
	r, ok = s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'c', r)
	_, ok = s.Next()
	assert.False(t, ok)
	assert.True(t, s.AtEnd())
}

func TestRuneCursorPeekWithoutConsuming(t *testing.T) {
	s := NewRuneCursor("xy")
	r1, _ := s.Peek()
	r2, _ := s.Peek()
	assert.Equal(t, 'x', r1)
	assert.Equal(t, 'x', r2, "expected Peek to not advance")
	assert.Equal(t, 0, s.Mark())
}

func TestRuneCursorIsNameChar(t *testing.T) {
	s := NewRuneCursor("")
	assert.True(t, s.IsNameChar('_'))
	assert.True(t, s.IsNameChar('a'))
	assert.True(t, s.IsNameChar('Z'))
	assert.True(t, s.IsNameChar('0'))
	assert.False(t, s.IsNameChar('-'))
}

func TestRuneCursorIsName(t *testing.T) {
	s := NewRuneCursor("")
	assert.True(t, s.IsName("hello"))
	assert.True(t, s.IsName("_foo"))
	assert.False(t, s.IsName(""))
	assert.False(t, s.IsName("123abc"))
}

func TestRuneCursorLocationAt(t *testing.T) {
	s := NewRuneCursorFromSource("calc.ebnf", "hello\nworld", 0)
	loc := s.LocationAt(6)
	assert.Equal(t, "calc.ebnf", loc.Source)
	assert.Equal(t, 2, loc.Line)
	assert.Equal(t, 1, loc.Col)
}

func TestRuneCursorCloneIndependence(t *testing.T) {
	s := NewRuneCursor("hello world")
	c := s.Clone()
	s.Next()
	assert.Equal(t, 1, s.Mark())
	assert.Equal(t, 0, c.Mark(), "expected clone to be independent")
}

func TestRuneCursorConfigureIgnoreCase(t *testing.T) {
	s := NewRuneCursor("hello")
	s.Configure(Cfg{IgnoreCase: true})
	assert.True(t, s.IgnoreCase())
}

func TestRuneCursorConfigureNameGuard(t *testing.T) {
	s := NewRuneCursor("hello")
	s.Configure(Cfg{NameGuard: ptr(true)})
	assert.True(t, s.NameGuard())
}

func TestRuneCursorConfigureSource(t *testing.T) {
	s := NewRuneCursor("hello")
	s.Configure(Cfg{Source: "test.ebnf"})
	assert.Equal(t, "test.ebnf", s.InputSource())
}

func TestRuneCursorMatchEOLAfterReset(t *testing.T) {
	s := NewRuneCursor("hello\nworld")
	for range 6 {
		s.Next()
	}
	s.Reset(6)
	got := s.MatchEOL()
	assert.False(t, got, "expected MatchEOL to fail after reset to middle of line")
	assert.Equal(t, 6, s.Mark(), "expected mark unchanged after failed MatchEOL")
}

func TestRuneCursorNextWithNonASCII(t *testing.T) {
	s := NewRuneCursor("héllo")
	r, ok := s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'h', r)
	r, ok = s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'é', r)
	assert.Equal(t, 2, s.Mark())
	r, ok = s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'l', r)
}

func TestRuneCursorResetToRuneBoundary(t *testing.T) {
	s := NewRuneCursor("héllo")
	s.Next()
	s.Next()
	s.Reset(1)
	assert.Equal(t, 1, s.Mark())
	r, ok := s.Peek()
	assert.True(t, ok)
	assert.Equal(t, 'é', r)
}

func TestRuneCursorLineCountEmpty(t *testing.T) {
	s := NewRuneCursor("")
	assert.Equal(t, 0, s.LineCount())
}

func TestRuneCursorLineCountTrailingNewline(t *testing.T) {
	s := NewRuneCursor("hello\n")
	assert.Equal(t, 1, s.LineCount())
}

func TestRuneCursorLineAtTrailingNewline(t *testing.T) {
	s := NewRuneCursor("hello\n")
	assert.Equal(t, "hello\n", s.LineAt(0))
	assert.Equal(t, "", s.LineAt(1))
}

func TestRuneCursorLinesAtTrailingNewline(t *testing.T) {
	s := NewRuneCursor("hello\n")
	lines := s.LinesAt(0, 2)
	assert.Equal(t, 1, len(lines))
	assert.Equal(t, "hello\n", lines[0])
}

func TestRuneCursorLinesAtOutOfRange(t *testing.T) {
	s := NewRuneCursor("hello\nworld\nfoo")
	assert.Zero(t, s.LinesAt(3, 5))
}

func TestRuneCursorPosAtZero(t *testing.T) {
	s := NewRuneCursor("hello")
	line, col := s.PosAt(0)
	assert.Equal(t, 1, line)
	assert.Equal(t, 1, col)
}

func TestRuneCursorPosAtEnd(t *testing.T) {
	s := NewRuneCursor("hi")
	line, col := s.PosAt(2)
	assert.Equal(t, 1, line)
	assert.Equal(t, 2, col)
}

func TestRuneCursorMatchEOLWithNewlineOnly(t *testing.T) {
	s := NewRuneCursor("a\nb")
	s.Next() // past 'a'
	assert.True(t, s.MatchEOL(), "expected EOL after 'a'")
	assert.Equal(t, 2, s.Mark())
	r, ok := s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'b', r)
}

func TestRuneCursorMatchEOLAfterTrailingSpaces(t *testing.T) {
	p, _ := pyre.Compile(`[ \t]+`)
	s := NewRuneCursor("a   \nb")
	s.Next() // past 'a'
	s.SetPatterns(&TokenizingPatterns{Wsp: p})
	assert.True(t, s.MatchEOL(), "expected EOL to skip trailing spaces")
	assert.Equal(t, 5, s.Mark())
}

func TestRuneCursorPeekTokenNonASCII(t *testing.T) {
	s := NewRuneCursor("héllo")
	assert.True(t, s.PeekToken("hé"), "expected PeekToken to match non-ASCII")
}

func TestRuneCursorMatchTokenNonASCII(t *testing.T) {
	s := NewRuneCursor("héllo world")
	assert.True(t, s.MatchToken("hé"), "expected MatchToken to match non-ASCII")
	assert.Equal(t, 2, s.Mark())
}
