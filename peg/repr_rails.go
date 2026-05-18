// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package peg

import (
	"fmt"
	"strings"
)

type rails []string

// String returns the string representation of the rails.
func (r rails) String() string {
	lines := make([]string, len(r))
	for i, s := range r {
		lines[i] = strings.TrimRight(s, " \t")
	}
	return strings.Join(lines, "\n")
}

const (
	etx = "\x03"
	eol = "\x0a"
)

func ulen(s string) int {
	// ASCII-only for now; full Unicode width support would need a library
	return len(s)
}

func pad(s string, c rune, width int) string {
	n := ulen(s)
	if n >= width {
		return s
	}
	return s + strings.Repeat(string(c), width-n)
}

func railpad(s string, maxl int) string { return pad(s, '─', maxl) }

func blankpad(s string, maxl int) string { return pad(s, ' ', maxl) }

func weld(a, b []string) []string {
	if len(a) == 0 {
		out := make([]string, len(b))
		copy(out, b)
		return out
	}
	if len(b) == 0 {
		out := make([]string, len(a))
		copy(out, a)
		return out
	}
	for _, s := range a {
		if strings.Contains(s, etx) {
			out := make([]string, len(a))
			copy(out, a)
			return out
		}
	}

	lenA := ulen(a[0])
	lenB := ulen(b[0])
	height := len(a)
	if len(b) > height {
		height = len(b)
	}
	common := len(a)
	if len(b) < common {
		common = len(b)
	}

	out := make([]string, height)
	for i := 0; i < height; i++ {
		switch {
		case i < common:
			out[i] = a[i] + b[i]
		case i < len(a):
			out[i] = a[i] + strings.Repeat(" ", lenB)
		default:
			out[i] = strings.Repeat(" ", lenA) + b[i]
		}
	}
	return out
}

func layOut(tracks [][]string) []string {
	if len(tracks) == 0 {
		return nil
	}
	if len(tracks) == 1 {
		out := make([]string, len(tracks[0]))
		copy(out, tracks[0])
		return out
	}

	maxl := 0
	for _, track := range tracks {
		if len(track) > 0 {
			if l := ulen(track[0]); l > maxl {
				maxl = l
			}
		}
	}

	var out []string

	for ti, track := range tracks {
		if len(track) == 0 {
			continue
		}
		isLast := ti == len(tracks)-1
		joint := track[0]

		if !isLast {
			if !strings.Contains(joint, etx) {
				out = append(out, fmt.Sprintf("  ├─%s─┤ ", railpad(joint, maxl)))
			} else {
				out = append(out, fmt.Sprintf("  ├─%s │ ", blankpad(joint, maxl)))
			}
			for _, rail := range track[1:] {
				out = append(out, fmt.Sprintf("  │ %s │ ", blankpad(rail, maxl)))
			}
		} else {
			if !strings.Contains(joint, etx) {
				// Check if we can connect to previous track's bottom
				if len(out) > 0 {
					last := out[len(out)-1]
					if strings.HasSuffix(last, "─┤ ") {
						out[len(out)-1] = strings.TrimSuffix(last, "─┤ ") + "─┘ "
					}
				}
				out = append(out, fmt.Sprintf("  └─%s─┘ ", railpad(joint, maxl)))
			} else {

				out = append(out, fmt.Sprintf("  └─%s   ", blankpad(joint, maxl)))
			}
			for _, rail := range track[1:] {
				out = append(out, fmt.Sprintf("    %s   ", blankpad(rail, maxl)))
			}
		}
	}

	if len(out) > 0 && len(tracks) > 0 {
		firstTrack := tracks[0]
		if len(firstTrack) > 0 {
			joint := firstTrack[0]
			if !strings.Contains(joint, etx) {
				out[0] = fmt.Sprintf("──┬─%s─┬─", railpad(joint, maxl))
			} else {
				out[0] = fmt.Sprintf("──┬─%s ┬─", blankpad(joint, maxl))
			}
		}
	}

	return out
}

