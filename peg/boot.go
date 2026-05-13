package peg

func LoadBootGrammar(data []byte) (*Grammar, error) {
	g, err := ParseGrammar(data)
	if err != nil {
		return nil, err
	}
	if err := g.Initialize(); err != nil {
		return nil, err
	}
	return g, nil
}
