package peg

import (
	"encoding/json"
	"fmt"
	"testing"

	asjson "github.com/neogeny/ogopego/json"
)

func TestRepro(t *testing.T) {
	n := &Node{
		Ast: map[string]any{"key": "value", "num": float64(42)},
	}
	result := asjson.AsJSON(n)
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("JSON: %s\n", b)

	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	for k, v := range out {
		fmt.Printf("  key=%q type=%T value=%v\n", k, v, v)
	}
}
