package pyre

type Match interface {
	Group(i int) (string, bool)
	Groups() []*string
	Start() int
	End() int
	Span() (int, int)
	GroupName(name string) (string, bool)
	GroupDict() map[string]*string
	Expand(template string) string
}

type Pattern interface {
	Match(text string) (Match, bool)
	Search(text string) (Match, bool)
	FullMatch(text string) (Match, bool)
	Split(text string, maxSplit int) []string
	FindAll(text string) [][]string
	FindIter(text string) []Match
	Sub(repl, text string, count int) string
	SubN(repl, text string, count int) (string, int)
	Pattern() string
	MatchesEmpty() bool
	IsEmpty() bool
	GroupIndex() map[string]int
	GroupsCount() int
}
