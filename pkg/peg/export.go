package peg

import (
	"bytes"
	"encoding/json"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func ModelToJSONStr(v Model) string {
	out := ModelToJSON(v)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	_ = enc.Encode(out)

	return string(bytes.TrimRight(buf.Bytes(), "\n"))
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
		return mapClass("Void", "ast", "()")
	case *NULL:
		return mapClass("Null")
	case *Cut:
		return mapClass("Cut")
	case *Synth:
		return mapClass("Synth", "exp", ModelToJSON(m.Exp))
	case *NameMeta:
		return mapClass("NameMeta")
	case *IntMeta:
		return mapClass("IntMeta")
	case *UIntMeta:
		return mapClass("UIntMeta")
	case *FloatMeta:
		return mapClass("FloatMeta")
	case *BoolMeta:
		return mapClass("BoolMeta")

	case *EmptyClosure:
		return mapClass("EmptyClosure", "ast", []any{})
	case *SkipTo:
		return mapClass("SkipTo", "exp", ModelToJSON(m.Exp))
	case *Call:
		return mapClass("Call", "name", m.Name)
	case *RuleInclude:
		return mapClass("RuleInclude", "name", m.Name)
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

func mapClass(class string, kv ...any) *orderedmap.OrderedMap[string, any] {
	out := orderedmap.New[string, any]()
	out.Set("__class__", class)
	for i := 0; i < len(kv); i += 2 {
		k := kv[i].(string)
		v := kv[i+1]
		out.Set(k, v)
	}
	return out
}

func serializeGrammar(g *Grammar) any {
	out := orderedmap.New[string, any]()
	out.Set("__class__", "Grammar")
	out.Set("name", g.Name)

	dirs := orderedmap.New[string, any]()
	for _, d := range g.Directives {
		v := d[1]
		var value any = v
		switch d[1] {
		case "true", "True":
			value = true
		case "false", "False":
			value = false
		case "null", "None":
			value = nil
		}
		dirs.Set(d[0], value)
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

func serializeRule(r *Rule) *orderedmap.OrderedMap[string, any] {
	r.normalize()
	out := orderedmap.New[string, any]()
	out.Set("__class__", "Rule")
	out.Set("name", r.Name)
	out.Set("exp", ModelToJSON(r.Exp))
	if r.Params == nil {
		out.Set("params", []string{})
	} else {
		out.Set("params", r.Params)
	}
	if r.KWParams == nil {
		out.Set("kwparams", map[string]any{})
	} else {
		out.Set("kwparams", r.KWParams)
	}
	out.Set("decorators", r.Decorators)
	out.Set("base", nil)
	out.Set("is_name", r.IsName)
	out.Set("is_tokn", r.IsTokn)
	out.Set("no_memo", r.NoMemo)
	out.Set("no_stak", r.NoStak)
	out.Set("is_memo", r.IsMemo)
	out.Set("is_lrec", r.IsLrec)
	return out
}
