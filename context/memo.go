package context

import "github.com/neogeny/ogopego/trees"

type MemoKey struct {
	Mark    int
	Name    string
	CanMemo bool
}

type Memo struct {
	Tree trees.Tree
	Mark int
}
