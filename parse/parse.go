// Package parse reads ledger input.
package parse

import (
	"strconv"
	"strings"

	"example.com/toy/ledger"
)

// ParseAmount parses an amount such as "12.34" into cents.
func ParseAmount(s string) (int64, error) {
	s = strings.TrimSpace(s)
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	whole, frac, _ := strings.Cut(s, ".")
	w, _ := strconv.ParseInt(whole, 10, 64)
	var f int64
	if frac != "" {
		if len(frac) == 1 {
			frac += "0"
		}
		f, _ = strconv.ParseInt(frac[:2], 10, 64)
	}
	cents := w*100 + f
	if neg {
		cents = -cents
	}
	return cents, nil
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
