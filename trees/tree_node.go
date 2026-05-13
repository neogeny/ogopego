package trees

type Node struct {
	TypeName string
	Tree     Tree
}

func (*Node) tree()                         {}
func (r *Node) fold(gather *treeMerge) Tree { return r }
