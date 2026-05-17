// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// Synth synthesizes a value from its nested expression; it is a thin wrapper
// around the nested expression used for model-level transformations.
type Synth struct {
	Box
}

func (s *Synth) Parse(ctx Ctx) (Tree, error) {
	return s.Exp.Parse(ctx)
}

func (t *Synth) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Synth) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Synth) AsJSONStr() string   { return t.AsJSONStrOf(t) }
