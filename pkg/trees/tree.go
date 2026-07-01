// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"fmt"
	"reflect"

	"github.com/neogeny/ogopego/pkg/util"
)

const (
	keyAt        = "@"
	keyNamed     = ":"
	keyListAt    = "@+"
	keyListNamed = "+:"
)

type Tree interface {
	As_JSON_(seen map[uintptr]bool) any
	isTree()
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
		case *treeSeq:
			var out any = nil
			for _, item := range t.Items {
				out = MergeTrees(out, fold(ast, item))
			}
			return out
		default:
			var f any
			switch n := t.(type) {
			case *namedTreeT:
				f := fold(ast, n.Tree)
				mapTree(ast, n.Name(), f)
			case *namedTreeSeqT:
				f := fold(ast, n.Tree)
				mapTreeAsSeq(ast, n.Name(), f)
			case *ovrTreeT:
				f := fold(ast, n.Tree)
				mapTree(ast, keyAt, f)
			case *overrideTreeSeqT:
				f := fold(ast, n.Tree)
				mapTreeAsSeq(ast, keyAt, f)
			default:
				panic(fmt.Sprintf("fold: unexpected Tree type %T", t))
			}
			return f
		}

	case string, bool, nil,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64:
		return val

	case *util.OrderedMap:
		out := make(map[string]any, val.Len())
		for _, k := range util.OrderedMapKeys(val) {
			if item, ok := val.Get(k); ok {
				out[k] = fold(ast, item)
			}
		}
		return fold(ast, out)

	case map[string]any:
		out := make(map[string]any, len(val))
		for k, item := range val {
			out[k] = fold(ast, item)
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
		if _, isAtSet := ast[keyAt]; isAtSet {
			return ast[keyAt]
		}
		return ast
	}
	return closed(base)
}

func closed(t any) any {
	if s, ok := t.(*treeSeq); ok {
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
		sa, aIsSeq := a.(*treeSeq)
		sb, bIsSeq := b.(*treeSeq)
		switch {
		case aIsSeq && bIsSeq:
			items := make([]any, len(sa.Items)+len(sb.Items))
			copy(items, sa.Items)
			copy(items[len(sa.Items):], sb.Items)
			return &treeSeq{Items: items}
		case aIsSeq:
			items := make([]any, len(sa.Items)+1)
			copy(items, sa.Items)
			items[len(sa.Items)] = b
			return &treeSeq{Items: items}
		case bIsSeq:
			items := make([]any, 1+len(sb.Items))
			items[0] = a
			copy(items[1:], sb.Items)
			return &treeSeq{Items: items}
		default:
			return &treeSeq{Items: []any{a, b}}
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
		if s, ok := a.(*treeSeq); ok {
			items := make([]any, len(s.Items)+1)
			copy(items, s.Items)
			items[len(s.Items)] = b
			return &treeSeq{Items: items}
		}
		return &treeSeq{Items: []any{a, b}}
	}
}

func appendAsSeq(a, b any) any {
	if isNil(a) {
		return &treeSeq{Items: []any{b}}
	}
	if s, aIsSeq := a.(*treeSeq); aIsSeq {
		items := make([]any, len(s.Items)+1)
		copy(items, s.Items)
		items[len(s.Items)] = b
		return &treeSeq{Items: items}
	}
	return &treeSeq{Items: []any{a, b}}
}

func mapTree(m map[string]any, key string, val any) {
	existing, ok := m[key]
	if !ok {
		m[key] = val
	} else {
		m[key] = appendTree(existing, val)
	}
}

func mapTreeAsSeq(m map[string]any, key string, val any) {
	existing, ok := m[key]
	if !ok {
		m[key] = &treeSeq{Items: []any{val}}
	} else {
		m[key] = appendAsSeq(existing, val)
	}
}

func isNil(t any) bool {
	return t == nil
}
