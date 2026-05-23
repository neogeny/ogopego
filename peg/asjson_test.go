// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"testing"

	asjson "github.com/neogeny/ogopego/json"
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
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if cls, _ := om["__class__"]; cls != "peg.Node" {
		t.Errorf("expected __class__ peg.Node, got %v", cls)
	}
	parseinfoRaw, hasPI := om["parse_info"]
	if !hasPI {
		t.Fatal("expected non-nil ParseInfo")
	}
	_ = parseinfoRaw
}
