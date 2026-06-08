// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
)

// optimizeExpr recursively simplifies a grammar model tree.
// It returns the optimized model, which may be a different concrete type
// (e.g., *Group is eliminated, returning its inner expression).
// The original tree is not modified.
func optimizeExpr(m Model) Model {
	switch e := m.(type) {
	// --- Leaves ---
	case *Dot, *Cut, *Void, *Fail, *EOF, *EOL,
		*Token, *Pattern, *Constant, *Alert,
		*EmptyClosure, *NULL,
		*MetaExp:
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
		return &Synth{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp)}

	// --- Binary containers: clone and recurse into both ---
	case *Join:
		return &Join{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Sep: optimizeExpr(e.Sep)}
	case *PositiveJoin:
		return &PositiveJoin{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Sep: optimizeExpr(e.Sep)}
	case *Gather:
		return &Gather{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Sep: optimizeExpr(e.Sep)}
	case *PositiveGather:
		return &PositiveGather{ModelBase: e.ModelBase, Exp: optimizeExpr(e.Exp), Sep: optimizeExpr(e.Sep)}

	// --- Call: always clone (re-linked by Initialize) ---
	case *Call:
		return &Call{ModelBase: e.ModelBase, Name: e.Name}

	// --- RuleInclude: copy and recurse into inner exp ---
	case *RuleInclude:
		return &RuleInclude{ModelBase: e.ModelBase, Name: e.Name, exp: optimizeExpr(e.exp)}

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
//
//  1. Recursively optimize the expression tree.
//  2. Unwrap single-element Sequences that result from optimization.
//  3. If the optimized expression is a single Call with a linked target rule,
//     inline the target rule's optimized body (alias resolution).
func (r *Rule) Optimized() *Rule {
	opt := *r
	opt.Exp = optimizeExpr(r.Exp)

	// Unwrap single-element Sequences (matching Python's Rule.optimized)
	for seq, ok := opt.Exp.(*Sequence); ok && len(seq.Sequence) == 1; {
		opt.Exp = seq.Sequence[0]
		seq, ok = opt.Exp.(*Sequence)
	}

	// If the body is a single Call to another rule, inline that rule's body
	if call, ok := opt.Exp.(*Call); ok && call.rule != nil {
		opt.Exp = optimizeExpr(call.rule.Exp)
	}

	return &opt
}

// Optimized simplifies all rules in the grammar and returns a new grammar.
// The original grammar is not modified.
func (g *Grammar) Optimized() (*Grammar, error) {
	if g.optimized {
		return g, nil
	}
	opt := *g
	opt.Rules = make([]*Rule, len(g.Rules))
	for i, r := range g.Rules {
		opt.Rules[i] = r.Optimized()
	}
	err := opt.Initialize()
	if err != nil {
		return nil, err
	}
	opt.optimized = true
	return &opt, nil
}
