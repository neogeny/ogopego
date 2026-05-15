// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type Cut struct {
	ModelBase
}

func (c *Cut) Parse(ctx Ctx) (Tree, error) {
	ctx.Cut()
	t := &trees.Nil{}
	return t, nil
}

func (t *Cut) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Cut) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Cut) AsJSONStr() string   { return t.AsJSONStrOf(t) }
