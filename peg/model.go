// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	gojson "encoding/json"
	"errors"

	"github.com/neogeny/ogopego/context"
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

// Ctx is an alias for the parse-time context type.
type Ctx = context.Ctx

// MemoKey is an alias for the type used as memoization keys.
type MemoKey = context.MemoKey

// Tree is an alias for the tree node interface used in parse results.
type Tree = trees.Tree

// OrderedMap is an alias for the JSON ordered map type used by models.
type OrderedMap = asjson.OrderedMap

// NIL is the sentinel nil-like tree value used to represent empty results.
var NIL = trees.NIL

// Model is the interface implemented by all grammar model nodes.
type Model interface {
	Parse(ctx Ctx) (Tree, error)
	Link(rules map[string]*Rule) error
	ValidateLinked() error
	followRef() *ModelBase

	PrettyPrint() string
	Railroads() string
}

type ModelBase struct {
	Model
	Node
}

// ModelBase provides common embedding for model nodes to carry identity and
// node data.

func (m *ModelBase) followRef() *ModelBase { return m }

func (m *ModelBase) Parse(ctx Ctx) (Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("Parse not implemented"))
}
func (m *ModelBase) Link(rules map[string]*Rule) error { return nil }
func (m *ModelBase) ValidateLinked() error             { return nil }
func (m *ModelBase) PrettyPrint() string               { return "" }
func (m *ModelBase) Railroads() string                 { return "" }

func (m *ModelBase) MarshalJSON() ([]byte, error) { return gojson.Marshal(m.AsJSON()) }

// modelMarshalJSON is a helper for model types to implement json.Marshaler
// by dispatching through AsJSONMixin to the correct AsJSON() override.
func modelMarshalJSON(v any) ([]byte, error) {
	if mixin, ok := v.(asjson.AsJSONMixin); ok {
		return gojson.MarshalIndent(mixin.AsJSON(), "", "  ")
	}
	return gojson.Marshal(nil)
}
