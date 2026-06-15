package test

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/asjson"
	"github.com/neogeny/ogopego/pkg/peg"
	"github.com/neogeny/ogopego/pkg/tool"
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

	// 4. Validate structure via AsJSON
	j := asjson.AsJSON(tree).(map[string]any)
	assert.Equal(t, "CompilationUnit", j["__class__"])

	// 5. Validate types of each entry
	t.Run("package_field", func(t *testing.T) {
		v, ok := j["package"]
		assert.True(t, ok, "missing 'package' entry")
		pn, ok := v.(map[string]any)
		assert.True(t, ok, "package: expected map, got %T", v)
		assert.Equal(t, "PackageDeclaration", pn["__class__"], "package TypeName")
		name, ok := pn["name"]
		assert.True(t, ok, "PackageDeclaration missing 'name' entry")
		qn, ok := name.(map[string]any)
		assert.True(t, ok, "PackageDeclaration.name: expected map, got %T", name)
		assert.Equal(t, "QualifiedName", qn["__class__"], "PackageDeclaration.name TypeName")
	})

	t.Run("imports_field", func(t *testing.T) {
		v, ok := j["imports"]
		assert.True(t, ok, "missing 'imports' entry")
		impList, ok := v.([]any)
		assert.True(t, ok, "imports: expected list, got %T", v)
		assert.Equal(t, 1, len(impList), "imports: expected 1 item")
		impNode, ok := impList[0].(map[string]any)
		assert.True(t, ok, "imports[0]: expected map, got %T", impList[0])
		assert.Equal(t, "ImportDeclaration", impNode["__class__"], "imports[0] TypeName")
	})

	t.Run("declarations_field", func(t *testing.T) {
		v, ok := j["declarations"]
		assert.True(t, ok, "missing 'declarations' entry")
		items, ok := v.([]any)
		assert.True(t, ok, "declarations: expected list, got %T", v)
		assert.Equal(t, 1, len(items), "declarations: expected 1 item")
		classNode, ok := items[0].(map[string]any)
		assert.True(t, ok, "declarations[0]: expected map, got %T", items[0])
		assert.Equal(t, "ClassDeclaration", classNode["__class__"], "declarations[0] TypeName")
	})

	t.Run("linecount_field", func(t *testing.T) {
		v, ok := j["linecount"]
		assert.True(t, ok, "missing 'linecount' entry")
		assert.True(t, v == nil, "linecount: expected nil, got %T", v)
	})

	// 6. Test FromTree with hand-authored model types
	t.Run("from_tree_identifier", func(t *testing.T) {
		id, err := identifierFromTree(map[string]any{
			"__class__": "Identifier",
			"value":     "Hello",
		})
		assert.NoError(t, err, "IdentifierFromTree")
		assert.Equal(t, "Hello", id.Value, "Identifier.Value")
	})

	t.Run("from_tree_qualified_name", func(t *testing.T) {
		qn, err := qualifiedNameFromTree(map[string]any{
			"__class__": "QualifiedName",
			"qualifiers": []any{
				map[string]any{"__class__": "Identifier", "value": "com"},
				map[string]any{"__class__": "Identifier", "value": "example"},
			},
			"name": map[string]any{"__class__": "Identifier", "value": "List"},
		})
		assert.NoError(t, err, "QualifiedNameFromTree")
		assert.Equal(t, 2, len(qn.Qualifiers))
		assert.Equal(t, "com", qn.Qualifiers[0].Value, "Qualifiers[0].Value")
		assert.Equal(t, "List", qn.Name.Value, "Name.Value")
	})

	t.Run("from_tree_package_via_optional", func(t *testing.T) {
		pkgTree := map[string]any{
			"__class__":   "PackageDeclaration",
			"annotations": []any{},
			"name": map[string]any{
				"__class__": "QualifiedName",
				"qualifiers": []any{
					map[string]any{"__class__": "Identifier", "value": "com"},
				},
				"name": map[string]any{"__class__": "Identifier", "value": "example"},
			},
		}

		cu, err := compilationUnitFromTree(map[string]any{
			"__class__":    "CompilationUnit",
			"package":      pkgTree,
			"imports":      []any{},
			"declarations": []any{},
			"linecount":    nil,
		})
		assert.NoError(t, err, "CompilationUnitFromTree")
		pkg, ok := cu.Package.(*PackageDeclaration)
		assert.True(t, ok, "Package: expected *PackageDeclaration, got %T", cu.Package)
		assert.NotZero(t, pkg.Name, "Package.Name is nil")
	})

	t.Run("from_tree_no_package", func(t *testing.T) {
		cu, err := compilationUnitFromTree(map[string]any{
			"__class__":    "CompilationUnit",
			"package":      nil,
			"imports":      []any{},
			"declarations": []any{},
			"linecount":    nil,
		})
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
	m, ok := tree.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("IdentifierFromTree: expected map, got %T", tree)
	}
	var result Identifier
	if v, ok := m["value"]; ok {
		result.Value = v
	}
	return &result, nil
}

