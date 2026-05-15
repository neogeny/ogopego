package ogopego_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/neogeny/ogopego/test"
)

func TestPrivateGrammars(t *testing.T) {
	if os.Getenv("XONSH_VERSION") == "" {
		t.Skip("XONSH_VERSION not set — local test only")
	}
	dir := "../grammar"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Skipf("cannot read %s: %v", dir, err)
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".ebnf" {
			continue
		}
		t.Run(e.Name(), func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(dir, e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			g := ogopego.Compile(t, string(data), nil)
			if len(g.Rules) == 0 {
				t.Fatal("expected at least one rule")
			}
		})
	}
}
