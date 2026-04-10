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
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
