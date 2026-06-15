// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// Void represents a model node that consumes no value but indicates a
// voided/ignored result (used for suppressed nodes).
type Void struct {
	ModelBase
}

func (v *Void) Link(g *Grammar) error {
	return nil
}

// Parse implements the Model interface for Void.
func (v *Void) Parse(ctx Ctx) (any, error) {
	ctx.Void()
	return nil, nil
}
