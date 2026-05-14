// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

type Void struct {
	ModelBase
}

func (v *Void) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	err := ctx.Void()
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	return NIL, nil
}

func (t *Void) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Void) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Void) AsJSONStr() string   { return t.AsJSONStrOf(t) }
