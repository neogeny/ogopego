package peg

type Model struct{}

func (m *Model) followRef() *Model { return m }
