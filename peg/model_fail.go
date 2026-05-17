// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/util"
)

// Fail represents a parsing failure.
type Fail struct {
	ModelBase
}

// Parse implements the Model interface for Fail.
func (f *Fail) Parse(ctx Ctx) (Tree, error) {
	return nil, ctx.Fail()
}

// PubMap returns an ordered map of the Fail's public fields.
func (t *Fail) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Fail.
func (t *Fail) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Fail.
func (t *Fail) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }
