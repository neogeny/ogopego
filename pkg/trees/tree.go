// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import (
	"reflect"

	"github.com/neogeny/ogopego/pkg/util"
)

type Tree interface {
	As_JSON_() any
	tree()
	fold(gather *treeMerge) any
}

type treeMerge struct {
	Root any
	Map  map[string]any
}

func Fold(tree any) any {
	if tree == nil {
		return NIL
	}
	g := &treeMerge{Map: make(map[string]any)}
	return finish(g, fold(g, tree))
}

func fold(g *treeMerge, tree any) any {
	switch val := tree.(type) {
	case Tree:
		return val.fold(g)

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
		return out

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
			if key == "@+" {
				g.Root = appendAsSeq(g.Root, tree)
				return tree
			}
			if key == "@" {
				g.Root = appendTree(g.Root, tree)
				return tree
			}
			if key[0:2] == ":+" {
				g.insertAsList(key[2:], tree)
				return tree
			}
			if key[0:1] == ":" {
				g.insert(key[2:], tree)
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

func finish(g *treeMerge, base any) any {
	switch g.Root.(type) {
	case *Nil, nil:
		{
		}
	default:
		return closed(g.Root)
	}
	if len(g.Map) > 0 {
		for k, v := range g.Map {
			g.Map[k] = closed(v)
		}
		return g.Map
	}
	return closed(base)
}

func closed(t any) any {
	if s, ok := t.(*Seq); ok {
		return s.Items
	}
	return t
}

func merge(a, b any) any {
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

func (m *treeMerge) insert(key string, val any) {
	existing, ok := m.Map[key]
	if !ok {
		m.Map[key] = val
	} else {
		m.Map[key] = appendTree(existing, val)
	}
}

func (m *treeMerge) insertAsList(key string, val any) {
	existing, ok := m.Map[key]
	if !ok {
		m.Map[key] = &Seq{Items: []any{val}}
	} else {
		m.Map[key] = appendAsSeq(existing, val)
	}
}

func isNil(t any) bool {
	if t == nil {
		return true
	}
	_, ok := t.(*Nil)
	return ok || t == nil
}
