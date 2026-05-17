// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/util"
)

// Constant represents a literal constant value.
type Constant struct {
	ModelBase
	Literal string
}

// Alert represents an alert with a literal message and a level.
type Alert struct {
	Constant
	Level int
}

// Parse implements the Model interface for Constant.
func (c *Constant) Parse(ctx Ctx) (Tree, error) {
	return ctx.Constant(c.Literal)
}

// PubMap returns an ordered map of the Constant's public fields.
func (t *Constant) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Constant.
func (t *Constant) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Constant.
func (t *Constant) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }

// Parse implements the Model interface for Alert.
func (a *Alert) Parse(ctx Ctx) (Tree, error) {
	return ctx.Constant(a.Literal)
}

// PubMap returns an ordered map of the Alert's public fields.
func (t *Alert) PubMap() *OrderedMap { return util.PubMapOf(t) }

// AsJSON returns a JSON-compatible representation of the Alert.
func (t *Alert) AsJSON() any { return asjson.AsJSONOf(t) }

// AsJSONStr returns a JSON string representation of the Alert.
func (t *Alert) AsJSONStr() string { return asjson.AsJSONStr(t.AsJSON()) }
