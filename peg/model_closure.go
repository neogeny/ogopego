// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

type Closure struct {
	Box
}

type PositiveClosure struct {
	Closure
}

func (c *Closure) Parse(ctx Ctx) (Tree, error) {
	return repeat(ctx, c.Exp, false)
}

func (p *PositiveClosure) Parse(ctx Ctx) (Tree, error) {
	return repeat(ctx, p.Exp, true)
}

func (t *Closure) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Closure) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Closure) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (t *PositiveClosure) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *PositiveClosure) AsJSON() any         { return t.AsJSONOf(t) }
func (t *PositiveClosure) AsJSONStr() string   { return t.AsJSONStrOf(t) }
