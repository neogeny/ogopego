// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type Named struct {
	NamedBox
}

type NamedList struct {
	Named
}

func (n *Named) Parse(ctx Ctx) (Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Named{Name: n.Name, Value: result}, nil
}

func (n *NamedList) Parse(ctx Ctx) (Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.NamedAsList{Name: n.Name, Value: result}, nil
}

func (t *Named) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Named) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Named) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (t *NamedList) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *NamedList) AsJSON() any         { return t.AsJSONOf(t) }
func (t *NamedList) AsJSONStr() string   { return t.AsJSONStrOf(t) }
