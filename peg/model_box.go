// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import "github.com/neogeny/ogopego/util"

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

// PubMap returns an ordered map of the Box's public fields.
func (t *Box) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Box.
func (t *Box) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Box.
func (t *Box) AsJSONStr() string { return t.AsJSONStrOf(t) }

// PubMap returns an ordered map of the NamedBox's public fields.
func (t *NamedBox) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the NamedBox.
func (t *NamedBox) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the NamedBox.
func (t *NamedBox) AsJSONStr() string { return t.AsJSONStrOf(t) }
