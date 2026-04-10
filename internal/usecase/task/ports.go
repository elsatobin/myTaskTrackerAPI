package task

import (
	"context"

	"example.com/taskservice/internal/domain/recurrence"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository interface {
	Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	CreateMany(ctx context.Context, tasks []*taskdomain.Task) ([]*taskdomain.Task, error)
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
	// Expand materializes individual task instances for each occurrence of a
	// recurring task in [from, to] and persists them. Returns the created tasks.
	Expand(ctx context.Context, id int64, from, to string) ([]*taskdomain.Task, error)
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
