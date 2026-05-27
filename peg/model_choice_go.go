package peg

import (
	"errors"
	"fmt"
	"strings"
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

	// 1. Launch concurrent branches
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

	// 2. Evaluate in strict priority order
	for i := 0; i < numOptions; i++ {
		res := <-chans[i]

		if res.Err == nil {
			// Even if res.Tree is nil, a nil error confirms this branch won
			ctx.Reset(clones[i].Mark())
			return res.Tree, nil
		}
	}

	// 3. Fallback: handle failures...
	msg := "no option matched"
	if len(c.la) > 0 {
		msg = fmt.Sprintf("expecting %s", strings.Join(c.la, ", "))
	}
	lastErr := ctx.Failure(startMark, errors.New(msg))
	ctx.Reset(startMark)
	return nil, lastErr
}
