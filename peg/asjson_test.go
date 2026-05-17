// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

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
	om, ok := result.(*asjson.OrderedMap)
	if !ok {
		t.Fatalf("expected *OrderedMap, got %T", result)
	}
	if cls, _ := om.Get("__class__"); cls != "Node" {
		t.Errorf("expected __class__ Node, got %v", cls)
	}
	// FIXME
	//if ast, hasAst := om.Get("ast"); hasAst {
	//	t.Errorf("expected no 'ast', got %v", ast)
	//}
	parseinfoRaw, _ := om.Get("parse_info")
	if parseinfoRaw == nil {
		t.Fatal("expected non-nil ParseInfo")
	}
}
