package toy_test

import (
	"os"
	"strings"
	"testing"
)

const (
	cdr0008Step1   = "1. `git clone https://github.com/sondalab-ai/ads-toyrepo.git && cd ads-toyrepo`"
	cdr0008Step2   = "2. `ads init`: the onboarding wizard. It sets up the signing key, ssh-agent, the hooks, and the GitHub rulesets."
	cdr0008Step3   = "3. `ads doctor`: checks the setup and changes nothing."
	cdr0008Step4   = "4. Open Claude Code in the repository and type `fix #<n>`, where `<n>` is a GitHub issue number."
	cdr0008GoTest  = "`go test ./...` must pass."
	cdr0008Heading = "## Try ads here"
)

func cdr0008TrySection(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	start := -1
	for i, l := range lines {
		if strings.TrimRight(l, " \t") == cdr0008Heading {
			start = i + 1
			break
		}
	}
	if start < 0 {
		t.Fatalf("README.md has no %q section", cdr0008Heading)
	}
	end := len(lines)
	for i := start; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "#") {
			end = i
			break
		}
	}
	return lines[start:end]
}

func cdr0008IndexOf(lines []string, want string) int {
	for i, l := range lines {
		if strings.TrimRight(l, " \t") == want {
			return i
		}
	}
	return -1
}

// Steps 1-3 are unchanged and consecutive.
func TestCDR0008_Steps1To3Unchanged(t *testing.T) {
	sec := cdr0008TrySection(t)
	i1 := cdr0008IndexOf(sec, cdr0008Step1)
	if i1 < 0 {
		t.Fatalf("step 1 changed or missing; want %q", cdr0008Step1)
	}
	if i1+2 >= len(sec) || sec[i1+1] != cdr0008Step2 || sec[i1+2] != cdr0008Step3 {
		t.Fatalf("steps 2-3 changed or not directly after step 1; want %q then %q", cdr0008Step2, cdr0008Step3)
	}
}

// The `go test ./...` line is unchanged and still follows the list.
func TestCDR0008_GoTestLineUnchanged(t *testing.T) {
	sec := cdr0008TrySection(t)
	ig := cdr0008IndexOf(sec, cdr0008GoTest)
	if ig < 0 {
		t.Fatalf("README.md 'Try ads here' section lost the line %q", cdr0008GoTest)
	}
	i3 := cdr0008IndexOf(sec, cdr0008Step3)
	if i3 < 0 || ig < i3 {
		t.Fatalf("%q must come after the numbered list", cdr0008GoTest)
	}
}

// If a step 4 exists, it sits directly after step 3 with the exact text.
// The list has no step beyond 4.
func TestCDR0008_Step4TextAndPlacement(t *testing.T) {
	sec := cdr0008TrySection(t)
	i3 := cdr0008IndexOf(sec, cdr0008Step3)
	if i3 < 0 {
		t.Fatalf("step 3 missing; want %q", cdr0008Step3)
	}
	for i, l := range sec {
		if strings.HasPrefix(strings.TrimSpace(l), "4.") && i != i3+1 {
			t.Fatalf("step 4 must be directly after step 3; found at offset %d: %q", i-i3, l)
		}
		if strings.HasPrefix(strings.TrimSpace(l), "5.") {
			t.Fatalf("unexpected step 5: %q", l)
		}
	}
	if i3+1 < len(sec) {
		next := sec[i3+1]
		if strings.HasPrefix(strings.TrimSpace(next), "4.") && next != cdr0008Step4 {
			t.Fatalf("step 4 text wrong:\n got: %q\nwant: %q", next, cdr0008Step4)
		}
	}
}
