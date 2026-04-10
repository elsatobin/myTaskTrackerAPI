package recurrence

import "time"


type RecurrenceRule struct {
	ID string

	Type RecurrenceType

	Interval int

	DaysOfMonth []int

	SpecificDates []time.Time

	EvenOdd string // "even" or "odd"

	StartDate time.Time

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

func (r *RecurrenceRule) NextRun(from time.Time) *time.Time {
	switch r.Type {

	case Daily:
		next := from.AddDate(0, 0, r.Interval)
		return &next

	case Monthly:
		if len(r.DaysOfMonth) == 0 {
			return nil
		}
		day := r.DaysOfMonth[0]
		next := time.Date(from.Year(), from.Month()+1, day, 0, 0, 0, 0, from.Location())
		return &next

	case Specific:
		for _, d := range r.SpecificDates {
			if d.After(from) {
				return &d
			}
		}
		return nil

	case EvenOdd:
		next := from
		for {
			next = next.AddDate(0, 0, 1)

			if r.EvenOdd == "even" && next.Day()%2 == 0 {
				return &next
			}
			if r.EvenOdd == "odd" && next.Day()%2 == 1 {
				return &next
			}
		}
	}

	return nil
}