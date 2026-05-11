package peg

type Comment struct {
	Model
	Comment string
}

type EOLComment struct {
	Comment
}
