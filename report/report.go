// Package report renders ledger entries as text.
package report

import (
	"fmt"

	"example.com/toy/ledger"
)

// Summary returns one line "account: cents" per account, in the order of
// ledger.Accounts.
func Summary(entries []ledger.Entry) []string {
	var out []string
	for _, a := range ledger.Accounts(entries) {
		out = append(out, fmt.Sprintf("%s: %d", a, ledger.Balance(entries, a)))
	}
	return out
}
