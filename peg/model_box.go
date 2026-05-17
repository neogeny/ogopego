// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// Box is a model wrapper that contains a nested expression (Exp).
type Box struct {
	ModelBase
	Exp Model
}

// NamedBox is a Box that carries a name for the nested expression.
type NamedBox struct {
	Box
	Name string
}

func (t *Box) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Box) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Box) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (t *NamedBox) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *NamedBox) AsJSON() any         { return t.AsJSONOf(t) }
func (t *NamedBox) AsJSONStr() string   { return t.AsJSONStrOf(t) }
