package pyre

func Compile(pattern string) (Pattern, error) {
	return NewRegexp2Pattern(pattern)
}

func MustCompile(pattern string) Pattern {
	p, err := NewRegexp2Pattern(pattern)
	if err != nil {
		panic(err)
	}
	return p
}
