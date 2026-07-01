package container

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestOrdered_ZeroCapacity(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	n := 100
	for i := range n {
		assert.NoError(t, bm.Set(intToKey(i), i))
	}
	assert.Equal(t, n, bm.Len())
	for i := range n {
		key := intToKey(i)
		v, err := bm.Get(key)
		assert.NoError(t, err, "Get(%q)", key)
		assert.Equal(t, i, v, "Get(%q)", key)
	}
}

func TestOrdered_NegativeCapacity(t *testing.T) {
	bm := NewBoundedMap[string, int](-1)
	n := 100
	for i := range n {
		assert.NoError(t, bm.Set(intToKey(i), i))
	}
	assert.Equal(t, n, bm.Len())
}

func TestOrdered_NoEviction(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	n := 1000
	for i := range n {
		assert.NoError(t, bm.Set(intToKey(i), i))
	}
	assert.Equal(t, n, bm.Len())
}

func TestOrdered_EntriesAll(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	n := 26
	for i := range n {
		assert.NoError(t, bm.Set(intToKey(i), i))
	}
	count := 0
	seen := make(map[string]bool)
	for k, v := range bm.Entries() {
		assert.False(t, seen[k], "duplicate key %q in Entries", k)
		seen[k] = true
		assert.Equal(t, keyToInt(k), v, "Entries key %q", k)
		count++
	}
	assert.Equal(t, n, count)
}

func TestOrdered_Delete(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	for i := range 5 {
		assert.NoError(t, bm.Set(intToKey(i), i))
	}
	bm.Delete(intToKey(2))
	assert.Equal(t, 4, bm.Len())
	_, err := bm.Get(intToKey(2))
	assert.Error(t, err, "Get(deleted key) should miss")
	for _, i := range []int{0, 1, 3, 4} {
		v, err := bm.Get(intToKey(i))
		assert.NoError(t, err, "Get(%q) should hit", intToKey(i))
		assert.Equal(t, i, v, "Get(%q)", intToKey(i))
	}
}

func TestOrdered_Retain(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	for i := range 10 {
		assert.NoError(t, bm.Set(intToKey(i), i))
	}
	bm.Retain(func(k string, v int) bool {
		return v%2 == 0
	})
	assert.Equal(t, 5, bm.Len())
	for i := range 10 {
		_, err := bm.Get(intToKey(i))
		if i%2 == 0 {
			assert.NoError(t, err, "even key %q should survive Retain", intToKey(i))
		} else {
			assert.Error(t, err, "odd key %q should be removed by Retain", intToKey(i))
		}
	}
}

func TestOrdered_UpdateNoDupe(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	assert.NoError(t, bm.Set("k", 1))
	assert.NoError(t, bm.Set("k", 2))
	assert.Equal(t, 1, bm.Len())
	v, err := bm.Get("k")
	assert.NoError(t, err, "Get(k) missing after update")
	assert.Equal(t, 2, v, "Get(k)")
}

func TestOrdered_EntriesOrderAfterSet(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	assert.NoError(t, bm.Set("a", 1))
	assert.NoError(t, bm.Set("b", 2))
	assert.NoError(t, bm.Set("c", 3))
	keys := collectKeys(bm)
	want := []string{"a", "b", "c"}
	assert.Equal(t, want, keys, "Entries order (insertion order)")
}

func TestOrdered_GetDoesNotReorderWhenUnbounded(t *testing.T) {
	bm := NewBoundedMap[string, int](0)
	assert.NoError(t, bm.Set("a", 1))
	assert.NoError(t, bm.Set("b", 2))
	assert.NoError(t, bm.Set("c", 3))
	bm.Get("a")
	keys := collectKeys(bm)
	want := []string{"a", "b", "c"}
	assert.Equal(t, want, keys, "Entries order after Get (Get should not reorder with capacity<=0)")
}

func TestBounded_EvictionAtCapacity(t *testing.T) {
	bm := NewBoundedMap[string, int](3)
	for i := range 5 {
		assert.NoError(t, bm.Set(intToKey(i), i))
	}
	assert.Equal(t, 3, bm.Len())
	for _, i := range []int{0, 1} {
		_, err := bm.Get(intToKey(i))
		assert.Error(t, err, "old key %q should have been evicted", intToKey(i))
	}
	for _, i := range []int{2, 3, 4} {
		_, err := bm.Get(intToKey(i))
		assert.NoError(t, err, "recent key %q should survive", intToKey(i))
	}
}

func TestBounded_LruPromotion(t *testing.T) {
	bm := NewBoundedMap[string, int](3)
	assert.NoError(t, bm.Set("a", 1))
	assert.NoError(t, bm.Set("b", 2))
	assert.NoError(t, bm.Set("c", 3))
	bm.Get("a")
	assert.NoError(t, bm.Set("d", 4))
	_, err := bm.Get("b")
	assert.Error(t, err, "Get(b) should miss — 'b' should be evicted as LRU")
	for _, k := range []string{"a", "c", "d"} {
		_, err := bm.Get(k)
		assert.NoError(t, err, "Get(%s) should hit", k)
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
