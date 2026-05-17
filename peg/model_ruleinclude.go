// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

// RuleInclude represents a reference to another rule whose expression is
// included into the current grammar model once resolved.
type RuleInclude struct {
	ModelBase
	Name string
	Exp  Model
}

// Parse implements the Model interface for RuleInclude.
func (r *RuleInclude) Parse(ctx Ctx) (Tree, error) {
	if r.Exp == nil {
		return nil, ctx.Failure(ctx.Mark(), fmt.Errorf("RuleInclude %q has not been resolved", r.Name))
	}
	return r.Exp.Parse(ctx)
}
