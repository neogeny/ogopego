package peg

import "fmt"

const (
	stateFirst = iota
	stateVisiting
	stateVisited
)

func isNullable(exp Model) bool {
	switch e := exp.(type) {
	case *Call:
		return false
	case *RuleInclude:
		if e.Exp != nil {
			return isNullable(e.Exp)
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

func callableRuleIDs(exp Model, ruleIndex map[*Rule]int) []int {
	switch e := exp.(type) {
	case *Call:
		if e.Target != nil {
			if id, ok := ruleIndex[e.Target]; ok {
				return []int{id}
			}
		}
		return nil

	case *RuleInclude:
		if e.Exp != nil {
			return callableRuleIDs(e.Exp, ruleIndex)
		}
		return nil

	case *Group, *SkipGroup, *Lookahead, *NegativeLookahead,
		*Override, *OverrideList, *Named, *NamedList,
		*Synth, *SkipTo, *Option, *Optional,
		*Closure, *PositiveClosure:
		return callableRuleIDs(unboxExp(e), ruleIndex)

	case *Join, *PositiveJoin, *Gather, *PositiveGather:
		return callableRuleIDs(unboxExp(e), ruleIndex)

	case *Choice:
		var result []int
		for _, opt := range e.Options {
			result = append(result, callableRuleIDs(opt, ruleIndex)...)
		}
		return result

	case *Sequence:
		var result []int
		for _, item := range e.Sequence {
			if _, ok := item.(*Cut); ok {
				continue
			}
			result = append(result, callableRuleIDs(item, ruleIndex)...)
			if !isNullable(item) {
				break
			}
		}
		return result

	case *Token, *Pattern, *Dot, *EOF, *EOL, *Void, *Fail,
		*NULL, *EmptyClosure, *Cut, *Constant, *Alert:
		return nil

	default:
		panic(fmt.Sprintf("callableRuleIDs: unhandled model type %T", exp))
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

func (g *Grammar) markLeftRecursion() {
	ruleIndex := make(map[*Rule]int)
	for i, rule := range g.Rules {
		rule.IsLrec = false
		rule.IsMemo = !rule.NoMemo
		ruleIndex[rule] = i
	}

	edges := make([][]int, len(g.Rules))
	for i, rule := range g.Rules {
		edges[i] = callableRuleIDs(rule.Exp, ruleIndex)
	}

	state := make([]int, len(g.Rules))
	stack := make([]int, 0, len(g.Rules))

	var dfs func(int)
	dfs = func(ruleID int) {
		switch state[ruleID] {
		case stateVisiting, stateVisited:
			return
		}
		state[ruleID] = stateVisiting
		stack = append(stack, ruleID)

		for _, childID := range edges[ruleID] {
			switch state[childID] {
			case stateFirst:
				dfs(childID)
			case stateVisiting:
				g.Rules[childID].IsLrec = true
				g.Rules[childID].IsMemo = false
				for _, id := range stack {
					g.Rules[id].IsMemo = false
				}
			}
		}

		stack = stack[:len(stack)-1]
		state[ruleID] = stateVisited
	}

	for i := range g.Rules {
		dfs(i)
	}
}
