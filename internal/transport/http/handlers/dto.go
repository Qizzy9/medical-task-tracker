package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type recurrenceDTO struct {
	Type          string      `json:"type"`
	Interval      *int        `json:"interval,omitempty"`
	DayOfMonth    *int        `json:"day_of_month,omitempty"`
	SpecificDates []time.Time `json:"specific_dates,omitempty"`
	Parity        *string     `json:"parity,omitempty"`
}

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	ScheduledAt *time.Time        `json:"scheduled_at,omitempty"`
	Recurrence  *recurrenceDTO    `json:"recurrence,omitempty"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	ParentID    *int64            `json:"parent_id,omitempty"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	ScheduledAt *time.Time        `json:"scheduled_at,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Recurrence  *recurrenceDTO    `json:"recurrence,omitempty"`
}

func newTaskDTO(t *taskdomain.Task) taskDTO {
	dto := taskDTO{
		ID:          t.ID,
		ParentID:    t.ParentID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		ScheduledAt: t.ScheduledAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	if t.RecurrenceConfig != nil {
		dto.Recurrence = &recurrenceDTO{
			Type:          string(t.RecurrenceConfig.Type),
			Interval:      t.RecurrenceConfig.Interval,
			DayOfMonth:    t.RecurrenceConfig.DayOfMonth,
			SpecificDates: t.RecurrenceConfig.SpecificDates,
			Parity:        t.RecurrenceConfig.Parity,
		}
	}

	return dto
}
