package peg

type Constant struct {
	Model
	Literal string
}

type Alert struct {
	Constant
	Level int
}
