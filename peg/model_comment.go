package peg

import (
	"errors"
)

type Comment struct {
	ModelBase
	Comment string
}

type EOLComment struct {
	Comment
}

func (c *Comment) Parse(ctx Ctx) (Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("exp Comment not yet implemented"))
}

func (t *Comment) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *Comment) AsJSON() any         { return t.AsJSONOf(t) }
func (t *Comment) AsJSONStr() string   { return t.AsJSONStrOf(t) }

func (e *EOLComment) Parse(ctx Ctx) (Tree, error) {
	return nil, ctx.Failure(ctx.Mark(), errors.New("exp EOLComment not yet implemented"))
}

func (t *EOLComment) PubMap() *OrderedMap { return t.PubMapOf(t) }
func (t *EOLComment) AsJSON() any         { return t.AsJSONOf(t) }
func (t *EOLComment) AsJSONStr() string   { return t.AsJSONStrOf(t) }
