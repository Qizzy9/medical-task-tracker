package task

import (
	"testing"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func TestCalculateDates(t *testing.T) {
	interval := 2
	config := taskdomain.RecurrenceConfig{
		Type:     taskdomain.RecurrenceDaily,
		Interval: &interval,
	}
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	dates := CalculateDates(config, start)

	if len(dates) == 0 {
		t.Error("Expected generated dates, got 0")
	}

	if dates[0].Day() != 3 {
		t.Errorf("Expected first day to be 3, got %d", dates[0].Day())
	}
}
