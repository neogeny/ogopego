package peg

type Rule struct {
	NamedBox
	Params     []string
	KWParams   map[string]any
	Decorators []string
	Base       string
	IsName     bool
	IsTokn     bool
	NoMemo     bool
	NoStak     bool
	IsMemo     bool
	IsLrec     bool
}
