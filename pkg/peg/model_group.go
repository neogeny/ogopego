// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// Group represents an explicit grouping of an expression; it has no
// semantics beyond grouping (used for precedence and organization).
type Group struct {
	ModelBase
	Exp Model
}

// Parse implements the Model interface for Group.
func (g *Group) Parse(ctx Ctx) (any, error) {
	return g.Exp.Parse(ctx)
}
