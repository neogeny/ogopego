// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"reflect"
	"unsafe"
)

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
