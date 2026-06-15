// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// EmptyClosure represents a closure that always matches an empty sequence, yielding an empty list.
type EmptyClosure struct {
	ModelBase
}

// Parse implements the Model interface for EmptyClosure.
func (e *EmptyClosure) Parse(ctx Ctx) (any, error) {
	return []any{}, nil
}
