// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	_ "embed"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego"
)

func TestImportCalcJSON(t *testing.T) {
	data := ogopego.CalcJSON
	g, err := LoadGrammarFromJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, "CALC", g.Name)
	assert.NotZero(t, len(g.Rules), "expected at least one rule")
	startRule := g.Rules[0]
	assert.Equal(t, "start", startRule.Name)
	seq, ok := startRule.Exp.(*Sequence)
	assert.True(t, ok, "expected start rule exp to be Sequence, got %T", startRule.Exp)
	assert.Equal(t, 2, len(seq.Sequence), "expected 2 items in start sequence, got %d", len(seq.Sequence))
	call, ok := seq.Sequence[0].(*Call)
	assert.True(t, ok, "expected first item to be Call, got %T", seq.Sequence[0])
	assert.Equal(t, "expression", call.Name)
	eof, ok := seq.Sequence[1].(*EOF)
	assert.True(t, ok, "expected second item to be EOF, got %T", seq.Sequence[1])
	_ = eof
}

func TestImportTatsuJSON(t *testing.T) {
	data := ogopego.TatSuGrammarJSON
	g, err := LoadGrammarFromJSON(data)
	assert.NoError(t, err)
	assert.NotZero(t, g.Name, "expected non-empty grammar name")
	assert.True(t, len(g.Rules) >= 10, "expected at least 10 rules in TatSu grammar, got %d", len(g.Rules))
	// Check for key rules
	names := make(map[string]bool)
	for _, r := range g.Rules {
		names[r.Name] = true
	}
	for _, want := range []string{"start", "grammar", "rule", "expre", "choice", "sequence"} {
		assert.True(t, names[want], "missing rule %q", want)
	}
}

func TestImportInvalidJSON(t *testing.T) {
	_, err := LoadGrammarFromJSON([]byte(`{"__class__": "NotAGrammar"}`))
	assert.Error(t, err, "expected error for non-Grammar root")
}

func TestImportMalformedJSON(t *testing.T) {
	_, err := LoadGrammarFromJSON([]byte(`{not valid json`))
	assert.Error(t, err, "expected error for malformed JSON")
}
