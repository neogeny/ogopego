// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package util

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const blackLineLength = 88

// FoldOption configures line-wrapping behavior in Fold.
type FoldOption struct {
	AddLevels int
	Amount    int
}

// Fold joins parts with commas, wraps in brackets, and handles line-breaking.
// Pass a single-element parts for field-level wrapping ("prefix + body").
func Fold(prefix string, parts []string, lbrack, rbrack string, opts ...FoldOption) string {
	opt := FoldOption{Amount: 2}
	if len(opts) > 0 {
		opt = opts[0]
	}

	single := prefix + lbrack + strings.Join(parts, ", ") + rbrack
	if fitsfmt(single, opt.AddLevels, opt.Amount) {
		return single
	}

	indent := strings.Repeat(" ", opt.Amount)
	return prefix + lbrack + "\n" + indent + strings.Join(parts, ",\n"+indent) + "\n" + rbrack
}

func fitsfmt(line string, addLevels, amount int) bool {
	if strings.Contains(line, "\n") {
		return false
	}
	return len(line)+addLevels*amount <= blackLineLength
}

// Repr returns a Go-composite-literal representation of v by consuming
// the PubMap protocol on AsJSONMixin types and delegating containers
// and scalars to the appropriate Go literal syntax.
func Repr(v any) string {
	if v == nil {
		return "nil"
	}

	rv := reflect.ValueOf(v)

	// Dereference pointers not handled directly by the type switch
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		if _, ok := v.(PubMapper); !ok {
			if _, ok := v.(*OrderedMap); !ok {
				return Repr(rv.Elem().Interface())
			}
		}
	}

	switch val := v.(type) {
	case *OrderedMap:
		var parts []string
		for _, k := range val.Keys() {
			item, _ := val.Get(k)
			parts = append(parts, fmt.Sprintf("%q: %s", k, Repr(item)))
		}
		prefix := reflect.TypeOf(val).String()
		return Fold(prefix, parts, "{", "}")

	case map[string]any:
		if val == nil {
			return "nil"
		}
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var parts []string
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%q: %s", k, Repr(val[k])))
		}
		return Fold("map[string]any", parts, "{", "}")

	case []any:
		if val == nil {
			return "nil"
		}
		parts := make([]string, len(val))
		for i, e := range val {
			parts[i] = Repr(e)
		}
		return Fold("[]any", parts, "{", "}")

	default:
		// AsJSONMixin check
		if m, ok := v.(PubMapper); ok {
			return reprMixin(m)
		}

		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			n := rv.Len()
			typ := rv.Type().String()
			if n == 0 {
				return typ + "{}"
			}
			parts := make([]string, n)
			for i := range n {
				parts[i] = Repr(rv.Index(i).Interface())
			}
			return Fold(typ, parts, "{", "}")

		case reflect.Map:
			typ := rv.Type().String()
			if rv.IsNil() {
				return typ + "{}"
			}
			keys := rv.MapKeys()
			sort.Slice(keys, func(i, j int) bool {
				return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
			})
			parts := make([]string, len(keys))
			for i, k := range keys {
				parts[i] = Repr(k.Interface()) + ": " + Repr(rv.MapIndex(k).Interface())
			}
			return Fold(typ, parts, "{", "}")

		case reflect.Struct:
			if m, ok := v.(PubMapper); ok {
				return reprMixin(m)
			}
		}

		return fmt.Sprintf("%#v", v)
	}
}

func reprMixin(m PubMapper) string {
	pub := m.PubMap()
	if pub == nil {
		return fmt.Sprintf("%#v", m)
	}

	// Build type name with optional & prefix for pointer receivers
	t := reflect.TypeOf(m)
	isPtr := false
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		isPtr = true
	}
	typeName := t.Name()
	if isPtr {
		typeName = "&" + typeName
	}

	// Filter nil values, collect formatted parts
	var parts []string
	for _, k := range pub.Keys() {
		val, ok := pub.Get(k)
		if !ok || isReprNil(val) {
			continue
		}
		parts = append(parts, k+": "+Repr(val))
	}

	if len(parts) == 0 {
		return typeName + "{}"
	}
	return Fold(typeName, parts, "{", "}")
}

// isReprNil returns true for values that should be treated as "nothing"
// in the representation, matching Python's None-skip semantics.
// Nil slices and maps are NOT skipped; they render as empty containers.
// Weak references are not representable here because PubMapOf excludes
// unexported fields.
func isReprNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}
