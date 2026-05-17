package context

import (
	"github.com/neogeny/ogopego/trees"
)

// MemoKey represents a key for memoization, combining the input mark and rule name.
type MemoKey struct {
	Mark    int    // Mark is the input position at which the rule was attempted.
	Name    string // Name is the name of the rule.
	CanMemo bool   // CanMemo indicates if the rule can be memoized.
}

// Memo stores the result of a memoized parse operation.
type Memo struct {
	Tree trees.Tree // Tree is the parse tree produced by the rule.
	Mark int        // Mark is the input position after the rule successfully parsed.
}

// IsBottomEntry checks if the memo entry represents a "bottom" (failed) parse.
func (m Memo) IsBottomEntry() bool {
	_, isBottom := m.Tree.(*trees.Bottom)
	return isBottom
}

// pruneCacheInPlace removes entries with Mark before cutpoint from cache.
func _(cache map[MemoKey]Memo, cutpoint int) {
	for k := range cache {
		if k.Mark < cutpoint {
			delete(cache, k)
		}
	}
}

// pruneCacheWithCopy returns a new map with only entries at or after cutpoint.
// Two-pass: count survivors for preallocation, then copy.
func pruneCacheWithCopy(cache map[MemoKey]Memo, cutpoint int) map[MemoKey]Memo {
	if cache == nil {
		return nil
	}

	n := 0
	for k := range cache {
		if k.Mark >= cutpoint || cache[k].IsBottomEntry() {
			n++
		}
	}

	result := make(map[MemoKey]Memo, n)
	for k, v := range cache {
		if k.Mark >= cutpoint || cache[k].IsBottomEntry() {
			result[k] = v
		}
	}
	return result
}
