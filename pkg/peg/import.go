// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package peg

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/neogeny/ogopego/pkg/asjson"
)

// helper is a utility struct for parsing JSON.
type helper struct {
	value map[string]any
	path  []string
}

// newHelper creates a new helper instance.
func newHelper(value map[string]any) *helper {
	return &helper{value: value, path: []string{}}
}

// push adds a label to the helper's path.
func (h *helper) push(label string) *helper {
	p := make([]string, len(h.path)+1)
	copy(p, h.path)
	p[len(h.path)] = label
	return &helper{value: h.value, path: p}
}

// withValue sets the helper's value to a new map.
func (h *helper) withValue(v any) *helper {
	obj, ok := v.(map[string]any)
	if !ok {
		return h
	}
	return &helper{value: obj, path: h.path}
}

// error creates a JsonError with the helper's path.
func (h *helper) error(msg string) *asjson.JsonError {
	s := strings.Join(h.path, " -> ")
	if s == "" {
		return asjson.NewJsonError(msg)
	}
	return asjson.NewJsonError(fmt.Sprintf("%s at %s", msg, s))
}

// getClass retrieves the "__class__" field from the helper's value.
func (h *helper) getClass() (string, error) {
	raw, ok := h.value["__class__"]
	if !ok {
		return "", h.error("Missing __class__")
	}
	s, ok := raw.(string)
	if !ok {
		return "", h.error("__class__ is not a string")
	}
	return s, nil
}

// getString retrieves a string field from the helper's value.
func (h *helper) getString(field string) (string, error) {
	raw, ok := h.value[field]
	if !ok {
		return "", h.error(fmt.Sprintf("Missing field: %s", field))
	}
	s, ok := raw.(string)
	if !ok {
		return "", h.error(fmt.Sprintf("Field %s is not a string", field))
	}
	return s, nil
}

// optString retrieves an optional string field from the helper's value.
func (h *helper) optString(field string) string {
	raw, ok := h.value[field]
	if !ok {
		return ""
	}
	if s, ok := raw.(string); ok {
		return s
	}
	return ""
}

// optBool retrieves an optional boolean field from the helper's value.
func (h *helper) optBool(field string, def bool) bool {
	raw, ok := h.value[field]
	if !ok {
		return def
	}
	if b, ok := raw.(bool); ok {
		return b
	}
	return def
}

// optFloat retrieves an optional float64 field from the helper's value.
func (h *helper) optFloat(field string) (float64, bool) {
	raw, ok := h.value[field]
	if !ok {
		return 0, false
	}
	if n, ok := raw.(float64); ok {
		return n, true
	}
	return 0, false
}

// getNested retrieves a nested object as a new helper instance.
func (h *helper) getNested(field string) (*helper, error) {
	raw, ok := h.value[field]
	if !ok {
		return nil, h.error(fmt.Sprintf("Missing field: %s", field))
	}
	obj, ok := raw.(map[string]any)
	if !ok {
		return nil, h.error(fmt.Sprintf("Field %s is not an object", field))
	}
	label := field
	if cls, ok := obj["__class__"].(string); ok {
		label = fmt.Sprintf("%s:%s", field, cls)
	}
	return h.push(label).withValue(obj), nil
}

// getArray retrieves an array of objects as a slice of helper instances.
func (h *helper) getArray(field string) ([]*helper, error) {
	raw, ok := h.value[field]
	if !ok {
		return nil, h.error(fmt.Sprintf("Missing field: %s", field))
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, h.error(fmt.Sprintf("Field %s is not an array", field))
	}
	var result []*helper
	for i, v := range arr {
		obj, ok := v.(map[string]any)
		if !ok {
			continue
		}
		label := fmt.Sprintf("%s[%d]", field, i)
		if cls, ok := obj["__class__"].(string); ok {
			label = fmt.Sprintf("%s[%d]:%s", field, i, cls)
		}
		result = append(result, h.push(label).withValue(obj))
	}
	return result, nil
}

// LoadGrammarFromJSON parses grammar data from JSON bytes.
func LoadGrammarFromJSON(data []byte) (*Grammar, error) {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, asjson.NewJsonError(fmt.Sprintf("JSON parse error: %v", err))
	}
	return grammarFromJSON(newHelper(raw.(map[string]any)))
}

