// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/util"
)

// Void represents a model node that consumes no value but indicates a
// voided/ignored result (used for suppressed nodes).
type Void struct {
	ModelBase
}

// Parse implements the Model interface for Void.
func (v *Void) Parse(ctx Ctx) (Tree, error) {
	ctx.Void()
	return NIL, nil
}

// PubMap returns an ordered map of the Void's public fields.
func (t *Void) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Void.
func (t *Void) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Void.
func (t *Void) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }
