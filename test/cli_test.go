package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
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
	if os.Getenv("XONSH_VERSION") == "" {
		t.Skip("XONSH_VERSION not set — local test only")
	}
	out, err := runOgo(t, "boot", "--pretty")
	if err != nil {
		t.Fatalf("ogo boot: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestCLIGrammar(t *testing.T) {
	if os.Getenv("XONSH_VERSION") == "" {
		t.Skip("XONSH_VERSION not set — local test only")
	}
	// Construct a grammar on the fly: write a temp file, run ogo on it
	ebnf := "@@grammar :: Test\nstart := 'hello' ;\n"
	tmp := filepath.Join(t.TempDir(), "test.ebnf")
	if err := os.WriteFile(tmp, []byte(ebnf), 0644); err != nil {
		t.Fatal(err)
	}
	out, err := runOgo(t, "grammar", tmp)
	if err != nil {
		t.Fatalf("ogo grammar: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty output")
	}
}

// codegenCompiles runs ogo grammar with the given flags, captures the
// generated Go source, and verifies it compiles via go vet.
func codegenCompiles(t *testing.T, args ...string) {
	out, err := runOgo(t, args...)
	if err != nil {
		t.Fatalf("ogo grammar: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty output")
	}

	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "gen.go"), []byte(out), 0644); err != nil {
		t.Fatal(err)
	}
	modContent := fmt.Sprintf(`module ogocodegentest

go 1.26.3

require github.com/neogeny/ogopego v0.0.0

replace github.com/neogeny/ogopego => %s
`, root)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
	cmd = exec.Command("go", "vet", "./...")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go vet failed:\n%s", out)
	}
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
	if err != nil {
		t.Fatalf("ogo grammar --parser: %v\n%s", err, out)
	}

	root, _ := filepath.Abs("..")
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "tatsuparser.go"), []byte(out), 0644); err != nil {
		t.Fatal(err)
	}

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
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainCode), 0644); err != nil {
		t.Fatal(err)
	}

	modContent := fmt.Sprintf(`module ogotest

go 1.26.3

require github.com/neogeny/ogopego v0.0.0

replace github.com/neogeny/ogopego => %s
`, root)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
	cmd = exec.Command("go", "build", "-o", filepath.Join(dir, "roundtrip"), ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	cmd = exec.Command(filepath.Join(dir, "roundtrip"), filepath.Join(root, "grammar", "tatsu.ebnf"))
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("roundtrip failed: %v\n%s", err, out)
	}
}

func TestCLIRun(t *testing.T) {
	if os.Getenv("XONSH_VERSION") == "" {
		t.Skip("XONSH_VERSION not set — local test only")
	}
	ebnf := "@@grammar :: Test\nstart := 'hello' ;\n"
	grm := filepath.Join(t.TempDir(), "test.ebnf")
	if err := os.WriteFile(grm, []byte(ebnf), 0644); err != nil {
		t.Fatal(err)
	}
	inp := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(inp, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	out, err := runOgo(t, "run", grm, inp)
	if err != nil {
		t.Fatalf("ogo run: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Fatal("expected non-empty output")
	}
}
