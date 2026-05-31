package api

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/neogeny/ogopego/pkg/asjson"
	"github.com/neogeny/ogopego/pkg/config"
	"github.com/neogeny/ogopego/pkg/peg"
)

type GrammarHandle int64

var (
	pyMu      sync.RWMutex
	pyHandles               = make(map[GrammarHandle]*peg.Grammar)
	pyNextID  GrammarHandle = 1
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
	pyMu.Lock()
	defer pyMu.Unlock()
	h := pyNextID
	pyNextID++
	pyHandles[h] = g
	return h
}

func lookupGrammar(h GrammarHandle) (*peg.Grammar, error) {
	pyMu.RLock()
	defer pyMu.RUnlock()
	g, ok := pyHandles[h]
	if !ok {
		return nil, fmt.Errorf("unknown grammar handle: %d", h)
	}
	return g, nil
}

func PyCompile(grammar string, cfgJSON string) (GrammarHandle, error) {
	cfg := decodeConfig(cfgJSON)
	g, err := Compile(grammar, cfg)
	if err != nil {
		return 0, err
	}
	return storeGrammar(g), nil
}

func PyLoadGrammar(path string, cfgJSON string) (GrammarHandle, error) {
	cfg := decodeConfig(cfgJSON)
	g, err := LoadGrammar(path, cfg)
	if err != nil {
		return 0, err
	}
	return storeGrammar(g), nil
}

func PyParseToJSONString(grammar string, text string, cfgJSON string) (string, error) {
	cfg := decodeConfig(cfgJSON)
	g, err := Compile(grammar, cfg)
	if err != nil {
		return "", err
	}
	return ParseInputToJSONString(g, text, cfg)
}

func PyParseInputToJSONString(handle GrammarHandle, text string, cfgJSON string) (string, error) {
	g, err := lookupGrammar(handle)
	if err != nil {
		return "", err
	}
	return ParseInputToJSONString(g, text, decodeConfig(cfgJSON))
}

func PyGrammarToJSONString(handle GrammarHandle) (string, error) {
	g, err := lookupGrammar(handle)
	if err != nil {
		return "", err
	}
	return asjson.AsJSONStr(g), nil
}

func PyFreeGrammar(handle GrammarHandle) {
	pyMu.Lock()
	defer pyMu.Unlock()
	delete(pyHandles, handle)
}

func PyBootGrammarToJSONString(cfgJSON string) (string, error) {
	g, err := BootGrammar()
	if err != nil {
		return "", err
	}
	_ = cfgJSON
	out, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func PyParseGrammarToJSONString(grammar string, cfgJSON string) (string, error) {
	return ParseGrammarToJSONStr(grammar, decodeConfig(cfgJSON))
}
