package tree

type Named struct {
	Name  string
	Value Tree
}

func (*Named) tree() {}
func (n *Named) fold(gather *treeMerge) Tree {
	val := n.Value.fold(gather)
	gather.insert(n.Name, val)
	return val
}

type NamedAsList struct {
	Name  string
	Value Tree
}

func (*NamedAsList) tree() {}
func (n *NamedAsList) fold(gather *treeMerge) Tree {
	val := n.Value.fold(gather)
	gather.insertAsList(n.Name, val)
	return val
}

type Override struct{ Value Tree }

func (*Override) tree() {}
func (o *Override) fold(gather *treeMerge) Tree {
	val := o.Value.fold(gather)
	gather.Root = appendTree(gather.Root, val)
	return val
}

type OverrideAsList struct{ Value Tree }

func (*OverrideAsList) tree() {}
func (o *OverrideAsList) fold(gather *treeMerge) Tree {
	val := o.Value.fold(gather)
	gather.Root = appendAsList(gather.Root, val)
	return val
}
