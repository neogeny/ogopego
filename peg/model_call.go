package peg

import (
	"fmt"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Call struct {
	ModelBase
	Name   string
	Target *Rule
}

func (c *Call) Parse(ctx Ctx) (trees.Tree, error) {
	if c.Target == nil {
		return nil, ctx.Failure(ctx.Mark(), fmt.Errorf("call to %q has not been linked", c.Name))
	}

	name := c.Name
	rule := c.Target
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

func (c *Call) PubMap() *asjson.OrderedMap { return c.PubMapOf(c) }
func (c *Call) AsJSON() any                { return c.AsJSONOf(c) }
func (c *Call) AsJSONStr() string          { return c.AsJSONStrOf(c) }
