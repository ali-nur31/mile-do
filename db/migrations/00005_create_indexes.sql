-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_goals_user ON goals(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_calendar ON tasks(user_id, scheduled_date);
CREATE INDEX IF NOT EXISTS idx_tasks_inbox ON tasks(user_id)
    WHERE scheduled_date IS NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_analytics ON tasks(user_id, scheduled_date, duration_minutes)
    WHERE is_done = TRUE;
CREATE INDEX IF NOT EXISTS idx_tasks_goal ON tasks(goal_id);
CREATE INDEX IF NOT EXISTS idx_recurring_user ON recurring_tasks_templates(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_goals_user;
DROP INDEX IF EXISTS idx_tasks_calendar;
DROP INDEX IF EXISTS idx_tasks_inbox;
DROP INDEX IF EXISTS idx_tasks_analytics;
DROP INDEX IF EXISTS idx_tasks_goal;
DROP INDEX IF EXISTS idx_recurring_user;
-- +goose StatementEnd
