// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestBootGrammarCfgFromDirectives(t *testing.T) {
	g, err := BootGrammar()
	assert.NoError(t, err, "LoadBootGrammar")

	cfg := *g.CfgFromDirectives()

	assert.Equal(t, "TatSuBootstrap", cfg.Grammar)
	assert.Equal(t, "", cfg.Name)
	assert.Equal(t, nil, cfg.Whitespace, "expected Whitespace pattern")
	// NOTE Now TatSu relies on the default whitespace definition
	// assert.Equal(t, `(?m)\s+`, *cfg.Whitespace)
	assert.NotZero(t, cfg.Comments, "expected Comments pattern")
	assert.NotZero(t, cfg.EolComments, "expected EolComments pattern")
	assert.True(t, cfg.ParseInfo, "expected ParseInfo to be true")
	assert.True(t, cfg.NoLeftRecursion, "expected NoLeftRecursion to be true (from left_recursion: false)")
	assert.False(t, cfg.IgnoreCase, "expected IgnoreCase to be false")
	assert.Equal(t, "", cfg.Source)
	assert.Zero(t, cfg.Keywords, "expected Keywords to be nil")
	assert.NotZero(t, cfg.Semantics, "expected Semantics to be non-nil")
}
