// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package context

import (
	"strings"
	"testing"

	"github.com/neogeny/ogopego/input"
)

func TestMementoNew(t *testing.T) {
	cursor := input.NewStrCursor("hello world\nsecond line")
	cs := []string{"rule1", "rule2"}

	m := NewMemento(0, "expected token", cursor, cs)
	if m.Msg != "expected token" {
		t.Errorf("Msg = %q, want %q", m.Msg, "expected token")
	}
	if m.InputSource() != "some input" {
		t.Errorf("InputSource = %q, want %q", m.InputSource(), "some input")
	}
	if len(m.CallStack) != 2 {
		t.Errorf("expected 2 callstack entries, got %d", len(m.CallStack))
	}
}

func TestMementoErrorContainsMsg(t *testing.T) {
	cursor := input.NewStrCursor("hello world")
	m := NewMemento(0, "unexpected character", cursor, nil)
	err := m.Error()
	if !strings.Contains(err, "unexpected character") {
		t.Errorf("expected msg in error output, got: %s", err)
	}
}

func TestMementoErrorShowsSourceText(t *testing.T) {
	cursor := input.NewStrCursor("hello world")
	m := NewMemento(0, "test error", cursor, nil)
	err := m.Error()
	if !strings.Contains(err, "hello world") {
		t.Errorf("expected source text in error output, got: %s", err)
	}
}

func TestMementoErrorShowsPosition(t *testing.T) {
	cursor := input.NewStrCursor("hello world")
	// Advance past "hello "
	for i := 0; i < 6; i++ {
		cursor.Next()
	}
	m := NewMemento(6, "expected 'world'", cursor, nil)
	err := m.Error()
	if !strings.Contains(err, "world") {
		t.Errorf("expected position context in error output, got: %s", err)
	}
}

func TestMementoErrorWithCallstack(t *testing.T) {
	cursor := input.NewStrCursor("x")
	cs := []string{"expr", "term", "factor"}
	m := NewMemento(0, "parse error", cursor, cs)
	err := m.Error()
	for _, call := range cs {
		if !strings.Contains(err, call) {
			t.Errorf("expected callstack entry %q in output, got: %s", call, err)
		}
	}
}

func TestMementoErrorMultiLineSource(t *testing.T) {
	src := "first line\nsecond line\nthird line\nfourth line"
	cursor := input.NewStrCursor(src)
	// Position at "fourth line" area - byte offset after "third line\n"
	idx := strings.Index(src, "fourth")
	cursor.Reset(idx)
	m := NewMemento(idx, "error here", cursor, nil)
	err := m.Error()
	// Should show "fourth line" context near the error
	if !strings.Contains(err, "fourth line") {
		t.Errorf("expected error line context, got: %s", err)
	}
	if !strings.Contains(err, "^") {
		t.Errorf("expected caret marker, got: %s", err)
	}
}

func TestMementoString(t *testing.T) {
	cursor := input.NewStrCursor("test input")
	m := NewMemento(0, "error msg", cursor, []string{"rule"})
	if m.String() != m.Error() {
		t.Error("String() and Error() should match")
	}
}

func TestMementoEmptyCallstack(t *testing.T) {
	cursor := input.NewStrCursor("test")
	m := NewMemento(0, "err", cursor, nil)
	err := m.Error()
	if !strings.Contains(err, "err") {
		t.Errorf("expected error message in output, got: %s", err)
	}
}
