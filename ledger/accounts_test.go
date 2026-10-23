package ledger

import "testing"

func TestAccountsDistinct(t *testing.T) {
	got := Accounts([]Entry{{"a", 1}, {"a", 2}})
	if len(got) != 1 || got[0] != "a" {
		t.Fatalf("Accounts = %v", got)
	}
}
