package peg

type Call struct {
	ModelBase
	Name   string
	Target *Rule
}
