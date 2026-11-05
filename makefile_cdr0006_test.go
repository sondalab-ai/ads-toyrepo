package toy

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// These tests use `make -n` (dry run) and never execute the recipes. A real
// `make test` here would run `go test ./...`, which would run these tests again.

// dryRun runs `make -n` with the given targets in dir and returns its output.
// It drops the make variables of any outer make, so a run under `make test`
// behaves the same as a direct `go test`.
func dryRun(t *testing.T, dir string, targets ...string) string {
	t.Helper()
	cmd := exec.Command("make", append([]string{"-n"}, targets...)...)
	cmd.Dir = dir
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		switch name {
		case "MAKEFLAGS", "MFLAGS", "MAKELEVEL", "GNUMAKEFLAGS":
			continue
		}
		cmd.Env = append(cmd.Env, kv)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make -n %v: %v\n%s", targets, err, out)
	}
	return string(out)
}

func wantRecipe(t *testing.T, output, recipe string) {
	t.Helper()
	if !strings.Contains(output, recipe) {
		t.Errorf("recipe %q not in make -n output:\n%s", recipe, output)
	}
}

func TestMakefileExistsAtRoot(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	// Compare exact names: a stat would accept "makefile" on a
	// case-insensitive file system.
	for _, e := range entries {
		if e.Name() == "Makefile" && e.Type().IsRegular() {
			return
		}
	}
	t.Fatal("no regular file named Makefile at the repository root")
}

func TestMakeTestRunsGoTest(t *testing.T) {
	wantRecipe(t, dryRun(t, ".", "test"), "go test ./...")
}

func TestMakeVetRunsGoVet(t *testing.T) {
	wantRecipe(t, dryRun(t, ".", "vet"), "go vet ./...")
}

func TestMakeDefaultRunsTestAndVet(t *testing.T) {
	out := dryRun(t, ".")
	wantRecipe(t, out, "go test ./...")
	wantRecipe(t, out, "go vet ./...")
}

// A .PHONY target runs even when a file with its name exists. Without .PHONY,
// make reports the target as up to date and prints no recipe.
func TestTestAndVetArePhony(t *testing.T) {
	makefile, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, content := range map[string][]byte{
		"Makefile": makefile,
		"test":     nil,
		"vet":      nil,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	wantRecipe(t, dryRun(t, dir, "test"), "go test ./...")
	wantRecipe(t, dryRun(t, dir, "vet"), "go vet ./...")
}
