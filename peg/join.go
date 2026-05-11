package peg

type Join struct {
	Box
	Sep any
}

type PositiveJoin struct {
	Join
}

type Gather struct {
	Join
}

type PositiveGather struct {
	Gather
}
