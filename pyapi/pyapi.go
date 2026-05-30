package pyapi

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/asjson"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/peg"
)

type GrammarHandle int64

var (
	mu      sync.RWMutex
	handles               = make(map[GrammarHandle]*peg.Grammar)
	nextID  GrammarHandle = 1
)

func decodeConfig(cfgJSON string) *config.Cfg {
	if cfgJSON == "" {
		return nil
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(cfgJSON), &raw); err != nil {
		return nil
	}
	cfg := config.DefaultCfg()
	for k, v := range raw {
		switch k {
		case "trace":
			if b, ok := v.(bool); ok {
				cfg.Trace = b
			}
		case "ignorecase":
			if b, ok := v.(bool); ok {
				cfg.IgnoreCase = b
			}
		case "left_recursion":
			if b, ok := v.(bool); ok {
				cfg.NoLeftRecursion = !b
			}
		case "parseinfo":
			if b, ok := v.(bool); ok {
				cfg.ParseInfo = b
			}
		case "memoization":
			if b, ok := v.(bool); ok {
				cfg.NoMemo = !b
			}
		case "prune_memos_on_cut":
			if b, ok := v.(bool); ok {
				cfg.NoPruneMemosOnCut = !b
			}
		case "colorize", "color":
			if b, ok := v.(bool); ok {
				cfg.Colorize = b
			}
		case "namechars":
			if s, ok := v.(string); ok {
				cfg.NameChars = s
			}
		case "nameguard":
			if b, ok := v.(bool); ok {
				cfg.NameGuard = &b
			}
		case "whitespace":
			if s, ok := v.(string); ok {
				cfg.Whitespace = &s
			}
		case "comments":
			if s, ok := v.(string); ok {
				cfg.Comments = s
			}
		case "eol_comments":
			if s, ok := v.(string); ok {
				cfg.EolComments = s
			}
		case "keywords":
			switch val := v.(type) {
			case []any:
				for _, item := range val {
					if s, ok := item.(string); ok {
						cfg.Keywords = append(cfg.Keywords, s)
					}
				}
			case []string:
				cfg.Keywords = val
			}
		case "name":
			if s, ok := v.(string); ok {
				cfg.Name = s
			}
		case "source":
			if s, ok := v.(string); ok {
				cfg.Source = s
			}
		case "start":
			if s, ok := v.(string); ok {
				cfg.Start = s
			}
		case "grammar":
			if s, ok := v.(string); ok {
				cfg.Grammar = s
			}
		case "perlinememos":
			if f, ok := v.(float64); ok {
				cfg.PerLineMemos = f
			}
		}
	}
	return cfg
}

func storeGrammar(g *peg.Grammar) GrammarHandle {
	mu.Lock()
	defer mu.Unlock()
	h := nextID
	nextID++
	handles[h] = g
	return h
}

func lookupGrammar(h GrammarHandle) (*peg.Grammar, error) {
	mu.RLock()
	defer mu.RUnlock()
	g, ok := handles[h]
	if !ok {
		return nil, fmt.Errorf("unknown grammar handle: %d", h)
	}
	return g, nil
}

func PyCompile(grammar string, cfgJSON string) (GrammarHandle, error) {
	cfg := decodeConfig(cfgJSON)
	g, err := api.Compile(grammar, cfg)
	if err != nil {
		return 0, err
	}
	return storeGrammar(g), nil
}

func PyParseToJSONString(grammar string, text string, cfgJSON string) (string, error) {
	cfg := decodeConfig(cfgJSON)
	g, err := api.Compile(grammar, cfg)
	if err != nil {
		return "", err
	}
	return api.ParseInputToJSONString(g, text, cfg)
}

func PyParseInputToJSONString(handle GrammarHandle, text string, cfgJSON string) (string, error) {
	g, err := lookupGrammar(handle)
	if err != nil {
		return "", err
	}
	return api.ParseInputToJSONString(g, text, decodeConfig(cfgJSON))
}

func toJSONStr(v any) (string, error) {
	bts, err := json.MarshalIndent(asjson.AsJSON(v), "", "  ")
	if err != nil {
		return "", err
	}
	return string(bts), nil
}

func PyGrammarToJSONString(handle GrammarHandle) (string, error) {
	g, err := lookupGrammar(handle)
	if err != nil {
		return "", err
	}
	return toJSONStr(g)
}

func PyFreeGrammar(handle GrammarHandle) {
	mu.Lock()
	defer mu.Unlock()
	delete(handles, handle)
}

func PyBootGrammarToJSONString(cfgJSON string) (string, error) {
	g, err := api.BootGrammar()
	if err != nil {
		return "", err
	}
	_ = cfgJSON
	return toJSONStr(g)
}

func PyParseGrammarToJSONString(grammar string, cfgJSON string) (string, error) {
	val, err := api.ParseGrammarToJSON(grammar, decodeConfig(cfgJSON))
	if err != nil {
		return "", err
	}
	return toJSONStr(val)
}
