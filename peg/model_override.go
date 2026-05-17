// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util"
)

// Override marks an expression whose value overrides surrounding values.
type Override struct {
	Box
}

// OverrideList marks an expression whose override value should be treated
// as a list.
type OverrideList struct {
	Box
}

// Parse implements the Model interface for Override.
func (o *Override) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Override{Value: result}, nil
}

// Parse implements the Model interface for OverrideList.
func (o *OverrideList) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.OverrideAsList{Value: result}, nil
}

// PubMap returns an ordered map of the Override's public fields.
func (t *Override) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Override.
func (t *Override) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Override.
func (t *Override) AsJSONStr() string { return t.AsJSONStrOf(t) }

// PubMap returns an ordered map of the OverrideList's public fields.
func (t *OverrideList) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the OverrideList.
func (t *OverrideList) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the OverrideList.
func (t *OverrideList) AsJSONStr() string { return t.AsJSONStrOf(t) }
