// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import "fmt"

// optimizeExpr recursively simplifies a grammar model tree.
// It returns the optimized model, which may be a different concrete type
// (e.g., *Group is eliminated, returning its inner expression).
func optimizeExpr(m Model) Model {
	switch e := m.(type) {
	// --- Leaves ---
	case *Dot, *Cut, *Void, *Fail, *EOF, *EOL,
		*Token, *Pattern, *Constant, *Alert,
		*EmptyClosure, *NULL, *Call, *RuleInclude:
		return e

	// --- Unary containers: recurse into Exp, keep self ---
	case *Optional:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *Closure:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *PositiveClosure:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *Lookahead:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *NegativeLookahead:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *SkipGroup:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *SkipTo:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *Override:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *OverrideList:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *Named:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *NamedList:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *Option:
		e.Exp = optimizeExpr(e.Exp)
		return e
	case *Synth:
		return optimizeExpr(e.Exp)

	// --- Binary containers: recurse into both ---
	case *Join:
		e.Exp = optimizeExpr(e.Exp)
		e.Sep = optimizeExpr(e.Sep)
		return e
	case *PositiveJoin:
		e.Exp = optimizeExpr(e.Exp)
		e.Sep = optimizeExpr(e.Sep)
		return e
	case *Gather:
		e.Exp = optimizeExpr(e.Exp)
		e.Sep = optimizeExpr(e.Sep)
		return e
	case *PositiveGather:
		e.Exp = optimizeExpr(e.Exp)
		e.Sep = optimizeExpr(e.Sep)
		return e

	// --- Eliminated: Group unwraps to its inner expression ---
	case *Group:
		return optimizeExpr(e.Exp)

	// --- Collections: recurse into children ---
	case *Choice:
		for _, o := range e.Options {
			o.Exp = optimizeExpr(o.Exp)
		}
		if len(e.Options) == 1 {
			return e.Options[0].Exp
		}
		return e

	case *Sequence:
		for i, s := range e.Sequence {
			e.Sequence[i] = optimizeExpr(s)
		}
		if len(e.Sequence) == 1 {
			return e.Sequence[0]
		}
		return e

	default:
		panic(fmt.Sprintf("optimizeExpr: unhandled model type %T", m))
	}
}

// Optimize simplifies the rule's expression tree.
func (r *Rule) Optimize() {
	r.Exp = optimizeExpr(r.Exp)
}

// Optimize simplifies all rules in the grammar.
func (g *Grammar) Optimize() {
	for _, r := range g.Rules {
		r.Optimize()
	}
}
