// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/trees"
)

func TestParseGrammar(t *testing.T) {
	result, err := ParseGrammar("@@grammar :: Test start := 'x'", nil)
	assert.NoError(t, err, "ParseGrammar error")
	assert.NotZero(t, result, "expected non-nil Tree")
	_, ok := result.(*trees.Nil)
	assert.False(t, ok, "unexpected Nil tree")
}

func TestCompile(t *testing.T) {
	src := "@@grammar :: EBNFTest\nstart := expression $\nexpression := expression '+' term | expression '-' term | term\nterm := term '*' factor | term '/' factor | factor\nfactor := '(' expression ')' | number\nnumber := /\\d+/\n"
	g, err := Compile(src, nil)
	assert.NoError(t, err, "Compile error")
	assert.Equal(t, "EBNFTest", g.Name, "expected name 'EBNFTest', got %q", g.Name)
	assert.True(t, g.Analyzed, "expected analyzed grammar")
	assert.Equal(t, 5, len(g.Rules), "expected 5 rules, got %d", len(g.Rules))
}

func TestCompileToJSON(t *testing.T) {
	src := `start := 'x'`
	json, err := CompileToJSON(src, nil)
	assert.NoError(t, err, "CompileToJSON error")
	assert.NotZero(t, json, "expected non-nil json")
}

func TestCompileToJSONString(t *testing.T) {
	src := `start := 'x'`
	s, err := CompileToJSONString(src, nil)
	assert.NoError(t, err, "CompileToJSONString error")
	assert.NotZero(t, s, "expected non-empty string")
}

func TestBootGrammar(t *testing.T) {
	g, err := BootGrammar()
	assert.NoError(t, err, "BootGrammar error")
	assert.True(t, g.Analyzed, "expected analyzed grammar")
}
