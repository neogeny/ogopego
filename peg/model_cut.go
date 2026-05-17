// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util"
)

// Cut represents a cut operator that prunes alternative backtracking paths.
type Cut struct {
	ModelBase
}

// Parse implements the Model interface for Cut.
func (c *Cut) Parse(ctx Ctx) (Tree, error) {
	ctx.Cut()
	t := &trees.Nil{}
	return t, nil
}

// PubMap returns an ordered map of the Cut's public fields.
func (t *Cut) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Cut.
func (t *Cut) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Cut.
func (t *Cut) AsJSONStr() string { return t.AsJSONStrOf(t) }
