package peg

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestOptimizePreservesOriginalRule(t *testing.T) {
	inner := &Token{Token: "hello"}
	group := &Group{Exp: inner}
	rule := &Rule{
		Name: "start",
		Exp:  group,
	}
	g := &Grammar{
		Name:  "test",
		Rules: []*Rule{rule},
	}

	orig := rule.Exp

	_ = g.Optimized()

	assert.Equal(t, orig, rule.Exp,
		"Optimize should not mutate the original rule's Exp")
}

func TestOptimizeReturnsOptimizedGrammar(t *testing.T) {
	inner := &Token{Token: "hello"}
	group := &Group{Exp: inner}
	rule := &Rule{
		Name: "start",
		Exp:  group,
	}
	g := &Grammar{
		Name:  "test",
		Rules: []*Rule{rule},
	}

	got := g.Optimized()

	gotRule := got.Rules[0]
	_, isGroup := gotRule.Exp.(*Group)
	assert.False(t, isGroup, "Group should be eliminated in the returned grammar")
	assert.Equal(t, inner, gotRule.Exp.(*Token))
}
