// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/pkg/trees"
)

// Override marks an expression whose value overrides surrounding values.
type Override struct {
	ModelBase
	Exp Model
}

// OverrideList marks an expression whose override value should be treated
// as a list.
type OverrideList struct {
	ModelBase
	Exp Model
}

// Parse implements the Model interface for Override.
func (o *Override) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.Override{Value: result}, nil
}

// Parse implements the Model interface for OverrideList.
func (o *OverrideList) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return &trees.OverrideAsList{Value: result}, nil
}
