package peg

import (
	"errors"
	"fmt"
)

type Result struct {
	Tree Tree
	Err  error
	Cut  bool
}

func (c *Choice) ParsePar(ctx Ctx) (Tree, error) {
	startMark := ctx.Mark()
	numOptions := len(c.Options)

	chans := make([]chan Result, numOptions)
	clones := make([]Ctx, numOptions)

	for i, opt := range c.Options {
		chans[i] = make(chan Result, 1)
		clones[i] = ctx.Clone()

		go func(ch chan Result, o *Option, wCtx Ctx) {
			wCtx.CutStackPush()
			tree, err := o.Exp.Parse(wCtx)
			cutSeen := wCtx.CutStackPop()

			ch <- Result{Tree: tree, Err: err, Cut: cutSeen}
		}(chans[i], opt, clones[i])
	}

	for i := 0; i < numOptions; i++ {
		res := <-chans[i]

		if res.Err == nil {
			ctx.Merge(clones[i])
			return res.Tree, nil
		} else if res.Cut {
			return nil, res.Err
		}
	}
	ctx.Reset(startMark)
	msg := "no option matched"
	if len(c.la) > 0 {
		msg = fmt.Sprintf("expecting one of: %s", c.LookAheadStr())
	}
	lastErr := ctx.Failure(startMark, errors.New(msg))
	return nil, lastErr
}
