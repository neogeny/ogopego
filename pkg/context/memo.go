package context

import (
	"github.com/neogeny/ogopego/pkg/trees"
)

// MemoKey represents a key for memoization, combining the input mark and rule name.
type MemoKey struct {
	Mark    int    // Mark is the input position at which the rule was attempted.
	Name    string // Name is the name of the rule.
	CanMemo bool   // CanMemo indicates if the rule can be memoized.
}

// Memo stores the result of a memoized parse operation.
type Memo struct {
	Tree any // Tree is the parse tree produced by the rule.
	Mark int // Mark is the input position after the rule successfully parsed.
}

type MemoCache = map[MemoKey]Memo

func NewMemoMache(_capacity int) MemoCache {
	return make(MemoCache)
}

func PruneMemoCache(cache MemoCache, cutpoint int) {
	Retain(cache, func(key MemoKey, memo Memo) bool {
		// NOTE
		// 	Keep trees.BOTTOM for the sake of left recursion
		//  To keep them is easier than to calculate when they can be prunned
		return key.Mark >= cutpoint || memo.Tree == trees.BOTTOM
	})
}

// IsBottomEntry checks if the memo entry represents a "bottom" (failed) parse.
func (m Memo) IsBottomEntry() bool {
	return m.Tree == trees.BOTTOM
}

// Retain removes all key-value pairs that do not satisfy the keep function.
func Retain[K comparable, V any](m map[K]V, keep func(K, V) bool) {
	for k, v := range m {
		if !keep(k, v) {
			delete(m, k)
		}
	}
}
