// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

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

// Parse implements the Model interface for Alert.
func (a *Alert) Parse(ctx Ctx) (Tree, error) {
	return ctx.Constant(a.Literal)
}
