// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import "github.com/neogeny/ogopego/util"

// NULL represents a parse model node that always succeeds without consuming
// input and yields a nil-ish tree value.
type NULL struct {
	ModelBase
}

// Parse implements the Model interface for NULL.
func (n *NULL) Parse(ctx Ctx) (Tree, error) {
	return NIL, nil
}

// PubMap returns an ordered map of the NULL's public fields.
func (t *NULL) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the NULL.
func (t *NULL) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the NULL.
func (t *NULL) AsJSONStr() string { return t.AsJSONStrOf(t) }
