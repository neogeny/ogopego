package context

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/trees"
)

func TestPruneCacheKeepsAfterCutpoint(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 0, Name: "a"}, Memo{})
	cache.Set(MemoKey{Mark: 5, Name: "b"}, Memo{})
	cache.Set(MemoKey{Mark: 10, Name: "c"}, Memo{})
	PruneMemoCache(cache, 5)
	_, ok := cache.Get(MemoKey{Mark: 0, Name: "a"})
	assert.False(t, ok, "expected entry at mark 0 to be removed")
	_, ok = cache.Get(MemoKey{Mark: 5, Name: "b"})
	assert.True(t, ok, "expected entry at mark 5 to be kept")
	_, ok = cache.Get(MemoKey{Mark: 10, Name: "c"})
	assert.True(t, ok, "expected entry at mark 10 to be kept")
}

func TestPruneCacheRemovesBeforeCutpoint(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 0, Name: "a"}, Memo{})
	cache.Set(MemoKey{Mark: 3, Name: "b"}, Memo{})
	PruneMemoCache(cache, 5)
	assert.Equal(t, 0, cache.Len(), "expected 0 entries")
}

func TestPruneCacheAtCutpoint(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 5, Name: "a"}, Memo{})
	PruneMemoCache(cache, 5)
	assert.Equal(t, 1, cache.Len(), "expected 1 entry")
	_, ok := cache.Get(MemoKey{Mark: 5, Name: "a"})
	assert.True(t, ok, "expected entry at mark 5 to be kept")
}

func TestPruneCacheEmpty(t *testing.T) {
	cache := NewMemoMache(64)
	PruneMemoCache(cache, 5)
	assert.Equal(t, 0, cache.Len(), "expected empty map")
}

func TestPruneCachePreservesValues(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 10, Name: "x"}, Memo{Mark: 20})
	PruneMemoCache(cache, 5)
	m, ok := cache.Get(MemoKey{Mark: 10, Name: "x"})
	assert.True(t, ok, "expected entry to be kept")
	assert.Equal(t, 20, m.Mark, "expected Mark=20")
}

func TestPruneCachePreservesBottom(t *testing.T) {
	t.Skip("pruning Bottom or not is optional and experimental")
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 0, Name: "a"}, Memo{Tree: trees.BOTTOM})
	cache.Set(MemoKey{Mark: 5, Name: "b"}, Memo{})
	PruneMemoCache(cache, 5)
	_, ok := cache.Get(MemoKey{Mark: 0, Name: "a"})
	assert.True(t, ok, "expected bottom entry at mark 0 to be preserved")
	_, ok = cache.Get(MemoKey{Mark: 5, Name: "b"})
	assert.True(t, ok, "expected entry at mark 5 to be kept")
	assert.Equal(t, 2, cache.Len(), "expected 2 entries")
}
