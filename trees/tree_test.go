package trees

import (
	"testing"

	asjson "github.com/neogeny/ogopego/json"
)

func text(s string) *Text      { return &Text{Value: s} }
func seq(items ...Tree) *Seq   { return &Seq{Items: items} }
func list(items ...Tree) *List { return &List{Items: items} }

func TestFoldNil(t *testing.T) {
	result := Fold(&Nil{})
	if _, ok := result.(*Nil); !ok {
		t.Errorf("expected Nil, got %T", result)
	}
}

func TestFoldBottom(t *testing.T) {
	result := Fold(&Bottom{})
	if _, ok := result.(*Bottom); !ok {
		t.Errorf("expected Bottom, got %T", result)
	}
}

func TestFoldGoNil(t *testing.T) {
	result := Fold(nil)
	if _, ok := result.(*Nil); !ok {
		t.Errorf("expected Nil, got %T", result)
	}
}

func TestFoldText(t *testing.T) {
	result := Fold(text("hello"))
	txt, ok := result.(*Text)
	if !ok || txt.Value != "hello" {
		t.Errorf("expected Text(hello), got %T %v", result, result)
	}
}

func TestFoldBool(t *testing.T) {
	result := Fold(&Bool{Value: true})
	b, ok := result.(*Bool)
	if !ok || b.Value != true {
		t.Errorf("expected Bool(true), got %T %v", result, result)
	}
}

func TestFoldNumber(t *testing.T) {
	result := Fold(&Number{Value: 42.5})
	n, ok := result.(*Number)
	if !ok || n.Value != 42.5 {
		t.Errorf("expected Number(42.5), got %T %v", result, result)
	}
}

func TestFoldSeqToSeq(t *testing.T) {
	// Seq with no Named/Override becomes List after fold (closed)
	result := Fold(seq(text("a"), text("b"), text("c")))
	l, ok := result.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", result)
	}
	if len(l.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(l.Items))
	}
	if l.Items[0].(*Text).Value != "a" {
		t.Errorf("expected 'a', got %v", l.Items[0])
	}
}

func TestFoldNamedToMap(t *testing.T) {
	result := Fold(&Named{Name: "x", Value: text("hello")})
	m, ok := result.(*MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", result)
	}
	if m.Entries["x"] == nil {
		t.Fatal("expected key 'x'")
	}
	if m.Entries["x"].(*Text).Value != "hello" {
		t.Errorf("expected 'hello', got %v", m.Entries["x"])
	}
}

func TestFoldOverride(t *testing.T) {
	result := Fold(&Override{Value: text("result")})
	txt, ok := result.(*Text)
	if !ok {
		t.Fatalf("expected Text, got %T", result)
	}
	if txt.Value != "result" {
		t.Errorf("expected 'result', got %v", txt.Value)
	}
}

func TestFoldMultipleNamed(t *testing.T) {
	result := Fold(seq(
		&Named{Name: "a", Value: text("1")},
		&Named{Name: "b", Value: text("2")},
	))
	m, ok := result.(*MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", result)
	}
	if m.Entries["a"].(*Text).Value != "1" {
		t.Errorf("expected '1', got %v", m.Entries["a"])
	}
	if m.Entries["b"].(*Text).Value != "2" {
		t.Errorf("expected '2', got %v", m.Entries["b"])
	}
}

func TestFoldNamedAccumulates(t *testing.T) {
	result := Fold(seq(
		&Named{Name: "x", Value: text("a")},
		&Named{Name: "x", Value: text("b")},
	))
	m, ok := result.(*MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", result)
	}
	s := m.Entries["x"].(*Seq)
	if s.Items[0].(*Text).Value != "a" {
		t.Errorf("expected 'a', got %v", s.Items[0])
	}
	if s.Items[1].(*Text).Value != "b" {
		t.Errorf("expected 'b', got %v", s.Items[1])
	}
}

