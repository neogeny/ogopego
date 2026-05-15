package ogopego_test

import (
	"testing"

	"github.com/neogeny/ogopego/util/pyre"
)

func TestMultilineDollarWithAnchor(t *testing.T) {
	p, err := pyre.Compile(`(?m)[ \t]*$`)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := p.Match("  \n")
	if !ok {
		t.Fatal("should match trailing whitespace")
	}
	if g, ok := m.Group(0); !ok || g != "  " {
		t.Errorf("expected '  ', got %q", g)
	}
}

func TestTatsuEOLPatterns(t *testing.T) {
	p1, err := pyre.Compile(`(?m)[ \t]*$`)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := pyre.Compile(`(?m)(?:\r?\n|\r)?`)
	if err != nil {
		t.Fatal(err)
	}

	text := "  \nNext line"

	m1, ok := p1.Match(text)
	if !ok {
		t.Fatal("p1 should match")
	}
	g1, _ := m1.Group(0)
	if g1 != "  " {
		t.Errorf("expected '  ', got %q", g1)
	}

	rest := text[m1.End():]
	m2, ok := p2.Match(rest)
	if !ok {
		t.Fatal("p2 should match")
	}
	g2, _ := m2.Group(0)
	if g2 != "\n" {
		t.Errorf("expected '\\n', got %q", g2)
	}
}

func TestMatchZeroWidthLookahead(t *testing.T) {
	p, err := pyre.Compile(`(?=\s*(?:\r?\n|\r)\S)`)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := p.Match("\nnext")
	if !ok {
		t.Fatal("lookahead should match at start")
	}
	if m.Start() != 0 || m.End() != 0 {
		t.Errorf("expected 0-width match at start, got start=%d end=%d", m.Start(), m.End())
	}
}

func TestMatchEndruleUnindentedBranch(t *testing.T) {
	p, err := pyre.Compile(`\s*[;]|(?=\s*(?:\r?\n|\r)\S)|(?:\s*(?:\r?\n|\r)){2,}[;]?`)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := p.Match("\nnext_rule")
	if !ok {
		t.Fatal("ENDRULE should match before an unindented next rule")
	}
	if m.Start() != 0 || m.End() != 0 {
		t.Errorf("expected 0-width at start, got start=%d end=%d", m.Start(), m.End())
	}
}

func TestMatchEndruleBlankLine(t *testing.T) {
	p, err := pyre.Compile(`\s*[;]|(?=\s*(?:\r?\n|\r)\S)|(?:\s*(?:\r?\n|\r)){2,}[;]?`)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := p.Match("\n\n")
	if !ok {
		t.Fatal("ENDRULE should match blank line")
	}
}

func TestMatchEndruleCRLF(t *testing.T) {
	p, err := pyre.Compile(`\s*[;]|(?=\s*(?:\r?\n|\r)\S)|(?:\s*(?:\r?\n|\r)){2,}[;?]`)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := p.Match("\r\nnext_rule")
	if !ok {
		t.Fatal("ENDRULE should match CRLF before an unindented next rule")
	}
}

func TestMatchFirstIsStartPosition(t *testing.T) {
	p, err := pyre.Compile(`\d+`)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := p.Search("abc 123 def")
	if !ok {
		t.Fatal("should find digits via Search")
	}
	if m.Start() != 4 || m.End() != 7 {
		t.Errorf("expected start=4 end=7, got start=%d end=%d", m.Start(), m.End())
	}
	if g, ok := m.Group(0); !ok || g != "123" {
		t.Errorf("expected '123', got %q", g)
	}
}

func TestMatchRequiresStartAtZero(t *testing.T) {
	p, err := pyre.Compile(`\d+`)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := p.Match("x123")
	if ok {
		t.Fatal("Match should require match at position 0")
	}
}

func TestEatPatternZeroWidth(t *testing.T) {
	p, err := pyre.Compile(`(?=\s*(?:\r?\n|\r)\S)`)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := p.Match("\nrule")
	if !ok {
		t.Fatal("lookahead should match")
	}
	if m.End() != 0 {
		t.Errorf("expected zero-width match, got end=%d", m.End())
	}
}
