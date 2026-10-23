package ledger

import "errors"

var errOverLimit = errors.New("daily limit exceeded")

// Withdraw checks a withdrawal of amount cents against the daily limit of
// role, given the cents already withdrawn today.
func Withdraw(role string, already, amount int64) error {
	if already+amount > DailyLimit(role) {
		return errOverLimit
	}
	return nil
}

// WithdrawAll checks several withdrawals as one against the daily limit of
// role, given the cents already withdrawn today.
func WithdrawAll(role string, already int64, amounts []int64) error {
	var sum int64
	for _, a := range amounts {
		sum += a
	}
	if already+sum > DailyLimit(role) {
		return errOverLimit
	}
	return nil
}
