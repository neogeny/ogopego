package peg

import (
	gojson "encoding/json"
	"errors"

	"github.com/neogeny/ogopego/context"
	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Ctx = context.Ctx
type MemoKey = context.MemoKey
type Tree = trees.Tree
type OrderedMap = asjson.OrderedMap

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

func (m *ModelBase) MarshalJSON() ([]byte, error) { return gojson.Marshal(m.AsJSON()) }

// modelMarshalJSON is a helper for model types to implement json.Marshaler
// by dispatching through AsJSONMixin to the correct AsJSON() override.
func modelMarshalJSON(v any) ([]byte, error) {
	if mixin, ok := v.(asjson.AsJSONMixin); ok {
		return gojson.MarshalIndent(mixin.AsJSON(), "", "  ")
	}
	return gojson.Marshal(nil)
}
