// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"errors"
)

// Option represents a single alternative within a Choice expression.
type Option struct {
	ModelBase
	Exp Model
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
