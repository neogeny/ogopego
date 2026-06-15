// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// Join represents a sequence joined by a separator expression (e.g. a list
// of elements with separators).
type Join struct {
	ModelBase
	Exp Model
	Sep Model
}

// PositiveJoin is like Join but requires at least one element.
type PositiveJoin struct {
	ModelBase
	Exp Model
	Sep Model
}

// Gather behaves like Join but collects elements differently (semantics
// vary by model consumer).
type Gather struct {
	ModelBase
	Exp Model
	Sep Model
}

// PositiveGather is the one-or-more variant of Gather.
type PositiveGather struct {
	ModelBase
	Exp Model
	Sep Model
}

// Parse implements the Model interface for Join.
func (j *Join) Parse(ctx Ctx) (any, error) {
	return repeatWithSep(ctx, j.Exp, j.Sep, false, true)
}

// Parse implements the Model interface for PositiveJoin.
func (p *PositiveJoin) Parse(ctx Ctx) (any, error) {
	return repeatWithSep(ctx, p.Exp, p.Sep, true, true)
}

// Parse implements the Model interface for Gather.
func (g *Gather) Parse(ctx Ctx) (any, error) {
	return repeatWithSep(ctx, g.Exp, g.Sep, false, false)
}

// Parse implements the Model interface for PositiveGather.
func (p *PositiveGather) Parse(ctx Ctx) (any, error) {
	return repeatWithSep(ctx, p.Exp, p.Sep, true, false)
}
