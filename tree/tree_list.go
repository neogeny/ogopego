package tree

type Seq struct{ Items []Tree }

func (*Seq) tree() {}
func (s *Seq) fold(gather *treeMerge) Tree {
	var out Tree = &Nil{}
	for _, item := range s.Items {
		out = merge(out, item.fold(gather))
	}
	return out
}

type List struct{ Items []Tree }

func (*List) tree() {}
func (l *List) fold(gather *treeMerge) Tree {
	items := make([]Tree, len(l.Items))
	for i, item := range l.Items {
		items[i] = item.fold(gather)
	}
	return &List{Items: items}
}
