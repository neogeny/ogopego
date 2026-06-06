// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import "fmt"

// optimizeExpr recursively simplifies a grammar model tree.
// It returns the optimized model, which may be a different concrete type
// (e.g., *Group is eliminated, returning its inner expression).
// The original tree is not modified.
func optimizeExpr(m Model) Model {
	switch e := m.(type) {
	// --- Leaves ---
	case *Dot, *Cut, *Void, *Fail, *EOF, *EOL,
		*Token, *Pattern, *Constant, *Alert,
		*EmptyClosure, *NULL, *Call, *RuleInclude:
		return e

	// --- Unary containers: clone and recurse into Exp ---
	case *Optional:
		return &Optional{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *Closure:
		return &Closure{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *PositiveClosure:
		return &PositiveClosure{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *Lookahead:
		return &Lookahead{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *NegativeLookahead:
		return &NegativeLookahead{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *SkipGroup:
		return &SkipGroup{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *SkipTo:
		return &SkipTo{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *Override:
		return &Override{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *OverrideList:
		return &OverrideList{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *Named:
		return &Named{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Name: e.Name}
	case *NamedList:
		return &NamedList{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Name: e.Name}
	case *Option:
		return &Option{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}
	case *Synth:
		return optimizeExpr(e.Exp)

	// --- Binary containers: clone and recurse into both ---
	case *Join:
		return &Join{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Sep: optimizeExpr(e.Sep)}
	case *PositiveJoin:
		return &PositiveJoin{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Sep: optimizeExpr(e.Sep)}
	case *Gather:
		return &Gather{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Sep: optimizeExpr(e.Sep)}
	case *PositiveGather:
		return &PositiveGather{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Sep: optimizeExpr(e.Sep)}

	// --- Eliminated: Group unwraps to its inner expression ---
	case *Group:
		return optimizeExpr(e.Exp)

	// --- Collections: clone and recurse into children ---
	case *Choice:
		opts := make([]*Option, len(e.Options))
		for i, o := range e.Options {
			opts[i] = &Option{ModelBase: o.ModelBase, Exp: optimizeExpr(o.Exp)}
		}
		if len(opts) == 1 {
			return opts[0].Exp
		}
		return &Choice{ModelBase: e.ModelBase, Options: opts}

	case *Sequence:
		seq := make([]Model, len(e.Sequence))
		for i, s := range e.Sequence {
			seq[i] = optimizeExpr(s)
		}
		if len(seq) == 1 {
			return seq[0]
		}
		return &Sequence{ModelBase: e.ModelBase, Sequence: seq}

	default:
		panic(fmt.Sprintf("optimizeExpr: unhandled model type %T", m))
	}
}

// Optimized simplifies the rule's expression tree and returns a new rule.
// The original rule is not modified.
func (r *Rule) Optimized() *Rule {
	r2 := *r
	r2.Exp = optimizeExpr(r.Exp)
	return &r2
}

// Optimized simplifies all rules in the grammar and returns a new grammar.
// The original grammar is not modified.
func (g *Grammar) Optimized() (*Grammar, error) {
	g2 := *g
	g2.Rules = make([]*Rule, len(g.Rules))
	for i, r := range g.Rules {
		g2.Rules[i] = r.Optimized()
	}
	err := g2.Initialize()
	if err != nil {
		return nil, err
	}
	return &g2, nil
}
