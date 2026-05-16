package context

import (
	"github.com/neogeny/ogopego/trees"
)

type MemoKey struct {
	Mark    int
	Name    string
	CanMemo bool
}

type Memo struct {
	Tree trees.Tree
	Mark int
}

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
