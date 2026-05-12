package tree

type Text struct{ Value string }

func (*Text) tree()                         {}
func (t *Text) fold(gather *treeMerge) Tree { return t }

type Number struct{ Value float64 }

func (*Number) tree()                         {}
func (n *Number) fold(gather *treeMerge) Tree { return n }

type Bool struct{ Value bool }

func (*Bool) tree()                         {}
func (b *Bool) fold(gather *treeMerge) Tree { return b }

type Nil struct{}

func (*Nil) tree()                         {}
func (n *Nil) fold(gather *treeMerge) Tree { return n }

type Bottom struct{}

func (*Bottom) tree()                         {}
func (b *Bottom) fold(gather *treeMerge) Tree { return b }
