// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
	"github.com/neogeny/ogopego/util"
)

// EOL matches an end-of-line sequence.
type EOL struct {
	ModelBase
}

// Parse implements the Model interface for EOL.
func (e *EOL) Parse(ctx Ctx) (Tree, error) {
	if !ctx.MatchEOL() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOL"),
		)
	}
	return &trees.Nil{}, nil
}

// PubMap returns an ordered map of the EOL's public fields.
func (t *EOL) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the EOL.
func (t *EOL) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the EOL.
func (t *EOL) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }
