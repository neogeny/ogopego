package test

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/peg"
	"github.com/neogeny/ogopego/tool"
	"github.com/neogeny/ogopego/trees"
)

func TestJavaEndToEnd(t *testing.T) {
	if os.Getenv("XONSH_VERSION") == "" {
		t.Skip("XONSH_VERSION not set — local test only")
	}
	// 1. Load pre-compiled Java grammar
	data, err := os.ReadFile("../grammar/java.json")
	if err != nil {
		t.Fatalf("read java.json: %v", err)
	}
	g, err := peg.ParseGrammar(data)
	if err != nil {
		t.Fatalf("parse grammar: %v", err)
	}
	if err := g.Initialize(); err != nil {
		t.Fatalf("init grammar: %v", err)
	}

	// 2. Verify tool.ModelRepr produces valid Go source
	code := tool.ModelRepr(*g, "java")
	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, "", code, parser.AllErrors)
	if err != nil {
		t.Fatalf("generated code is not valid Go:\n%s\n\nError: %v", code, err)
	}

	// 3. Parse a Java source snippet
	javaSrc := "package com.example;\nimport java.util.List;\npublic class Hello {}\n"
	tree, err := api.ParseInput(g, javaSrc, nil)
	if err != nil {
		t.Fatalf("parse Java: %v", err)
	}

	// 4. Walk the tree manually to validate structure
	n, ok := tree.(*trees.Node)
	if !ok {
		t.Fatalf("expected *trees.Node, got %T", tree)
	}
	if n.TypeName != "CompilationUnit" {
		t.Fatalf("expected TypeName CompilationUnit, got %q", n.TypeName)
	}
	m, ok := n.Tree.(*trees.MapNode)
	if !ok {
		t.Fatalf("expected MapNode, got %T", n.Tree)
	}

	// 5. Validate concrete types of each MapNode entry
	t.Run("package_field", func(t *testing.T) {
		v, ok := m.Entries["package"]
		if !ok {
			t.Fatal("missing 'package' entry")
		}
		pn, ok := v.(*trees.Node)
		if !ok {
			t.Fatalf("package: expected *trees.Node, got %T", v)
		}
		if pn.TypeName != "PackageDeclaration" {
			t.Fatalf("package TypeName: expected PackageDeclaration, got %q", pn.TypeName)
		}
		pm, ok := pn.Tree.(*trees.MapNode)
		if !ok {
			t.Fatalf("PackageDeclaration.Tree: expected MapNode, got %T", pn.Tree)
		}
		// name field is *trees.Node{QualifiedName}
		name, ok := pm.Entries["name"]
		if !ok {
			t.Fatal("PackageDeclaration missing 'name' entry")
		}
		qn, ok := name.(*trees.Node)
		if !ok {
			t.Fatalf("PackageDeclaration.name: expected *trees.Node, got %T", name)
		}
		if qn.TypeName != "QualifiedName" {
			t.Fatalf("PackageDeclaration.name TypeName: expected QualifiedName, got %q", qn.TypeName)
		}
	})

	t.Run("imports_field", func(t *testing.T) {
		v, ok := m.Entries["imports"]
		if !ok {
			t.Fatal("missing 'imports' entry")
		}
		impList, ok := v.(*trees.List)
		if !ok {
			t.Fatalf("imports: expected *trees.List, got %T", v)
		}
		if len(impList.Items) != 1 {
			t.Fatalf("imports: expected 1 item, got %d", len(impList.Items))
		}
		impNode, ok := impList.Items[0].(*trees.Node)
		if !ok {
			t.Fatalf("imports[0]: expected *trees.Node, got %T", impList.Items[0])
		}
		if impNode.TypeName != "ImportDeclaration" {
			t.Fatalf("imports[0] TypeName: expected ImportDeclaration, got %q", impNode.TypeName)
		}
	})

	t.Run("declarations_field", func(t *testing.T) {
		v, ok := m.Entries["declarations"]
		if !ok {
			t.Fatal("missing 'declarations' entry")
		}
		seq, ok := v.(*trees.Seq)
		if !ok {
			t.Fatalf("declarations: expected *trees.Seq, got %T", v)
		}
		if len(seq.Items) != 1 {
			t.Fatalf("declarations: expected 1 item, got %d", len(seq.Items))
		}
		classNode, ok := seq.Items[0].(*trees.Node)
		if !ok {
			t.Fatalf("declarations[0]: expected *trees.Node, got %T", seq.Items[0])
		}
		if classNode.TypeName != "ClassDeclaration" {
			t.Fatalf("declarations[0] TypeName: expected ClassDeclaration, got %q", classNode.TypeName)
		}
	})

	t.Run("linecount_field", func(t *testing.T) {
		v, ok := m.Entries["linecount"]
		if !ok {
			t.Fatal("missing 'linecount' entry")
		}
		if _, ok := v.(*trees.Nil); !ok {
			t.Fatalf("linecount: expected *trees.Nil, got %T", v)
		}
	})

	// 6. Test FromTree with hand-authored model types
	t.Run("from_tree_identifier", func(t *testing.T) {
		id, err := identifierFromTree(
			&trees.Node{TypeName: "Identifier", Tree: &trees.MapNode{
				Entries: map[string]trees.Tree{
					"value": &trees.Text{Value: "Hello"},
				},
			}},
		)
		if err != nil {
			t.Fatalf("IdentifierFromTree: %v", err)
		}
		txt, ok := id.Value.(*trees.Text)
		if !ok {
			t.Fatalf("Identifier.Value: expected *trees.Text, got %T", id.Value)
		}
		if txt.Value != "Hello" {
			t.Fatalf("Identifier.Value = %q, want %q", txt.Value, "Hello")
		}
	})

	t.Run("from_tree_qualified_name", func(t *testing.T) {
		// Build a tree for: com.example.List
		qn, err := qualifiedNameFromTree(
			&trees.Node{TypeName: "QualifiedName", Tree: &trees.MapNode{
				Entries: map[string]trees.Tree{
					"qualifiers": &trees.Seq{Items: []trees.Tree{
						&trees.Node{TypeName: "Identifier", Tree: &trees.MapNode{
							Entries: map[string]trees.Tree{
								"value": &trees.Text{Value: "com"},
							},
						}},
						&trees.Node{TypeName: "Identifier", Tree: &trees.MapNode{
							Entries: map[string]trees.Tree{
								"value": &trees.Text{Value: "example"},
							},
						}},
					}},
					"name": &trees.Node{TypeName: "Identifier", Tree: &trees.MapNode{
						Entries: map[string]trees.Tree{
							"value": &trees.Text{Value: "List"},
						},
					}},
				},
			}},
		)
		if err != nil {
			t.Fatalf("QualifiedNameFromTree: %v", err)
		}

		// Validate concrete types
		if len(qn.Qualifiers) != 2 {
			t.Fatalf("Qualifiers: len=%d, want 2", len(qn.Qualifiers))
		}
		com := qn.Qualifiers[0]
		comTxt, ok := com.Value.(*trees.Text)
		if !ok {
			t.Fatalf("Qualifiers[0].Value: expected *trees.Text, got %T", com.Value)
		}
		if comTxt.Value != "com" {
			t.Fatalf("Qualifiers[0].Value = %q, want %q", comTxt.Value, "com")
		}

		list := qn.Name
		listTxt, ok := list.Value.(*trees.Text)
		if !ok {
			t.Fatalf("Name.Value: expected *trees.Text, got %T", list.Value)
		}
		if listTxt.Value != "List" {
			t.Fatalf("Name.Value = %q, want %q", listTxt.Value, "List")
		}
	})

	t.Run("from_tree_package_via_optional", func(t *testing.T) {
		// Build a CompilationUnit tree with package present
		pkgTree := &trees.Node{TypeName: "PackageDeclaration", Tree: &trees.MapNode{
			Entries: map[string]trees.Tree{
				"annotations": &trees.List{Items: []trees.Tree{}},
				"name": &trees.Node{TypeName: "QualifiedName", Tree: &trees.MapNode{
					Entries: map[string]trees.Tree{
						"qualifiers": &trees.Seq{Items: []trees.Tree{
							&trees.Node{TypeName: "Identifier", Tree: &trees.MapNode{
								Entries: map[string]trees.Tree{
									"value": &trees.Text{Value: "com"},
								},
							}},
						}},
						"name": &trees.Node{TypeName: "Identifier", Tree: &trees.MapNode{
							Entries: map[string]trees.Tree{
								"value": &trees.Text{Value: "example"},
							},
						}},
					},
				}},
			},
		}}

		cu, err := compilationUnitFromTree(
			&trees.Node{TypeName: "CompilationUnit", Tree: &trees.MapNode{
				Entries: map[string]trees.Tree{
					"package":      pkgTree,
					"imports":      &trees.List{Items: []trees.Tree{}},
					"declarations": &trees.Seq{Items: []trees.Tree{}},
					"linecount":    &trees.Nil{},
				},
			}},
		)
		if err != nil {
			t.Fatalf("CompilationUnitFromTree: %v", err)
		}

		pkg, ok := cu.Package.(*PackageDeclaration)
		if !ok {
			t.Fatalf("Package: expected *PackageDeclaration, got %T", cu.Package)
		}
		if pkg.Name == nil {
			t.Fatal("Package.Name is nil")
		}
	})

	t.Run("from_tree_no_package", func(t *testing.T) {
		// Build a CompilationUnit tree WITHOUT package (simulating optional-nil)
		cu, err := compilationUnitFromTree(
			&trees.Node{TypeName: "CompilationUnit", Tree: &trees.MapNode{
				Entries: map[string]trees.Tree{
					"package":      &trees.Nil{},
					"imports":      &trees.List{Items: []trees.Tree{}},
					"declarations": &trees.Seq{Items: []trees.Tree{}},
					"linecount":    &trees.Nil{},
				},
			}},
		)
		if err != nil {
			t.Fatalf("CompilationUnitFromTree (no package): %v", err)
		}
		if cu.Package != nil {
			t.Fatal("expected nil Package when entry is Nil")
		}
	})

	// 7. Validate field types via reflect against the generated code structure
	t.Run("field_types_in_generated_code", func(t *testing.T) {
		// Check that the generated code uses the correct types for key fields
		tests := []struct {
			field   string
			wantAny bool
		}{
			{"Package *PackageDeclaration", false},
			{"Imports []*ImportDeclaration", false},
			{"Declarations []any", false},
			{"Linecount any", false},
			{"Value any", false},
		}
		for _, tt := range tests {
			if !contains(code, tt.field) {
				t.Errorf("generated code should contain %q", tt.field)
			}
		}

		// Verify the generated code references the expected types via reflect
		if !contains(code, "*PackageDeclaration") {
			t.Error("expected *PackageDeclaration in generated code")
		}
		if !contains(code, "[]*ImportDeclaration") {
			t.Error("expected []*ImportDeclaration in generated code")
		}
		if !contains(code, "IdentifierFromTree") {
			t.Error("expected IdentifierFromTree in generated code")
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// --- Hand-authored model types matching the Java grammar subset ---

type Identifier struct {
	Value any
}

func identifierFromTree(tree trees.Tree) (*Identifier, error) {
	n, ok := tree.(*trees.Node)
	if !ok {
		return nil, fmt.Errorf("IdentifierFromTree: expected *trees.Node, got %T", tree)
	}
	m, ok := n.Tree.(*trees.MapNode)
	if !ok {
		return nil, fmt.Errorf("IdentifierFromTree: expected MapNode, got %T", n.Tree)
	}
	var result Identifier
	if v, ok := m.Entries["value"]; ok {
		result.Value = v
	}
	return &result, nil
}

type QualifiedName struct {
	Qualifiers []*Identifier
	Name       *Identifier
}

func qualifiedNameFromTree(tree trees.Tree) (*QualifiedName, error) {
	n, ok := tree.(*trees.Node)
	if !ok {
		return nil, fmt.Errorf("QualifiedNameFromTree: expected *trees.Node, got %T", tree)
	}
	m, ok := n.Tree.(*trees.MapNode)
	if !ok {
		return nil, fmt.Errorf("QualifiedNameFromTree: expected MapNode, got %T", n.Tree)
	}
	var result QualifiedName
	var err error
	if v, ok := m.Entries["qualifiers"]; ok {
		seq, ok := v.(*trees.Seq)
		if !ok {
			return nil, fmt.Errorf("QualifiedName.Qualifiers: expected Seq, got %T", v)
		}
		result.Qualifiers = make([]*Identifier, len(seq.Items))
		for i, item := range seq.Items {
			result.Qualifiers[i], err = identifierFromTree(item)
			if err != nil {
				return nil, fmt.Errorf("QualifiedName.Qualifiers[%d]: %w", i, err)
			}
		}
	}
	if v, ok := m.Entries["name"]; ok {
		if n, ok := v.(*trees.Node); ok {
			result.Name, err = identifierFromTree(n)
			if err != nil {
				return nil, fmt.Errorf("QualifiedName.Name: %w", err)
			}
		}
	}
	return &result, nil
}

type PackageDeclaration struct {
	Annotations any
	Name        *QualifiedName
}

func packageDeclarationFromTree(tree trees.Tree) (*PackageDeclaration, error) {
	n, ok := tree.(*trees.Node)
	if !ok {
		return nil, fmt.Errorf("PackageDeclarationFromTree: expected *trees.Node, got %T", tree)
	}
	m, ok := n.Tree.(*trees.MapNode)
	if !ok {
		return nil, fmt.Errorf("PackageDeclarationFromTree: expected MapNode, got %T", n.Tree)
	}
	var result PackageDeclaration
	var err error
	if v, ok := m.Entries["annotations"]; ok {
		result.Annotations = v
	}
	if v, ok := m.Entries["name"]; ok {
		if n, ok := v.(*trees.Node); ok {
			result.Name, err = qualifiedNameFromTree(n)
			if err != nil {
				return nil, fmt.Errorf("PackageDeclaration.Name: %w", err)
			}
		}
	}
	return &result, nil
}

type ImportDeclaration struct {
	Static any
	Name   *QualifiedName
	All    any
}

func importDeclarationFromTree(tree trees.Tree) (*ImportDeclaration, error) {
	n, ok := tree.(*trees.Node)
	if !ok {
		return nil, fmt.Errorf("ImportDeclarationFromTree: expected *trees.Node, got %T", tree)
	}
	m, ok := n.Tree.(*trees.MapNode)
	if !ok {
		return nil, fmt.Errorf("ImportDeclarationFromTree: expected MapNode, got %T", n.Tree)
	}
	var result ImportDeclaration
	var err error
	if v, ok := m.Entries["static"]; ok {
		result.Static = v
	}
	if v, ok := m.Entries["name"]; ok {
		if n, ok := v.(*trees.Node); ok {
			result.Name, err = qualifiedNameFromTree(n)
			if err != nil {
				return nil, fmt.Errorf("ImportDeclaration.Name: %w", err)
			}
		}
	}
	if v, ok := m.Entries["all"]; ok {
		result.All = v
	}
	return &result, nil
}

type CompilationUnit struct {
	Package      any
	Imports      any
	Declarations any
	Linecount    any
}

func compilationUnitFromTree(tree trees.Tree) (*CompilationUnit, error) {
	n, ok := tree.(*trees.Node)
	if !ok {
		return nil, fmt.Errorf("CompilationUnitFromTree: expected *trees.Node, got %T", tree)
	}
	m, ok := n.Tree.(*trees.MapNode)
	if !ok {
		return nil, fmt.Errorf("CompilationUnitFromTree: expected MapNode, got %T", n.Tree)
	}
	var result CompilationUnit
	var err error
	if v, ok := m.Entries["package"]; ok {
		if n, ok := v.(*trees.Node); ok {
			var pkg *PackageDeclaration
			pkg, err = packageDeclarationFromTree(n)
			if err != nil {
				return nil, fmt.Errorf("CompilationUnit.Package: %w", err)
			}
			result.Package = pkg
		}
	}
	if v, ok := m.Entries["imports"]; ok {
		list, ok := v.(*trees.List)
		if !ok {
			return nil, fmt.Errorf("CompilationUnit.Imports: expected List, got %T", v)
		}
		imports := make([]*ImportDeclaration, len(list.Items))
		for i, item := range list.Items {
			imports[i], err = importDeclarationFromTree(item)
			if err != nil {
				return nil, fmt.Errorf("CompilationUnit.Imports[%d]: %w", i, err)
			}
		}
		result.Imports = imports
	}
	if v, ok := m.Entries["declarations"]; ok {
		result.Declarations = v
	}
	if v, ok := m.Entries["linecount"]; ok {
		result.Linecount = v
	}
	return &result, nil
}
