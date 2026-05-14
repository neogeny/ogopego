// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

type Fail struct {
	ModelBase
}

func (f *Fail) Parse(ctx Ctx) (Tree, error) {
	return nil, ctx.Fail()
}

func (t *Fail) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Fail) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Fail) AsJSONStr() string   { return t.AsJSONStrOf(t) }
