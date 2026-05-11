package trees

type Text struct{ Value string }

func (*Text) tree()                                    {}
func (t *Text) Fold(gather *Merge) Tree                { return t }

type Number struct{ Value float64 }

func (*Number) tree()                                    {}
func (n *Number) Fold(gather *Merge) Tree               { return n }

type Bool struct{ Value bool }

func (*Bool) tree()                                     {}
func (b *Bool) Fold(gather *Merge) Tree                 { return b }

type Nil struct{}

func (*Nil) tree()                                     {}
func (n *Nil) Fold(gather *Merge) Tree                 { return n }

type Bottom struct{}

func (*Bottom) tree()                                   {}
func (b *Bottom) Fold(gather *Merge) Tree               { return b }
