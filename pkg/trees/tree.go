// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"fmt"
	"reflect"

	"github.com/neogeny/ogopego/pkg/util"
)

const (
	AtKey        = "@"
	AtListKey    = "+@"
	NamedKey     = ":"
	NamedListKey = "+:"
)

type Tree interface {
	As_JSON_() any
	tree()
}

type FoldGather struct {
	At  any
	Ast map[string]any
}

func Fold(tree any) any {
	if tree == nil {
		return nil
	}
	g := &FoldGather{Ast: make(map[string]any)}
	return finish(g, fold(g, tree))
}

func fold(g *FoldGather, tree any) any {
	if tree == BOTTOM {
		return tree
	}
	switch val := tree.(type) {
	case Tree:
		switch t := val.(type) {
		case *Node:
			return t
		case *Seq:
			var out any = nil
			for _, item := range t.Items {
				out = MergeTrees(out, fold(g, item))
			}
			return out
		case *treeNamed:
			v := fold(g, t.Value)
			g.insert(t.Name, v)
			return v
		case *treeNamedAsList:
			v := fold(g, t.Value)
			g.insertAsList(t.Name, v)
			return v
		case *Override:
			v := fold(g, t.Value)
			g.At = appendTree(g.At, v)
			return v
		case *OverrideAsList:
			v := fold(g, t.Value)
			g.At = appendAsSeq(g.At, v)
			return v
		default:
			panic(fmt.Sprintf("fold: unexpected Tree type %T", t))
		}

	case string, bool, nil,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64:
		return val

	case *util.OrderedMap:
		out := make(map[string]any, val.Len())
		for _, k := range val.Keys() {
			item, _ := val.Get(k)
			out[k] = fold(g, item)
		}
		return fold(g, out)

	case map[string]any:
		out := make(map[string]any, len(val))
		for k, item := range val {
			out[k] = fold(g, item)
		}
		if len(out) != 1 {
			return out
		}
		for key := range out {
			tree := val[key]
			if key == AtListKey {
				g.At = appendAsSeq(g.At, tree)
				return tree
			}
			if key == AtKey {
				g.At = appendTree(g.At, tree)
				return tree
			}
			if len(key) > 2 && key[0:2] == NamedListKey {
				g.insertAsList(key[2:], tree)
				return tree
			}
			if len(key) > 1 && key[0:1] == NamedKey {
				g.insert(key[1:], tree)
				return tree
			}
		}
		return out

	case []any:
		out := make([]any, 0, len(val))
		for _, item := range val {
			out = append(out, fold(g, item))
		}
		return out
	default:
		rv := reflect.ValueOf(tree)
		if !rv.IsValid() {
			return nil
		}

		switch rv.Kind() {
		case reflect.Chan, reflect.Func:
			return nil

		case reflect.Slice, reflect.Array:
			length := rv.Len()
			out := make([]any, 0, length)
			for i := range length {
				out = append(out, fold(g, rv.Index(i).Interface()))
			}
			return out

		default:
			return val
		}
	}
}

func finish(g *FoldGather, base any) any {
	switch g.At.(type) {
	case nil:
		{
		}
	default:
		return closed(g.At)
	}
	if len(g.Ast) > 0 {
		for k, v := range g.Ast {
			g.Ast[k] = closed(v)
		}
		return g.Ast
	}
	return closed(base)
}

func closed(t any) any {
	if s, ok := t.(*Seq); ok {
		return s.Items
	}
	return t
}

func MergeTrees(a, b any) any {
	switch {
	case isNil(a):
		return b
	case isNil(b):
		return a
	default:
		sa, aIsSeq := a.(*Seq)
		sb, bIsSeq := b.(*Seq)
		switch {
		case aIsSeq && bIsSeq:
			items := make([]any, len(sa.Items)+len(sb.Items))
			copy(items, sa.Items)
			copy(items[len(sa.Items):], sb.Items)
			return &Seq{Items: items}
		case aIsSeq:
			items := make([]any, len(sa.Items)+1)
			copy(items, sa.Items)
			items[len(sa.Items)] = b
			return &Seq{Items: items}
		case bIsSeq:
			items := make([]any, 1+len(sb.Items))
			items[0] = a
			copy(items[1:], sb.Items)
			return &Seq{Items: items}
		default:
			return &Seq{Items: []any{a, b}}
		}
	}
}

func appendTree(a, b any) any {
	switch {
	case isNil(a):
		return b
	case isNil(b):
		return a
	default:
		if s, ok := a.(*Seq); ok {
			items := make([]any, len(s.Items)+1)
			copy(items, s.Items)
			items[len(s.Items)] = b
			return &Seq{Items: items}
		}
		return &Seq{Items: []any{a, b}}
	}
}

func appendAsSeq(a, b any) any {
	if isNil(a) {
		return &Seq{Items: []any{b}}
	}
	if s, aIsSeq := a.(*Seq); aIsSeq {
		items := make([]any, len(s.Items)+1)
		copy(items, s.Items)
		items[len(s.Items)] = b
		return &Seq{Items: items}
	}
	return &Seq{Items: []any{a, b}}
}

func (m *FoldGather) insert(key string, val any) {
	existing, ok := m.Ast[key]
	if !ok {
		m.Ast[key] = val
	} else {
		m.Ast[key] = appendTree(existing, val)
	}
}

func (m *FoldGather) insertAsList(key string, val any) {
	existing, ok := m.Ast[key]
	if !ok {
		m.Ast[key] = &Seq{Items: []any{val}}
	} else {
		m.Ast[key] = appendAsSeq(existing, val)
	}
}

func isNil(t any) bool {
	return t == nil
}
