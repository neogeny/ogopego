package peg

import (
	asjson "github.com/neogeny/ogopego/json"
	"unicode"
)

type Rule struct {
	NamedBox
	Params     []string
	KWParams   map[string]any
	Decorators []string
	Base       string
	IsName     bool
	IsTokn     bool
	NoMemo     bool
	NoStak     bool
	IsMemo     bool
	IsLrec     bool
}

func (r *Rule) IsToken() bool {
	if r.IsTokn {
		return true
	}
	for _, c := range r.Name {
		if c != '_' {
			return unicode.IsUpper(c)
		}
	}
	return false
}

func (r *Rule) IsLeftRecursive() bool { return r.IsLrec }

func (r *Rule) IsMemoizable() bool {
	return r.IsLrec || (r.IsMemo && !r.NoMemo)
}

func (r *Rule) ShouldTrace() bool {
	return !r.NoStak && !r.IsToken()
}

func (r *Rule) PubMap() *asjson.OrderedMap { return r.PubMapOf(r) }
func (r *Rule) AsJSON() any                { return r.AsJSONOf(r) }
func (r *Rule) AsJSONStr() string          { return r.AsJSONStrOf(r) }
