// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// Closure represents the Kleene-closure (zero-or-more) of an expression.
type Closure struct {
	ModelBase
	Exp Model
}

// PositiveClosure represents the positive closure (one-or-more) of an
// expression.
type PositiveClosure struct {
	ModelBase
	Exp Model
}

// Parse implements the Model interface for Closure.
func (c *Closure) Parse(ctx Ctx) (Tree, error) {
	return repeat(ctx, c.Exp, false)
}

// Parse implements the Model interface for PositiveClosure.
func (p *PositiveClosure) Parse(ctx Ctx) (Tree, error) {
	return repeat(ctx, p.Exp, true)
}
