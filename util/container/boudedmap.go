package container

import "container/list"

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
		bm.evictList.MoveToFront(elem)
		return elem.Value.(*cacheEntry[K, V]).value, true
	}
	var zero V
	return zero, false
}

func (bm *BoundedMap[K, V]) Set(key K, value V) {
	// 1. Update if it already exists
	if elem, exists := bm.items[key]; exists {
		bm.evictList.MoveToFront(elem)
		elem.Value.(*cacheEntry[K, V]).value = value
		return
	}

	// 2. Evict if we hit the boundary limit
	if bm.evictList.Len() >= bm.capacity {
		oldest := bm.evictList.Back()
		if oldest != nil {
			bm.evictList.Remove(oldest)
			kv := oldest.Value.(*cacheEntry[K, V])
			delete(bm.items, kv.key) // Free map slot
		}
	}

	// 3. Insert new entry
	entry := &cacheEntry[K, V]{key: key, value: value}
	elem := bm.evictList.PushFront(entry)
	bm.items[key] = elem
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
