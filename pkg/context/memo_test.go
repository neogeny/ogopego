package context

import (
	"testing"

	"github.com/neogeny/ogopego/pkg/trees"
)

func TestPruneCacheKeepsAfterCutpoint(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 0, Name: "a"}, Memo{})
	cache.Set(MemoKey{Mark: 5, Name: "b"}, Memo{})
	cache.Set(MemoKey{Mark: 10, Name: "c"}, Memo{})
	PruneMemoCache(cache, 5)
	if _, ok := cache.Get(MemoKey{Mark: 0, Name: "a"}); ok {
		t.Error("expected entry at mark 0 to be removed")
	}
	if _, ok := cache.Get(MemoKey{Mark: 5, Name: "b"}); !ok {
		t.Error("expected entry at mark 5 to be kept")
	}
	if _, ok := cache.Get(MemoKey{Mark: 10, Name: "c"}); !ok {
		t.Error("expected entry at mark 10 to be kept")
	}
}

func TestPruneCacheRemovesBeforeCutpoint(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 0, Name: "a"}, Memo{})
	cache.Set(MemoKey{Mark: 3, Name: "b"}, Memo{})
	PruneMemoCache(cache, 5)
	if cache.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", cache.Len())
	}
}

func TestPruneCacheAtCutpoint(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 5, Name: "a"}, Memo{})
	PruneMemoCache(cache, 5)
	if cache.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", cache.Len())
	}
	if _, ok := cache.Get(MemoKey{Mark: 5, Name: "a"}); !ok {
		t.Error("expected entry at mark 5 to be kept")
	}
}

func TestPruneCacheEmpty(t *testing.T) {
	cache := NewMemoMache(64)
	PruneMemoCache(cache, 5)
	if cache.Len() != 0 {
		t.Fatalf("expected empty map, got %d entries", cache.Len())
	}
}

func TestPruneCachePreservesValues(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 10, Name: "x"}, Memo{Mark: 20})
	PruneMemoCache(cache, 5)
	m, ok := cache.Get(MemoKey{Mark: 10, Name: "x"})
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if m.Mark != 20 {
		t.Errorf("expected Mark=20, got %d", m.Mark)
	}
}

func TestPruneCachePreservesBottom(t *testing.T) {
	cache := NewMemoMache(64)
	cache.Set(MemoKey{Mark: 0, Name: "a"}, Memo{Tree: trees.BOTTOM})
	cache.Set(MemoKey{Mark: 5, Name: "b"}, Memo{})
	PruneMemoCache(cache, 5)
	if _, ok := cache.Get(MemoKey{Mark: 0, Name: "a"}); !ok {
		t.Error("expected bottom entry at mark 0 to be preserved")
	}
	if _, ok := cache.Get(MemoKey{Mark: 5, Name: "b"}); !ok {
		t.Error("expected entry at mark 5 to be kept")
	}
	if cache.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", cache.Len())
	}
}
