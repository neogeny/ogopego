// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/util"
)

// Lookahead performs a positive lookahead: it succeeds if the nested
// expression matches without consuming input.
type Lookahead struct {
	Box
}

// NegativeLookahead succeeds when the nested expression does not match.
type NegativeLookahead struct {
	Box
}

// Parse implements the Model interface for Lookahead.
func (l *Lookahead) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	_, err := l.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err != nil {
		return nil, err
	}
	return NIL, nil
}

// Parse implements the Model interface for NegativeLookahead.
func (n *NegativeLookahead) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	_, err := n.Exp.Parse(ctx)
	ctx.Reset(mark)
	if err == nil {
		return nil, ctx.Failure(
			mark,
			fmt.Errorf(
				"negative lookahead matched:%v",
				n.Exp,
			),
		)
	}
	return NIL, nil
}

// PubMap returns an ordered map of the Lookahead's public fields.
func (t *Lookahead) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Lookahead.
func (t *Lookahead) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Lookahead.
func (t *Lookahead) AsJSONStr() string { return asjson.AsJSONStrOf(t) }

// PubMap returns an ordered map of the NegativeLookahead's public fields.
func (t *NegativeLookahead) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the NegativeLookahead.
func (t *NegativeLookahead) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the NegativeLookahead.
func (t *NegativeLookahead) AsJSONStr() string { return asjson.AsJSONStrOf(t) }
