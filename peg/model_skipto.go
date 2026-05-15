// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
)

type SkipTo struct {
	Box
}

func (s *SkipTo) Parse(ctx Ctx) (Tree, error) {
	for !ctx.Eof() {
		mark := ctx.Mark()
		result, err := s.Exp.Parse(ctx)
		if err == nil {
			return result, nil
		}
		ctx.Reset(mark)
		if _, ok := ctx.Next(); !ok {
			return nil, ctx.Failure(mark, fmt.Errorf("skip_to: target not found"))
		}
	}
	return nil, ctx.Failure(
		ctx.Mark(),
		fmt.Errorf("skip_to: target not found"),
	)
}

func (t *SkipTo) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *SkipTo) AsJSON() any                { return t.AsJSONOf(t) }
func (t *SkipTo) AsJSONStr() string          { return t.AsJSONStrOf(t) }
