package peg

type Model struct {
	*Node
}

func (m *Model) followRef() *Model { return m }
