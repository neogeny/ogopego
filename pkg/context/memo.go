package context

import (
	"unique"

	"github.com/neogeny/ogopego/pkg/trees"
)

// MemoKey represents a key for memoization, combining the input mark and rule name.
type MemoKey struct {
	Mark       int // Mark is the input position at which the rule was attempted.
	NameHandle unique.Handle[string]
	CanMemo    bool // CanMemo indicates if the rule can be memoized.
}

// Memo stores the result of a memoized parse operation.
type Memo struct {
	Tree any // Tree is the parse tree produced by the rule.
	Mark int // Mark is the input position after the rule successfully parsed.
}

type MemoEntries = map[MemoKey]Memo

type MemoCache struct {
	entries  MemoEntries
	oldKeys  []MemoKey
	oldMemos []Memo
}

func NewMemoCache(_capacity int) MemoCache {
	return MemoCache{entries: make(MemoEntries)}
}

func _newMemoKey(mark int, name string, canMemo bool) MemoKey {
	return MemoKey{
		Mark:       mark,
		NameHandle: unique.Make(name),
		CanMemo:    canMemo,
	}
}

func (k MemoKey) Name() string {
	return k.NameHandle.Value()
}

func (cache *MemoCache) NewKey(mark int, name string, canMemo bool) MemoKey {
	if len(cache.oldKeys) == 0 {
		return _newMemoKey(mark, name, canMemo)
	}
	p := len(cache.oldKeys) - 1
	key := cache.oldKeys[p]
	cache.oldKeys = cache.oldKeys[0:p]
	key.Mark = mark
	key.NameHandle = unique.Make(name)
	key.CanMemo = canMemo
	return key
}

func (cache *MemoCache) Put(key MemoKey, tree any, mark int) {
	var memo Memo
	if len(cache.oldMemos) == 0 {
		memo = Memo{tree, mark}
	} else {
		p := len(cache.oldMemos) - 1
		memo = cache.oldMemos[p]
		memo.Tree = tree
		memo.Mark = mark
		cache.oldMemos = cache.oldMemos[0:p]
	}
	cache.entries[key] = memo
}

func (cache *MemoCache) Get(key MemoKey) (Memo, bool) {
	memo, ok := cache.entries[key]
	return memo, ok
}

func (cache *MemoCache) Prune(cutpoint int) {
	// FIXME
	// WARNING diaabled to experiment with CPU and MEM profiling
	// cache.Retain(func(key MemoKey, memo Memo) bool {
	// 	// NOTE
	// 	// 	Keep trees.BOTTOM for the sake of left recursion
	// 	//  To keep them is easier than to calculate when they can be prunned
	// 	return key.Mark >= cutpoint || memo.Tree == trees.BOTTOM
	// })
}

// IsBottomEntry checks if the memo entry represents a "bottom" (failed) parse.
func (m Memo) IsBottomEntry() bool {
	return m.Tree == trees.BOTTOM
}

// Retain removes all key-value pairs that do not satisfy the keep function.
func (cache *MemoCache) Retain(keep func(MemoKey, Memo) bool) {
	for k, v := range cache.entries {
		if !keep(k, v) {
			delete(cache.entries, k)
			cache.oldKeys = append(cache.oldKeys, k)
			cache.oldMemos = append(cache.oldMemos, v)
		}
	}
}
