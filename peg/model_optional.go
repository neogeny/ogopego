// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

// Optional represents an optional expression that may succeed with a
// Nil result if the nested expression fails without a cut.
type Optional struct {
	ModelBase
	Exp Model
}

// Parse implements the Model interface for Optional.
func (o *Optional) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()

	ctx.CutStackPush()
	result, err := o.Exp.Parse(ctx)
	cutSeen := ctx.CutStackPop()

	if err != nil {
		ctx.Reset(mark)
		if cutSeen {
			return nil, err
		}
		return &trees.Nil{}, nil
	}
	return result, nil
}
