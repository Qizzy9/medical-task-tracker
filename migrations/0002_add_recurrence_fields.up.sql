
ALTER TABLE tasks ADD COLUMN scheduled_at TIMESTAMPTZ;

ALTER TABLE tasks ADD COLUMN parent_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE;

ALTER TABLE tasks ADD COLUMN recurrence_type TEXT;           
ALTER TABLE tasks ADD COLUMN recurrence_interval INTEGER;   
ALTER TABLE tasks ADD COLUMN recurrence_day_of_month INTEGER;
ALTER TABLE tasks ADD COLUMN specific_dates DATE[];          
ALTER TABLE tasks ADD COLUMN recurrence_parity TEXT;         

CREATE INDEX idx_tasks_scheduled_at ON tasks (scheduled_at);
CREATE INDEX idx_tasks_parent_id ON tasks (parent_id);