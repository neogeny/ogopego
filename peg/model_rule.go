// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"unicode"

	"github.com/neogeny/ogopego/trees"
)

// Rule represents a named grammar rule with optional parameters and
// metadata used during parsing and code generation.
type Rule struct {
	NamedBox
	// Params are the parameters for the rule.
	Params []string
	// KWParams are keyword parameters for the rule.
	KWParams map[string]any
	// Decorators are decorators applied to the rule.
	Decorators []string
	// Base is the base rule name.
	Base string
	// IsName indicates if the rule is a name rule.
	IsName bool
	// IsTokn indicates if the rule is a token rule.
	IsTokn bool
	// NoMemo disables memoization for the rule.
	NoMemo bool
	// NoStak disables tracing for the rule.
	NoStak bool
	// IsMemo indicates if the rule is memoizable.
	IsMemo bool
	// IsLrec indicates if the rule is left-recursive.
	IsLrec bool
}

// Parse implements the Model interface for Rule.
func (r *Rule) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	result, err := r.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	folded := trees.Fold(result)
	if len(r.Params) == 0 || r.Params[0] == "bool" {
		return folded, nil
	}
	return &trees.Node{TypeName: r.Params[0], Tree: folded}, nil
}

// IsToken returns true if the rule is a token rule.
func (r *Rule) IsToken() bool {
	if r.IsTokn {
		return true
	}
	for _, c := range r.Name {
		if c != '_' {
			return unicode.IsUpper(c)
		}
	}
	return false
}

// IsLeftRecursive returns true if the rule is left-recursive.
func (r *Rule) IsLeftRecursive() bool { return r.IsLrec }

// IsMemoizable returns true if the rule is memoizable.
func (r *Rule) IsMemoizable() bool {
	return r.IsLrec || (r.IsMemo && !r.NoMemo)
}

// ShouldTrace returns true if the rule should be traced.
func (r *Rule) ShouldTrace() bool {
	return !r.NoStak && !r.IsToken()
}
