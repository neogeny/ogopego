// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import "github.com/neogeny/ogopego/util"

// Closure represents the Kleene-closure (zero-or-more) of an expression.
type Closure struct {
	Box
}

// PositiveClosure represents the positive closure (one-or-more) of an
// expression.
type PositiveClosure struct {
	Closure
}

// Parse implements the Model interface for Closure.
func (c *Closure) Parse(ctx Ctx) (Tree, error) {
	return repeat(ctx, c.Exp, false)
}

// Parse implements the Model interface for PositiveClosure.
func (p *PositiveClosure) Parse(ctx Ctx) (Tree, error) {
	return repeat(ctx, p.Exp, true)
}

// PubMap returns an ordered map of the Closure's public fields.
func (t *Closure) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Closure.
func (t *Closure) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Closure.
func (t *Closure) AsJSONStr() string { return t.AsJSONStrOf(t) }

// PubMap returns an ordered map of the PositiveClosure's public fields.
func (t *PositiveClosure) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the PositiveClosure.
func (t *PositiveClosure) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the PositiveClosure.
func (t *PositiveClosure) AsJSONStr() string { return t.AsJSONStrOf(t) }
