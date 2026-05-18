// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

// SkipTo scans forward until the nested expression matches, returning the
// nested result; used to skip to a target token or pattern.
type SkipTo struct {
	ModelBase
	Exp Model
}

// Parse implements the Model interface for SkipTo.
func (s *SkipTo) Parse(ctx Ctx) (Tree, error) {
	for !ctx.Eof() {
		mark := ctx.Mark()
		result, err := s.Exp.Parse(ctx)
		if err == nil {
			return result, nil
		}
		ctx.Reset(mark)
		if _, ok := ctx.Next(); !ok {
			return nil, ctx.Failure(mark, fmt.Errorf("skip_to: target not found"))
		}
	}
	return nil, ctx.Failure(
		ctx.Mark(),
		fmt.Errorf("skip_to: target not found"),
	)
}
