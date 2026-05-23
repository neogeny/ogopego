// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const BlackLineLength = 88

// FoldOption configures line-wrapping behavior in Fold.
type FoldOption struct {
	AddLevels int
	Amount    int
}

// Repr returns a Go-composite-literal representation of v.
// Structs are rendered via PubMapOf, giving each one its type name
// and exported fields. Slices, maps, and scalars follow Go literal syntax.
func Repr(v any) string {
	return reprValue(v, make(map[uintptr]bool))
}

// Fold joins parts with commas, wraps in brackets, and handles line-breaking.
// Pass a single-element parts for field-level wrapping ("prefix + body").
func Fold(prefix string, parts []string, lbrack, rbrack string, opts ...FoldOption) string {
	opt := FoldOption{Amount: 2}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if len(parts) == 0 {
		return prefix + lbrack + rbrack
	}

	single := prefix + lbrack + strings.Join(parts, ", ") + rbrack
	if fitsfmt(single, opt.AddLevels, opt.Amount) {
		return single
	}

	indent := strings.Repeat(" ", opt.Amount)
	for i, p := range parts {
		parts[i] = indent + strings.ReplaceAll(p, "\n", "\n"+indent)
	}
	return prefix + lbrack + "\n" + strings.Join(parts, ",\n") + ",\n" + prefix + rbrack
}

func fitsfmt(line string, addLevels, amount int) bool {
	if strings.Contains(line, "\n") {
		return false
	}
	return len(line)+addLevels*amount <= BlackLineLength
}

// classKeys extracts __class__ from a map and returns the type name
// plus remaining keys sorted alphabetically.
func classKeys(m map[string]any) (string, []string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	typeName := ""
	if len(keys) > 0 && keys[0] == "__class__" {
		typeName = fmt.Sprint(m["__class__"])
		keys = keys[1:]
	}
	return typeName, keys
}

func reprFold(parts []string, typeName string) string {
	if typeName == "" {
		return Fold("", parts, "map[string]any{", "}")
	}
	return Fold("", parts, typeName+"{", "}")
}

func reprValue(v any, seen map[uintptr]bool) string {
	id := Id(v)
	if id != 0 {
		if seen[id] {
			return "nil"
		}
		seen[id] = true
		defer func() { seen[id] = false }()
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() && rv.Elem().Kind() == reflect.Struct {
		if _, ok := v.(*OrderedMap); !ok {
			return "&" + reprValue(rv.Elem().Interface(), seen)
		}
	}

	v = PubMapOf(v)
	switch val := v.(type) {
	case nil:
		return "nil"
	case string:
		return fmt.Sprintf("%q", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64:
		return fmt.Sprint(val)
	case []any:
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = reprValue(item, seen)
		}
		return Fold("", parts, "[]any{", "}")
	case map[string]any:
		typeName, keys := classKeys(val)
		parts := make([]string, len(keys))
		for i, k := range keys {
			parts[i] = fmt.Sprintf("%s: %s", k, reprValue(val[k], seen))
		}
		return reprFold(parts, typeName)
	case *OrderedMap:
		return reprOrderedMap(val, seen)
	default:
		rv := reflect.ValueOf(val)
		if !rv.IsValid() {
			return "nil"
		}
		switch rv.Kind() {
		case reflect.Ptr:
			if rv.IsNil() {
				return "nil"
			}
			return "&" + reprValue(rv.Elem().Interface(), seen)
		case reflect.Slice, reflect.Array:
			n := rv.Len()
			parts := make([]string, n)
			for i := range n {
				parts[i] = reprValue(rv.Index(i).Interface(), seen)
			}
			return Fold("", parts, rv.Type().String()+"{", "}")
		default:
			return fmt.Sprintf("%v", val)
		}
	}
}

func reprOrderedMap(om *OrderedMap, seen map[uintptr]bool) string {
	keys := om.Keys()
	typeName := ""
	if len(keys) > 0 && keys[0] == "__class__" {
		if cls, ok := om.Get("__class__"); ok {
			typeName = fmt.Sprint(cls)
		}
		keys = keys[1:]
	}
	if len(keys) == 0 {
		if typeName == "" {
			return "map[string]any{}"
		}
		return typeName + "{}"
	}
	parts := make([]string, len(keys))
	for i, k := range keys {
		item, _ := om.Get(k)
		parts[i] = fmt.Sprintf("%s: %s", k, reprValue(item, seen))
	}
	if typeName == "" {
		return Fold("", parts, "map[string]any{", "}")
	}
	return reprFold(parts, typeName)
}
