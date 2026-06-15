// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"

	"github.com/neogeny/ogopego/pkg/trees"
)

// Call represents a call to a grammar rule.
type Call struct {
	ModelBase
	Name string
	rule *Rule
}

// Parse implements the Model interface for Call.
func (c *Call) Parse(ctx Ctx) (any, error) {
	if c.rule == nil {
		return nil, ctx.Failure(ctx.Mark(), fmt.Errorf("call to %q has not been linked", c.Name))
	}

	name := c.Name
	rule := c.rule
	return call(ctx, name, rule)
}

func call(ctx Ctx, name string, rule *Rule) (any, error) {
	start := ctx.Mark()

	if !rule.IsToken() {
		ctx.NextToken()
	}

	key := ctx.Key(name, rule.IsMemoizable())

	if rule.ShouldTrace() {
		ctx.Tracer().TraceEntry(ctx)
	}

	ctx.Enter(name)
	result, err := ruleCall(ctx, name, rule, key, start)
	ctx.Leave()

	if err != nil {
		ctx.Tracer().TraceFailure(ctx, name)
		ctx.Memoize(key, trees.BOTTOM, start)
		return nil, err
	}

	if rule.IsName {
		if text, ok := result.(string); ok && ctx.IsKeyword(text) {
			ctx.Memoize(key, trees.BOTTOM, start)
			if rule.ShouldTrace() {
				ctx.Tracer().TraceFailure(ctx, text)
			}
			return nil, ctx.Failure(start, fmt.Errorf("'%s' is a reserved word", text))
		}
	}

	ctx.Memoize(key, result, ctx.Mark())
	if rule.ShouldTrace() {
		ctx.Tracer().TraceSuccess(ctx)
	}
	ctx.HeartbeatTick()

	return result, nil
}

func ruleCall(ctx Ctx, name string, rule *Rule, key MemoKey, start int) (any, error) {
	if memo, ok := ctx.Memo(key); ok {
		ctx.Reset(memo.Mark)
		if memo.Tree == trees.BOTTOM {
			return nil, ctx.Failure(start, fmt.Errorf("failed parsing %q", name))
		}
		return memo.Tree, nil
	}

	if rule.IsLeftRecursive() {
		return recursiveCall(ctx, name, rule, key, start)
	}

	return rule.Parse(ctx)
}

func recursiveCall(ctx Ctx, name string, rule *Rule, key MemoKey, start int) (any, error) {
	ctx.Tracer().TraceRecursion(ctx)

	ctx.Memoize(key, trees.BOTTOM, start)

	var lastMark = start
	var lastTree any
	var lastErr error

	for {
		ctx.Reset(start)

		if err := ctx.TrackRecursionDepth(key); err != nil {
			return nil, ctx.Failure(start, err)
		}

		result, err := rule.Parse(ctx)

		ctx.Untrack(key)

		if err != nil {
			lastErr = err
			break
		}

		endMark := ctx.Mark()
		if endMark <= lastMark {
			break
		}

		lastMark = endMark
		lastTree = result
		ctx.Memoize(key, lastTree, lastMark)
	}

	ctx.Reset(lastMark)
	ctx.Memoize(key, lastTree, lastMark)

	if lastTree == nil {
		return nil, lastErr
	}
	if lastTree == trees.BOTTOM {
		return nil, lastErr
	}
	return lastTree, nil
}
