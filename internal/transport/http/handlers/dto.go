package handlers

import (
	"fmt"
	"time"

	"example.com/taskservice/internal/domain/recurrence"
	taskdomain "example.com/taskservice/internal/domain/task"
	taskusecase "example.com/taskservice/internal/usecase/task"
)

// ── Request DTOs ──────────────────────────────────────────────────────────────

type recurrenceRequestDTO struct {
	Type          string   `json:"type"`
	Interval      int      `json:"interval,omitempty"`
	DaysOfMonth   []int    `json:"days_of_month,omitempty"`
	SpecificDates []string `json:"specific_dates,omitempty"` // "2006-01-02"
	EvenOdd       string   `json:"even_odd,omitempty"`
	StartDate     string   `json:"start_date"`          // "2006-01-02"
	EndDate       *string  `json:"end_date,omitempty"`  // "2006-01-02"
}

func (d *recurrenceRequestDTO) toDomain() (*recurrence.Recurrence, error) {
	start, err := time.Parse("2006-01-02", d.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	var endDate *time.Time
	if d.EndDate != nil {
		t, err := time.Parse("2006-01-02", *d.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date: %w", err)
		}
		endDate = &t
	}

	var specificDates []time.Time
	for _, s := range d.SpecificDates {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return nil, fmt.Errorf("invalid specific_date %q: %w", s, err)
		}
		specificDates = append(specificDates, t)
	}

	return &recurrence.Recurrence{
		Type:          recurrence.Type(d.Type),
		Interval:      d.Interval,
		DaysOfMonth:   d.DaysOfMonth,
		SpecificDates: specificDates,
		EvenOdd:       recurrence.EvenOdd(d.EvenOdd),
		StartDate:     start,
		EndDate:       endDate,
	}, nil
}

type taskMutationDTO struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Status      taskdomain.Status     `json:"status"`
	Recurrence  *recurrenceRequestDTO `json:"recurrence,omitempty"`
}

func (d *taskMutationDTO) toCreateInput() (taskusecase.CreateInput, error) {
	input := taskusecase.CreateInput{
		Title:       d.Title,
		Description: d.Description,
		Status:      d.Status,
	}

	if d.Recurrence != nil {
		rec, err := d.Recurrence.toDomain()
		if err != nil {
			return taskusecase.CreateInput{}, err
		}
		input.Recurrence = rec
	}

	return input, nil
}

func (d *taskMutationDTO) toUpdateInput() (taskusecase.UpdateInput, error) {
	input := taskusecase.UpdateInput{
		Title:       d.Title,
		Description: d.Description,
		Status:      d.Status,
	}

	if d.Recurrence != nil {
		rec, err := d.Recurrence.toDomain()
		if err != nil {
			return taskusecase.UpdateInput{}, err
		}
		input.Recurrence = rec
	}

	return input, nil
}

// ── Response DTOs ─────────────────────────────────────────────────────────────

type recurrenceResponseDTO struct {
	Type          recurrence.Type    `json:"type"`
	Interval      int                `json:"interval,omitempty"`
	DaysOfMonth   []int              `json:"days_of_month,omitempty"`
	SpecificDates []string           `json:"specific_dates,omitempty"`
	EvenOdd       recurrence.EvenOdd `json:"even_odd,omitempty"`
	StartDate     string             `json:"start_date"`
	EndDate       *string            `json:"end_date,omitempty"`
}

type taskDTO struct {
	ID          int64                  `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Status      taskdomain.Status      `json:"status"`
	Recurrence  *recurrenceResponseDTO `json:"recurrence,omitempty"`
	DueDate     *string                `json:"due_date,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	dto := taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if task.DueDate != nil {
		s := task.DueDate.Format("2006-01-02")
		dto.DueDate = &s
	}

	if task.Recurrence != nil {
		r := task.Recurrence
		rdto := &recurrenceResponseDTO{
			Type:        r.Type,
			Interval:    r.Interval,
			DaysOfMonth: r.DaysOfMonth,
			EvenOdd:     r.EvenOdd,
			StartDate:   r.StartDate.Format("2006-01-02"),
		}

		for _, d := range r.SpecificDates {
			rdto.SpecificDates = append(rdto.SpecificDates, d.Format("2006-01-02"))
		}

		if r.EndDate != nil {
			s := r.EndDate.Format("2006-01-02")
			rdto.EndDate = &s
		}

		dto.Recurrence = rdto
	}

	return dto
}
