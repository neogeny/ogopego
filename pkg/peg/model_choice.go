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
	Options []Model
}

// Parse implements the Model interface for Option.
func (o *Option) Parse(ctx Ctx) (any, error) {
	return o.Exp.Parse(ctx)
}

func (c *Choice) Parse(ctx Ctx) (any, error) {
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
		result, err := opt.Parse(ctx)
		cutSeen := ctx.CutStackPop()

		if err == nil {
			return result, nil
		}

		if cutSeen {
			ctx.Reset(startMark)
			return nil, err
		}
	}
	ctx.Reset(startMark)
	return nil, failure()
}
