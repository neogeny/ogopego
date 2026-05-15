package pyre

func Compile(pattern string) (Pattern, error) {
	return NewPCRE2CgoPattern(pattern)
}

func MustCompile(pattern string) Pattern {
	p, err := NewPCRE2CgoPattern(pattern)
	if err != nil {
		panic(err)
	}
	return p
}
