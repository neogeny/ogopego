// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// Synth synthesizes a value from its nested expression; it is a thin wrapper
// around the nested expression used for model-level transformations.
type Synth struct {
	Box
}

// Parse implements the Model interface for Synth.
func (s *Synth) Parse(ctx Ctx) (Tree, error) {
	return s.Exp.Parse(ctx)
}
