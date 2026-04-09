package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type RecurrenceType string

const (
	RecurrenceDaily   RecurrenceType = "daily"
	RecurrenceMonthly RecurrenceType = "monthly"
	RecurrenceDates   RecurrenceType = "dates"
	RecurrenceParity  RecurrenceType = "parity"
)

type RecurrenceConfig struct {
	Type          RecurrenceType `json:"type"`
	Interval      *int           `json:"interval,omitempty"`
	DayOfMonth    *int           `json:"day_of_month,omitempty"`
	SpecificDates []time.Time    `json:"specific_dates,omitempty"`
	Parity        *string        `json:"parity,omitempty"`
}

type Task struct {
	ID               int64             `json:"id"`
	ParentID         *int64            `json:"parent_id,omitempty"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Status           Status            `json:"status"`
	ScheduledAt      *time.Time        `json:"scheduled_at,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	RecurrenceConfig *RecurrenceConfig `json:"recurrence_config,omitempty"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
