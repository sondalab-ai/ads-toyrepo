package ledger

import "testing"

func TestWithdrawLimit(t *testing.T) {
	if err := Withdraw("teller", 0, 100); err != nil {
		t.Fatal(err)
	}
	if err := Withdraw("teller", 49950, 100); err == nil {
		t.Fatal("over the limit")
	}
}

func TestWithdrawAllLimit(t *testing.T) {
	if err := WithdrawAll("teller", 0, []int64{100, 200}); err != nil {
		t.Fatal(err)
	}
	if err := WithdrawAll("teller", 0, []int64{30000, 30000}); err == nil {
		t.Fatal("over the limit")
	}
}
