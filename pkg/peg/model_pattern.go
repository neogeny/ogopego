// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

// Pattern matches input according to a configured pattern and returns text.
type Pattern struct {
	ModelBase
	Pattern string
}

// Parse implements the Model interface for Pattern.
func (p *Pattern) Parse(ctx Ctx) (any, error) {
	matched, err := ctx.MatchPattern(p.Pattern)
	if err != nil {
		return nil, err
	}
	return matched, nil
}
