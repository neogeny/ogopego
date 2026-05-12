package tree

type TreeNode struct {
	TypeName string
	Tree     Tree
}

func (*TreeNode) tree()                         {}
func (r *TreeNode) fold(gather *treeMerge) Tree { return r }
