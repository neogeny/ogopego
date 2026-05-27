// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"errors"
	"fmt"
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

func (c *Choice) Parse(ctx Ctx) (Tree, error) {
	startMark := ctx.Mark()
	failure := func() error {
		msg := "no option matched"
		if len(c.la) > 0 {
			msg = fmt.Sprintf("expecteing one of: %s", c.LookAheadStr())
		}
		return ctx.Failure(startMark, errors.New(msg))
	}
	for _, opt := range c.Options {
		ctx.Reset(startMark)

		ctx.CutStackPush()
		result, err := opt.Exp.Parse(ctx)
		cutSeen := ctx.CutStackPop()

		if err == nil {
			return result, nil
		}

		if cutSeen {
			ctx.Reset(startMark)
			return nil, failure()
		}
	}
	ctx.Reset(startMark)
	return nil, failure()
}
