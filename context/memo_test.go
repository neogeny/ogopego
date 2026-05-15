package context

import (
	"testing"
)

func TestPruneCacheKeepsAfterCutpoint(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 0, Name: "a"}:  {},
		{Mark: 5, Name: "b"}:  {},
		{Mark: 10, Name: "c"}: {},
	}
	result := pruneCacheWithCopy(cache, 5)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if _, ok := result[MemoKey{Mark: 0, Name: "a"}]; ok {
		t.Error("expected entry at mark 0 to be removed")
	}
	if _, ok := result[MemoKey{Mark: 5, Name: "b"}]; !ok {
		t.Error("expected entry at mark 5 to be kept")
	}
	if _, ok := result[MemoKey{Mark: 10, Name: "c"}]; !ok {
		t.Error("expected entry at mark 10 to be kept")
	}
}

func TestPruneCacheRemovesBeforeCutpoint(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 0, Name: "a"}: {},
		{Mark: 3, Name: "b"}: {},
	}
	result := pruneCacheWithCopy(cache, 5)
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestPruneCacheAtCutpoint(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 5, Name: "a"}: {},
	}
	result := pruneCacheWithCopy(cache, 5)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if _, ok := result[MemoKey{Mark: 5, Name: "a"}]; !ok {
		t.Error("expected entry at mark 5 to be kept")
	}
}

func TestPruneCacheEmpty(t *testing.T) {
	result := pruneCacheWithCopy(map[MemoKey]Memo{}, 5)
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(result))
	}
}

func TestPruneCacheNil(t *testing.T) {
	result := pruneCacheWithCopy(nil, 5)
	if result != nil {
		t.Fatal("expected nil")
	}
}

func TestPruneCachePreservesValues(t *testing.T) {
	cache := map[MemoKey]Memo{
		{Mark: 10, Name: "x"}: {Mark: 20},
	}
	result := pruneCacheWithCopy(cache, 5)
	m, ok := result[MemoKey{Mark: 10, Name: "x"}]
	if !ok {
		t.Fatal("expected entry to be kept")
	}
	if m.Mark != 20 {
		t.Errorf("expected Mark=20, got %d", m.Mark)
	}
}
