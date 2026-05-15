package test

import (
	"testing"

	"github.com/neogeny/ogopego/api"
)

func TestMissingRule(t *testing.T) {
	_, err := api.Compile(`
		@@grammar :: TestGrammar
		block = test
	`, nil)
	if err == nil {
		t.Fatal("expected error for missing rule 'test'")
	}
}
