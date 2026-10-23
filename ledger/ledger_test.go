package ledger

import "testing"

func TestBalance(t *testing.T) {
	es := []Entry{{"a", 100}, {"b", 50}, {"a", -30}}
	if got := Balance(es, "a"); got != 70 {
		t.Fatalf("Balance(a) = %d", got)
	}
	if got := Balance(es, "c"); got != 0 {
		t.Fatalf("Balance(c) = %d", got)
	}
}

func TestLastAll(t *testing.T) {
	es := []Entry{{"a", 1}, {"a", 2}}
	if got := Last(es, 5); len(got) != 2 {
		t.Fatalf("Last(5) = %v", got)
	}
	if got := Last(es, 0); got != nil {
		t.Fatalf("Last(0) = %v", got)
	}
}

func TestApplyFeeExact(t *testing.T) {
	if got := ApplyFee(10000, 25); got != 25 {
		t.Fatalf("ApplyFee = %d", got)
	}
}

func TestSplitExact(t *testing.T) {
	got := Split(9, 3)
	if len(got) != 3 || got[0] != 3 || got[1] != 3 || got[2] != 3 {
		t.Fatalf("Split = %v", got)
	}
}
