package peg

import "fmt"

func isNullable(exp Model) bool {
	switch e := exp.(type) {
	case *Call:
		return false
	case *RuleInclude:
		if e.exp != nil {
			return isNullable(e.exp)
		}
		return false

	case *Group, *SkipGroup, *Lookahead, *NegativeLookahead,
		*Override, *OverrideList, *Named, *NamedList,
		*Synth, *Option:
		return isNullable(unboxExp(e))

	case *Optional:
		return true

	case *Closure:
		return true
	case *PositiveClosure:
		return isNullable(e.Exp)

	case *Join, *Gather:
		return true
	case *PositiveJoin:
		return isNullable(e.Exp)
	case *PositiveGather:
		return isNullable(e.Exp)

	case *Choice:
		for _, opt := range e.Options {
			if isNullable(opt) {
				return true
			}
		}
		return false

	case *Sequence:
		for _, item := range e.Sequence {
			if !isNullable(item) {
				return false
			}
		}
		return true

	case *EOL, *Void, *NULL, *EmptyClosure, *Cut, *Constant, *Alert:
		return true

	case *Token, *Pattern, *Dot, *EOF, *Fail, *SkipTo:
		return false

	default:
		panic(fmt.Sprintf("isNullable: unhandled model type %T", exp))
	}
}

func unboxExp(e Model) Model {
	switch e2 := e.(type) {
	case *Group:
		return e2.Exp
	case *SkipGroup:
		return e2.Exp
	case *Lookahead:
		return e2.Exp
	case *NegativeLookahead:
		return e2.Exp
	case *Override:
		return e2.Exp
	case *OverrideList:
		return e2.Exp
	case *Named:
		return e2.Exp
	case *NamedList:
		return e2.Exp
	case *Synth:
		return e2.Exp
	case *SkipTo:
		return e2.Exp
	case *Option:
		return e2.Exp
	case *Optional:
		return e2.Exp
	case *Closure:
		return e2.Exp
	case *PositiveClosure:
		return e2.Exp
	case *Join:
		return e2.Exp
	case *PositiveJoin:
		return e2.Exp
	case *Gather:
		return e2.Exp
	case *PositiveGather:
		return e2.Exp
	default:
		panic(fmt.Sprintf("unboxExp: unhandled model type %T", e))
	}
}
