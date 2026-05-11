package peg

type Comment struct {
	ModelBase
	Comment string
}

type EOLComment struct {
	Comment
}
