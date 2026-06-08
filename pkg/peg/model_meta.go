// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"github.com/neogeny/ogopego/pkg/trees"
)

// MetaExp is the base type for meta-expressions.
type MetaExp struct {
	ModelBase
}

// NameMeta matches a @name meta-expression.
type NameMeta struct{ MetaExp }

// IntMeta matches a @int meta-expression.
type IntMeta struct{ MetaExp }

// UIntMeta matches a @uint meta-expression.
type UIntMeta struct{ MetaExp }

// FloatMeta matches a @float meta-expression.
type FloatMeta struct{ MetaExp }

// BoolMeta matches a @bool meta-expression.
type BoolMeta struct{ MetaExp }

func (m *NameMeta) Parse(ctx Ctx) (Tree, error) {
	s, err := ctx.MatchName()
	if err != nil {
		return nil, err
	}
	return &trees.Text{Value: s}, nil
}

func (m *IntMeta) Parse(ctx Ctx) (Tree, error) {
	n, err := ctx.MatchInt()
	if err != nil {
		return nil, err
	}
	return &trees.Number{Value: float64(n)}, nil
}

func (m *UIntMeta) Parse(ctx Ctx) (Tree, error) {
	n, err := ctx.MatchUInt()
	if err != nil {
		return nil, err
	}
	return &trees.Number{Value: float64(n)}, nil
}

func (m *FloatMeta) Parse(ctx Ctx) (Tree, error) {
	f, err := ctx.MatchFloat()
	if err != nil {
		return nil, err
	}
	return &trees.Number{Value: f}, nil
}

func (m *BoolMeta) Parse(ctx Ctx) (Tree, error) {
	b, err := ctx.MatchBool()
	if err != nil {
		return nil, err
	}
	return &trees.Bool{Value: b}, nil
}
