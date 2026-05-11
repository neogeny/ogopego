package peg

import (
	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/trees"
)

type Model interface {
	Parse(ctx context.Ctx) (trees.Tree, error)
	Link(rules map[string]*Rule) error
	ValidateLinked() error
	followRef() *ModelBase
}

type ModelBase struct {
	*Node
}

func (m *ModelBase) followRef() *ModelBase { return m }

func (m *ModelBase) Link(rules map[string]*Rule) error             { return nil }
func (m *ModelBase) ValidateLinked() error                         { return nil }
