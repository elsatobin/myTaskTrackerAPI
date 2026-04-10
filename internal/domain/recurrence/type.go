package recurrence

type RecurrenceType string

const (
	Daily        RecurrenceType = "daily"
	Monthly      RecurrenceType = "monthly"
	Specific     RecurrenceType = "specific_dates"
	EvenOdd      RecurrenceType = "even_odd"
)