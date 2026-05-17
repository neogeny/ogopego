// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

// Named wraps an expression result with a name, producing a Named tree node.
type Named struct {
	NamedBox
}

// NamedList wraps an expression result into a Named-as-list node.
type NamedList struct {
	Named
}

// Parse implements the Model interface for Named.
func (n *Named) Parse(ctx Ctx) (Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Named{Name: n.Name, Value: result}, nil
}

// Parse implements the Model interface for NamedList.
func (n *NamedList) Parse(ctx Ctx) (Tree, error) {
	result, err := n.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.NamedAsList{Name: n.Name, Value: result}, nil
}
