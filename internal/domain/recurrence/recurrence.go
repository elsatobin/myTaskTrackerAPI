package recurrence

import (
	"fmt"
	"time"
)

// Type defines the recurrence strategy.
type Type string

const (
	// TypeDaily repeats every N days (N >= 1).
	TypeDaily Type = "daily"
	// TypeMonthly repeats on specific days-of-month (1–30).
	TypeMonthly Type = "monthly"
	// TypeSpecificDates repeats only on the listed dates.
	TypeSpecificDates Type = "specific_dates"
	// TypeEvenOdd repeats on even or odd days of the month.
	TypeEvenOdd Type = "even_odd"
)

// EvenOdd selects even or odd days.
type EvenOdd string

const (
	Even EvenOdd = "even"
	Odd  EvenOdd = "odd"
)

// Recurrence holds the full recurrence configuration for a task.
type Recurrence struct {
	Type Type `json:"type"`

	// Daily: interval in days (>= 1).
	Interval int `json:"interval,omitempty"`

	// Monthly: days of month to trigger on (1–30).
	DaysOfMonth []int `json:"days_of_month,omitempty"`

	// SpecificDates: exact dates (date part only, time is ignored).
	SpecificDates []time.Time `json:"specific_dates,omitempty"`

	// EvenOdd: "even" or "odd".
	EvenOdd EvenOdd `json:"even_odd,omitempty"`

	// StartDate: first possible occurrence date (inclusive).
	StartDate time.Time `json:"start_date"`

	// EndDate: last possible occurrence date (inclusive). Nil means no end.
	EndDate *time.Time `json:"end_date,omitempty"`
}

// Validate checks that the recurrence config is self-consistent.
func (r *Recurrence) Validate() error {
	if r.StartDate.IsZero() {
		return fmt.Errorf("start_date is required")
	}

	if r.EndDate != nil && !r.EndDate.After(r.StartDate) {
		return fmt.Errorf("end_date must be after start_date")
	}

	switch r.Type {
	case TypeDaily:
		if r.Interval < 1 {
			return fmt.Errorf("daily recurrence requires interval >= 1")
		}
	case TypeMonthly:
		if len(r.DaysOfMonth) == 0 {
			return fmt.Errorf("monthly recurrence requires at least one day_of_month")
		}
		for _, d := range r.DaysOfMonth {
			if d < 1 || d > 30 {
				return fmt.Errorf("day_of_month must be between 1 and 30, got %d", d)
			}
		}
	case TypeSpecificDates:
		if len(r.SpecificDates) == 0 {
			return fmt.Errorf("specific_dates recurrence requires at least one date")
		}
	case TypeEvenOdd:
		if r.EvenOdd != Even && r.EvenOdd != Odd {
			return fmt.Errorf("even_odd must be 'even' or 'odd'")
		}
	default:
		return fmt.Errorf("unknown recurrence type: %q", r.Type)
	}

	return nil
}

// Occurrences returns all dates in [from, to] that match this recurrence.
// Dates are returned as midnight UTC.
func (r *Recurrence) Occurrences(from, to time.Time) []time.Time {
	from = truncateToDay(from)
	to = truncateToDay(to)

	start := truncateToDay(r.StartDate)
	if from.Before(start) {
		from = start
	}

	if r.EndDate != nil {
		end := truncateToDay(*r.EndDate)
		if to.After(end) {
			to = end
		}
	}

	if from.After(to) {
		return nil
	}

	var result []time.Time

	switch r.Type {
	case TypeDaily:
		result = r.dailyOccurrences(from, to)
	case TypeMonthly:
		result = r.monthlyOccurrences(from, to)
	case TypeSpecificDates:
		result = r.specificOccurrences(from, to)
	case TypeEvenOdd:
		result = r.evenOddOccurrences(from, to)
	}

	return result
}

func (r *Recurrence) dailyOccurrences(from, to time.Time) []time.Time {
	var result []time.Time
	start := truncateToDay(r.StartDate)

	// find first occurrence >= from
	cur := start
	if cur.Before(from) {
		diff := int(from.Sub(start).Hours() / 24)
		steps := diff / r.Interval
		cur = start.AddDate(0, 0, steps*r.Interval)
		if cur.Before(from) {
			cur = cur.AddDate(0, 0, r.Interval)
		}
	}

	for !cur.After(to) {
		result = append(result, cur)
		cur = cur.AddDate(0, 0, r.Interval)
	}

	return result
}

func (r *Recurrence) monthlyOccurrences(from, to time.Time) []time.Time {
	var result []time.Time

	// iterate month by month
	cur := time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !cur.After(end) {
		for _, day := range r.DaysOfMonth {
			// clamp to last day of month
			lastDay := lastDayOfMonth(cur.Year(), cur.Month())
			if day > lastDay {
				continue
			}
			d := time.Date(cur.Year(), cur.Month(), day, 0, 0, 0, 0, time.UTC)
			if !d.Before(from) && !d.After(to) {
				result = append(result, d)
			}
		}
		cur = cur.AddDate(0, 1, 0)
	}

	return result
}

func (r *Recurrence) specificOccurrences(from, to time.Time) []time.Time {
	var result []time.Time
	for _, d := range r.SpecificDates {
		d = truncateToDay(d)
		if !d.Before(from) && !d.After(to) {
			result = append(result, d)
		}
	}
	return result
}

func (r *Recurrence) evenOddOccurrences(from, to time.Time) []time.Time {
	var result []time.Time
	cur := from
	for !cur.After(to) {
		day := cur.Day()
		if (r.EvenOdd == Even && day%2 == 0) || (r.EvenOdd == Odd && day%2 != 0) {
			result = append(result, cur)
		}
		cur = cur.AddDate(0, 0, 1)
	}
	return result
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
