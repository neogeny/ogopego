package peg

import (
	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/trees"
)

type Model interface {
	Parse(ctx context.Ctx) (trees.Tree, error)
	followRef() *ModelBase
}

type ModelBase struct {
	*Node
}

func (m *ModelBase) followRef() *ModelBase { return m }
