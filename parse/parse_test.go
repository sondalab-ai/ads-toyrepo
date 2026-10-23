package parse

import "testing"

func TestParseAmount(t *testing.T) {
	for in, want := range map[string]int64{"12.34": 1234, "7": 700, "0.5": 50, "-1.50": -150} {
		got, err := ParseAmount(in)
		if err != nil || got != want {
			t.Fatalf("ParseAmount(%q) = %d, %v; want %d", in, got, err, want)
		}
	}
}

func TestParseLine(t *testing.T) {
	e, err := ParseLine("acct , 3.25")
	if err != nil || e.Account != "acct" || e.Cents != 325 {
		t.Fatalf("ParseLine = %+v, %v", e, err)
	}
}
