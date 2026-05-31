// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/asjson"
)

func TestNodeAsJSON(t *testing.T) {
	n := &Node{
		Ast: "test-ast",
		ParseInfo: &ParseInfo{
			Source:  "test",
			Rule:    "start",
			Start:   0,
			Mark:    4,
			Line:    1,
			EndLine: 1,
		},
	}
	result := asjson.AsJSON(n)
	om, ok := result.(map[string]any)
	assert.True(t, ok, "expected map, got %T", result)

	cls, _ := om["__class__"]
	assert.Equal(t, "peg.Node", cls)

	_, hasPI := om["parse_info"]
	assert.True(t, hasPI, "expected non-nil ParseInfo")
}
