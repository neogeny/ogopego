package trees

type Named struct{ Name string; Value Tree }

func (*Named) tree() {}
func (n *Named) Fold(gather *Merge) Tree {
	val := n.Value.Fold(gather)
	gather.insert(n.Name, val)
	return val
}

type NamedAsList struct{ Name string; Value Tree }

func (*NamedAsList) tree() {}
func (n *NamedAsList) Fold(gather *Merge) Tree {
	val := n.Value.Fold(gather)
	gather.insertAsList(n.Name, val)
	return val
}

type Override struct{ Value Tree }

func (*Override) tree() {}
func (o *Override) Fold(gather *Merge) Tree {
	val := o.Value.Fold(gather)
	gather.Root = appendTree(gather.Root, val)
	return val
}

type OverrideAsList struct{ Value Tree }

func (*OverrideAsList) tree() {}
func (o *OverrideAsList) Fold(gather *Merge) Tree {
	val := o.Value.Fold(gather)
	gather.Root = appendAsList(gather.Root, val)
	return val
}
