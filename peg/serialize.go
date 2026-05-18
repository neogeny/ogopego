package peg

import (
	"encoding/json"

	"github.com/iancoleman/orderedmap"
)

func ModelToJSONStr(v Model) string {
	out := ModelToJSON(v)
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b)
}

// ModelToJSON returns an any that can be serialized to JSON,
// matching the format expected by modelFromJSON in import.go.
func ModelToJSON(v Model) any {
	switch m := v.(type) {
	case nil:
		return nil

	case *Grammar:
		return serializeGrammar(m)
	case *Token:
		return mapClass("Token", "token", m.Token)
	case *Pattern:
		return mapClass("Pattern", "pattern", m.Pattern)
	case *Constant:
		return mapClass("Constant", "literal", m.Literal)
	case *Alert:
		return mapClass("Alert", "literal", m.Literal, "level", m.Level)
	case *Dot:
		return mapClass("Dot")
	case *EOF:
		return mapClass("EOF")
	case *EOL:
		return mapClass("EOL")
	case *Fail:
		return mapClass("Fail")
	case *Void:
		return mapClass("Void")
	case *NULL:
		return mapClass("Null")
	case *Cut:
		return mapClass("Cut")
	case *Synth:
		return mapClass("Synth", "exp", ModelToJSON(m.Exp))
	case *EmptyClosure:
		return mapClass("EmptyClosure")
	case *SkipTo:
		return mapClass("SkipTo", "exp", ModelToJSON(m.Exp))

	case *Call:
		return mapClass("Call", "name", m.Name)
	case *RuleInclude:
		return mapClass("RuleInclude", "name", m.Name, "exp", ModelToJSON(m.exp))

	case *Group:
		return mapClass("Group", "exp", ModelToJSON(m.Exp))
	case *SkipGroup:
		return mapClass("SkipGroup", "exp", ModelToJSON(m.Exp))
	case *Lookahead:
		return mapClass("Lookahead", "exp", ModelToJSON(m.Exp))
	case *NegativeLookahead:
		return mapClass("NegativeLookahead", "exp", ModelToJSON(m.Exp))
	case *Override:
		return mapClass("Override", "exp", ModelToJSON(m.Exp))
	case *OverrideList:
		return mapClass("OverrideList", "exp", ModelToJSON(m.Exp))
	case *Option:
		return mapClass("Option", "exp", ModelToJSON(m.Exp))
	case *Optional:
		return mapClass("Optional", "exp", ModelToJSON(m.Exp))
	case *Closure:
		return mapClass("Closure", "exp", ModelToJSON(m.Exp))
	case *PositiveClosure:
		return mapClass("PositiveClosure", "exp", ModelToJSON(m.Exp))

	case *Named:
		return mapClass("Named", "name", m.Name, "exp", ModelToJSON(m.Exp))
	case *NamedList:
		return mapClass("NamedList", "name", m.Name, "exp", ModelToJSON(m.Exp))

	case *Join:
		return mapClass("Join", "exp", ModelToJSON(m.Exp), "sep", ModelToJSON(m.Sep))
	case *PositiveJoin:
		return mapClass("PositiveJoin", "exp", ModelToJSON(m.Exp), "sep", ModelToJSON(m.Sep))
	case *Gather:
		return mapClass("Gather", "exp", ModelToJSON(m.Exp), "sep", ModelToJSON(m.Sep))
	case *PositiveGather:
		return mapClass("PositiveGather", "exp", ModelToJSON(m.Exp), "sep", ModelToJSON(m.Sep))

	case *Sequence:
		items := make([]any, len(m.Sequence))
		for i, item := range m.Sequence {
			items[i] = ModelToJSON(item)
		}
		return mapClass("Sequence", "sequence", items)
	case *Choice:
		opts := make([]any, len(m.Options))
		for i, opt := range m.Options {
			opts[i] = ModelToJSON(opt)
		}
		return mapClass("Choice", "options", opts)

	default:
		panic("unknown model type")
	}
}

func mapClass(class string, kv ...any) *orderedmap.OrderedMap {
	out := orderedmap.New()
	out.Set("__class__", class)
	for i := 0; i < len(kv); i += 2 {
		k := kv[i].(string)
		v := kv[i+1]
		out.Set(k, v)
	}
	return out
}

// serializeGrammar returns a JSON string representation of g
// that can be read back by ParseGrammar.
func serializeGrammar(g *Grammar) any {
	out := orderedmap.New()
	out.Set("__class__", "Grammar")
	out.Set("name", g.Name)

	dirs := orderedmap.New()
	if g.Directives != nil {
		for _, k := range g.Directives.Keys() {
			v, _ := g.Directives.Get(k)
			dirs.Set(k, v)
		}
	}
	out.Set("directives", dirs)

	kw := make([]string, len(g.Keywords))
	copy(kw, g.Keywords)
	out.Set("keywords", kw)

	rules := make([]any, len(g.Rules))
	for i, rule := range g.Rules {
		rules[i] = serializeRule(rule)
	}
	out.Set("rules", rules)
	return out
}

func serializeRule(r *Rule) *orderedmap.OrderedMap {
	out := orderedmap.New()
	out.Set("__class__", "Rule")
	out.Set("name", r.Name)
	if r.Params == nil {
		out.Set("params", []string{})
	} else {
		out.Set("params", r.Params)
	}
	out.Set("is_name", r.IsName)
	out.Set("is_tokn", r.IsTokn)
	out.Set("no_memo", r.NoMemo)
	out.Set("no_stak", r.NoStak)
	out.Set("is_memo", r.IsMemo)
	out.Set("is_lrec", r.IsLrec)
	out.Set("exp", ModelToJSON(r.Exp))
	return out
}
