package task

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

func CalculateDates(config taskdomain.RecurrenceConfig, start time.Time) []time.Time {
	var dates []time.Time

	limit := start.AddDate(0, 0, 30)

	switch config.Type {
	case taskdomain.RecurrenceDaily:
		step := 1
		if config.Interval != nil && *config.Interval > 0 {
			step = *config.Interval
		}
		for d := start.AddDate(0, 0, step); d.Before(limit); d = d.AddDate(0, 0, step) {
			dates = append(dates, d)
		}

	case taskdomain.RecurrenceMonthly:
		if config.DayOfMonth == nil {
			return nil
		}
		for i := 1; i <= 3; i++ {
			d := time.Date(start.Year(), start.Month()+time.Month(i), *config.DayOfMonth, start.Hour(), start.Minute(), 0, 0, time.UTC)
			dates = append(dates, d)
		}

	case taskdomain.RecurrenceDates:
		return config.SpecificDates

	case taskdomain.RecurrenceParity:
		if config.Parity == nil {
			return nil
		}
		for d := start.AddDate(0, 0, 1); d.Before(limit); d = d.AddDate(0, 0, 1) {
			isEven := d.Day()%2 == 0
			if (*config.Parity == "even" && isEven) || (*config.Parity == "odd" && !isEven) {
				dates = append(dates, d)
			}
		}
	}
	return dates
}
