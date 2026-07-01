// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package context

import (
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/input"
	"github.com/neogeny/ogopego/pkg/util"
)

func TestMementoNew(t *testing.T) {
	cursor := input.NewStrCursor("hello world\nsecond line")
	cs := util.NewTokenStack()
	cs.Push("rule1")
	cs.Push("rule2")

	m := NewMemento(0, "expected token", cursor, cs)
	assert.Equal(t, "expected token", m.Msg)
	assert.Equal(t, "some input", m.InputSource())
	assert.Equal(t, 2, m.CallStack.Len(), "expected 2 callstack entries")
}

func TestMementoErrorContainsMsg(t *testing.T) {
	cursor := input.NewStrCursor("hello world")
	m := NewMemento(0, "unexpected character", cursor, util.NewTokenStack())
	err := m.Error()
	assert.True(t, strings.Contains(err, "unexpected character"), "expected msg in error output, got: %s", err)
}

func TestMementoErrorShowsSourceText(t *testing.T) {
	cursor := input.NewStrCursor("hello world")
	m := NewMemento(0, "test error", cursor, util.NewTokenStack())
	err := m.Error()
	assert.True(t, strings.Contains(err, "hello world"), "expected source text in error output, got: %s", err)
}

func TestMementoErrorShowsPosition(t *testing.T) {
	cursor := input.NewStrCursor("hello world")
	// Advance past "hello "
	for i := 0; i < 6; i++ {
		cursor.Next()
	}
	m := NewMemento(6, "expected 'world'", cursor, util.NewTokenStack())
	err := m.Error()
	assert.True(t, strings.Contains(err, "world"), "expected position context in error output, got: %s", err)
}

func TestMementoErrorWithCallstack(t *testing.T) {
	cursor := input.NewStrCursor("x")
	cs := util.NewTokenStack()
	cs.Push("expr")
	cs.Push("term")
	cs.Push("factor")
	m := NewMemento(0, "parse error", cursor, cs)
	err := m.Error()
	// Check each callstack entry appears in the error output
	assert.True(t, strings.Contains(err, "expr"), "expected callstack entry %q in output, got: %s", "expr", err)
	assert.True(t, strings.Contains(err, "term"), "expected callstack entry %q in output, got: %s", "term", err)
	assert.True(t, strings.Contains(err, "factor"), "expected callstack entry %q in output, got: %s", "factor", err)
}

func TestMementoErrorMultiLineSource(t *testing.T) {
	src := "first line\nsecond line\nthird line\nfourth line"
	cursor := input.NewStrCursor(src)
	idx := strings.Index(src, "fourth")
	cursor.Reset(idx)
	m := NewMemento(idx, "error here", cursor, util.NewTokenStack())
	err := m.Error()
	assert.True(t, strings.Contains(err, "fourth line"), "expected error line context, got: %s", err)
	assert.True(t, strings.Contains(err, "⌃"), "expected caret marker, got: %s", err)
}

func TestMementoString(t *testing.T) {
	cursor := input.NewStrCursor("test input")
	m := NewMemento(0, "error msg", cursor, util.Atom("rule"))
	assert.Equal(t, m.Error(), m.String(), "String() and Error() should match")
}

func TestMementoEmptyCallstack(t *testing.T) {
	cursor := input.NewStrCursor("test")
	m := NewMemento(0, "err", cursor, util.NewTokenStack())
	err := m.Error()
	assert.True(t, strings.Contains(err, "err"), "expected error message in output, got: %s", err)
}
