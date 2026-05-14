// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/iancoleman/orderedmap"
)

type OrderedMap = orderedmap.OrderedMap

type AsJSONMixin interface {
	PubMap() *OrderedMap
	AsJSON() any
	AsJSONStr() string
}

type AsJSONBase struct{}

func AsJSON(v any) any {
	seen := make(map[uintptr]bool)
	return asjson(reflect.ValueOf(v), seen)
}

func AsJSONs(v any) string {
	b, err := json.MarshalIndent(AsJSON(v), "", "  ")
	if err != nil {
		return fmt.Sprintf("!json:%v", err)
	}
	return string(b)
}

func ToJSONString(v any) string {
	return toJSONString(v, "", "  ")
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

func (b *AsJSONBase) AsJSONStrOf(ref any) string {
	val := b.AsJSONOf(ref)
	if val == nil {
		return ""
	}
	bts, err := json.MarshalIndent(val, "", "  ")
	if err != nil {
		return fmt.Sprintf("!json:%v", err)
	}
	return string(bts)
}

func (b *AsJSONBase) AsJSONOf(ref any) any {
	pub := b.PubMapOf(ref)
	if pub == nil {
		return nil
	}

	t := reflect.TypeOf(ref)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	typename := t.Name()
	out := orderedmap.New()
	out.Set("__class__", typename)
	for _, k := range pub.Keys() {
		if k == "__class__" {
			continue
		}
		val, _ := pub.Get(k)
		out.Set(k, val)
	}
	return out
}

func (b *AsJSONBase) PubMapOf(ref any) *OrderedMap {
	if ref == nil {
		return nil
	}
	v := reflect.ValueOf(ref)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := orderedmap.New()
	t := v.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		out.Set(f.Name, v.Field(i).Interface())
	}
	return out
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
		if m, ok := v.Interface().(AsJSONMixin); ok {
			return mixinToJSON(m, seen, v.Type().Elem().Name())
		}
		if m, ok := v.Interface().(json.Marshaler); ok {
			return marshalToAny(m)
		}
		return asjson(v.Elem(), seen)

	case reflect.Struct:
		if v.CanAddr() {
			if m, ok := v.Addr().Interface().(AsJSONMixin); ok {
				return mixinToJSON(m, seen, v.Type().Name())
			}
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

func mixinToJSON(m AsJSONMixin, seen map[uintptr]bool, className string) any {
	return m.AsJSON()
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
