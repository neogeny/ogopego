// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"errors"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/util"
)

// Option represents a single alternative within a Choice expression.
type Option struct {
	Box
}

// Choice represents a PEG choice expression (ordered choice).
type Choice struct {
	ModelBase
	Options []*Option
}

// Parse implements the Model interface for Option.
func (o *Option) Parse(ctx Ctx) (Tree, error) {
	result, err := o.Exp.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Parse implements the Model interface for Choice.
func (c *Choice) Parse(ctx Ctx) (Tree, error) {
	startMark := ctx.Mark()
	var lastErr error
	for _, opt := range c.Options {
		mark := ctx.Mark()
		ctx.CutStackPush()
		result, err := opt.Parse(ctx)
		ctx.CutStackPop()

		if err == nil {
			return result, nil
		}

		ctx.Reset(mark)
		//if cutSeen {
		//	return nil, err
		//}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = ctx.Failure(startMark, errors.New("no option matched"))
	}
	ctx.Reset(startMark)
	return nil, lastErr
}

// PubMap returns an ordered map of the Choice's public fields.
func (c *Choice) PubMap() *OrderedMap { return util.PubMapOf(c) }

// AsJSON returns a JSON-compatible representation of the Choice.
func (c *Choice) AsJSON() any { return asjson.AsJSONOf(c) }

// AsJSONStr returns a JSON string representation of the Choice.
func (c *Choice) AsJSONStr() string { return asjson.AsJSONStr(c.AsJSON()) }

// PubMap returns an ordered map of the Option's public fields.
func (o *Option) PubMap() *OrderedMap { return util.PubMapOf(o) }

// AsJSON returns a JSON-compatible representation of the Option.
func (o *Option) AsJSON() any { return asjson.AsJSONOf(o) }

// AsJSONStr returns a JSON string representation of the Option.
func (o *Option) AsJSONStr() string { return asjson.AsJSONStr(o.AsJSON()) }
