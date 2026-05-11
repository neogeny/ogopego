package json

import "github.com/neogeny/ogopego/peg"

func LoadBootGrammar(data []byte) (*peg.Grammar, error) {
	g, err := ParseGrammar(data)
	if err != nil {
		return nil, err
	}
	if err := g.Initialize(); err != nil {
		return nil, err
	}
	return g, nil
}
