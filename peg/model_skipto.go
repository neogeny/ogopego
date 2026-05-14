// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	asjson "github.com/neogeny/ogopego/json"
)

type SkipTo struct {
	Box
}

func (s *SkipTo) Parse(ctx Ctx) (Tree, error) {
	for {
		_, err := s.Exp.Parse(ctx)
		if err == nil {
			return NIL, nil
		}
		return nil, err
	}
}

func (t *SkipTo) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *SkipTo) AsJSON() any                { return t.AsJSONOf(t) }
func (t *SkipTo) AsJSONStr() string          { return t.AsJSONStrOf(t) }
