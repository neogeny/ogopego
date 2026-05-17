// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/util"
)

// Join represents a sequence joined by a separator expression (e.g. a list
// of elements with separators).
type Join struct {
	Box
	Sep Model
}

// PositiveJoin is like Join but requires at least one element.
type PositiveJoin struct {
	Join
}

// Gather behaves like Join but collects elements differently (semantics
// vary by model consumer).
type Gather struct {
	Join
}

// PositiveGather is the one-or-more variant of Gather.
type PositiveGather struct {
	Gather
}

// Parse implements the Model interface for Join.
func (j *Join) Parse(ctx Ctx) (Tree, error) {
	return repeatWithSep(ctx, j.Exp, j.Sep, false, true)
}

// PubMap returns an ordered map of the Join's public fields.
func (t *Join) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Join.
func (t *Join) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Join.
func (t *Join) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }

// Parse implements the Model interface for PositiveJoin.
func (p *PositiveJoin) Parse(ctx Ctx) (Tree, error) {
	return repeatWithSep(ctx, p.Exp, p.Sep, true, true)
}

// PubMap returns an ordered map of the PositiveJoin's public fields.
func (t *PositiveJoin) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the PositiveJoin.
func (t *PositiveJoin) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the PositiveJoin.
func (t *PositiveJoin) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }

// Parse implements the Model interface for Gather.
func (g *Gather) Parse(ctx Ctx) (Tree, error) {
	return repeatWithSep(ctx, g.Exp, g.Sep, false, false)
}

// PubMap returns an ordered map of the Gather's public fields.
func (t *Gather) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Gather.
func (t *Gather) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Gather.
func (t *Gather) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }

// Parse implements the Model interface for PositiveGather.
func (p *PositiveGather) Parse(ctx Ctx) (Tree, error) {
	return repeatWithSep(ctx, p.Exp, p.Sep, true, false)
}

// PubMap returns an ordered map of the PositiveGather's public fields.
func (t *PositiveGather) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the PositiveGather.
func (t *PositiveGather) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the PositiveGather.
func (t *PositiveGather) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }
