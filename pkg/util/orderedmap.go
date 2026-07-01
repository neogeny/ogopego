// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"iter"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// OrderedMap is an insertion-ordered map, backed by github.com/wk8/go-ordered-map/v2.
type OrderedMap = orderedmap.OrderedMap[string, any]

// OrderedMapEntries returns an iterator over key-value pairs of an ordered map in insertion order.
func OrderedMapEntries[K comparable, V any](om *orderedmap.OrderedMap[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for pair := om.Oldest(); pair != nil; pair = pair.Next() {
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	}
}

// OrderedMapKeys returns the keys of an ordered map in insertion order.
func OrderedMapKeys[K comparable, V any](om *orderedmap.OrderedMap[K, V]) []K {
	keys := make([]K, 0, om.Len())
	for pair := om.Oldest(); pair != nil; pair = pair.Next() {
		keys = append(keys, pair.Key)
	}
	return keys
}
