package peg

type Constant struct {
	ModelBase
	Literal string
}

type Alert struct {
	Constant
	Level int
}
