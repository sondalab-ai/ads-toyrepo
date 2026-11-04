package parse

import "testing"

func assertAmountErrorCDR0002(t *testing.T, in string) {
	t.Helper()
	got, err := ParseAmount(in)
	if err == nil {
		t.Fatalf("ParseAmount(%q) = %d, nil; want non-nil error", in, got)
	}
}

func TestParseAmountRejectsLettersCDR0002(t *testing.T) {
	assertAmountErrorCDR0002(t, "abc")
}

func TestParseAmountRejectsTrailingGarbageCDR0002(t *testing.T) {
	assertAmountErrorCDR0002(t, "12.3x")
}

func TestParseAmountRejectsEmptyCDR0002(t *testing.T) {
	assertAmountErrorCDR0002(t, "")
}

func TestParseAmountRejectsMoreThanTwoDecimalsCDR0002(t *testing.T) {
	assertAmountErrorCDR0002(t, "1.234")
}

func TestParseAmountValidValuesCDR0002(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"12.34", 1234},
		{"-5", -500},
		{"0.5", 50},
	}
	for _, c := range cases {
		got, err := ParseAmount(c.in)
		if err != nil || got != c.want {
			t.Errorf("ParseAmount(%q) = %d, %v; want %d, nil", c.in, got, err, c.want)
		}
	}
}

func TestParseLineInvalidAmountCDR0002(t *testing.T) {
	for _, line := range []string{"acct,abc", "acct,12.3x", "acct,", "acct,1.234"} {
		e, err := ParseLine(line)
		if err == nil {
			t.Errorf("ParseLine(%q) = %+v, nil; want non-nil error", line, e)
		}
	}
}
