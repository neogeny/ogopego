// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/util"
)

// EOF matches the end of input.
type EOF struct {
	ModelBase
}

// Parse implements the Model interface for EOF.
func (e *EOF) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	ctx.NextToken()
	if !ctx.Eof() {
		ctx.Reset(mark)
		return nil, ctx.Failure(
			mark,
			fmt.Errorf("expected EOF"),
		)
	}
	return NIL, nil
}

// PubMap returns an ordered map of the EOF's public fields.
func (t *EOF) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the EOF.
func (t *EOF) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the EOF.
func (t *EOF) AsJSONStr() string { return t.AsJSONStrOf(t) }
