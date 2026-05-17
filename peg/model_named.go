// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

// Named wraps an expression result with a name, producing a Named tree node.
type Named struct {
	NamedBox
}

// NamedList wraps an expression result into a Named-as-list node.
type NamedList struct {
	Named
}

// Parse implements the Model interface for Named.
func (n *Named) Parse(ctx Ctx) (Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Named{Name: n.Name, Value: result}, nil
}

// Parse implements the Model interface for NamedList.
func (n *NamedList) Parse(ctx Ctx) (Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.NamedAsList{Name: n.Name, Value: result}, nil
}

// PubMap returns an ordered map of the Named's public fields.
func (t *Named) PubMap() *OrderedMap { return t.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Named.
func (t *Named) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Named.
func (t *Named) AsJSONStr() string { return t.AsJSONStrOf(t) }

// PubMap returns an ordered map of the NamedList's public fields.
func (t *NamedList) PubMap() *OrderedMap { return t.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the NamedList.
func (t *NamedList) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the NamedList.
func (t *NamedList) AsJSONStr() string { return t.AsJSONStrOf(t) }
