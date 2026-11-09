// Package makefiletest holds acceptance tests for the repository Makefile.
package makefiletest

import (
	"bytes"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

const (
	repoRoot = ".."

	// nestedEnv marks the `go test` started by `make cover`. That run covers
	// this package too, so without the marker the test would start itself
	// forever.
	nestedEnv = "ADS_CDR0007_NESTED_MAKE_COVER"
)

// makeEnv returns the environment for a child make, without the variables a
// parent make would pass down (a parent `make -n` must not change the child).
func makeEnv(extra ...string) []string {
	env := []string{}
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		switch name {
		case "MAKEFLAGS", "MFLAGS", "MAKELEVEL", "MAKEOVERRIDES":
			continue
		}
		env = append(env, kv)
	}
	return append(env, extra...)
}

func runMake(t *testing.T, env []string, args ...string) (string, error) {
	t.Helper()
	if _, err := exec.LookPath("make"); err != nil {
		t.Fatalf("make is not on PATH: %v", err)
	}
	cmd := exec.Command("make", args...)
	cmd.Dir = repoRoot
	cmd.Env = env
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func makefileLines(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(repoRoot + "/Makefile")
	if err != nil {
		t.Fatalf("read Makefile: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		line, _, _ = strings.Cut(line, "#")
		lines[i] = line
	}
	return lines
}

// prerequisites returns the prerequisites of every rule line for target, and
// whether the Makefile has such a line.
func prerequisites(lines []string, target string) ([]string, bool) {
	var prereqs []string
	found := false
	for _, line := range lines {
		rest, ok := strings.CutPrefix(line, target+":")
		if !ok || strings.HasPrefix(rest, "=") {
			continue
		}
		found = true
		prereqs = append(prereqs, strings.Fields(rest)...)
	}
	return prereqs, found
}

func TestMakeCoverRunsCoverageAndExitsZero(t *testing.T) {
	if os.Getenv(nestedEnv) != "" {
		t.Skip("already inside make cover")
	}

	out, err := runMake(t, makeEnv(nestedEnv+"=1"), "cover")
	if err != nil {
		t.Fatalf("make cover failed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "coverage:") {
		t.Errorf("make cover printed no coverage report, so it did not run go test -cover:\n%s", out)
	}
}

func TestPhonyDeclaresCover(t *testing.T) {
	phony, found := prerequisites(makefileLines(t), ".PHONY")
	if !found {
		t.Fatal("Makefile has no .PHONY line")
	}
	if !slices.Contains(phony, "cover") {
		t.Errorf(".PHONY = %v, want it to include cover", phony)
	}
}

func TestAllStillDependsOnlyOnTestAndVetOnceCoverExists(t *testing.T) {
	lines := makefileLines(t)

	// The criterion is about `all` after the cover target is added.
	if _, found := prerequisites(lines, "cover"); !found {
		t.Fatal("Makefile has no cover target, so all cannot be checked against it")
	}

	got, found := prerequisites(lines, "all")
	if !found {
		t.Fatal("Makefile has no all target")
	}
	slices.Sort(got)
	if want := []string{"test", "vet"}; !slices.Equal(got, want) {
		t.Errorf("all depends on %v, want only %v", got, want)
	}
}

func TestMakeDryRunCoverPrintsCommand(t *testing.T) {
	out, err := runMake(t, makeEnv(), "-n", "cover")
	if err != nil {
		t.Fatalf("make -n cover failed: %v\n%s", err, out)
	}

	const want = "go test -cover ./..."
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == want {
			return
		}
	}
	t.Errorf("make -n cover did not print %q:\n%s", want, out)
}
