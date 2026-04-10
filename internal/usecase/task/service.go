package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	domain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*domain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	now := s.now()
	model := &domain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Recurrence:  normalized.Recurrence,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.Create(ctx, model)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*domain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &domain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Recurrence:  normalized.Recurrence,
		UpdatedAt:   s.now(),
	}

	return s.repo.Update(ctx, model)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]domain.Task, error) {
	return s.repo.List(ctx)
}

// Occurrences returns all scheduled dates for a task in the given date range.
// from and to are expected as "2006-01-02" (date only).
func (s *Service) Occurrences(ctx context.Context, id int64, from, to string) ([]string, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if task.Recurrence == nil {
		return []string{}, nil
	}

	fromT, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid from date: %s", ErrInvalidInput, from)
	}

	toT, err := time.Parse("2006-01-02", to)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid to date: %s", ErrInvalidInput, to)
	}

	if toT.Before(fromT) {
		return nil, fmt.Errorf("%w: to must be >= from", ErrInvalidInput)
	}

	dates := task.Recurrence.Occurrences(fromT, toT)
	result := make([]string, len(dates))
	for i, d := range dates {
		result[i] = d.Format("2006-01-02")
	}

	return result, nil
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = domain.StatusNew
	}

	if !input.Status.IsValid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.Recurrence != nil {
		if err := input.Recurrence.Validate(); err != nil {
			return CreateInput{}, fmt.Errorf("%w: %s", ErrInvalidInput, err)
		}
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = domain.StatusNew
	}

	if !input.Status.IsValid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if input.Recurrence != nil {
		if err := input.Recurrence.Validate(); err != nil {
			return UpdateInput{}, fmt.Errorf("%w: %s", ErrInvalidInput, err)
		}
	}

	return input, nil
}
