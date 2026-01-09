-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS recurring_tasks_templates (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    goal_id INT NOT NULL,
    FOREIGN KEY (goal_id) REFERENCES goals(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    scheduled_datetime TIMESTAMP NOT NULL,
    has_time BOOLEAN NOT NULL DEFAULT true,
    duration_minutes INT NOT NULL DEFAULT 15,
    recurrence_rrule VARCHAR NOT NULL,
    last_generated_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS recurring_tasks_templates;
-- +goose StatementEnd
