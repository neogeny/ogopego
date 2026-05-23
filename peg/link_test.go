// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"testing"
)

func mustLink(t *testing.T, g *Grammar) {
	t.Helper()
	if err := g.LinkGrammar(); err != nil {
		t.Fatalf("Link() error: %v", err)
	}
}

func mustValidateLinked(t *testing.T, g *Grammar) {
	t.Helper()
	if err := g.ValidateLinked(); err != nil {
		t.Fatalf("ValidateLinked() error: %v", err)
	}
}

func expectLinkError(t *testing.T, g *Grammar) {
	t.Helper()
	if err := g.LinkGrammar(); err == nil {
		t.Fatal("expected Link() error, got nil")
	}
}

func expectValidateError(t *testing.T, g *Grammar) {
	t.Helper()
	if err := g.ValidateLinked(); err == nil {
		t.Fatal("expected ValidateLinked() error, got nil")
	}
}

// linkTestGrammar builds a simple grammar:
//
//	expr = atom ('+' atom)*
//	atom = 'x'
//
// Uses Call to reference atom from expr.
func linkTestGrammar() *Grammar {
	atom := &Rule{
		Exp:  &Token{Token: "x"},
		Name: "atom",
	}
	expr := &Rule{
		Exp: &Sequence{
			Sequence: []Model{
				&Call{Name: "atom"},
				&Closure{Exp: &Sequence{
					Sequence: []Model{
						&Token{Token: "+"},
						&Call{Name: "atom"},
					},
				}},
			},
		},
		Name: "expr",
	}
	return &Grammar{Name: "test", Rules: []*Rule{expr, atom}}
}

func TestLinkCall(t *testing.T) {
	g := linkTestGrammar()
	mustLink(t, g)
	mustValidateLinked(t, g)

	seq, ok := g.Rules[0].Exp.(*Sequence)
	if !ok {
		t.Fatal("expr.Exp is not *Sequence")
	}
	call, ok := seq.Sequence[0].(*Call)
	if !ok {
		t.Fatal("seq[0] is not *Call")
	}
	if call.rule == nil {
		t.Fatal("Call.rule is nil after Link")
	}
	if call.rule.Name != "atom" {
		t.Fatalf("Call.rule.Name = %q, want %q", call.rule.Name, "atom")
	}
}

func TestLinkCallUnlinked(t *testing.T) {
	g := linkTestGrammar()
	expectValidateError(t, g)
}

func TestLinkCallUndefined(t *testing.T) {
	g := linkTestGrammar()
	g.Rules[1].Exp = &Sequence{
		Sequence: []Model{
			&Call{Name: "nonexistent"},
		},
	}
	expectLinkError(t, g)
}

// includeTestGrammar builds a grammar with RuleInclude:
//
//	rule_a = 'a'
//	rule_b = >rule_a   (includes rule_a)
func includeTestGrammar() *Grammar {
	ruleA := &Rule{
		Exp:  &Token{Token: "a"},
		Name: "rule_a",
	}
	ruleB := &Rule{
		Exp:  &RuleInclude{Name: "rule_a"},
		Name: "rule_b",
	}
	return &Grammar{Name: "test", Rules: []*Rule{ruleB, ruleA}}
}

func TestLinkRuleInclude(t *testing.T) {
	g := includeTestGrammar()
	mustLink(t, g)
	mustValidateLinked(t, g)

	ri, ok := g.Rules[0].Exp.(*RuleInclude)
	if !ok {
		t.Fatal("rule_b.Exp is not *RuleInclude")
	}
	if ri.exp == nil {
		t.Fatal("RuleInclude.exp is nil after Link")
	}
	if _, ok := ri.exp.(*Token); !ok {
		t.Fatalf("RuleInclude.exp is %T, want *Token", ri.exp)
	}
	if ri.exp.(*Token).Token != "a" {
		t.Fatalf("RuleInclude.exp.Token = %q, want %q", ri.exp.(*Token).Token, "a")
	}
}

func TestLinkRuleIncludeUnlinked(t *testing.T) {
	g := includeTestGrammar()
	expectValidateError(t, g)
}

func TestLinkRuleIncludeUndefined(t *testing.T) {
	g := includeTestGrammar()
	g.Rules[0].Exp = &RuleInclude{Name: "nonexistent"}
	expectLinkError(t, g)
}

// deepTestGrammar builds a deeply nested grammar exercising all Link recursion paths.
func deepTestGrammar() *Grammar {
	prim := &Rule{
		Exp:  &Token{Token: "p"},
		Name: "prim",
	}
	deep := &Rule{
		Exp: &Choice{
			Options: []*Option{
				{Exp: &Call{Name: "prim"}},
				{Exp: &Group{Exp: &Call{Name: "prim"}}},
				{Exp: &Override{Exp: &Call{Name: "prim"}}},
				{Exp: &OverrideList{Exp: &Call{Name: "prim"}}},
				{Exp: &Lookahead{Exp: &Call{Name: "prim"}}},
				{Exp: &NegativeLookahead{Exp: &Call{Name: "prim"}}},
				{Exp: &Closure{Exp: &Call{Name: "prim"}}},
				{Exp: &PositiveClosure{Exp: &Call{Name: "prim"}}},
				{Exp: &Optional{Exp: &Call{Name: "prim"}}},
				{Exp: &SkipGroup{Exp: &Call{Name: "prim"}}},
				{Exp: &SkipTo{Exp: &Call{Name: "prim"}}},
				{Exp: &Named{Exp: &Call{Name: "prim"}, Name: "n"}},
				{Exp: &NamedList{Exp: &Call{Name: "prim"}, Name: "nl"}},
				{Exp: &Join{Exp: &Call{Name: "prim"}, Sep: &Token{Token: ","}}},
				{Exp: &PositiveJoin{Exp: &Call{Name: "prim"}, Sep: &Token{Token: ","}}},
				{Exp: &Gather{Exp: &Call{Name: "prim"}, Sep: &Token{Token: ","}}},
				{Exp: &PositiveGather{Exp: &Call{Name: "prim"}, Sep: &Token{Token: ","}}},
			},
		},
		Name: "deep",
	}
	return &Grammar{Name: "test", Rules: []*Rule{deep, prim}}
}

func TestLinkDeep(t *testing.T) {
	g := deepTestGrammar()
	mustLink(t, g)
	mustValidateLinked(t, g)
}

func TestValidateCatchesSingleUnlinked(t *testing.T) {
	g := linkTestGrammar()
	// Remove atom rule so expr's Call to "atom" remains unlinked
	g.Rules = g.Rules[:1]
	expectLinkError(t, g) // undefined Call
}

func TestGrammarValidateOk(t *testing.T) {
	g := linkTestGrammar()
	mustLink(t, g)
	mustValidateLinked(t, g)
}