func TestFoldNamedAsList(t *testing.T) {
	result := Fold(&NamedAsList{Name: "items", Value: text("x")})
	m, ok := result.(*MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", result)
	}
	s := m.Entries["items"].(*Seq)
	if len(s.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(s.Items))
	}
	if s.Items[0].(*Text).Value != "x" {
		t.Errorf("expected 'x', got %v", s.Items[0])
	}
}

func TestFoldNamedAsListAccumulates(t *testing.T) {
	result := Fold(seq(
		&NamedAsList{Name: "items", Value: text("a")},
		&NamedAsList{Name: "items", Value: text("b")},
	))
	m, ok := result.(*MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", result)
	}
	s := m.Entries["items"].(*Seq)
	if len(s.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(s.Items))
	}
	if s.Items[0].(*Text).Value != "a" {
		t.Errorf("expected 'a', got %v", s.Items[0])
	}
	if s.Items[1].(*Text).Value != "b" {
		t.Errorf("expected 'b', got %v", s.Items[1])
	}
}

func TestFoldOverrideWins(t *testing.T) {
	result := Fold(seq(
		&Named{Name: "x", Value: text("ignored")},
		text("also ignored"),
		&Override{Value: text("result")},
	))
	txt, ok := result.(*Text)
	if !ok {
		t.Fatalf("expected Text, got %T", result)
	}
	if txt.Value != "result" {
		t.Errorf("expected 'result', got %v", txt.Value)
	}
}

func TestFoldOverrideAsList(t *testing.T) {
	result := Fold(seq(
		&OverrideAsList{Value: text("a")},
		&OverrideAsList{Value: text("b")},
	))
	l, ok := result.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", result)
	}
	if len(l.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(l.Items))
	}
}

func TestFoldNestedNamed(t *testing.T) {
	result := Fold(&Named{
		Name: "x",
		Value: seq(
			&Named{Name: "a", Value: text("1")},
			&Named{Name: "b", Value: text("2")},
		),
	})
	m, ok := result.(*MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", result)
	}
	if _, exists := m.Entries["x"]; !exists {
		t.Fatal("expected key 'x'")
	}
	if _, exists := m.Entries["a"]; !exists {
		t.Fatal("expected key 'a'")
	}
	if _, exists := m.Entries["b"]; !exists {
		t.Fatal("expected key 'b'")
	}
}

func TestFoldSeqWithNil(t *testing.T) {
	result := Fold(seq(text("a"), &Nil{}, text("b")))
	l, ok := result.(*List)
	if !ok {
		t.Fatalf("expected List, got %T", result)
	}
	if len(l.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(l.Items))
	}
}

func TestTextAsJSON(t *testing.T) {
	result := asjson.AsJSON(&Text{Value: "hello"})
	s, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	if s != "hello" {
		t.Errorf("expected 'hello', got %v", s)
	}
}

func TestNumberAsJSON(t *testing.T) {
	result := asjson.AsJSON(&Number{Value: 42.5})
	f, ok := result.(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", result)
	}
	if f != 42.5 {
		t.Errorf("expected 42.5, got %v", f)
	}
}

func TestNodeAsJSONTree(t *testing.T) {
	result := asjson.AsJSON(&Node{TypeName: "expr", Tree: text("42")})
	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if cls := m["__class__"]; cls != "expr" {
		t.Errorf("expected __class__=expr, got %v", cls)
	}
	if ast := m["ast"]; ast != "42" {
		t.Errorf("expected ast='42', got %v", ast)
	}
}

func TestFoldRuleNode(t *testing.T) {
	result := Fold(&Node{TypeName: "expr", Tree: text("42")})
	r, ok := result.(*Node)
	if !ok {
		t.Fatalf("expected RuleNode, got %T", result)
	}
	if r.TypeName != "expr" {
		t.Errorf("expected 'expr', got %q", r.TypeName)
	}
}
