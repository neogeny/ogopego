package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
)

var ogoBin string

func ogoPath(t *testing.T) string {
	if ogoBin == "" {
		root := filepath.Join("..")
		bin := filepath.Join(root, "target/debug", "ogo")
		if _, err := os.Stat(bin); err != nil {
			t.Skipf("ogo binary not found at %s — run 'just build' first", bin)
		}
		ogoBin, _ = filepath.Abs(bin)
	}
	return ogoBin
}

func runOgo(t *testing.T, args ...string) (string, error) {
	cmd := exec.Command(ogoPath(t), args...)
	cmd.Dir = filepath.Join("..")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestCLIBoot(t *testing.T) {
	out, err := runOgo(t, "boot", "--pretty")
	assert.NoError(t, err, "ogo boot:\n%s", out)
	assert.NotZero(t, len(out), "expected non-empty output")
}

func TestCLIGrammar(t *testing.T) {
	// Construct a grammar on the fly: write a temp file, run ogo on it
	ebnf := "@@grammar :: Test\nstart := 'hello' ;\n"
	tmp := filepath.Join(t.TempDir(), "test.ebnf")
	err := os.WriteFile(tmp, []byte(ebnf), 0644)
	assert.NoError(t, err)
	out, err := runOgo(t, "grammar", tmp)
	assert.NoError(t, err, "ogo grammar:\n%s", out)
	assert.NotZero(t, len(out), "expected non-empty output")
}

// codegenCompiles runs ogo grammar with the given flags, captures the
// generated Go source, and verifies it compiles via go vet.
func codegenCompiles(t *testing.T, args ...string) {
	out, err := runOgo(t, args...)
	assert.NoError(t, err, "ogo grammar:\n%s", out)
	assert.NotZero(t, len(out), "expected non-empty output")

	root, err := filepath.Abs("..")
	assert.NoError(t, err)

	dir := t.TempDir()
	err = os.WriteFile(filepath.Join(dir, "gen.go"), []byte(out), 0644)
	assert.NoError(t, err)
	modContent := fmt.Sprintf(`module ogocodegentest

go 1.26.3

require github.com/neogeny/ogopego v0.0.0

replace github.com/neogeny/ogopego => %s
`, root)
	err = os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644)
	assert.NoError(t, err)

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	raw, err := cmd.CombinedOutput()
	assert.NoError(t, err, "go mod tidy:\n%s", raw)
	cmd = exec.Command("go", "vet", "./...")
	cmd.Dir = dir
	raw, err = cmd.CombinedOutput()
	assert.NoError(t, err, "go vet failed:\n%s", raw)
}

func TestCLIParserCodegen(t *testing.T) {
	codegenCompiles(t, "grammar", "--color", "never", "-x", "ogotest", "grammar/calc.ebnf")
}

func TestCLIModelGenCodegen(t *testing.T) {
	codegenCompiles(t, "grammar", "--color", "never", "-g", "ogotest", "grammar/calc_model.tatsu")
}

func TestCLIRoundtrip(t *testing.T) {
	// Generate parser code from tatsu.ebnf, build it, and run it
	// on tatsu.ebnf to verify the generated parser works end-to-end.
	out, err := runOgo(t, "grammar", "--color", "never", "-x", "main", "grammar/tatsu.ebnf")
	assert.NoError(t, err, "ogo grammar --parser:\n%s", out)

	root, _ := filepath.Abs("..")
	dir := t.TempDir()
	err = os.WriteFile(filepath.Join(dir, "tatsuparser.go"), []byte(out), 0644)
	assert.NoError(t, err)

	mainCode := `package main

import (
	"fmt"
	"os"

	"github.com/neogeny/ogopego/api"
	"github.com/neogeny/ogopego/pkg/config"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gram := TatSuParser()
	if err := gram.Initialize(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cfg := &config.Cfg{}
	_, err = api.ParseInput(&gram, string(data), cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
`
	err = os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainCode), 0644)
	assert.NoError(t, err)

	modContent := fmt.Sprintf(`module ogotest

go 1.26.3

require github.com/neogeny/ogopego v0.0.0

replace github.com/neogeny/ogopego => %s
`, root)
	err = os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644)
	assert.NoError(t, err)

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	raw, err := cmd.CombinedOutput()
	assert.NoError(t, err, "go mod tidy:\n%s", raw)
	cmd = exec.Command("go", "build", "-o", filepath.Join(dir, "roundtrip"), ".")
	cmd.Dir = dir
	raw, err = cmd.CombinedOutput()
	assert.NoError(t, err, "go build:\n%s", raw)
	cmd = exec.Command(filepath.Join(dir, "roundtrip"), filepath.Join(root, "grammar", "tatsu.ebnf"))
	raw, err = cmd.CombinedOutput()
	assert.NoError(t, err, "roundtrip failed:\n%s", raw)
}

func TestCLIRun(t *testing.T) {
	ebnf := "@@grammar :: Test\nstart := 'hello' ;\n"
	grm := filepath.Join(t.TempDir(), "test.ebnf")
	err := os.WriteFile(grm, []byte(ebnf), 0644)
	assert.NoError(t, err)
	inp := filepath.Join(t.TempDir(), "input.txt")
	err = os.WriteFile(inp, []byte("hello"), 0644)
	assert.NoError(t, err)
	out, err := runOgo(t, "run", grm, inp)
	assert.NoError(t, err, "ogo run:\n%s", out)
	assert.NotZero(t, len(out), "expected non-empty output")
}
