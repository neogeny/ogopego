package container

import (
	"encoding/json"
	"iter"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// BoundedMap is an ordered map that optionally acts as an LRU cache.
//
// When created with capacity <= 0, it's an unbounded insertion-ordered map.
// When created with capacity > 0, it's an LRU cache: the least recently used
// entry is evicted when the map reaches capacity.
//
// The map preserves insertion/access order for iteration, JSON marshaling, etc.
// It is safe for use with any comparable key type and any value type.
type BoundedMap[K comparable, V any] struct {
	capacity int
	items    *orderedmap.OrderedMap[K, V]
}

// NewBoundedMap creates a new BoundedMap with the given capacity.
// capacity <= 0 means unbounded.
func NewBoundedMap[K comparable, V any](capacity int) BoundedMap[K, V] {
	return BoundedMap[K, V]{
		capacity: capacity,
		items:    orderedmap.New[K, V](),
	}
}

// Get retrieves the value for the given key.
// With capacity > 0, the entry is promoted to most recently used.
// Returns KeyNotFoundError if the key is not present.
func (bm *BoundedMap[K, V]) Get(key K) (V, error) {
	if bm.capacity > 0 {
		return bm.items.GetAndMoveToBack(key)
	}
	pair := bm.items.GetPair(key)
	if pair == nil {
		var zero V
		return zero, &orderedmap.KeyNotFoundError[K]{MissingKey: key}
	}
	return pair.Value, nil
}

// Set inserts or updates a key-value pair.
// With capacity > 0, the entry is promoted to most recently used,
// and the least recently used entry is evicted if at capacity.
func (bm *BoundedMap[K, V]) Set(key K, value V) error {
	// Update existing entry
	if _, present := bm.items.Get(key); present {
		bm.items.Set(key, value)
		if bm.capacity > 0 {
			return bm.items.MoveToBack(key)
		}
		return nil
	}

	// Evict LRU entry if at capacity
	if bm.capacity > 0 && bm.items.Len() >= bm.capacity {
		if oldest := bm.items.Oldest(); oldest != nil {
			bm.items.Delete(oldest.Key)
		}
	}

	// Insert new entry
	bm.items.Set(key, value)
	return nil
}

// Keys returns the keys in order (oldest to newest).
func (bm *BoundedMap[K, V]) Keys() []K {
	keys := make([]K, 0, bm.items.Len())
	for pair := bm.items.Oldest(); pair != nil; pair = pair.Next() {
		keys = append(keys, pair.Key)
	}
	return keys
}

// Entries returns an iterator over key-value pairs in order (oldest to newest).
func (bm *BoundedMap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for pair := bm.items.Oldest(); pair != nil; pair = pair.Next() {
			if !yield(pair.Key, pair.Value) {
				return
			}
		}
	}
}

// Delete removes the key-value pair for the given key.
func (bm *BoundedMap[K, V]) Delete(key K) {
	bm.items.Delete(key)
}

// Retain removes all key-value pairs that do not satisfy the keep function.
func (bm *BoundedMap[K, V]) Retain(keep func(K, V) bool) {
	var toDelete []K
	for pair := bm.items.Oldest(); pair != nil; pair = pair.Next() {
		if !keep(pair.Key, pair.Value) {
			toDelete = append(toDelete, pair.Key)
		}
	}
	for _, key := range toDelete {
		bm.items.Delete(key)
	}
}

// Len returns the number of entries in the map.
func (bm *BoundedMap[K, V]) Len() int {
	return bm.items.Len()
}

// MarshalJSON implements json.Marshaler.
func (bm *BoundedMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(bm.items)
}
