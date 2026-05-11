package peg

type Grammar struct {
	ModelBase
	Name       string
	Directives map[string]any
	Keywords   []string
	Rules      []*Rule
}
