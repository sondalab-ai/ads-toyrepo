package makefilecheck_test

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot walks up from the test's working directory to the directory holding go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found above test directory")
		}
		dir = parent
	}
}

func makefilePath(t *testing.T) string {
	t.Helper()
	p := filepath.Join(repoRoot(t), "Makefile")
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("Makefile must exist at repository root: %v", err)
	}
	if info.IsDir() {
		t.Fatalf("Makefile at repository root is a directory")
	}
	return p
}

// makeDryRun runs `make -n [target]` at the repo root and returns the trimmed output lines.
func makeDryRun(t *testing.T, target string) []string {
	t.Helper()
	makefilePath(t)
	if _, err := exec.LookPath("make"); err != nil {
		t.Fatalf("make not found on PATH: %v", err)
	}
	args := []string{"-n"}
	if target != "" {
		args = append(args, target)
	}
	cmd := exec.Command("make", args...)
	cmd.Dir = repoRoot(t)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		lines = append(lines, strings.TrimSpace(sc.Text()))
	}
	return lines
}

func hasLine(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}

func TestMakefileExistsAtRoot(t *testing.T) {
	makefilePath(t)
}

func TestMakeTestPrintsGoTest(t *testing.T) {
	lines := makeDryRun(t, "test")
	if !hasLine(lines, "go test ./...") {
		t.Fatalf("make -n test must print `go test ./...`; got:\n%s", strings.Join(lines, "\n"))
	}
}

func TestMakeVetPrintsGoVet(t *testing.T) {
	lines := makeDryRun(t, "vet")
	if !hasLine(lines, "go vet ./...") {
		t.Fatalf("make -n vet must print `go vet ./...`; got:\n%s", strings.Join(lines, "\n"))
	}
}

func TestMakeDefaultPrintsVetAndTest(t *testing.T) {
	lines := makeDryRun(t, "")
	if !hasLine(lines, "go vet ./...") {
		t.Errorf("make -n (default) must print `go vet ./...`; got:\n%s", strings.Join(lines, "\n"))
	}
	if !hasLine(lines, "go test ./...") {
		t.Errorf("make -n (default) must print `go test ./...`; got:\n%s", strings.Join(lines, "\n"))
	}
}

func TestTestAndVetArePhony(t *testing.T) {
	data, err := os.ReadFile(makefilePath(t))
	if err != nil {
		t.Fatalf("read Makefile: %v", err)
	}
	phony := map[string]bool{}
	// Join backslash-continued lines before scanning.
	text := strings.ReplaceAll(string(data), "\\\n", " ")
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, ".PHONY") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(trimmed, ".PHONY"))
		if !strings.HasPrefix(rest, ":") {
			continue
		}
		for _, name := range strings.Fields(strings.TrimPrefix(rest, ":")) {
			phony[name] = true
		}
	}
	for _, target := range []string{"test", "vet"} {
		if !phony[target] {
			t.Errorf("target %q must be declared .PHONY", target)
		}
	}
}
