// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/util"
)

// SkipGroup represents a group that is parsed but whose result is discarded.
type SkipGroup struct {
	Box
}

// Parse implements the Model interface for SkipGroup.
func (s *SkipGroup) Parse(ctx Ctx) (Tree, error) {
	_, err := s.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

// PubMap returns an ordered map of the SkipGroup's public fields.
func (t *SkipGroup) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the SkipGroup.
func (t *SkipGroup) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the SkipGroup.
func (t *SkipGroup) AsJSONStr() string { return asjson.AsJSONStrOf(t) }
