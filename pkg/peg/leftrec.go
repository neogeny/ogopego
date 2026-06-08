// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"fmt"
	"sort"
)

func callableRuleIDs(exp Model, ruleIndex map[*Rule]int) []int {
	switch e := exp.(type) {
	case *Call:
		if e.rule != nil {
			if id, ok := ruleIndex[e.rule]; ok {
				return []int{id}
			}
		}
		return nil

	case *RuleInclude:
		if e.exp != nil {
			return callableRuleIDs(e.exp, ruleIndex)
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
		*NULL, *EmptyClosure, *Cut, *Constant, *Alert,
		*MetaExp:
		return nil

	default:
		panic(fmt.Sprintf("callableRuleIDs: unhandled model type %T", exp))
	}
}

func tarjanSCC(edges [][]int) [][]int {
	n := len(edges)
	index := make([]int, n)
	lowlink := make([]int, n)
	onStack := make([]bool, n)
	stack := make([]int, 0, n)
	for i := range index {
		index[i] = -1
	}
	var currentIndex int
	var sccs [][]int

	var strongconnect func(v int)
	strongconnect = func(v int) {
		index[v] = currentIndex
		lowlink[v] = currentIndex
		currentIndex++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range edges[v] {
			if index[w] == -1 {
				strongconnect(w)
				if lowlink[w] < lowlink[v] {
					lowlink[v] = lowlink[w]
				}
			} else if onStack[w] {
				if index[w] < lowlink[v] {
					lowlink[v] = index[w]
				}
			}
		}

		if lowlink[v] == index[v] {
			var scc []int
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}
			sccs = append(sccs, scc)
		}
	}

	for v := 0; v < n; v++ {
		if index[v] == -1 {
			strongconnect(v)
		}
	}
	return sccs
}

func findCyclesInSCC(edges [][]int, scc []int, start int) [][]int {
	sccSet := make(map[int]bool)
	for _, id := range scc {
		sccSet[id] = true
	}

	var cycles [][]int
	var dfs func(node int, path []int)
	dfs = func(node int, path []int) {
		for _, p := range path {
			if p == node {
				cycles = append(cycles, append([]int{}, path...))
				return
			}
		}
		path = append(path, node)
		for _, child := range edges[node] {
			if sccSet[child] {
				newPath := append([]int{}, path...)
				dfs(child, newPath)
			}
		}
	}

	dfs(start, []int{})
	return cycles
}

func (g *Grammar) markLeftRecursion() {
	ruleIndex := make(map[*Rule]int)
	for i, rule := range g.Rules {
		rule.IsLrec = false
		rule.IsMemo = !rule.NoMemo
		ruleIndex[rule] = i
	}

	edges := make([][]int, len(g.Rules))
	ruleNames := make([]string, len(g.Rules))
	for i, rule := range g.Rules {
		edges[i] = callableRuleIDs(rule.Exp, ruleIndex)
		ruleNames[i] = rule.Name
	}

	sccs := tarjanSCC(edges)

	for _, scc := range sccs {
		if len(scc) > 1 {
			for _, id := range scc {
				g.Rules[id].IsMemo = false
			}

			leaders := make(map[int]bool)
			for _, id := range scc {
				leaders[id] = true
			}

			for _, start := range scc {
				cycles := findCyclesInSCC(edges, scc, start)
				for _, cycle := range cycles {
					cycleSet := make(map[int]bool)
					for _, id := range cycle {
						cycleSet[id] = true
					}
					for id := range leaders {
						if !cycleSet[id] {
							delete(leaders, id)
						}
					}
					if len(leaders) == 0 {
						break
					}
				}
				if len(leaders) == 0 {
					break
				}
			}

			if len(leaders) == 0 {
				for _, id := range scc {
					leaders[id] = true
				}
			}

			var leaderIDs []int
			for id := range leaders {
				leaderIDs = append(leaderIDs, id)
			}
			sort.Slice(leaderIDs, func(i, j int) bool {
				return ruleNames[leaderIDs[i]] < ruleNames[leaderIDs[j]]
			})
			g.Rules[leaderIDs[0]].IsLrec = true

		} else if len(scc) == 1 {
			id := scc[0]
			for _, child := range edges[id] {
				if child == id {
					g.Rules[id].IsLrec = true
					g.Rules[id].IsMemo = false
					break
				}
			}
		}
	}
}
