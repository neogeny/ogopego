// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type Optional struct {
	Box
}

func (o *Optional) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		if ctx.CutStackPop() {
			return nil, err
		}
		return &trees.Nil{}, nil
	}
	return result, nil
}

func (t *Optional) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Optional) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Optional) AsJSONStr() string   { return t.AsJSONStrOf(t) }
