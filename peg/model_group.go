// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

type Group struct {
	Box
}

func (g *Group) Parse(ctx Ctx) (Tree, error) {
	return g.Exp.Parse(ctx)
}

func (t *Group) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Group) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Group) AsJSONStr() string   { return t.AsJSONStrOf(t) }
