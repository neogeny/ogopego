// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// NULL represents a parse model node that always succeeds without consuming
// input and yields a nil-ish tree value.
type NULL struct {
	ModelBase
}

func (n *NULL) Link(g *Grammar) error {
	return nil
}

// Parse implements the Model interface for NULL.
func (n *NULL) Parse(ctx Ctx) (Tree, error) {
	return NIL, nil
}
