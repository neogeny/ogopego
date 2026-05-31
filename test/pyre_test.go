package test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/util/pyre"
)

func TestMultilineDollarWithAnchor(t *testing.T) {
	p, err := pyre.Compile(`(?m)[ \t]*$`)
	assert.NoError(t, err)
	m, ok := p.Match("  \n")
	assert.True(t, ok, "should match trailing whitespace")
	g, ok := m.Group(0)
	assert.True(t, ok, "expected group 0")
	assert.Equal(t, "  ", g)
}

func TestTatsuEOLPatterns(t *testing.T) {
	p1, err := pyre.Compile(`(?m)[ \t]*$`)
	assert.NoError(t, err)
	p2, err := pyre.Compile(`(?m)(?:\r?\n|\r)?`)
	assert.NoError(t, err)

	text := "  \nNext line"

	m1, ok := p1.Match(text)
	assert.True(t, ok, "p1 should match")
	g1, _ := m1.Group(0)
	assert.Equal(t, "  ", g1)

	rest := text[m1.End():]
	m2, ok := p2.Match(rest)
	assert.True(t, ok, "p2 should match")
	g2, _ := m2.Group(0)
	assert.Equal(t, "\n", g2)
}

func TestMatchFirstIsStartPosition(t *testing.T) {
	p, err := pyre.Compile(`\d+`)
	assert.NoError(t, err)
	m, ok := p.Search("abc 123 def")
	assert.True(t, ok, "should find digits via Search")
	assert.Equal(t, 4, m.Start())
	assert.Equal(t, 7, m.End())
	g, ok := m.Group(0)
	assert.True(t, ok, "expected group 0")
	assert.Equal(t, "123", g)
}

func TestMatchRequiresStartAtZero(t *testing.T) {
	p, err := pyre.Compile(`\d+`)
	assert.NoError(t, err)
	_, ok := p.Match("x123")
	assert.False(t, ok, "Match should require match at position 0")
}

func TestMatchEndruleBlankLine(t *testing.T) {
	pat := `\s*[;]|(?:\s*(?:\r?\n|\r)){2,}[;]?`
	if pyre.LookaheadSupport {
		pat = `\s*[;]|(?=\s*(?:\r?\n|\r)\S)|(?:\s*(?:\r?\n|\r)){2,}[;]?`
	}
	p, err := pyre.Compile(pat)
	assert.NoError(t, err)
	_, ok := p.Match("\n\n")
	assert.True(t, ok, "ENDRULE should match blank line")
}
