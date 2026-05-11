package trees

type Seq struct{ Items []Tree }

func (*Seq) tree() {}
func (s *Seq) Fold(gather *Merge) Tree {
	var out Tree = &Nil{}
	for _, item := range s.Items {
		out = merge(out, item.Fold(gather))
	}
	return out
}

type List struct{ Items []Tree }

func (*List) tree() {}
func (l *List) Fold(gather *Merge) Tree {
	items := make([]Tree, len(l.Items))
	for i, item := range l.Items {
		items[i] = item.Fold(gather)
	}
	return &List{Items: items}
}

type MapNode struct{ Entries map[string]Tree }

func (*MapNode) tree() {}
func (m *MapNode) Fold(gather *Merge) Tree { return m }

type RuleNode struct{ TypeName string; Tree Tree }

func (*RuleNode) tree() {}
func (r *RuleNode) Fold(gather *Merge) Tree { return r }
