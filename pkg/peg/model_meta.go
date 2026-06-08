// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/pkg/trees"
)

// MetaExp matches a meta-expression (@name, @int, @uint, @float, @bool).
type MetaExp struct {
	ModelBase
	Kind string
}

// Parse implements the Model interface for MetaExp by dispatching on Kind.
func (m *MetaExp) Parse(ctx Ctx) (Tree, error) {
	switch m.Kind {
	case "name":
		s, err := ctx.MatchName()
		if err != nil {
			return nil, err
		}
		return &trees.Text{Value: s}, nil

	case "int":
		n, err := ctx.MatchInt()
		if err != nil {
			return nil, err
		}
		return &trees.Number{Value: float64(n)}, nil

	case "uint":
		n, err := ctx.MatchUInt()
		if err != nil {
			return nil, err
		}
		return &trees.Number{Value: float64(n)}, nil

	case "float":
		f, err := ctx.MatchFloat()
		if err != nil {
			return nil, err
		}
		return &trees.Number{Value: f}, nil

	case "bool":
		b, err := ctx.MatchBool()
		if err != nil {
			return nil, err
		}
		return &trees.Bool{Value: b}, nil

	default:
		panic("unknown meta kind: " + m.Kind)
	}
}