// grammarFromJSON converts a JSON helper to a Grammar object.
func grammarFromJSON(h *helper) (*Grammar, error) {
	cls, err := h.getClass()
	if err != nil {
		return nil, err
	}
	if cls != "Grammar" {
		return nil, h.error("Expected Grammar root")
	}

	name, err := h.getString("name")
	if err != nil {
		return nil, err
	}

	ruleHelpers, err := h.getArray("rules")
	if err != nil {
		return nil, err
	}

	var rules []*Rule
	for i, rh := range ruleHelpers {
		r, err := ruleFromJSON(rh)
		if err != nil {
			return nil, asjson.NewJsonError(fmt.Sprintf("rules[%d]: %v", i, err))
		}
		rules = append(rules, r)
	}

	directives := parseDirectives(h.value)

	var keywords []string
	if kwRaw, ok := h.value["keywords"]; ok {
		if kwArr, ok := kwRaw.([]any); ok {
			for _, v := range kwArr {
				if s, ok := v.(string); ok {
					keywords = append(keywords, s)
				}
			}
		}
	}

	g := &Grammar{
		Name:       name,
		Directives: directives,
		Keywords:   keywords,
		Rules:      rules,
	}
	return g, nil
}

// parseDirectives parses directives from a map.
func parseDirectives(obj map[string]any) [][]string {
	raw, ok := obj["directives"]
	if !ok {
		return nil
	}
	dirObj, ok := raw.(map[string]any)
	if !ok {
		return nil
	}
	var result [][]string
	for k, v := range dirObj {
		switch val := v.(type) {
		case string:
			result = append(result, []string{k, val})
		case bool:
			if val {
				result = append(result, []string{k, "true"})
			} else {
				result = append(result, []string{k, "false"})
			}
		case float64:
			result = append(result, []string{k, fmt.Sprintf("%v", val)})
		default:
			result = append(result, []string{k, fmt.Sprintf("%v", val)})
		}
	}
	return result
}

// ruleFromJSON converts a JSON helper to a Rule object.
func ruleFromJSON(h *helper) (*Rule, error) {
	cls, err := h.getClass()
	if err != nil {
		return nil, err
	}
	if cls != "Rule" {
		return nil, h.error("Expected Rule")
	}

	name, err := h.getString("name")
	if err != nil {
		return nil, err
	}

	exp, err := modelFromJSON(h.getNested("exp"))
	if err != nil {
		return nil, err
	}

	var params []string
	if pRaw, ok := h.value["params"]; ok {
		if pArr, ok := pRaw.([]any); ok {
			for _, v := range pArr {
				if s, ok := v.(string); ok {
					params = append(params, s)
				}
			}
		}
	}

	var kwparams = make(map[string]string)
	if pRaw, ok := h.value["kwparams"]; ok {
		if pMap, ok := pRaw.(map[any]any); ok {
			for k, v := range pMap {
				if ks, ok := k.(string); ok {
					if vs, ok := v.(string); ok {
						kwparams[ks] = vs
					}
				}
			}
		}
	}

	var decorators []string
	if pRaw, ok := h.value["decorators"]; ok {
		if pArr, ok := pRaw.([]any); ok {
			for _, v := range pArr {
				if s, ok := v.(string); ok {
					decorators = append(decorators, s)
				}
			}
		}
	}

	noMemo := h.optBool("no_memo", false) ||
		slices.Contains(decorators, "nomemo")
	noStak := h.optBool("no_stak", false) ||
		slices.Contains(decorators, "nostak")
	isName := h.optBool("is_name", false) ||
		slices.Contains(decorators, "name") ||
		slices.Contains(decorators, "isname")
	isTokn := h.optBool("is_tokn", false)
	isMemo := h.optBool("is_memo", true)
	isLrec := h.optBool("is_lrec", false)

	r := &Rule{
		Exp:        exp,
		Name:       name,
		Params:     params,
		KWParams:   kwparams,
		Decorators: decorators,
		IsName:     isName,
		IsTokn:     isTokn,
		NoMemo:     noMemo,
		NoStak:     noStak,
		IsMemo:     isMemo,
		IsLrec:     isLrec,
	}
	return r, nil
}

