// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util"
)

// EmptyClosure represents a closure that always matches an empty sequence, yielding an empty list.
type EmptyClosure struct {
	ModelBase
}

// Parse implements the Model interface for EmptyClosure.
func (e *EmptyClosure) Parse(ctx Ctx) (Tree, error) {
	return &trees.List{Items: nil}, nil
}

// PubMap returns an ordered map of the EmptyClosure's public fields.
func (t *EmptyClosure) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the EmptyClosure.
func (t *EmptyClosure) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the EmptyClosure.
func (t *EmptyClosure) AsJSONStr() string { return asjson.AsJSONStrOf(t) }
