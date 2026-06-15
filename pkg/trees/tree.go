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
	m := make(map[string]any)
	return finish(m, fold(m, tree))
}

func fold(ast map[string]any, tree any) any {
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
				out = MergeTrees(out, fold(ast, item))
			}
			return out
		case *treeNamed:
			v := fold(ast, t.Value)
			insert(ast, t.Name, v)
			return v
		case *treeNamedAsList:
			v := fold(ast, t.Value)
			insertAsSeq(ast, t.Name, v)
			return v
		case *treeOverride:
			v := fold(ast, t.Value)
			insert(ast, AtKey, v)
			return v
		case *treeOverrideAsList:
			v := fold(ast, t.Value)
			insertAsSeq(ast, AtKey, v)
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
			out[k] = fold(ast, item)
		}
		return fold(ast, out)

	case map[string]any:
		out := make(map[string]any, len(val))
		for k, item := range val {
			out[k] = fold(ast, item)
		}
		if len(out) != 1 {
			return out
		}
		for key := range out {
			tree := val[key]
			if key == AtListKey {
				insertAsSeq(ast, AtKey, tree)
				return tree
			}
			if key == AtKey {
				insert(ast, AtKey, tree)
				return tree
			}
			if len(key) > 2 && key[0:2] == NamedListKey {
				insertAsSeq(ast, key[2:], tree)
				return tree
			}
			if len(key) > 1 && key[0:1] == NamedKey {
				insert(ast, key[1:], tree)
				return tree
			}
		}
		return out

	case []any:
		out := make([]any, 0, len(val))
		for _, item := range val {
			out = append(out, fold(ast, item))
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
				out = append(out, fold(ast, rv.Index(i).Interface()))
			}
			return out

		default:
			return val
		}
	}
}

func finish(ast map[string]any, base any) any {
	if len(ast) > 0 {
		for k, v := range ast {
			ast[k] = closed(v)
		}
		if _, isAtSet := ast[AtKey]; isAtSet {
			return ast[AtKey]
		}
		return ast
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

func insert(m map[string]any, key string, val any) {
	existing, ok := m[key]
	if !ok {
		m[key] = val
	} else {
		m[key] = appendTree(existing, val)
	}
}

func insertAsSeq(m map[string]any, key string, val any) {
	existing, ok := m[key]
	if !ok {
		m[key] = &Seq{Items: []any{val}}
	} else {
		m[key] = appendAsSeq(existing, val)
	}
}

func isNil(t any) bool {
	return t == nil
}
