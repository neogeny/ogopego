package context

type MemoKey struct {
	Mark    int
	Name    string
	CanMemo bool
}

type Memo struct {
	Tree interface{}
	Mark int
}
