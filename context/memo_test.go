package context

import (
	"testing"

	"github.com/neogeny/ogopego/trees"
)

func TestPruneCacheKeepsAfterCutpoint(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 0, Name: "a"}:  {},
		{Mark: 5, Name: "b"}:  {},
		{Mark: 10, Name: "c"}: {},
	}
	pruneCacheInPlace(cache, 5)
	if len(cache) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cache))
	}
	if _, ok := cache[MemoKey{Mark: 0, Name: "a"}]; ok {
		t.Error("expected entry at mark 0 to be removed")
	}
	if _, ok := cache[MemoKey{Mark: 5, Name: "b"}]; !ok {
		t.Error("expected entry at mark 5 to be kept")
	}
	if _, ok := cache[MemoKey{Mark: 10, Name: "c"}]; !ok {
		t.Error("expected entry at mark 10 to be kept")
	}
}

func TestPruneCacheRemovesBeforeCutpoint(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 0, Name: "a"}: {},
		{Mark: 3, Name: "b"}: {},
	}
	pruneCacheInPlace(cache, 5)
	if len(cache) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(cache))
	}
}

func TestPruneCacheAtCutpoint(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 5, Name: "a"}: {},
	}
	pruneCacheInPlace(cache, 5)
	if len(cache) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(cache))
	}
	if _, ok := cache[MemoKey{Mark: 5, Name: "a"}]; !ok {
		t.Error("expected entry at mark 5 to be kept")
	}
}

func TestPruneCacheEmpty(t *testing.T) {
	cache := map[MemoKey]Memo{}
	pruneCacheInPlace(cache, 5)
	if len(cache) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(cache))
	}
}

func TestPruneCacheNil(t *testing.T) {
	pruneCacheInPlace(nil, 5)
}

func TestPruneCachePreservesValues(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 10, Name: "x"}: {Mark: 20},
	}
	pruneCacheInPlace(cache, 5)
	m, ok := cache[MemoKey{Mark: 10, Name: "x"}]
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if m.Mark != 20 {
		t.Errorf("expected Mark=20, got %d", m.Mark)
	}
}

func TestPruneCachePreservesBottom(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 0, Name: "a"}: {Tree: trees.BOTTOM},
		{Mark: 5, Name: "b"}: {},
	}
	pruneCacheInPlace(cache, 5)
	if _, ok := cache[MemoKey{Mark: 0, Name: "a"}]; !ok {
		t.Error("expected bottom entry at mark 0 to be preserved")
	}
	if _, ok := cache[MemoKey{Mark: 5, Name: "b"}]; !ok {
		t.Error("expected entry at mark 5 to be kept")
	}
	if len(cache) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cache))
	}
}
