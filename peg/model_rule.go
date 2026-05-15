// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"encoding/json"
	"unicode"

	"github.com/neogeny/ogopego/trees"
)

type Rule struct {
	NamedBox
	Params     []string
	KWParams   map[string]any
	Decorators []string
	Base       string
	IsName     bool
	IsTokn     bool
	NoMemo     bool
	NoStak     bool
	IsMemo     bool
	IsLrec     bool
}

func (r *Rule) Parse(ctx Ctx) (Tree, error) {
	mark := ctx.Mark()
	result, err := r.Exp.Parse(ctx)
	if err != nil {
		ctx.Reset(mark)
		return nil, err
	}
	folded := trees.Fold(result)
	if len(r.Params) == 0 || r.Params[0] == "bool" {
		return folded, nil
	}
	return &trees.Node{TypeName: r.Params[0], Tree: folded}, nil
}

func (r *Rule) IsToken() bool {
	if r.IsTokn {
		return true
	}
	for _, c := range r.Name {
		if c != '_' {
			return unicode.IsUpper(c)
		}
	}
	return false
}

func (r *Rule) IsLeftRecursive() bool { return r.IsLrec }

func (r *Rule) IsMemoizable() bool {
	return r.IsLrec || (r.IsMemo && !r.NoMemo)
}

func (r *Rule) ShouldTrace() bool {
	return !r.NoStak && !r.IsToken()
}

func (r *Rule) PubMap() *OrderedMap { return r.PubMapOf(r) }
func (r *Rule) AsJSON() any         { return r.AsJSONOf(r) }
func (r *Rule) AsJSONStr() string   { return r.AsJSONStrOf(r) }

func (r *Rule) MarshalJSON() ([]byte, error) { return json.Marshal(r.AsJSON()) }
