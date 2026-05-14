// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

type Synth struct {
	Box
}

func (s *Synth) Parse(ctx Ctx) (Tree, error) {
	return s.Exp.Parse(ctx)
}

func (t *Synth) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Synth) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Synth) AsJSONStr() string   { return t.AsJSONStrOf(t) }
