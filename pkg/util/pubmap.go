// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"reflect"

	ctn "github.com/neogeny/ogopego/pkg/util/container"
)

const MaxPointerDerefDepth = 4

type OrderedMap = ctn.BoundedMap[string, any]

func DeRef(ref any) any {
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
	return v.Interface()
}

// PubMapOf returns an OrderedMap containing the public fields of the given reference.
func PubMapOf(ref any) any {
	if ref == nil {
		return nil
	}
	if _, ok := ref.(*OrderedMap); ok {
		return ref
	}

	depth := 0
	v := reflect.ValueOf(ref)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		depth += 1
		if depth > MaxPointerDerefDepth {
			panic(v)
		}
	}
	if v.Kind() != reflect.Struct {
		return ref
	}

	t := v.Type()
	typeName := t.String()
	out := ctn.NewBoundedMap[string, any](0)
	_ = out.Set("__class__", typeName)
	flattenFields(t, v, &out)
	return &out
}

func flattenFields(t reflect.Type, v reflect.Value, out *OrderedMap) {
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		if f.Anonymous {
			if f.Type.Kind() == reflect.Struct && f.Type.Name() != "Node" {
				flattenFields(f.Type, v.Field(i), out)
			}
			continue
		}
		_ = out.Set(f.Name, v.Field(i).Interface())
	}
}
