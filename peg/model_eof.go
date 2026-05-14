// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
)

type EOF struct {
	ModelBase
}

func (e *EOF) Parse(ctx Ctx) (Tree, error) {
	if !ctx.Eof() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOF"),
		)
	}
	return NIL, nil
}

func (t *EOF) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *EOF) AsJSON() any         { return t.AsJSONOf(t) }
func (t *EOF) AsJSONStr() string   { return t.AsJSONStrOf(t) }
