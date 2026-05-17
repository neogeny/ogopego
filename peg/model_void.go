// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// Void represents a model node that consumes no value but indicates a
// voided/ignored result (used for suppressed nodes).
type Void struct {
	ModelBase
}

func (v *Void) Parse(ctx Ctx) (Tree, error) {
	ctx.Void()
	return NIL, nil
}

func (t *Void) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Void) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Void) AsJSONStr() string   { return t.AsJSONStrOf(t) }
