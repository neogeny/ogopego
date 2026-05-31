// Disabled tests: these require regexp lookaheads (regexp2 backend) and are
// skipped when using the default Go regexp backend.
package pyre

import (
	"testing"
)

func TestMatchZeroWidthLookahead(t *testing.T) {
	if !LookaheadSupport {
		t.Skip("requires regexp lookaheads")
	}
	p, err := Compile(`(?=\s*(?:\r?\n|\r)\S)`)
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
	if !LookaheadSupport {
		t.Skip("requires regexp lookaheads")
	}
	p, err := Compile(`\s*[;]|(?=\s*(?:\r?\n|\r)\S)|(?:\s*(?:\r?\n|\r)){2,}[;]?`)
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

func TestMatchEndruleCRLF(t *testing.T) {
	if !LookaheadSupport {
		t.Skip("requires regexp lookaheads")
	}
	p, err := Compile(`\s*[;]|(?=\s*(?:\r?\n|\r)\S)|(?:\s*(?:\r?\n|\r)){2,}[;?]`)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := p.Match("\r\nnext_rule")
	if !ok {
		t.Fatal("ENDRULE should match CRLF before an unindented next rule")
	}
}

func TestEatPatternZeroWidth(t *testing.T) {
	if !LookaheadSupport {
		t.Skip("requires regexp lookaheads")
	}
	p, err := Compile(`(?=\s*(?:\r?\n|\r)\S)`)
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
