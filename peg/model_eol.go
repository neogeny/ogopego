// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/trees"
)

type EOL struct {
	ModelBase
}

func (e *EOL) Parse(ctx Ctx) (Tree, error) {
	if !ctx.MatchEOL() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOL"),
		)
	}
	return &trees.Nil{}, nil
}

func (t *EOL) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *EOL) AsJSON() any         { return t.AsJSONOf(t) }
func (t *EOL) AsJSONStr() string   { return t.AsJSONStrOf(t) }
