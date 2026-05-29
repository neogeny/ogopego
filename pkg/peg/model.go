// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"errors"
	"fmt"
	"strings"

	context2 "github.com/neogeny/ogopego/pkg/context"
	asjson "github.com/neogeny/ogopego/pkg/json"
	"github.com/neogeny/ogopego/pkg/trees"
)

// Ctx is an alias for the parse-time context type.
type Ctx = context2.Ctx

// MemoKey is an alias for the type used as memoization keys.
type MemoKey = context2.MemoKey

// Tree is an alias for the tree node interface used in parse results.
type Tree = trees.Tree

// OrderedMap is an alias for the JSON ordered map type used by models.
type OrderedMap = asjson.OrderedMap

// NIL is the sentinel nil-like tree value used to represent empty results.
var NIL = trees.NIL

// Model is the interface implemented by all grammar model nodes.
type Model interface {
	// Parse attempts to parse the input using the model.
	Parse(ctx Ctx) (Tree, error)
	// Link resolves rule references within the model.
	Link(rules map[string]*Rule) error
	// ValidateLinked checks if all rule references are resolved.
	ValidateLinked() error
	// followRef returns the underlying ModelBase.
	followRef() *ModelBase

	// PrettyPrint returns a pretty-printed string representation of the model.
	PrettyPrint() string
	// Railroads returns a railroad diagram representation of the model.
	Railroads() string
}

// ModelBase provides common embedding for model nodes to carry identity and
// node data.
type ModelBase struct {
	Model
	Node
	la []string
}

// followRef returns the underlying ModelBase.
func (m *ModelBase) followRef() *ModelBase { return m }

// Parse is a placeholder for Model implementations.
func (m *ModelBase) Parse(ctx Ctx) (Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("method Parse() not implemented"))
}

// Link is a placeholder for Model implementations.
func (m *ModelBase) Link(rules map[string]*Rule) error { return nil }

// ValidateLinked is a placeholder for Model implementations.
func (m *ModelBase) ValidateLinked() error { return nil }

// PrettyPrint is a placeholder for Model implementations.
func (m *ModelBase) PrettyPrint() string { return "" }

// Railroads is a placeholder for Model implementations.
func (m *ModelBase) Railroads() string { return "" }

func (m *ModelBase) LookAheadStr() string {
	reprs := make([]string, len(m.la))
	for i, item := range m.la {
		reprs[i] = fmt.Sprintf("`%v`", item)
	}
	return strings.Join(reprs, " ")
}
