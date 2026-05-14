// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

type Box struct {
	ModelBase
	Exp Model
}

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
