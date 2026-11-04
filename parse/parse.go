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
	w, err := strconv.ParseInt(whole, 10, 64)
	if err != nil || w > (math.MaxInt64-99)/100 {
		return 0, fmt.Errorf("parse: amount %q out of range", orig)
	}
	var f int64
	if frac != "" {
		if len(frac) == 1 {
			frac += "0"
		}
		f, _ = strconv.ParseInt(frac, 10, 64)
	}
	cents := w*100 + f
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
