// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
)

// Token matches a literal string token.
type Token struct {
	ModelBase
	Token string
}

// Parse implements the Model interface for Token.
func (t *Token) Parse(ctx Ctx) (any, error) {
	mark := ctx.Mark()
	if !ctx.MatchToken(t.Token) {
		ctx.Reset(mark)
		return nil, ctx.Failure(mark, fmt.Errorf("expected: %q", t.Token))
	}
	return t.Token, nil
}
