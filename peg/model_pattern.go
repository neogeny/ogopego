// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

// Pattern matches input according to a configured pattern and returns text.
type Pattern struct {
	ModelBase
	Pattern string
}

// Parse implements the Model interface for Pattern.
func (p *Pattern) Parse(ctx Ctx) (Tree, error) {
	matched, err := ctx.MatchPattern(p.Pattern)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
}

// PubMap returns an ordered map of the Pattern's public fields.
func (t *Pattern) PubMap() *OrderedMap { return t.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Pattern.
func (t *Pattern) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Pattern.
func (t *Pattern) AsJSONStr() string { return t.AsJSONStrOf(t) }
