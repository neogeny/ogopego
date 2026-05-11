package json

import (
	"encoding/json"
	"fmt"
	"reflect"
)

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
