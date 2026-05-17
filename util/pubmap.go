// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package util

import (
	"reflect"

	"github.com/iancoleman/orderedmap"
)

type OrderedMap = orderedmap.OrderedMap

type PubMapper interface {
	PubMap() *OrderedMap
}

func NewOrderedMap() *OrderedMap {
	return orderedmap.New()
}

// PubMapOf returns an OrderedMap containing the public fields of the given reference.
func PubMapOf(ref any) *OrderedMap {
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
		if !f.IsExported() || f.Anonymous {
			continue
		}
		out.Set(f.Name, v.Field(i).Interface())
	}
	return out
}
