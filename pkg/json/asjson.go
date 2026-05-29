// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	util2 "github.com/neogeny/ogopego/pkg/util"
)

type OrderedMap = util2.OrderedMap

func AsJSON(v any) any {
	return toJSONValue(v, make(map[uintptr]bool))
}

// AsJSONStr returns a JSON string representation of the given value.
func AsJSONStr(v any) string {
	bts, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("!json:%v", err)
	}
	return string(bts)
}

func toJSONValue(v any, seen map[uintptr]bool) any {
	id := util2.Id(v)
	if in, ok := seen[id]; ok && in {
		return fmt.Sprintf("%T@%p", v, v)
	}
	seen[id] = true
	defer func() {
		seen[id] = false
	}()

	v = util2.PubMapOf(v)
	switch val := v.(type) {
	case *OrderedMap:
		out := make(map[string]any, val.Len())
		for _, k := range val.Keys() {
			item, _ := val.Get(k)
			out[PythonizeName(k)] = toJSONValue(item, seen)
		}
		return out

	case map[string]any:
		out := make(map[string]any, len(val))
		for k, item := range val {
			out[PythonizeName(k)] = toJSONValue(item, seen)
		}
		return out

	case []any:
		out := make([]any, 0, len(val))
		for _, item := range val {
			out = append(out, toJSONValue(item, seen))
		}
		return out

	case string, bool, nil,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64:
		return val

	default:
		rv := reflect.ValueOf(v)
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
				out = append(out, toJSONValue(rv.Index(i).Interface(), seen))
			}
			return out

		default:
			return val
		}
	}
}

// PythonizeName converts a Go field name to a Python-style snake_case name.
func PythonizeName(s string) string {
	if len(s) == 0 {
		return s
	}
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rune(s[i-1])
				if unicode.IsLower(prev) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) {
					result.WriteByte('_')
				}
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
