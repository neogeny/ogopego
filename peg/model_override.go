// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
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

func (o *Override) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Override{Value: result}, nil
}

func (o *OverrideList) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.OverrideAsList{Value: result}, nil
}

func (t *Override) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Override) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Override) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (t *OverrideList) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *OverrideList) AsJSON() any         { return t.AsJSONOf(t) }
func (t *OverrideList) AsJSONStr() string   { return t.AsJSONStrOf(t) }
