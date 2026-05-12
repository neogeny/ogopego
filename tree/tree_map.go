package tree

type MapNode struct{ Entries map[string]Tree }

func (*MapNode) tree()                         {}
func (m *MapNode) fold(gather *treeMerge) Tree { return m }