func loopRails(rs []string) []string {
	if len(rs) == 0 {
		return []string{"───>───"}
	}
	maxl := 0
	for _, s := range rs {
		if l := ulen(s); l > maxl {
			maxl = l
		}
	}
	out := []string{fmt.Sprintf("──┬→%s─┬──", railpad("", maxl))}
	out = append(out, fmt.Sprintf("  ├→%s─┤  ", railpad(rs[0], maxl)))
	for _, rail := range rs[1:] {
		out = append(out, fmt.Sprintf("  │ %s │  ", blankpad(rail, maxl)))
	}
	out = append(out, fmt.Sprintf("  └─%s<┘  ", railpad("", maxl)))
	return out
}

func stopnloop(rs []string) []string {
	if len(rs) == 0 {
		return []string{"───>───"}
	}
	maxl := 0
	for _, s := range rs {
		if l := ulen(s); l > maxl {
			maxl = l
		}
	}
	out := []string{fmt.Sprintf("──┬─%s─┬──", railpad(rs[0], maxl))}
	for _, rail := range rs[1:] {
		out = append(out, fmt.Sprintf("  │ %s │  ", blankpad(rail, maxl)))
	}
	out = append(out, fmt.Sprintf("  └─%s<┘  ", railpad("", maxl)))
	return out
}

func joinLists(tracks [][]string) []string {
	var out []string
	for _, t := range tracks {
		out = append(out, t...)
	}
	return out
}

func walkRule(r *Rule) []string {
	leftrec := ""
	if r.IsLrec {
		leftrec = "⟳"
	} else if r.NoMemo {
		leftrec = "⊬"
	}
	out := []string{fmt.Sprintf("%s%s ●─", r.Name, leftrec)}
	out = weld(out, walkExp(r.Exp))
	out = weld(out, []string{"─■"})

	len0 := ulen(out[0])
	padding := strings.Repeat(" ", len0)
	for i, s := range out {
		out[i] = s + padding
	}
	out = append(out, strings.Repeat(" ", ulen(out[0])))
	return out
}

func walkGrammar(g *Grammar) []string {
	var tracks [][]string
	for _, rule := range g.Rules {
		tracks = append(tracks, walkRule(rule))
	}
	return joinLists(tracks)
}

