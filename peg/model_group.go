// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/context"
)

type Group struct {
	Box
}

func (g *Group) Parse(ctx Ctx) (Tree, error) {
	result, err := g.Exp.Parse(ctx)
	if err != nil {
		if pf, ok := err.(*context.Nope); ok {
			pf.CutSeen = false
		}
		return nil, err
	}
	result.TakeCutSeen()
	return result, nil
}

func (t *Group) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Group) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Group) AsJSONStr() string   { return t.AsJSONStrOf(t) }
