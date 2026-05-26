package container

import (
	"slices"
	"testing"
)

func TestOrdered_ZeroCapacity(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	n := 100
	for i := range n {
		bm.Set(intToKey(i), i)
	}
	if got := bm.Len(); got != n {
		t.Errorf("Len() = %d, want %d", got, n)
	}
	for i := range n {
		key := intToKey(i)
		v, ok := bm.Get(key)
		if !ok {
			t.Errorf("Get(%q) missing", key)
		}
		if v != i {
			t.Errorf("Get(%q) = %d, want %d", key, v, i)
		}
	}
}

func TestOrdered_NegativeCapacity(t *testing.T) {
	bm := NewBoundedMap[string, int](-1)
	n := 100
	for i := range n {
		bm.Set(intToKey(i), i)
	}
	if got := bm.Len(); got != n {
		t.Errorf("Len() = %d, want %d", got, n)
	}
}

func TestOrdered_NoEviction(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	n := 1000
	for i := range n {
		bm.Set(intToKey(i), i)
	}
	if got := bm.Len(); got != n {
		t.Errorf("Len() = %d, want %d", got, n)
	}
}

func TestOrdered_EntriesAll(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	n := 26
	for i := range n {
		bm.Set(intToKey(i), i)
	}
	count := 0
	seen := make(map[string]bool)
	for k, v := range bm.Entries() {
		if seen[k] {
			t.Errorf("duplicate key %q in Entries", k)
		}
		seen[k] = true
		if v != keyToInt(k) {
			t.Errorf("Entries key %q has value %d, want %d", k, v, keyToInt(k))
		}
		count++
	}
	if count != n {
		t.Errorf("Entries yielded %d items, want %d", count, n)
	}
}

func TestOrdered_Delete(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	for i := range 5 {
		bm.Set(intToKey(i), i)
	}
	bm.Delete(intToKey(2))
	if got := bm.Len(); got != 4 {
		t.Errorf("Len() after Delete = %d, want 4", got)
	}
	if _, ok := bm.Get(intToKey(2)); ok {
		t.Error("Get(deleted key) should miss")
	}
	for _, i := range []int{0, 1, 3, 4} {
		v, ok := bm.Get(intToKey(i))
		if !ok {
			t.Errorf("Get(%q) should hit", intToKey(i))
		}
		if v != i {
			t.Errorf("Get(%q) = %d, want %d", intToKey(i), v, i)
		}
	}
}

func TestOrdered_Retain(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	for i := range 10 {
		bm.Set(intToKey(i), i)
	}
	bm.Retain(func(k string, v int) bool {
		return v%2 == 0
	})
	if got := bm.Len(); got != 5 {
		t.Errorf("Len() after Retain = %d, want 5", got)
	}
	for i := range 10 {
		_, ok := bm.Get(intToKey(i))
		if i%2 == 0 && !ok {
			t.Errorf("even key %q should survive Retain", intToKey(i))
		}
		if i%2 != 0 && ok {
			t.Errorf("odd key %q should be removed by Retain", intToKey(i))
		}
	}
}

func TestOrdered_UpdateNoDupe(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	bm.Set("k", 1)
	bm.Set("k", 2)
	if got := bm.Len(); got != 1 {
		t.Errorf("Len() after update = %d, want 1", got)
	}
	v, ok := bm.Get("k")
	if !ok {
		t.Fatal("Get(k) missing after update")
	}
	if v != 2 {
		t.Errorf("Get(k) = %d, want 2", v)
	}
}

func TestOrdered_EntriesOrderAfterSet(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	bm.Set("a", 1)
	bm.Set("b", 2)
	bm.Set("c", 3)
	keys := collectKeys(bm)
	want := []string{"a", "b", "c"}
	if !slices.Equal(keys, want) {
		t.Errorf("Entries order = %v, want %v (insertion order)", keys, want)
	}
}

func TestOrdered_GetDoesNotReorderWhenUnbounded(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	bm.Set("a", 1)
	bm.Set("b", 2)
	bm.Set("c", 3)
	bm.Get("a")
	keys := collectKeys(bm)
	want := []string{"a", "b", "c"}
	if !slices.Equal(keys, want) {
		t.Errorf("Entries order after Get = %v, want %v (Get should not reorder with capacity<=0)", keys, want)
	}
}

func TestBounded_EvictionAtCapacity(t *testing.T) {
	bm := NewBoundedMap[string, int](3)
	for i := range 5 {
		bm.Set(intToKey(i), i)
	}
	if got := bm.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}
	for _, i := range []int{0, 1} {
		if _, ok := bm.Get(intToKey(i)); ok {
			t.Errorf("old key %q should have been evicted", intToKey(i))
		}
	}
	for _, i := range []int{2, 3, 4} {
		if _, ok := bm.Get(intToKey(i)); !ok {
			t.Errorf("recent key %q should survive", intToKey(i))
		}
	}
}

func TestBounded_LruPromotion(t *testing.T) {
	bm := NewBoundedMap[string, int](3)
	bm.Set("a", 1)
	bm.Set("b", 2)
	bm.Set("c", 3)
	bm.Get("a")
	bm.Set("d", 4)
	if _, ok := bm.Get("b"); ok {
		t.Error("Get(b) should miss — 'b' should be evicted as LRU")
	}
	for _, k := range []string{"a", "c", "d"} {
		if _, ok := bm.Get(k); !ok {
			t.Errorf("Get(%s) should hit", k)
		}
	}
}

// helpers

func intToKey(i int) string {
	return string(rune('a' + i))
}

func keyToInt(k string) int {
	return int(rune(k[0]) - 'a')
}

func collectKeys(bm BoundedMap[string, int]) []string {
	var keys []string
	for k := range bm.Entries() {
		keys = append(keys, k)
	}
	return keys
}
