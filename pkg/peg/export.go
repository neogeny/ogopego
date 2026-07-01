package peg

import (
	"bytes"
	"encoding/json"

	ctn "github.com/neogeny/ogopego/pkg/util/container"
)

func ModelToJSONStr(v Model) string {
	out := ModelToJSON(v)

	// Temporarily print the concrete types inside your map/struct
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false) // Prevents escaping Unicode/HTML

	_ = enc.Encode(out)

	// Trim the trailing newline so it behaves exactly like MarshalIndent
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

func mapClass(class string, kv ...any) *ctn.BoundedMap[string, any] {
	out := ctn.NewBoundedMap[string, any](0)
	_ = out.Set("__class__", class)
	for i := 0; i < len(kv); i += 2 {
		k := kv[i].(string)
		v := kv[i+1]
		_ = out.Set(k, v)
	}
	return &out
}

// serializeGrammar returns a JSON string representation of g
// that can be read back by LoadGrammarFromJSON.
func serializeGrammar(g *Grammar) any {
	out := ctn.NewBoundedMap[string, any](0)
	_ = out.Set("__class__", "Grammar")
	_ = out.Set("name", g.Name)

	dirs := ctn.NewBoundedMap[string, any](0)
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
		_ = dirs.Set(d[0], value)
	}
	_ = out.Set("directives", &dirs)

	kw := make([]string, len(g.Keywords))
	copy(kw, g.Keywords)
	_ = out.Set("keywords", kw)

	rules := make([]any, len(g.Rules))
	for i, rule := range g.Rules {
		rules[i] = serializeRule(rule)
	}
	_ = out.Set("rules", rules)
	return &out
}

func serializeRule(r *Rule) *ctn.BoundedMap[string, any] {
	r.normalize()
	out := ctn.NewBoundedMap[string, any](0)
	_ = out.Set("__class__", "Rule")
	_ = out.Set("name", r.Name)
	// NOTE This is the field order used by TatSu @ 2026-05-27
	_ = out.Set("exp", ModelToJSON(r.Exp))
	if r.Params == nil {
		_ = out.Set("params", []string{})
	} else {
		_ = out.Set("params", r.Params)
	}
	if r.KWParams == nil {
		_ = out.Set("kwparams", map[string]any{})
	} else {
		_ = out.Set("kwparams", r.KWParams)
	}
	_ = out.Set("decorators", r.Decorators)
	_ = out.Set("base", nil)
	_ = out.Set("is_name", r.IsName)
	_ = out.Set("is_tokn", r.IsTokn)
	_ = out.Set("no_memo", r.NoMemo)
	_ = out.Set("no_stak", r.NoStak)
	_ = out.Set("is_memo", r.IsMemo)
	_ = out.Set("is_lrec", r.IsLrec)
	return &out
}
