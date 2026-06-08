// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/util/pyre"
)

func TestStrCursorMatchPattern(t *testing.T) {
	s := NewStrCursor("abc123def")
	m, ok := s.MatchPattern(`\d+`)
	assert.False(t, ok, "expected no match at start, got %q", m)
}

func TestStrCursorMatchToken(t *testing.T) {
	s := NewStrCursor("hello world")
	assert.True(t, s.MatchToken("hello"), "expected MatchToken to succeed")
	assert.Equal(t, 5, s.offset)
}

func TestStrCursorPeekToken(t *testing.T) {
	s := NewStrCursor("hello")
	assert.True(t, s.PeekToken("hello"), "expected PeekToken to succeed")
	assert.Equal(t, 0, s.offset, "expected offset unchanged")
}

func TestStrCursorNextPeek(t *testing.T) {
	s := NewStrCursor("ab")
	r, ok := s.Peek()
	assert.True(t, ok)
	assert.Equal(t, 'a', r)
	r, ok = s.Next()
	assert.True(t, ok)
	assert.Equal(t, 'a', r)
	assert.Equal(t, 1, s.offset)
}

func TestStrCursorAtEnd(t *testing.T) {
	s := NewStrCursor("a")
	assert.False(t, s.AtEnd(), "expected not at end")
	s.Next()
	assert.True(t, s.AtEnd(), "expected at end")
}

func TestStrCursorPos(t *testing.T) {
	s := NewStrCursor("hello\nworld")
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

func TestStrCursorClone(t *testing.T) {
	s := NewStrCursor("hello")
	c := s.Clone()
	s.Next()
	assert.True(t, s.Mark() != c.Mark(), "expected clone to have independent offset")
}

func TestMatchPatternSuccess(t *testing.T) {
	s := NewStrCursor("abc123def")
	// skip past "abc"
	s.Next()
	s.Next()
	s.Next()
	m, ok := s.MatchPattern(`\d+`)
	assert.True(t, ok, "expected match")
	assert.Equal(t, "123", m)
	assert.Equal(t, 6, s.offset)
}

func TestMatchPatternWithGroup(t *testing.T) {
	s := NewStrCursor("foo bar baz")
	m, ok := s.MatchPattern(`(foo)\s+(bar)`)
	assert.True(t, ok, "expected match")
	assert.Equal(t, "foo", m, "expected 'foo' from group 1")
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
		assert.Equal(t, tt.want, got, "MatchEOL(%q)", tt.input)
		assert.Equal(t, tt.off, s.offset, "MatchEOL(%q) offset", tt.input)
	}
}

func TestMatchEOLWithWhitespace(t *testing.T) {
	p, err := pyre.Compile(`[ \t]+`)
	assert.NoError(t, err)
	s := NewStrCursor("   \nnext")
	s.SetPatterns(&TokenizingPatterns{Wsp: p})
	assert.True(t, s.MatchEOL(), "expected MatchEOL to skip whitespace and match newline")
	assert.Equal(t, 4, s.offset, "expected offset 4 (whitespace+newline)")
}

