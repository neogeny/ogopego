// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/trees"
)

type Token struct {
	ModelBase
	Token string
}

func (t *Token) Parse(ctx Ctx) (Tree, error) {
	matched, err := ctx.Token(t.Token)
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: matched}, nil
}

func (t *Token) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Token) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Token) AsJSONStr() string   { return t.AsJSONStrOf(t) }
