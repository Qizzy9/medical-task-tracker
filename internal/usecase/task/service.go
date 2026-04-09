package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
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

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		return nil, fmt.Errorf("%w: title required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	now := s.now()
	model := &taskdomain.Task{
		Title: input.Title, Description: input.Description, Status: input.Status,
		ScheduledAt: input.ScheduledAt, RecurrenceConfig: input.Recurrence,
		CreatedAt: now, UpdatedAt: now,
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	if input.Recurrence != nil {
		start := now
		if input.ScheduledAt != nil {
			start = *input.ScheduledAt
		}
		dates := CalculateDates(*input.Recurrence, start)
		instances := make([]taskdomain.Task, 0, len(dates))
		for _, d := range dates {
			dCopy := d
			instances = append(instances, taskdomain.Task{
				ParentID: &created.ID, Title: created.Title, Description: created.Description,
				Status: taskdomain.StatusNew, ScheduledAt: &dCopy, CreatedAt: now, UpdatedAt: now,
			})
		}
		_ = s.repo.CreateMany(ctx, instances)
	}
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *Service) List(ctx context.Context, filterDate *time.Time) ([]taskdomain.Task, error) {
	return s.repo.List(ctx, filterDate)
}
func (s *Service) Delete(ctx context.Context, id int64) error { return s.repo.Delete(ctx, id) }
func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	model := &taskdomain.Task{ID: id, Title: input.Title, Description: input.Description, Status: input.Status, ScheduledAt: input.ScheduledAt, UpdatedAt: s.now()}
	return s.repo.Update(ctx, model)
}