func TestPosAt(t *testing.T) {
	s := NewStrCursor("hello\nworld\nfoo")

	testPos := func(pos int, eline int, ecol int) string {
		line, col := s.PosAt(pos)
		if line != eline || col != ecol {
			return fmt.Sprintf(
				"Expexted (%d, %d) at pos %d, got (%d, %d) ",
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

func TestPosAtPastEnd(t *testing.T) {
	s := NewStrCursor("hi")
	line, col := s.PosAt(100)
	assert.Equal(t, 1, line)
	assert.Equal(t, 2, col)
}

func TestLocation(t *testing.T) {
	s := NewStrCursorFromSource("test.txt", "hello\nworld", 0)
	loc := s.Location()
	assert.Equal(t, "test.txt", loc.Source)
	assert.Equal(t, 1, loc.Line)
	assert.Equal(t, 1, loc.Col)
	s.Next()
	s.Next()
	s.Next()
	s.Next()
	s.Next()
	s.Next() // past newline
	loc = s.Location()
	assert.Equal(t, "test.txt", loc.Source)
	assert.Equal(t, 2, loc.Line)
	assert.Equal(t, 1, loc.Col)
}

func TestLocationAt(t *testing.T) {
	s := NewStrCursorFromSource("src", "abc\ndef", 0)
	loc := s.LocationAt(4)
	assert.Equal(t, "src", loc.Source)
	assert.Equal(t, 2, loc.Line)
	assert.Equal(t, 1, loc.Col)
}

func TestInputSource(t *testing.T) {
	s := NewStrCursorFromSource("myfile.ebnf", "grammar", 0)
	assert.Equal(t, "myfile.ebnf", s.InputSource())
}

func TestAsStr(t *testing.T) {
	s := NewStrCursor("some text")
	assert.Equal(t, "some text", s.AsStr())
}

func TestAsRef(t *testing.T) {
	s := NewStrCursor("ref text")
	assert.Equal(t, "ref text", s.AsRef())
}

func TestIgnoreCase(t *testing.T) {
	s := NewStrCursor("HELLO")
	assert.False(t, s.IgnoreCase(), "expected IgnoreCase false by default")
}

func TestNameGuard(t *testing.T) {
	s := NewStrCursor("hello")
	assert.False(t, s.NameGuard(), "expected NameGuard false by default")
}

func TestLookahead(t *testing.T) {
	s := NewStrCursor("hello world\nnext line")
	la := s.Lookahead(0)
	assert.Equal(t, "hello world", la)
	la = s.Lookahead(12)
	assert.Equal(t, "next line", la)
	la = s.Lookahead(100)
	assert.Equal(t, "", la)
}

func TestLookaheadNoNewline(t *testing.T) {
	s := NewStrCursor("single line")
	la := s.Lookahead(0)
	assert.Equal(t, "single line", la)
}

func TestNewStrCursorFromSourceStart(t *testing.T) {
	s := NewStrCursorFromSource("src", "hello world", 6)
	assert.Equal(t, 6, s.offset)
	assert.Equal(t, "src", s.InputSource())
}

func TestNewStrCursorFromSourceClamp(t *testing.T) {
	s := NewStrCursorFromSource("src", "hi", 100)
	assert.Equal(t, 2, s.offset)
}

func TestSetPatternsAndNextToken(t *testing.T) {
	wsp, _ := pyre.Compile(`[ \t]+`)
	cmt, _ := pyre.Compile(`#[^\n]*`)
	eol, _ := pyre.Compile(`\r?\n`)
	s := NewStrCursor("  # comment \n  hello")
	s.SetPatterns(&TokenizingPatterns{Wsp: wsp, Cmt: cmt, Eol: eol})
	s.NextToken()
	assert.True(t, s.MatchToken("hello"), "expected to match 'hello' after skipping whitespace and comment")
}

func TestSetPatternsNoPatterns(t *testing.T) {
	s := NewStrCursor("hello")
	s.SetPatterns(nil)
	s.NextToken()
	assert.True(t, s.MatchToken("hello"), "expected to match 'hello' with nil patterns")
}

func TestReset(t *testing.T) {
	s := NewStrCursor("hello")
	s.offset = 3
	s.Reset(1)
	assert.Equal(t, 1, s.offset)
}

func TestCloneCursor(t *testing.T) {
	s := NewStrCursor("hello world")
	s.offset = 6
	c := s.Clone()
	assert.Equal(t, 6, c.Mark())
	s.offset = 0
	assert.Equal(t, 6, c.Mark(), "expected clone independent after original move")
}

func TestPosEmpty(t *testing.T) {
	s := NewStrCursor("")
	line, col := s.Pos()
	assert.Equal(t, 1, line)
	assert.Equal(t, 1, col)
}

func TestStrCursorGetPattern(t *testing.T) {
	s := NewStrCursor("")
	p1 := s.GetPattern(`\d+`)
	assert.NotZero(t, p1, "expected non-nil pattern")
	p2 := s.GetPattern(`\d+`)
	assert.True(t, p2 == p1, "expected cached pattern to be same instance")
}

func TestStrCursorGetPatternInvalid(t *testing.T) {
	s := NewStrCursor("")
	p := s.GetPattern(`[invalid`)
	assert.Zero(t, p, "expected nil for invalid pattern")
}

// Meta expression tests - ported from x/tatsu/tests/syntax/meta_test.py
// These test at the cursor unit level rather than the full grammar level.

func TestStrCursorMatchName(t *testing.T) {
	s := NewStrCursor("hello")
	n, ok := s.MatchName()
	assert.True(t, ok)
	assert.Equal(t, "hello", n)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchNameUnderscoreStart(t *testing.T) {
	s := NewStrCursor("_hello")
	n, ok := s.MatchName()
	assert.True(t, ok)
	assert.Equal(t, "_hello", n)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchNameRejectsDigitStart(t *testing.T) {
	s := NewStrCursor("1hello")
	_, ok := s.MatchName()
	assert.False(t, ok)
	assert.Equal(t, 0, s.Mark())
}

func TestStrCursorMatchNameRejectsEmpty(t *testing.T) {
	s := NewStrCursor("")
	_, ok := s.MatchName()
	assert.False(t, ok)
}

func TestStrCursorMatchIntSigned(t *testing.T) {
	s := NewStrCursor("+42")
	n, ok := s.MatchInt()
	assert.True(t, ok)
	assert.Equal(t, 42, n)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchIntNegative(t *testing.T) {
	s := NewStrCursor("-7")
	n, ok := s.MatchInt()
	assert.True(t, ok)
	assert.Equal(t, -7, n)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchIntUnsigned(t *testing.T) {
	s := NewStrCursor("99")
	n, ok := s.MatchInt()
	assert.True(t, ok)
	assert.Equal(t, 99, n)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchIntRejectsEmpty(t *testing.T) {
	s := NewStrCursor("")
	_, ok := s.MatchInt()
	assert.False(t, ok)
}

func TestStrCursorMatchUIntUnsigned(t *testing.T) {
	s := NewStrCursor("42")
	n, ok := s.MatchUInt()
	assert.True(t, ok)
	assert.Equal(t, uint64(42), n)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchUIntRejectsNegative(t *testing.T) {
	s := NewStrCursor("-5")
	_, ok := s.MatchUInt()
	assert.False(t, ok)
	assert.Equal(t, 0, s.Mark())
}

func TestStrCursorMatchUIntRejectsPlus(t *testing.T) {
	s := NewStrCursor("+5")
	_, ok := s.MatchUInt()
	assert.False(t, ok)
	assert.Equal(t, 0, s.Mark())
}

func TestStrCursorMatchUIntUnderscores(t *testing.T) {
	s := NewStrCursor("1_000_000")
	n, ok := s.MatchUInt()
	assert.True(t, ok)
	assert.Equal(t, uint64(1000000), n)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchFloatPi(t *testing.T) {
	piStr := strconv.FormatFloat(math.Pi, 'g', -1, 64)
	s := NewStrCursor(piStr)
	f, ok := s.MatchFloat()
	assert.True(t, ok)
	assert.Equal(t, math.Pi, f)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchFloatNegative(t *testing.T) {
	s := NewStrCursor("-2.5")
	f, ok := s.MatchFloat()
	assert.True(t, ok)
	assert.Equal(t, -2.5, f)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchFloatRejectsText(t *testing.T) {
	s := NewStrCursor("abc")
	_, ok := s.MatchFloat()
	assert.False(t, ok)
	assert.Equal(t, 0, s.Mark())
}

func TestStrCursorMatchBoolTrue(t *testing.T) {
	s := NewStrCursor("true")
	b, ok := s.MatchBool()
	assert.True(t, ok)
	assert.Equal(t, true, b)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchBoolCapitalizedTrue(t *testing.T) {
	s := NewStrCursor("True")
	b, ok := s.MatchBool()
	assert.True(t, ok)
	assert.Equal(t, true, b)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchBoolFalse(t *testing.T) {
	s := NewStrCursor("false")
	b, ok := s.MatchBool()
	assert.True(t, ok)
	assert.Equal(t, false, b)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchBoolCapitalizedFalse(t *testing.T) {
	s := NewStrCursor("False")
	b, ok := s.MatchBool()
	assert.True(t, ok)
	assert.Equal(t, false, b)
	assert.True(t, s.AtEnd())
}

func TestStrCursorMatchBoolRejectsEmpty(t *testing.T) {
	s := NewStrCursor("")
	_, ok := s.MatchBool()
	assert.False(t, ok)
}
