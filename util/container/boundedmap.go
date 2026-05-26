package container

import (
	"bytes"
	"container/list"
	"encoding/json"
	"iter"
)

type cacheEntry[K comparable, V any] struct {
	key   K
	value V
}

type BoundedMap[K comparable, V any] struct {
	capacity  int
	items     map[K]*list.Element
	evictList *list.List
}

func NewBoundedMap[K comparable, V any](capacity int) BoundedMap[K, V] {
	return BoundedMap[K, V]{
		capacity:  capacity,
		items:     make(map[K]*list.Element, capacity), // Pre-allocated capacity
		evictList: list.New(),
	}
}

func (bm *BoundedMap[K, V]) Get(key K) (V, bool) {
	if elem, exists := bm.items[key]; exists {
		if bm.capacity > 0 {
			bm.evictList.MoveToFront(elem)
		}
		return elem.Value.(*cacheEntry[K, V]).value, true
	}
	var zero V
	return zero, false
}

func (bm *BoundedMap[K, V]) Set(key K, value V) {
	// 1. Update if it already exists
	if elem, exists := bm.items[key]; exists {
		if bm.capacity > 0 {
			bm.evictList.MoveToFront(elem)
		}
		elem.Value.(*cacheEntry[K, V]).value = value
		return
	}

	// 2. Evict if we hit the boundary limit
	if bm.capacity > 0 && bm.evictList.Len() >= bm.capacity {
		oldest := bm.evictList.Back()
		if oldest != nil {
			bm.evictList.Remove(oldest)
			kv := oldest.Value.(*cacheEntry[K, V])
			delete(bm.items, kv.key) // Free map slot
		}
	}

	// 3. Insert new entry (front for LRU, back for ordered)
	entry := &cacheEntry[K, V]{key: key, value: value}
	var elem *list.Element
	if bm.capacity > 0 {
		elem = bm.evictList.PushFront(entry)
	} else {
		elem = bm.evictList.PushBack(entry)
	}
	bm.items[key] = elem
}

func (bm *BoundedMap[K, V]) Keys() []K {
	keys := make([]K, 0, bm.evictList.Len())
	for e := bm.evictList.Front(); e != nil; e = e.Next() {
		if pair, ok := e.Value.(*cacheEntry[K, V]); ok {
			keys = append(keys, pair.key)
		}
	}
	return keys
}

func (bm *BoundedMap[K, V]) Entries() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		// Traverse the doubly-linked list from front to back
		for e := bm.evictList.Front(); e != nil; e = e.Next() {
			// Type-assert the element's Value back to our internal entry struct
			if pair, ok := e.Value.(*cacheEntry[K, V]); ok {
				// yield passes the key/value to the for-range loop.
				// If yield returns false, the loop broke early, so we stop.
				if !yield(pair.key, pair.value) {
					return
				}
			}
		}
	}
}

func (bm *BoundedMap[K, V]) Delete(key K) {
	// 1. Check if the item exists
	elem, exists := bm.items[key]
	if !exists {
		return // No-op if the key isn't in the map
	}

	// 2. Remove the element from the LRU ordering list
	bm.evictList.Remove(elem)

	// 3. Delete the key from the underlying map to reclaim memory slot
	delete(bm.items, key)
}

func (bm *BoundedMap[K, V]) Retain(keep func(K, V) bool) {
	for key, elem := range bm.items {
		ent := elem.Value.(*cacheEntry[K, V])
		if !keep(ent.key, ent.value) {
			bm.evictList.Remove(elem)
			delete(bm.items, key)
		}
	}
}

func (bm *BoundedMap[K, V]) Len() int {
	return bm.evictList.Len()
}

func (bm *BoundedMap[K, V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	for e := bm.evictList.Front(); e != nil; e = e.Next() {
		pair := e.Value.(*cacheEntry[K, V])
		if !first {
			buf.WriteByte(',')
		}
		first = false
		keyJSON, err := json.Marshal(pair.key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyJSON)
		buf.WriteByte(':')
		valJSON, err := json.Marshal(pair.value)
		if err != nil {
			return nil, err
		}
		buf.Write(valJSON)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
