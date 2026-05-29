package peg

import (
	"fmt"
	"sort"
)

const sentinelEOF = "＄"

func computeLA(exp Model) []string {
	base := exp.followRef()
	if base.la != nil {
		return base.la
	}

	var la []string

	switch e := exp.(type) {
	case *Token:
		la = []string{e.Token}
	case *Pattern:
		la = []string{e.Pattern}
	case *Constant:
		la = []string{e.Literal}
	case *Alert:
		la = []string{e.Literal}
	case *EOF:
		la = []string{sentinelEOF}

	case *Group, *SkipGroup, *Lookahead, *NegativeLookahead,
		*Override, *OverrideList, *Named, *NamedList,
		*Synth, *SkipTo, *Option, *Optional,
		*Closure, *PositiveClosure:
		la = computeLA(unboxExp(e))

	case *Sequence:
		for _, item := range e.Sequence {
			if _, ok := item.(*Cut); ok {
				continue
			}
			la = mergeLA(la, computeLA(item))
			if !isNullable(item) {
				break
			}
		}

	case *Choice:
		for _, opt := range e.Options {
			la = mergeLA(la, computeLA(opt))
		}

	case *Join, *PositiveJoin, *Gather, *PositiveGather:
		la = computeLA(unboxExp(e))

	case *Call:
		if e.rule != nil {
			la = []string{fmt.Sprintf("→%s", e.Name)}
		}

	case *RuleInclude:
		la = computeLA(e.exp)

	case *NULL, *Void, *Cut, *EmptyClosure,
		*Fail, *Dot, *EOL:
	// no contribution

	default:
		panic(fmt.Sprintf("computeLA: unhandled model type %T", exp))
	}

	base.la = la
	return la
}

func mergeLA(a, b []string) []string {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	seen := make(map[string]struct{}, len(a)+len(b))
	for _, s := range a {
		seen[s] = struct{}{}
	}
	for _, s := range b {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			a = append(a, s)
		}
	}
	sort.Strings(a)
	return a
}
