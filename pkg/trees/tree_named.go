// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package trees

import "unicode"

func validateUserKeyName(name string) {
	for _, char := range name {
		if !unicode.IsLetter(char) && char != '_' {
			panic("invalid name: " + name)
		}
	}
}

// TreeNamed represents a named key/value pair folded from the parse tree.
func TreeNamed(name string, value any) any {
	validateUserKeyName(name)
	return map[string]any{keyNamed + name: value}
}

// TreeNamedAsList is like Named but its values are collected as a list.
func TreeNamedAsList(name string, value any) any {
	validateUserKeyName(name)
	return map[string]any{keyListNamed + name: value}
}

// TreeOverride indicates that the contained value should override other values
// when folding into the result.
func TreeOverride(value any) any {
	return map[string]any{keyAt: value}
}

// TreeOverrideAsList is a list-form override variant.
func TreeOverrideAsList(value any) any {
	return map[string]any{keyListAt: value}
}
