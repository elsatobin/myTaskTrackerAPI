package recurrence

import "time"

// RecurrenceRule defines how tasks are generated over time.
type RecurrenceRule struct {
	ID string

	Type RecurrenceType

	// For daily recurrence (every N days)
	Interval int

	// For monthly recurrence (e.g. 1, 15, 30)
	DaysOfMonth []int

	// For specific dates recurrence
	SpecificDates []time.Time

	// For even/odd day rules
	EvenOdd string // "even" or "odd"

	StartDate time.Time

	// Optional end date
	EndDate *time.Time
}

func (r RecurrenceRule) IsValid() bool {
	if r.Type == "" {
		return false
	}

	switch r.Type {

	case Daily:
		return r.Interval > 0

	case Monthly:
		return len(r.DaysOfMonth) > 0

	case Specific:
		return len(r.SpecificDates) > 0

	case EvenOdd:
		return r.EvenOdd == "even" || r.EvenOdd == "odd"
	}

	return false
}

func (r RecurrenceRule) IsActive(date time.Time) bool {
	if date.Before(r.StartDate) {
		return false
	}

	if r.EndDate != nil && date.After(*r.EndDate) {
		return false
	}

	return true
}