package ledger

// DailyLimit returns the daily withdrawal limit of role, in cents.
func DailyLimit(role string) int64 {
	switch role {
	case "admin":
		return 1000000
	case "teller":
		return 50000
	default:
		return 10000
	}
}
