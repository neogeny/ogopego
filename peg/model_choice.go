// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
	for i, opt := range c.Options {
		ctx.Reset(startMark)

		ctx.CutStackPush()
		result, err := opt.Exp.Parse(ctx)
		cutSeen := ctx.CutStackPop()

		if err == nil {
			return result, nil
		}

		if cutSeen {
			fmt.Fprintf(os.Stderr, "CUT: la=%q pos=%d opt=%d err=%v\n", c.la, startMark, i, err)
			return nil, err
		}
	}
	msg := "no option matched"
	if len(c.la) > 0 {
		msg = fmt.Sprintf("expecteing %s", strings.Join(c.la, ", "))
	}
	lastErr := ctx.Failure(startMark, errors.New(msg))
	ctx.Reset(startMark)
	return nil, lastErr
}
