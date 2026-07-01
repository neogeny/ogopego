package context

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/pkg/trees"
)

func TestPruneCacheKeepsAfterCutpoint(t *testing.T) {
	t.Skip("MemoCache.Prune() disabled")
	cache := NewMemoCache(64)
	cache.Put(cache.NewKey(0, "a", false), nil, 0)
	cache.Put(cache.NewKey(5, "b", false), nil, 0)
	cache.Put(cache.NewKey(10, "c", false), nil, 0)
	cache.Prune(5)
	_, ok := cache.Get(cache.NewKey(0, "a", false))
	assert.False(t, ok, "expected entry at mark 0 to be removed")
	_, ok = cache.Get(cache.NewKey(5, "b", false))
	assert.True(t, ok, "expected entry at mark 5 to be kept")
	_, ok = cache.Get(cache.NewKey(10, "c", false))
	assert.True(t, ok, "expected entry at mark 10 to be kept")
}

func TestPruneCacheRemovesBeforeCutpoint(t *testing.T) {
	t.Skip("MemoCache.Prune() disabled")
	cache := NewMemoCache(64)
	cache.Put(cache.NewKey(0, "a", false), nil, 0)
	cache.Put(cache.NewKey(3, "b", false), nil, 0)
	cache.Prune(5)
	assert.Equal(t, 0, len(cache.entries), "expected 0 entries")
}

func TestPruneCacheAtCutpoint(t *testing.T) {
	t.Skip("MemoCache.Prune() disabled")
	cache := NewMemoCache(64)
	cache.Put(cache.NewKey(5, "a", false), nil, 0)
	cache.Prune(5)
	assert.Equal(t, 1, len(cache.entries), "expected 1 entry")
	_, ok := cache.Get(cache.NewKey(5, "a", false))
	assert.True(t, ok, "expected entry at mark 5 to be kept")
}

func TestPruneCacheEmpty(t *testing.T) {
	t.Skip("MemoCache.Prune() disabled")
	cache := NewMemoCache(64)
	cache.Prune(5)
	assert.Equal(t, 0, len(cache.entries), "expected empty map")
}

func TestPruneCachePreservesValues(t *testing.T) {
	t.Skip("MemoCache.Prune() disabled")
	cache := NewMemoCache(64)
	cache.Put(cache.NewKey(10, "x", false), nil, 20)
	cache.Prune(5)
	m, ok := cache.Get(cache.NewKey(10, "x", false))
	assert.True(t, ok, "expected entry to be kept")
	assert.Equal(t, 20, m.Mark, "expected Mark=20")
}

func TestPruneCachePreservesBottom(t *testing.T) {
	t.Skip("MemoCache.Prune() disabled")
	cache := NewMemoCache(64)
	cache.Put(cache.NewKey(0, "a", false), trees.BOTTOM, 0)
	cache.Put(cache.NewKey(5, "b", false), nil, 0)
	cache.Prune(5)
	_, ok := cache.Get(cache.NewKey(0, "a", false))
	assert.True(t, ok, "expected bottom entry at mark 0 to be preserved")
	_, ok = cache.Get(cache.NewKey(5, "b", false))
	assert.True(t, ok, "expected entry at mark 5 to be kept")
	assert.Equal(t, 2, len(cache.entries), "expected 2 entries")
}
