package ledger

// Accounts returns the distinct accounts of entries in first-seen order.
func Accounts(entries []Entry) []string {
	var out []string
	seen := map[string]bool{}
	for _, e := range entries {
		if !seen[e.Account] {
			seen[e.Account] = true
			out = append(out, e.Account)
		}
	}
	return out
}
