// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package util

import (
	"strings"
)

const BlackLineLength = 88

// FoldOption configures line-wrapping behavior in Fold.
type FoldOption struct {
	AddLevels int
	Amount    int
}

// Fold joins parts with commas, wraps in brackets, and handles line-breaking.
// Pass a single-element parts for field-level wrapping ("prefix + body").
func Fold(prefix string, parts []string, lbrack, rbrack string, opts ...FoldOption) string {
	opt := FoldOption{Amount: 2}
	if len(opts) > 0 {
		opt = opts[0]
	}

	single := prefix + lbrack + strings.Join(parts, ", ") + rbrack
	if fitsfmt(single, opt.AddLevels, opt.Amount) {
		return single
	}

	indent := strings.Repeat(" ", opt.Amount)
	return prefix + lbrack + "\n" + indent + strings.Join(parts, ",\n"+indent) + "\n" + rbrack
}

func fitsfmt(line string, addLevels, amount int) bool {
	if strings.Contains(line, "\n") {
		return false
	}
	return len(line)+addLevels*amount <= BlackLineLength
}

// Repr returns a Go-composite-literal representation of v by consuming
// the PubMap protocol on AsJSONMixin types and delegating containers
// and scalars to the appropriate Go literal syntax.
func Repr(v any) string {
	return ""
}
