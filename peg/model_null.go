// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// NULL represents a parse model node that always succeeds without consuming
// input and yields a nil-ish tree value.
type NULL struct {
	ModelBase
}

func (n *NULL) Parse(ctx Ctx) (Tree, error) {
	return NIL, nil
}

func (t *NULL) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *NULL) AsJSON() any         { return t.AsJSONOf(t) }
func (t *NULL) AsJSONStr() string   { return t.AsJSONStrOf(t) }
