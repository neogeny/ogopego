package peg

type Grammar struct {
	Model
	Name       string
	Directives map[string]any
	Keywords   []string
	Rules      []*Rule
}
