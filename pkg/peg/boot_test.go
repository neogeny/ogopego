// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestLoadBootGrammar(t *testing.T) {
	g, err := BootGrammar()
	assert.NoError(t, err, "LoadBootGrammar")
	assert.NotZero(t, g.Name, "expected non-empty grammar name")
	assert.True(t, g.Analyzed, "expected grammar to be analyzed after Initialize")
	assert.True(t, len(g.Rules) >= 10, "expected at least 10 rules, got %d", len(g.Rules))
	required := []string{"start", "grammar", "rule", "expre", "choice", "sequence"}
	for _, name := range required {
		rule, err := g.GetRule(name)
		assert.NoError(t, err, "missing required rule %q", name)
		assert.NotZero(t, rule, "GetRule(%q) returned nil", name)
	}
}