// modelFromJSON converts a JSON helper to a Model object.
func modelFromJSON(h *helper, err error) (Model, error) {
	if err != nil {
		return nil, err
	}
	cls, err := h.getClass()
	if err != nil {
		return nil, err
	}

	switch cls {
	case "Sequence":
		items, err := h.getArray("sequence")
		if err != nil {
			return nil, err
		}
		var seq []Model
		for _, ih := range items {
			exp, err := modelFromJSON(ih, nil)
			if err != nil {
				return nil, err
			}
			seq = append(seq, exp)
		}
		return &Sequence{Sequence: seq}, nil

	case "Choice":
		items, err := h.getArray("options")
		if err != nil {
			return nil, err
		}
		var opts []Model
		for _, ih := range items {
			cls, _ := ih.getClass()
			var exp Model
			if cls == "Option" {
				exp, err = modelFromJSON(ih.getNested("exp"))
			} else {
				exp, err = modelFromJSON(ih, nil)
			}
			if err != nil {
				return nil, err
			}
			opts = append(opts, &Option{Exp: exp})
		}
		return &Choice{Options: opts}, nil

	case "Option":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &Option{Exp: exp}, nil

	case "Named":
		name, err := h.getString("name")
		if err != nil {
			return nil, err
		}
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &Named{Exp: exp, Name: name}, nil

	case "NamedList":
		name, err := h.getString("name")
		if err != nil {
			return nil, err
		}
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &NamedList{Exp: exp, Name: name}, nil

	case "Call":
		name, err := h.getString("name")
		if err != nil {
			return nil, err
		}
		return &Call{Name: name}, nil

	case "Token":
		tok, err := h.getString("token")
		if err != nil {
			return nil, err
		}
		return &Token{Token: tok}, nil

	case "Pattern":
		pat, err := h.getString("pattern")
		if err != nil {
			return nil, err
		}
		return &Pattern{Pattern: pat}, nil

	case "Constant":
		lit := h.optString("literal")
		return &Constant{Literal: lit}, nil

	case "Alert":
		lit := h.optString("literal")
		level := 0
		if n, ok := h.optFloat("level"); ok {
			level = int(n)
		}
		return &Alert{Constant: Constant{Literal: lit}, Level: level}, nil

	case "Group":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &Group{Exp: exp}, nil

	case "Optional":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &Optional{Exp: exp}, nil

	case "Closure":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &Closure{Exp: exp}, nil

	case "PositiveClosure":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &PositiveClosure{Exp: exp}, nil

	case "Lookahead":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &Lookahead{Exp: exp}, nil

	case "NegativeLookahead":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &NegativeLookahead{Exp: exp}, nil

	case "SkipGroup":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &SkipGroup{Exp: exp}, nil

	case "SkipTo":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &SkipTo{Exp: exp}, nil

	case "Override":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &Override{Exp: exp}, nil

	case "OverrideList":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &OverrideList{Exp: exp}, nil

	case "Join":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		sep, err := modelFromJSON(h.getNested("sep"))
		if err != nil {
			return nil, err
		}
		return &Join{Exp: exp, Sep: sep}, nil

	case "PositiveJoin":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		sep, err := modelFromJSON(h.getNested("sep"))
		if err != nil {
			return nil, err
		}
		return &PositiveJoin{Exp: exp, Sep: sep}, nil

	case "Gather":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		sep, err := modelFromJSON(h.getNested("sep"))
		if err != nil {
			return nil, err
		}
		return &Gather{Exp: exp, Sep: sep}, nil

	case "PositiveGather":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		sep, err := modelFromJSON(h.getNested("sep"))
		if err != nil {
			return nil, err
		}
		return &PositiveGather{Exp: exp, Sep: sep}, nil

	case "Void":
		return &Void{}, nil

	case "Null":
		return &NULL{}, nil

	case "Fail":
		return &Fail{}, nil

	case "Dot":
		return &Dot{}, nil

	case "Synth":
		exp, err := modelFromJSON(h.getNested("exp"))
		if err != nil {
			return nil, err
		}
		return &Synth{Exp: exp}, nil

	case "Cut":
		return &Cut{}, nil

	case "EOF":
		return &EOF{}, nil

	case "EOL":
		return &EOL{}, nil

	case "EmptyClosure":
		return &EmptyClosure{}, nil

	case "NameMeta":
		return &NameMeta{}, nil

	case "IntMeta":
		return &IntMeta{}, nil

	case "UIntMeta":
		return &UIntMeta{}, nil

	case "FloatMeta":
		return &FloatMeta{}, nil

	case "BoolMeta":
		return &BoolMeta{}, nil

	case "RuleInclude":
		name, err := h.getString("name")
		if err != nil {
			return nil, err
		}
		exp, _ := modelFromJSON(h.getNested("exp"))
		return &RuleInclude{Name: name, exp: exp}, nil

	default:
		return nil, h.error(fmt.Sprintf("Unsupported: %s", cls))
	}
}
