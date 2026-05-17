// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/trees"
)

// Token matches a literal string token.
type Token struct {
	ModelBase
	Token string
}

// Parse implements the Model interface for Token.
func (t *Token) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	if !ctx.MatchToken(t.Token) {
		ctx.Reset(mark)
		return nil, ctx.Failure(mark, fmt.Errorf("expected: %q", t.Token))
	}
	return &trees.Text{Value: t.Token}, nil
}

// PubMap returns an ordered map of the Token's public fields.
func (t *Token) PubMap() *OrderedMap { return t.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Token.
func (t *Token) AsJSON() any { return t.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Token.
func (t *Token) AsJSONStr() string { return t.AsJSONStrOf(t) }
