// Package parse reads ledger input.
package parse

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"example.com/toy/ledger"
)

// ParseAmount parses an amount such as "12.34" into cents.
func ParseAmount(s string) (int64, error) {
	orig := s
	s = strings.TrimSpace(s)
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	whole, frac, hasDot := strings.Cut(s, ".")
	if !isDigits(whole) {
		return 0, fmt.Errorf("parse: invalid amount %q", orig)
	}
	if hasDot && (!isDigits(frac) || len(frac) > 2) {
		return 0, fmt.Errorf("parse: invalid amount %q", orig)
	}
	// Compute the magnitude in cents as uint64 so the full int64 range,
	// including math.MinInt64, can be checked exactly.
	limit := uint64(math.MaxInt64)
	if neg {
		limit = uint64(math.MaxInt64) + 1
	}
	w, err := strconv.ParseUint(whole, 10, 64)
	if err != nil || w > limit/100 {
		return 0, fmt.Errorf("parse: amount %q out of range", orig)
	}
	var f uint64
	if frac != "" {
		if len(frac) == 1 {
			frac += "0"
		}
		f, _ = strconv.ParseUint(frac, 10, 64)
	}
	mag := w * 100
	if f > limit-mag {
		return 0, fmt.Errorf("parse: amount %q out of range", orig)
	}
	mag += f
	// For mag == 1<<63 (only reachable when neg), the conversion yields
	// math.MinInt64 and negation leaves it unchanged, which is correct.
	cents := int64(mag)
	if neg {
		cents = -cents
	}
	return cents, nil
}

// isDigits reports whether s is non-empty and made only of ASCII digits.
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// ParseLine parses "account,amount" into an entry.
func ParseLine(line string) (ledger.Entry, error) {
	acct, amt, _ := strings.Cut(line, ",")
	c, err := ParseAmount(amt)
	if err != nil {
		return ledger.Entry{}, err
	}
	return ledger.Entry{Account: strings.TrimSpace(acct), Cents: c}, nil
}
