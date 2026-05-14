// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type EmptyClosure struct {
	ModelBase
}

func (e *EmptyClosure) Parse(ctx Ctx) (Tree, error) {
	return &trees.List{Items: nil}, nil
}

func (t *EmptyClosure) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *EmptyClosure) AsJSON() any         { return t.AsJSONOf(t) }
func (t *EmptyClosure) AsJSONStr() string   { return t.AsJSONStrOf(t) }
