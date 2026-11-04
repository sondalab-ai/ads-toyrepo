package parse

import (
	"math"
	"testing"
)

// CDR-0003: ParseAmount must accept every amount whose value in cents fits
// in int64 and reject every amount whose value in cents does not.

func TestCDR0003MaxInt64CentsAccepted(t *testing.T) {
	got, err := ParseAmount("92233720368547758.07")
	if err != nil || got != math.MaxInt64 {
		t.Fatalf("ParseAmount(%q) = %d, %v; want %d, nil", "92233720368547758.07", got, err, int64(math.MaxInt64))
	}
}

func TestCDR0003NearMaxAccepted(t *testing.T) {
	for in, want := range map[string]int64{
		"92233720368547758":    9223372036854775800,
		"92233720368547758.0":  9223372036854775800,
		"92233720368547757.99": 9223372036854775799,
		"92233720368547758.06": 9223372036854775806,
	} {
		got, err := ParseAmount(in)
		if err != nil || got != want {
			t.Errorf("ParseAmount(%q) = %d, %v; want %d, nil", in, got, err, want)
		}
	}
}

func TestCDR0003OneCentOverMaxRejected(t *testing.T) {
	if got, err := ParseAmount("92233720368547758.08"); err == nil {
		t.Fatalf("ParseAmount(%q) = %d, nil; want an error", "92233720368547758.08", got)
	}
}

func TestCDR0003HugeWholeRejected(t *testing.T) {
	for _, in := range []string{
		"99999999999999999999",
		"92233720368547759",
		"92233720368547758.10",
		"-99999999999999999999",
	} {
		if got, err := ParseAmount(in); err == nil {
			t.Errorf("ParseAmount(%q) = %d, nil; want an error", in, got)
		}
	}
}

// The negative boundary is documented as ACCEPTED: "-92233720368547758.08"
// is exactly math.MinInt64 cents, which fits in int64, so it must parse.
// One cent beyond it must be rejected.
func TestCDR0003MinInt64CentsAccepted(t *testing.T) {
	got, err := ParseAmount("-92233720368547758.08")
	if err != nil || got != math.MinInt64 {
		t.Fatalf("ParseAmount(%q) = %d, %v; want %d, nil", "-92233720368547758.08", got, err, int64(math.MinInt64))
	}
	got, err = ParseAmount("-92233720368547758.07")
	if err != nil || got != -math.MaxInt64 {
		t.Fatalf("ParseAmount(%q) = %d, %v; want %d, nil", "-92233720368547758.07", got, err, int64(-math.MaxInt64))
	}
}

func TestCDR0003OneCentUnderMinRejected(t *testing.T) {
	if got, err := ParseAmount("-92233720368547758.09"); err == nil {
		t.Fatalf("ParseAmount(%q) = %d, nil; want an error", "-92233720368547758.09", got)
	}
}

func TestCDR0003ExistingBehaviorKept(t *testing.T) {
	for in, want := range map[string]int64{"12.34": 1234, "7": 700, "0.5": 50, "-1.50": -150, " 3.25 ": 325, "0": 0} {
		got, err := ParseAmount(in)
		if err != nil || got != want {
			t.Errorf("ParseAmount(%q) = %d, %v; want %d, nil", in, got, err, want)
		}
	}
	for _, in := range []string{"", "-", "abc", "1.234", "1.", ".5", "1.2x", "--1", "+1", "1,00"} {
		if got, err := ParseAmount(in); err == nil {
			t.Errorf("ParseAmount(%q) = %d, nil; want an error", in, got)
		}
	}
}
