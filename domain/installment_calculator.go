package domain

import (
	"errors"
	"math"
	"time"
)

// Installment represents a calculated installment with amount and due date.
type Installment struct {
	Number  int
	Amount  float64
	DueDate time.Time
}

// CalculateInstallments generates a list of installments based on the parameters.
// amount: The total amount (for split) or individual amount (for recurring).
// count: Number of installments.
// firstDueDate: The due date of the first installment.
// isRecurring: If true, the full amount is repeated for each installment. If false, the amount is split.
func CalculateInstallments(amount float64, count int, firstDueDate time.Time, isRecurring bool) ([]Installment, error) {
	if count < 1 {
		return nil, errors.New("installments count must be at least 1")
	}
	if amount < 0 {
		return nil, errors.New("amount must be non-negative")
	}

	installments := make([]Installment, count)

	// Round amount to 2 decimal places to avoid float precision issues initially
	amount = math.Round(amount*100) / 100

	var partAmount float64
	var remainder float64

	if !isRecurring {
		// Split logic
		// Example: 100 / 3 = 33.33, Remainder 0.01
		// 1st: 33.34, 2nd: 33.33, 3rd: 33.33

		// Use integer math for cents to ensure precision
		totalCents := int64(math.Round(amount * 100))
		partCents := totalCents / int64(count)
		remainderCents := totalCents % int64(count)

		partAmount = float64(partCents) / 100.0
		remainder = float64(remainderCents) / 100.0
	} else {
		// Recurring logic
		partAmount = amount
		remainder = 0
	}

	for i := 0; i < count; i++ {
		instAmount := partAmount
		if !isRecurring && i == 0 {
			instAmount += remainder
			// Fix precision again just in case float math added noise (e.g. 0.33 + 0.01)
			instAmount = math.Round(instAmount*100) / 100
		}

		// Calculate Due Date
		// AddDate(0, i, 0) adds 'i' months.
		// Go's AddDate handles rollover: Jan 31 + 1 month = March 3 (or Feb 28/29)?
		// Actually Go's Time.AddDate normalizes: Jan 31 + 1 month -> Feb 31 -> March 3 (non-leap) or March 2 (leap).
		// This is usually NOT what credit cards do. Credit cards typically snap to last day of month if overflow.
		// However, standard simplistic implementation often accepts AddDate behavior or snaps.
		// User requirement checks: "Handles month-end logic (31st -> 28th/29th)".
		// So we need specific logic to snap to end of month if day doesn't exist.
		dueDate := addMonths(firstDueDate, i)

		installments[i] = Installment{
			Number:  i + 1,
			Amount:  instAmount,
			DueDate: dueDate,
		}
	}

	return installments, nil
}

// addMonths adds months to a date, preserving the day of month unless it doesn't exist,
// in which case it snaps to the last day of that month.
// E.g. Jan 31 + 1 month -> Feb 28 (or 29).
func addMonths(t time.Time, months int) time.Time {
	// First calculate target year and month
	yyyy, mm, dd := t.Date()

	// Add months
	// We can use AddDate to get the target month/year easily, but we need to watch out for the day normalization.
	// Let's get the 1st of the target month
	targetFirstOfMonth := time.Date(yyyy, mm, 1, 0, 0, 0, 0, t.Location()).AddDate(0, months, 0)
	targetYear, targetMonth, _ := targetFirstOfMonth.Date()

	// Get the last day of the target month
	// (Next month 1st minus 1 day)
	targetLastOfMonth := targetFirstOfMonth.AddDate(0, 1, 0).Add(-24 * time.Hour)
	daysInTargetMonth := targetLastOfMonth.Day()

	// Determine resulting day
	// If original day (dd) is greater than days in target month, clamp to last day.
	targetDay := dd
	if targetDay > daysInTargetMonth {
		targetDay = daysInTargetMonth
	}

	// Construct result
	return time.Date(targetYear, targetMonth, targetDay, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}
