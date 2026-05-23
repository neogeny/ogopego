// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"encoding/json"
	"fmt"
	"maps"
)

type Tree interface {
	tree()
	fold(gather *treeMerge) Tree
}

type TreeBase struct{}

// TreeToJSONStr returns the JSON string representation of a tree.
func TreeToJSONStr(t Tree) string {
	v := TreeToJSON(t)
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("!json:%v", err)
	}
	return string(b)
}

func TreeToJSON(t Tree) any {
	switch v := t.(type) {
	case *Text:
		return v.Value
	case *Number:
		return v.Value
	case *Bool:
		return v.Value
	case *Nil:
		return nil
	case *Bottom:
		return nil
	case *TrueValue:
		return true
	case *FalseValue:
		return false
	case *NullValue:
		return nil
	case *Seq:
		items := make([]any, len(v.Items))
		for i, item := range v.Items {
			items[i] = TreeToJSON(item)
		}
		return items
	case *List:
		items := make([]any, len(v.Items))
		for i, item := range v.Items {
			items[i] = TreeToJSON(item)
		}
		return items
	case *MapNode:
		out := make(map[string]any, len(v.Entries))
		for k, val := range v.Entries {
			out[k] = TreeToJSON(val)
		}
		return out
	case *Named:
		return map[string]any{v.Name: TreeToJSON(v.Value)}
	case *NamedAsList:
		return map[string]any{v.Name: TreeToJSON(v.Value)}
	case *Override:
		return TreeToJSON(v.Value)
	case *OverrideAsList:
		return TreeToJSON(v.Value)
	case *Node:
		child := TreeToJSON(v.Tree)
		if m, ok := child.(map[string]any); ok {
			if _, has := m["__class__"]; !has {
				out := make(map[string]any, len(m)+1)
				out["__class__"] = v.TypeName
				maps.Copy(out, m)
				return out
			}
		}
		return map[string]any{"__class__": v.TypeName, "ast": child}
	default:
		return nil
	}
}

type treeMerge struct {
	Root Tree
	Map  map[string]Tree
}

func Fold(tree Tree) Tree {
	if tree == nil {
		return &Nil{}
	}
	g := &treeMerge{Map: make(map[string]Tree)}
	result := tree.fold(g)
	return finish(g, result)
}

func finish(g *treeMerge, base Tree) Tree {
	switch g.Root.(type) {
	case *Nil, nil:
	default:
		return closed(g.Root)
	}
	if len(g.Map) > 0 {
		return &MapNode{Entries: g.Map}
	}
	return closed(base)
}

func closed(t Tree) Tree {
	if _, ok := t.(*Seq); ok {
		return &List{Items: t.(*Seq).Items}
	}
	return t
}

func merge(a, b Tree) Tree {
	switch {
	case isNil(a):
		return b
	case isNil(b):
		return a
	default:
		sa, aOk := a.(*Seq)
		sb, bOk := b.(*Seq)
		switch {
		case aOk && bOk:
			items := make([]Tree, len(sa.Items)+len(sb.Items))
			copy(items, sa.Items)
			copy(items[len(sa.Items):], sb.Items)
			return &Seq{Items: items}
		case aOk:
			items := make([]Tree, len(sa.Items)+1)
			copy(items, sa.Items)
			items[len(sa.Items)] = b
			return &Seq{Items: items}
		case bOk:
			items := make([]Tree, 1+len(sb.Items))
			items[0] = a
			copy(items[1:], sb.Items)
			return &Seq{Items: items}
		default:
			return &Seq{Items: []Tree{a, b}}
		}
	}
}

func appendTree(a, b Tree) Tree {
	switch {
	case isNil(a):
		return b
	case isNil(b):
		return a
	default:
		if s, ok := a.(*Seq); ok {
			items := make([]Tree, len(s.Items)+1)
			copy(items, s.Items)
			items[len(s.Items)] = b
			return &Seq{Items: items}
		}
		return &Seq{Items: []Tree{a, b}}
	}
}

func appendAsList(a, b Tree) Tree {
	if isNil(a) {
		return &Seq{Items: []Tree{b}}
	}
	if s, ok := a.(*Seq); ok {
		items := make([]Tree, len(s.Items)+1)
		copy(items, s.Items)
		items[len(s.Items)] = b
		return &Seq{Items: items}
	}
	return &Seq{Items: []Tree{a, b}}
}

func (m *treeMerge) insert(key string, val Tree) {
	existing, ok := m.Map[key]
	if !ok {
		m.Map[key] = val
	} else {
		m.Map[key] = appendTree(existing, val)
	}
}

func (m *treeMerge) insertAsList(key string, val Tree) {
	existing, ok := m.Map[key]
	if !ok {
		m.Map[key] = &Seq{Items: []Tree{val}}
	} else {
		m.Map[key] = appendAsList(existing, val)
	}
}

func isNil(t Tree) bool {
	_, ok := t.(*Nil)
	return ok || t == nil
}
