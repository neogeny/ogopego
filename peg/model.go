package peg

import (
	"errors"

	"github.com/neogeny/ogopego/context"
	"github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Ctx = context.Ctx
type ParseFailure = context.ParseFailure
type MemoKey = context.MemoKey
type Tree = trees.Tree
type OrderedMap = json.OrderedMap

var NIL = trees.NIL

type Model interface {
	Parse(ctx Ctx) (Tree, error)
	Link(rules map[string]*Rule) error
	ValidateLinked() error
	followRef() *ModelBase

	PrettyPrint() string
	Railroads() string
}

type ModelBase struct {
	Model
	Node
}

func (m *ModelBase) followRef() *ModelBase { return m }

func (m *ModelBase) Parse(ctx Ctx) (Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("Parse not implemented"))
}
func (m *ModelBase) Link(rules map[string]*Rule) error { return nil }
func (m *ModelBase) ValidateLinked() error             { return nil }
func (m *ModelBase) PrettyPrint() string               { return "" }
func (m *ModelBase) Railroads() string                 { return "" }