func walkExp(m Model) []string {
	switch exp := m.(type) {
	case *EmptyClosure:
		return []string{"[]"}
	case *NULL:
		return []string{""}
	case *Cut:
		return []string{" ✂ "}
	case *Void:
		return []string{" ∅ "}
	case *Fail:
		return []string{" ⚠ "}
	case *Dot:
		return []string{" ∀ "}
	case *EOF:
		return []string{fmt.Sprintf("⇥%s ", etx)}
	case *EOL:
		return []string{fmt.Sprintf("⇥%s ", eol)}

	case *Token:
		return []string{fmt.Sprintf("%q", exp.Token)}
	case *Pattern:
		pat := strings.TrimPrefix(exp.Pattern, "r'")
		pat = strings.TrimSuffix(pat, "'")
		return []string{fmt.Sprintf("/%s/─", pat)}
	case *Constant:
		return []string{fmt.Sprintf("`%s`", exp.Literal)}
	case *Alert:
		return []string{fmt.Sprintf("%s^`%s`", strings.Repeat("^", exp.Level), exp.Literal)}

	case *Call:
		return []string{exp.Name}
	case *RuleInclude:
		return []string{fmt.Sprintf(" >(%s) ", exp.Name)}

	case *Optional:
		inner := walkExp(exp.Exp)
		withArrow := weld([]string{"→"}, inner)
		return layOut([][]string{withArrow, {"→"}})

	case *Closure:
		return loopRails(walkExp(exp.Exp))
	case *PositiveClosure:
		return stopnloop(walkExp(exp.Exp))

	case *Join:
		sep := weld(walkExp(exp.Sep), []string{" ✂ ─"})
		out := weld(sep, walkExp(exp.Exp))
		return loopRails(out)
	case *PositiveJoin:
		sep := weld(walkExp(exp.Sep), []string{" ✂ ─"})
		out := weld(sep, walkExp(exp.Exp))
		return stopnloop(out)
	case *Gather:
		sep := weld(walkExp(exp.Sep), []string{" │ "})
		out := weld(sep, walkExp(exp.Exp))
		return loopRails(out)
	case *PositiveGather:
		sep := weld(walkExp(exp.Sep), []string{" │ "})
		out := weld(sep, walkExp(exp.Exp))
		return stopnloop(out)

	case *SkipTo:
		return weld([]string{" ->("}, weld(walkExp(exp.Exp), []string{")"}))

	case *Sequence:
		if len(exp.Sequence) == 0 {
			return []string{""}
		}
		result := walkExp(exp.Sequence[0])
		for _, item := range exp.Sequence[1:] {
			result = weld(result, walkExp(item))
		}
		return result

	case *Choice:
		var tracks [][]string
		for _, opt := range exp.Options {
			tracks = append(tracks, walkExp(opt.Exp))
		}
		return layOut(tracks)

	case *Option:
		return walkExp(exp.Exp)

	case *Named:
		return weld(
			[]string{fmt.Sprintf(" %s=(", exp.Name)},
			weld(walkExp(exp.Exp), []string{") "}),
		)
	case *NamedList:
		return weld(
			[]string{fmt.Sprintf(" %s+=(", exp.Name)},
			weld(walkExp(exp.Exp), []string{") "}),
		)

	case *Group:
		return walkExp(exp.Exp)
	case *SkipGroup:
		return walkExp(exp.Exp)

	case *Lookahead:
		return weld(
			[]string{"─ &["},
			weld(walkExp(exp.Exp), []string{"]"}),
		)
	case *NegativeLookahead:
		return weld(
			[]string{"─ !["},
			weld(walkExp(exp.Exp), []string{"]"}),
		)

	case *Override:
		return weld(
			[]string{" =("},
			weld(walkExp(exp.Exp), []string{") "}),
		)
	case *OverrideList:
		return weld(
			[]string{" +=("},
			weld(walkExp(exp.Exp), []string{") "}),
		)

	case *Synth:
		return walkExp(exp.Exp)

	default:
		return []string{""}
	}
}

// Railroads implements the ToRailroad interface for Grammar.
func (m *Grammar) Railroads() string {
	return rails(walkGrammar(m)).String()
}

// Railroads implements the ToRailroad interface for Rule.
func (m *Rule) Railroads() string {
	return rails(walkRule(m)).String()
}

// Railroads implements the ToRailroad interface for NULL.
func (m *NULL) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Token.
func (m *Token) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Pattern.
func (m *Pattern) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Constant.
func (m *Constant) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Alert.
func (m *Alert) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Call.
func (m *Call) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for RuleInclude.
func (m *RuleInclude) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Cut.
func (m *Cut) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Dot.
func (m *Dot) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for EOF.
func (m *EOF) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for EOL.
func (m *EOL) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Fail.
func (m *Fail) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Void.
func (m *Void) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for EmptyClosure.
func (m *EmptyClosure) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Option.
func (m *Option) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Group.
func (m *Group) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for SkipGroup.
func (m *SkipGroup) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Lookahead.
func (m *Lookahead) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for NegativeLookahead.
func (m *NegativeLookahead) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for SkipTo.
func (m *SkipTo) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Optional.
func (m *Optional) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Closure.
func (m *Closure) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for PositiveClosure.
func (m *PositiveClosure) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Override.
func (m *Override) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for OverrideList.
func (m *OverrideList) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Synth.
func (m *Synth) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Named.
func (m *Named) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for NamedList.
func (m *NamedList) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Join.
func (m *Join) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for PositiveJoin.
func (m *PositiveJoin) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Gather.
func (m *Gather) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for PositiveGather.
func (m *PositiveGather) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Sequence.
func (m *Sequence) Railroads() string { return rails(walkExp(m)).String() }

// Railroads implements the ToRailroad interface for Choice.
func (m *Choice) Railroads() string { return rails(walkExp(m)).String() }
