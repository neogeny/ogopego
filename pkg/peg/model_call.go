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
func (c *Call) Parse(ctx Ctx) (trees.Tree, error) {
	if c.rule == nil {
		return nil, ctx.Failure(ctx.Mark(), fmt.Errorf("call to %q has not been linked", c.Name))
	}

	name := c.Name
	rule := c.rule
	start := ctx.Mark()

	if !rule.IsToken() {
		ctx.NextToken()
	}

	key := ctx.Key(name, rule.IsMemoizable())

	if rule.ShouldTrace() {
		ctx.Enter(name)
		ctx.Tracer().TraceEntry(ctx)
	}

	result, err := c.doCall(ctx, name, rule, key, start)

	if rule.ShouldTrace() {
		ctx.Leave()
	}

	if err != nil {
		ctx.Tracer().TraceFailure(ctx, name)
		ctx.Memoize(key, trees.BOTTOM, ctx.Mark())
		return nil, err
	}

	if rule.IsName {
		if text, ok := result.(*trees.Text); ok && ctx.IsKeyword(text.Value) {
			ctx.Memoize(key, trees.BOTTOM, ctx.Mark())
			ctx.Tracer().TraceFailure(ctx, text.Value)
			return nil, ctx.Failure(start, fmt.Errorf("'%s' is a reserved word", text.Value))
		}
	}

	ctx.Memoize(key, result, ctx.Mark())
	ctx.Tracer().TraceSuccess(ctx)
	ctx.HeartbeatTick()

	return result, nil
}

func (c *Call) doCall(ctx Ctx, name string, rule *Rule, key MemoKey, start int) (trees.Tree, error) {
	if memo, ok := ctx.Memo(key); ok {
		ctx.Reset(memo.Mark)
		if _, isBottom := memo.Tree.(*trees.Bottom); isBottom {
			return nil, ctx.Failure(start, fmt.Errorf("failed parsing %q", name))
		}
		return memo.Tree, nil
	}

	if rule.IsLeftRecursive() {
		return c.callRecursive(ctx, name, rule, key, start)
	}

	return rule.Parse(ctx)
}

func (c *Call) callRecursive(ctx Ctx, name string, rule *Rule, key MemoKey, start int) (trees.Tree, error) {
	ctx.Tracer().TraceRecursion(ctx)

	ctx.Memoize(key, trees.BOTTOM, start)

	var lastMark = start
	var lastTree trees.Tree
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
	if _, isBottom := lastTree.(*trees.Bottom); isBottom {
		return nil, lastErr
	}
	return lastTree, nil
}
