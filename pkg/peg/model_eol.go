// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
)

// EOL matches an end-of-line sequence.
type EOL struct {
	ModelBase
}

// Parse implements the Model interface for EOL.
func (e *EOL) Parse(ctx Ctx) (any, error) {
	if !ctx.MatchEOL() {
		return nil, ctx.Failure(
			ctx.Mark(),
			fmt.Errorf("expected EOL"),
		)
	}
	return nil, nil
}
