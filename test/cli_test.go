package test

import (
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
