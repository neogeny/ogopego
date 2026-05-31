package test

import (
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/neogeny/ogopego/api"
)

func TestMissingRule(t *testing.T) {
	_, err := api.Compile(`
		@@grammar :: TestGrammar
		block = test
	`, nil)
	assert.Error(t, err, "expected error for missing rule 'test'")
}
