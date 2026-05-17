// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util"
)

// Sequence represents an ordered sequence of model elements.
type Sequence struct {
	ModelBase
	Sequence []Model
}

// Parse implements the Model interface for Sequence.
func (s *Sequence) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	var items []Tree
	for _, el := range s.Sequence {
		result, err := el.Parse(ctx)
		if err != nil {
			ctx.Reset(mark)
			return nil, err
		}
		if _, ok := result.(*trees.Nil); !ok {
			items = append(items, result)
		}
	}
	var tree Tree = NIL
	switch len(items) {
	case 0:
	case 1:
		tree = items[0]
	default:
		tree = &trees.Seq{Items: items}
	}
	return tree, nil
}

// PubMap returns an ordered map of the Sequence's public fields.
func (t *Sequence) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Sequence.
func (t *Sequence) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Sequence.
func (t *Sequence) AsJSONStr() string { return asjson.AsJSONStrOf(t) }
