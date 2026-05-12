package peg

import (
	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/trees"
)

type Ctx = context.Ctx
type ParseError = context.ParseError
type DisasterReport = context.DisasterReport

type Model interface {
	Parse(ctx Ctx) (trees.Tree, error)
	Link(rules map[string]*Rule) error
	ValidateLinked() error
	followRef() *ModelBase

	PrettyPrint() string
	Railroads() string
}

type ModelBase struct {
	*Node
}

func (m *ModelBase) followRef() *ModelBase { return m }

func (m *ModelBase) Link(rules map[string]*Rule) error    { return nil }
func (m *ModelBase) ValidateLinked() error                { return nil }
func (m *ModelBase) PrettyPrint() string                  { return "" }
func (m *ModelBase) Railroads() string                    { return "" }
