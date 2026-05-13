package peg

import (
	"testing"

	asjson "github.com/neogeny/ogopego/json"
)

func TestNodeAsJSON(t *testing.T) {
	n := &Node{
		Ast: "test-ast",
		Pos: &ParseInfo{Source: "test", Rule: "start", Pos: 0, EndPos: 4, Line: 1, EndLine: 1},
	}
	result := asjson.AsJSON(n)
	om, ok := result.(*asjson.OrderedMap)
	if !ok {
		t.Fatalf("expected *OrderedMap, got %T", result)
	}
	if cls, _ := om.Get("__class__"); cls != "Node" {
		t.Errorf("expected __class__ Node, got %v", cls)
	}
	if ast, _ := om.Get("Ast"); ast != "test-ast" {
		t.Errorf("expected Ast test-ast, got %v", ast)
	}
	posRaw, _ := om.Get("Pos")
	if posRaw == nil {
		t.Fatal("expected non-nil Pos")
	}
}
