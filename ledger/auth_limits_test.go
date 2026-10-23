package ledger

import "testing"

func TestDailyLimit(t *testing.T) {
	if DailyLimit("admin") <= DailyLimit("teller") || DailyLimit("teller") <= DailyLimit("guest") {
		t.Fatal("limits are not ordered")
	}
}
