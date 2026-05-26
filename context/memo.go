package context

import (
	"github.com/neogeny/ogopego/trees"
	ctn "github.com/neogeny/ogopego/util/container"
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

type MemoCache = ctn.BoundedMap[MemoKey, Memo]

func NewMemoMache(capacity int) MemoCache {
	return ctn.NewBoundedMap[MemoKey, Memo](capacity)
}

func PruneMemoCache(cache MemoCache, cutpoint int) {
	cache.Retain(func(key MemoKey, memo Memo) bool {
		return key.Mark >= cutpoint || memo.Tree == trees.BOTTOM
	})
}

// IsBottomEntry checks if the memo entry represents a "bottom" (failed) parse.
func (m Memo) IsBottomEntry() bool {
	_, isBottom := m.Tree.(*trees.Bottom)
	return isBottom
}
