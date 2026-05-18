// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

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
		NamedBox: NamedBox{
			Box:  Box{Exp: &Token{Token: "x"}},
			Name: "atom",
		},
	}
	expr := &Rule{
		NamedBox: NamedBox{
			Box: Box{Exp: &Sequence{
				Sequence: []Model{
					&Call{Name: "atom"},
					&Closure{Box: Box{Exp: &Sequence{
						Sequence: []Model{
							&Token{Token: "+"},
							&Call{Name: "atom"},
						},
					}}},
				},
			}},
			Name: "expr",
		},
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
		NamedBox: NamedBox{
			Box:  Box{Exp: &Token{Token: "a"}},
			Name: "rule_a",
		},
	}
	ruleB := &Rule{
		NamedBox: NamedBox{
			Box:  Box{Exp: &RuleInclude{Name: "rule_a"}},
			Name: "rule_b",
		},
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
		NamedBox: NamedBox{
			Box:  Box{Exp: &Token{Token: "p"}},
			Name: "prim",
		},
	}
	deep := &Rule{
		NamedBox: NamedBox{
			Box: Box{Exp: &Choice{
				Options: []*Option{
					{Box: Box{Exp: &Call{Name: "prim"}}},
					{Box: Box{Exp: &Group{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &Override{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &OverrideList{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &Lookahead{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &NegativeLookahead{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &Closure{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &PositiveClosure{Closure: Closure{Box: Box{Exp: &Call{Name: "prim"}}}}}},
					{Box: Box{Exp: &Optional{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &SkipGroup{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &SkipTo{Box: Box{Exp: &Call{Name: "prim"}}}}},
					{Box: Box{Exp: &Named{NamedBox: NamedBox{Box: Box{Exp: &Call{Name: "prim"}}, Name: "n"}}}},
					{Box: Box{Exp: &NamedList{Named: Named{NamedBox: NamedBox{Box: Box{Exp: &Call{Name: "prim"}}, Name: "nl"}}}}},
					{Box: Box{Exp: &Join{Box: Box{Exp: &Call{Name: "prim"}}, Sep: &Token{Token: ","}}}},
					{Box: Box{Exp: &PositiveJoin{Join: Join{Box: Box{Exp: &Call{Name: "prim"}}, Sep: &Token{Token: ","}}}}},
					{Box: Box{Exp: &Gather{Join: Join{Box: Box{Exp: &Call{Name: "prim"}}, Sep: &Token{Token: ","}}}}},
					{Box: Box{Exp: &PositiveGather{Gather: Gather{Join: Join{Box: Box{Exp: &Call{Name: "prim"}}, Sep: &Token{Token: ","}}}}}},
				},
			}},
			Name: "deep",
		},
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
