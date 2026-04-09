package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `
		INSERT INTO tasks (title, description, status, scheduled_at, parent_id, 
		                   recurrence_type, recurrence_interval, recurrence_day_of_month, 
		                   specific_dates, recurrence_parity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, title, description, status, scheduled_at, parent_id, 
		          recurrence_type, recurrence_interval, recurrence_day_of_month, 
		          specific_dates, recurrence_parity, created_at, updated_at
	`
	var recType, recParity *string
	var recInt, recDay *int
	var specDates []time.Time
	if task.RecurrenceConfig != nil {
		t := string(task.RecurrenceConfig.Type)
		recType = &t
		recInt = task.RecurrenceConfig.Interval
		recDay = task.RecurrenceConfig.DayOfMonth
		recParity = task.RecurrenceConfig.Parity
		specDates = task.RecurrenceConfig.SpecificDates
	}
	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.ScheduledAt, task.ParentID, recType, recInt, recDay, specDates, recParity, task.CreatedAt, task.UpdatedAt)
	return scanTask(row)
}

func (r *Repository) CreateMany(ctx context.Context, tasks []taskdomain.Task) error {
	if len(tasks) == 0 {
		return nil
	}
	rows := [][]any{}
	for _, t := range tasks {
		rows = append(rows, []any{t.Title, t.Description, t.Status, t.ScheduledAt, t.ParentID, t.CreatedAt, t.UpdatedAt})
	}
	_, err := r.pool.CopyFrom(ctx, pgx.Identifier{"tasks"}, []string{"title", "description", "status", "scheduled_at", "parent_id", "created_at", "updated_at"}, pgx.CopyFromRows(rows))
	return err
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `SELECT id, title, description, status, scheduled_at, parent_id, recurrence_type, recurrence_interval, recurrence_day_of_month, specific_dates, recurrence_parity, created_at, updated_at FROM tasks WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}
		return nil, err
	}
	return found, nil
}

func (r *Repository) List(ctx context.Context, filterDate *time.Time) ([]taskdomain.Task, error) {
	var query string
	var args []any

	if filterDate != nil {
		query = `SELECT id, title, description, status, scheduled_at, parent_id, recurrence_type, 
		                recurrence_interval, recurrence_day_of_month, specific_dates, 
		                recurrence_parity, created_at, updated_at 
		         FROM tasks 
		         WHERE scheduled_at::date = $1 
		         ORDER BY scheduled_at ASC`
		args = append(args, filterDate.Format("2006-01-02"))
	} else {
		query = `SELECT id, title, description, status, scheduled_at, parent_id, recurrence_type, 
		                recurrence_interval, recurrence_day_of_month, specific_dates, 
		                recurrence_parity, created_at, updated_at 
		         FROM tasks 
		         ORDER BY id DESC`
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	const query = `UPDATE tasks SET title=$1, description=$2, status=$3, scheduled_at=$4, updated_at=$5 WHERE id=$6 
	               RETURNING id, title, description, status, scheduled_at, parent_id, recurrence_type, 
	                         recurrence_interval, recurrence_day_of_month, specific_dates, 
	                         recurrence_parity, created_at, updated_at`
	row := r.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.ScheduledAt, task.UpdatedAt, task.ID)
	return scanTask(row)
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`
	res, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}
	return nil
}

type taskScanner interface{ Scan(dest ...any) error }

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var t taskdomain.Task
	var status string
	var recType, recParity *string
	var recInt, recDay *int
	var specDates []time.Time
	err := scanner.Scan(&t.ID, &t.Title, &t.Description, &status, &t.ScheduledAt, &t.ParentID, &recType, &recInt, &recDay, &specDates, &recParity, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	t.Status = taskdomain.Status(status)
	if recType != nil {
		t.RecurrenceConfig = &taskdomain.RecurrenceConfig{Type: taskdomain.RecurrenceType(*recType), Interval: recInt, DayOfMonth: recDay, SpecificDates: specDates, Parity: recParity}
	}
	return &t, nil
}
