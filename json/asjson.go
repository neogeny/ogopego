// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"
	"unsafe"

	"github.com/iancoleman/orderedmap"
	"github.com/neogeny/ogopego/util"
)

type OrderedMap = orderedmap.OrderedMap

// Id returns a unique uintptr identifier for any Go value,
// mimicking Python's id() without panicking on unhashable types.
func Id(val any) uintptr {
	if val == nil {
		return 0
	}

	v := reflect.ValueOf(val)

	switch v.Kind() {
	// Reference types natively hold a pointer to their data block.
	// v.Pointer() extracts it safely without panicking.
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Chan, reflect.UnsafePointer:
		return v.Pointer()

	default:
		// For structs, arrays, and primitive scalars, we take the address
		// of the underlying concrete value held by the interface.
		if v.CanAddr() {
			return v.UnsafeAddr()
		}

		// Fallback for unaddressable values: extract the address of the
		// interface wrapper's data segment directly via unsafe.
		return uintptr(unsafe.Pointer(&val))
	}
}

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

// ToJSONString converts a Go value to a JSON string with optional prefix and indent.
func ToJSONString(v any) string {
	return toJSONString(v, "", "  ")
}

func toJSONValue(v any, seen map[uintptr]bool) any {
	id := Id(v)
	if in, ok := seen[id]; ok && in {
		return fmt.Sprintf("%T@%p", v, v)
	}
	seen[id] = true
	defer func() {
		seen[id] = false
	}()

	v = util.PubMapOf(v)
	switch val := v.(type) {
	// --- Fast-Path: Complex Types ---
	case *OrderedMap:
		out := make(map[string]any, len(val.Keys()))
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

	// --- Fast-Path: Safe Primitives (No reflection needed) ---
	case string, bool, nil,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64:
		return val

	// --- Fallback: Go-Specific and Dynamic Sequences ---
	default:
		rv := reflect.ValueOf(v)
		if !rv.IsValid() {
			return nil
		}

		switch rv.Kind() {
		case reflect.Chan, reflect.Func:
			// Safely intervene on channels and functions before JSON marshalling fails
			//return fmt.Sprintf("%T(%p)", v, v)
			return nil

		case reflect.Slice, reflect.Array:
			// Unpack typed slices/arrays (e.g., []func(), []chan string, []MyStruct)
			length := rv.Len()
			out := make([]any, 0, length)
			for i := range length {
				out = append(out, toJSONValue(rv.Index(i).Interface(), seen))
			}
			return out

		default:
			// Catch-all for underlying structs that escaped PubMapOf or custom types
			return val
		}
	}
}

func toJSONString(v any, prefix, indent string) string {
	switch val := v.(type) {
	case *OrderedMap:
		if len(val.Keys()) == 0 {
			return "{}"
		}
		inner := prefix + indent
		var buf strings.Builder
		buf.WriteString("{\n")
		for i, k := range val.Keys() {
			item, _ := val.Get(k)
			buf.WriteString(inner)
			buf.WriteString(fmt.Sprintf("%q: ", k))
			buf.WriteString(toJSONString(item, inner, indent))
			if i < len(val.Keys())-1 {
				buf.WriteString(",")
			}
			buf.WriteString("\n")
		}
		buf.WriteString(prefix + "}")
		return buf.String()
	case map[string]any:
		if len(val) == 0 {
			return "{}"
		}
		inner := prefix + indent
		var buf strings.Builder
		buf.WriteString("{\n")
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			buf.WriteString(inner)
			buf.WriteString(fmt.Sprintf("%q: ", k))
			buf.WriteString(toJSONString(val[k], inner, indent))
			if i < len(keys)-1 {
				buf.WriteString(",")
			}
			buf.WriteString("\n")
		}
		buf.WriteString(prefix + "}")
		return buf.String()
	case []any:
		if len(val) == 0 {
			return "[]"
		}
		inner := prefix + indent
		var buf strings.Builder
		buf.WriteString("[\n")
		for i, item := range val {
			buf.WriteString(inner)
			buf.WriteString(toJSONString(item, inner, indent))
			if i < len(val)-1 {
				buf.WriteString(",")
			}
			buf.WriteString("\n")
		}
		buf.WriteString(prefix + "]")
		return buf.String()
	case string:
		return fmt.Sprintf("%q", val)
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", val)
	}
}

// ToGoMap converts an OrderedMap or slice of any to a standard Go map or slice recursively.
func ToGoMap(v any) any {
	switch val := v.(type) {
	case *orderedmap.OrderedMap:
		m := make(map[string]any)
		for _, k := range val.Keys() {
			item, _ := val.Get(k)
			m[k] = ToGoMap(item)
		}
		return m
	case []any:
		s := make([]any, len(val))
		for i, item := range val {
			s[i] = ToGoMap(item)
		}
		return s
	default:
		return v
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

func asjson(val reflect.Value, seen map[uintptr]bool) any {
	if !val.IsValid() {
		return nil
	}

	v := val
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		if m, ok := v.Interface().(json.Marshaler); ok {
			return marshalToAny(m)
		}
		return asjson(v.Elem(), seen)

	case reflect.Struct:
		if v.CanAddr() {
			if m, ok := v.Addr().Interface().(json.Marshaler); ok {
				return marshalToAny(m)
			}
		}
		return structToJSON(v, seen)

	case reflect.Map:
		addr := v.Pointer()
		if seen[addr] {
			return fmt.Sprintf("%s@0x%X", v.Type().Name(), addr)
		}
		seen[addr] = true
		defer delete(seen, addr)

		out := make(map[string]any, v.Len())
		for _, key := range v.MapKeys() {
			k := fmt.Sprint(key.Interface())
			out[k] = asjson(v.MapIndex(key), seen)
		}
		return out

	case reflect.Slice, reflect.Array:
		n := v.Len()
		out := make([]any, 0, n)
		for i := range n {
			out = append(out, asjson(v.Index(i), seen))
		}
		return out

	case reflect.String:
		return v.String()

	case reflect.Bool:
		return v.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()

	case reflect.Float32, reflect.Float64:
		return v.Float()

	case reflect.Func, reflect.Chan:
		return nil

	default:
		return fmt.Sprint(v.Interface())
	}
}

func marshalToAny(m json.Marshaler) any {
	raw, err := m.MarshalJSON()
	if err != nil {
		return nil
	}
	var out any
	_ = json.Unmarshal(raw, &out)
	return out
}

func structToJSON(val reflect.Value, seen map[uintptr]bool) any {
	t := val.Type()

	if val.CanAddr() {
		addr := val.Addr().Pointer()
		if seen[addr] {
			return fmt.Sprintf("%s@0x%X", t.Name(), addr)
		}
		seen[addr] = true
		defer delete(seen, addr)
	}

	out := make(map[string]any)
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		out[f.Name] = asjson(val.Field(i), seen)
	}
	return out
}
