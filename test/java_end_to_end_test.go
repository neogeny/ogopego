package test

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/peg"
	"github.com/neogeny/ogopego/pkg/tool"
	"github.com/neogeny/ogopego/pkg/trees"
)

func TestJavaEndToEnd(t *testing.T) {
	if os.Getenv("XONSH_VERSION") == "" {
		t.Skip("XONSH_VERSION not set — local test only")
	}
	// 1. Load pre-compiled Java grammar
	data, err := os.ReadFile("../grammar/java.json")
	assert.NoError(t, err, "read java.json")
	g, err := peg.LoadGrammarFromJSON(data)
	assert.NoError(t, err, "parse grammar")
	assert.NoError(t, g.Initialize(), "init grammar")

	// 2. Verify tool.ModelRepr produces valid Go source
	code := tool.ModelRepr(*g, "java")
	fset := token.NewFileSet()
	_, err = parser.ParseFile(fset, "", code, parser.AllErrors)
	assert.NoError(t, err, "generated code is not valid Go:\n%s", code)

	// 3. Parse a Java source snippet
	javaSrc := "package com.example;\nimport java.util.List;\npublic class Hello {}\n"
	tree, err := api.ParseInput(g, javaSrc, nil)
	assert.NoError(t, err, "parse Java")

	// 4. Walk the tree manually to validate structure
	n, ok := tree.(*trees.Node)
	assert.True(t, ok, "expected *trees.Node, got %T", tree)
	assert.Equal(t, "CompilationUnit", n.TypeName)
	m, ok := n.Tree.(*trees.MapNode)
	assert.True(t, ok, "expected MapNode, got %T", n.Tree)

	// 5. Validate concrete types of each MapNode entry
	t.Run("package_field", func(t *testing.T) {
		v, ok := m.Entries["package"]
		assert.True(t, ok, "missing 'package' entry")
		pn, ok := v.(*trees.Node)
		assert.True(t, ok, "package: expected *trees.Node, got %T", v)
		assert.Equal(t, "PackageDeclaration", pn.TypeName, "package TypeName")
		pm, ok := pn.Tree.(*trees.MapNode)
		assert.True(t, ok, "PackageDeclaration.Tree: expected MapNode, got %T", pn.Tree)
		// name field is *trees.Node{QualifiedName}
		name, ok := pm.Entries["name"]
		assert.True(t, ok, "PackageDeclaration missing 'name' entry")
		qn, ok := name.(*trees.Node)
		assert.True(t, ok, "PackageDeclaration.name: expected *trees.Node, got %T", name)
		assert.Equal(t, "QualifiedName", qn.TypeName, "PackageDeclaration.name TypeName")
	})

	t.Run("imports_field", func(t *testing.T) {
		v, ok := m.Entries["imports"]
		assert.True(t, ok, "missing 'imports' entry")
		impList, ok := v.(*trees.Array)
		assert.True(t, ok, "imports: expected *trees.List, got %T", v)
		assert.Equal(t, 1, len(impList.Items), "imports: expected 1 item")
		impNode, ok := impList.Items[0].(*trees.Node)
		assert.True(t, ok, "imports[0]: expected *trees.Node, got %T", impList.Items[0])
		assert.Equal(t, "ImportDeclaration", impNode.TypeName, "imports[0] TypeName")
	})

	t.Run("declarations_field", func(t *testing.T) {
		v, ok := m.Entries["declarations"]
		assert.True(t, ok, "missing 'declarations' entry")
		seq, ok := v.(*trees.Seq)
		assert.True(t, ok, "declarations: expected *trees.Seq, got %T", v)
		assert.Equal(t, 1, len(seq.Items), "declarations: expected 1 item")
		classNode, ok := seq.Items[0].(*trees.Node)
		assert.True(t, ok, "declarations[0]: expected *trees.Node, got %T", seq.Items[0])
		assert.Equal(t, "ClassDeclaration", classNode.TypeName, "declarations[0] TypeName")
	})

	t.Run("linecount_field", func(t *testing.T) {
		v, ok := m.Entries["linecount"]
		assert.True(t, ok, "missing 'linecount' entry")
		_, ok = v.(*trees.Nil)
		assert.True(t, ok, "linecount: expected *trees.Nil, got %T", v)
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
		assert.NoError(t, err, "IdentifierFromTree")
		txt, ok := id.Value.(*trees.Text)
		assert.True(t, ok, "Identifier.Value: expected *trees.Text, got %T", id.Value)
		assert.Equal(t, "Hello", txt.Value, "Identifier.Value")
	})

	t.Run("from_tree_qualified_name", func(t *testing.T) {
		// Build a tree for: com.example.List
		qn, err := qualifiedNameFromTree(
			&trees.Node{TypeName: "QualifiedName", Tree: &trees.MapNode{
				Entries: map[string]trees.Tree{
					"qualifiers": &trees.Seq{Items: []any{
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
		assert.NoError(t, err, "QualifiedNameFromTree")

		// Validate concrete types
		assert.Equal(t, 2, len(qn.Qualifiers))
		com := qn.Qualifiers[0]
		comTxt, ok := com.Value.(*trees.Text)
		assert.True(t, ok, "Qualifiers[0].Value: expected *trees.Text, got %T", com.Value)
		assert.Equal(t, "com", comTxt.Value, "Qualifiers[0].Value")

		list := qn.Name
		listTxt, ok := list.Value.(*trees.Text)
		assert.True(t, ok, "Name.Value: expected *trees.Text, got %T", list.Value)
		assert.Equal(t, "List", listTxt.Value, "Name.Value")
	})

	t.Run("from_tree_package_via_optional", func(t *testing.T) {
		// Build a CompilationUnit tree with package present
		pkgTree := &trees.Node{TypeName: "PackageDeclaration", Tree: &trees.MapNode{
			Entries: map[string]trees.Tree{
				"annotations": &trees.Array{Items: []any{}},
				"name": &trees.Node{TypeName: "QualifiedName", Tree: &trees.MapNode{
					Entries: map[string]trees.Tree{
						"qualifiers": &trees.Seq{Items: []any{
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
					"imports":      &trees.Array{Items: []any{}},
					"declarations": &trees.Seq{Items: []any{}},
					"linecount":    &trees.Nil{},
				},
			}},
		)
		assert.NoError(t, err, "CompilationUnitFromTree")

		pkg, ok := cu.Package.(*PackageDeclaration)
		assert.True(t, ok, "Package: expected *PackageDeclaration, got %T", cu.Package)
		assert.NotZero(t, pkg.Name, "Package.Name is nil")
	})

	t.Run("from_tree_no_package", func(t *testing.T) {
		// Build a CompilationUnit tree WITHOUT package (simulating optional-nil)
		cu, err := compilationUnitFromTree(
			&trees.Node{TypeName: "CompilationUnit", Tree: &trees.MapNode{
				Entries: map[string]trees.Tree{
					"package":      &trees.Nil{},
					"imports":      &trees.Array{Items: []any{}},
					"declarations": &trees.Seq{Items: []any{}},
					"linecount":    &trees.Nil{},
				},
			}},
		)
		assert.NoError(t, err, "CompilationUnitFromTree (no package)")
		assert.Zero(t, cu.Package, "expected nil Package when entry is Nil")
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
			assert.True(t, contains(code, tt.field), "generated code should contain %q", tt.field)
		}

		// Verify the generated code references the expected types via reflect
		assert.True(t, contains(code, "*PackageDeclaration"), "expected *PackageDeclaration in generated code")
		assert.True(t, contains(code, "[]*ImportDeclaration"), "expected []*ImportDeclaration in generated code")
		assert.True(t, contains(code, "IdentifierFromTree"), "expected IdentifierFromTree in generated code")
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

func identifierFromTree(tree any) (*Identifier, error) {
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

func qualifiedNameFromTree(tree any) (*QualifiedName, error) {
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

func packageDeclarationFromTree(tree any) (*PackageDeclaration, error) {
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

func importDeclarationFromTree(tree any) (*ImportDeclaration, error) {
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

func compilationUnitFromTree(tree any) (*CompilationUnit, error) {
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
		list, ok := v.(*trees.Array)
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
