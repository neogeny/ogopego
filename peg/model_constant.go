// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

type Constant struct {
	ModelBase
	Literal string
}

type Alert struct {
	Constant
	Level int
}

func (c *Constant) Parse(ctx Ctx) (Tree, error) {
	return ctx.Constant(c.Literal)
}

func (t *Constant) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Constant) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Constant) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (a *Alert) Parse(ctx Ctx) (Tree, error) {
	return ctx.Constant(a.Literal)
}

func (t *Alert) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Alert) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Alert) AsJSONStr() string   { return t.AsJSONStrOf(t) }