type QualifiedName struct {
	Qualifiers []*Identifier
	Name       *Identifier
}

func qualifiedNameFromTree(tree any) (*QualifiedName, error) {
	m, ok := tree.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("QualifiedNameFromTree: expected map, got %T", tree)
	}
	var result QualifiedName
	var err error
	if v, ok := m["qualifiers"]; ok {
		items, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("QualifiedName.Qualifiers: expected list, got %T", v)
		}
		result.Qualifiers = make([]*Identifier, len(items))
		for i, item := range items {
			result.Qualifiers[i], err = identifierFromTree(item)
			if err != nil {
				return nil, fmt.Errorf("QualifiedName.Qualifiers[%d]: %w", i, err)
			}
		}
	}
	if v, ok := m["name"]; ok {
		result.Name, err = identifierFromTree(v)
		if err != nil {
			return nil, fmt.Errorf("QualifiedName.Name: %w", err)
		}
	}
	return &result, nil
}

type PackageDeclaration struct {
	Annotations any
	Name        *QualifiedName
}

func packageDeclarationFromTree(tree any) (*PackageDeclaration, error) {
	m, ok := tree.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("PackageDeclarationFromTree: expected map, got %T", tree)
	}
	var result PackageDeclaration
	var err error
	if v, ok := m["annotations"]; ok {
		result.Annotations = v
	}
	if v, ok := m["name"]; ok {
		result.Name, err = qualifiedNameFromTree(v)
		if err != nil {
			return nil, fmt.Errorf("PackageDeclaration.Name: %w", err)
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
	m, ok := tree.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("ImportDeclarationFromTree: expected map, got %T", tree)
	}
	var result ImportDeclaration
	var err error
	if v, ok := m["static"]; ok {
		result.Static = v
	}
	if v, ok := m["name"]; ok {
		result.Name, err = qualifiedNameFromTree(v)
		if err != nil {
			return nil, fmt.Errorf("ImportDeclaration.Name: %w", err)
		}
	}
	if v, ok := m["all"]; ok {
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
	m, ok := tree.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("CompilationUnitFromTree: expected map, got %T", tree)
	}
	var result CompilationUnit
	var err error
	if v, ok := m["package"]; ok {
		if pkgMap, ok := v.(map[string]any); ok {
			var pkg *PackageDeclaration
			pkg, err = packageDeclarationFromTree(pkgMap)
			if err != nil {
				return nil, fmt.Errorf("CompilationUnit.Package: %w", err)
			}
			result.Package = pkg
		}
	}
	if v, ok := m["imports"]; ok {
		items, ok := v.([]any)
		if !ok {
			return nil, fmt.Errorf("CompilationUnit.Imports: expected list, got %T", v)
		}
		imports := make([]*ImportDeclaration, len(items))
		for i, item := range items {
			imports[i], err = importDeclarationFromTree(item)
			if err != nil {
				return nil, fmt.Errorf("CompilationUnit.Imports[%d]: %w", i, err)
			}
		}
		result.Imports = imports
	}
	if v, ok := m["declarations"]; ok {
		result.Declarations = v
	}
	if v, ok := m["linecount"]; ok {
		result.Linecount = v
	}
	return &result, nil
}
