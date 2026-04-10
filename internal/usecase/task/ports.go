package task

import (
	"context"

	taskdomain "example.com/taskservice/internal/domain/task"
	"example.com/taskservice/internal/domain/recurrence"
)

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
	// Occurrences returns all scheduled dates for a task in [from, to].
	Occurrences(ctx context.Context, id int64, from, to string) ([]string, error)
}

type CreateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
	Recurrence  *recurrence.Recurrence
}

type UpdateInput struct {
	Title       string
	Description string
	Status      taskdomain.Status
	Recurrence  *recurrence.Recurrence
}
