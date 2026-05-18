// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

// SkipGroup represents a group that is parsed but whose result is discarded.
type SkipGroup struct {
	ModelBase
	Exp Model
}

// Parse implements the Model interface for SkipGroup.
func (s *SkipGroup) Parse(ctx Ctx) (Tree, error) {
	_, err := s.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}
