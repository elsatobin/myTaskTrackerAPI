package task

import (
	"time"

	"example.com/taskservice/internal/domain/recurrence"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      Status
	Recurrence  *recurrence.Recurrence
	// DueDate is set on tasks that were materialized from a recurrence rule.
	// It represents the specific date this task instance is scheduled for.
	DueDate   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
