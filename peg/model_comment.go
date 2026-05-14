package peg

import (
	"errors"

	asjson "github.com/neogeny/ogopego/json"
	"github.com/neogeny/ogopego/trees"
)

type Comment struct {
	ModelBase
	Comment string
}

func (c *Comment) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("exp Comment not yet implemented"))
}

func (t *Comment) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *Comment) AsJSON() any                { return t.AsJSONOf(t) }
func (t *Comment) AsJSONStr() string          { return t.AsJSONStrOf(t) }

type EOLComment struct {
	Comment
}

func (e *EOLComment) Parse(ctx Ctx) (trees.Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("exp EOLComment not yet implemented"))
}

func (t *EOLComment) PubMap() *asjson.OrderedMap { return t.PubMapOf(t) }
func (t *EOLComment) AsJSON() any                { return t.AsJSONOf(t) }
func (t *EOLComment) AsJSONStr() string          { return t.AsJSONStrOf(t) }
