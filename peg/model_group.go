// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// Group represents an explicit grouping of an expression; it has no
// semantics beyond grouping (used for precedence and organization).
type Group struct {
	Box
}

// Parse implements the Model interface for Group.
func (g *Group) Parse(ctx Ctx) (Tree, error) {
	return g.Exp.Parse(ctx)
}

// PubMap returns an ordered map of the Group's public fields.
func (t *Group) PubMap() *OrderedMap { return t.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Group.
func (t *Group) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Group.
func (t *Group) AsJSONStr() string { return t.AsJSONStrOf(t) }
