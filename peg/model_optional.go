// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util"
)

// Optional represents an optional expression that may succeed with a
// Nil result if the nested expression fails without a cut.
type Optional struct {
	Box
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

// PubMap returns an ordered map of the Optional's public fields.
func (t *Optional) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Optional.
func (t *Optional) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Optional.
func (t *Optional) AsJSONStr() string { return t.AsJSONStrOf(t) }
