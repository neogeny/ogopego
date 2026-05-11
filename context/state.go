package context

type CallStack []string

type ParseFailure struct {
	Message string
}

type DisasterReport struct {
	Start   int
	Failure ParseFailure
	CutSeen bool
}
