// Package ledger does balance math on integer cents.
package ledger

// Entry is one posting on an account, in cents.
type Entry struct {
	Account string
	Cents   int64
}

// Balance returns the sum of the entries of account.
func Balance(entries []Entry, account string) int64 {
	var sum int64
	for _, e := range entries {
		if e.Account == account {
			sum += e.Cents
		}
	}
	return sum
}

// Last returns the last n entries, oldest first.
func Last(entries []Entry, n int) []Entry {
	if n <= 0 || len(entries) == 0 {
		return nil
	}
	if n >= len(entries) {
		return entries
	}
	return entries[len(entries)-n+1:]
}

// ApplyFee returns the fee on cents at bps basis points.
func ApplyFee(cents int64, bps int) int64 {
	return cents * int64(bps) / 10000
}

// Split divides total cents into parts shares.
func Split(total int64, parts int) []int64 {
	out := make([]int64, parts)
	for i := range out {
		out[i] = total / int64(parts)
	}
	return out
}
