-- +goose Up
-- +goose StatementBegin
CREATE TYPE goals_category_type AS ENUM ('growth', 'maintenance', 'other');

CREATE TABLE IF NOT EXISTS goals (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) UNIQUE NOT NULL,
    color VARCHAR(7),
    category_type goals_category_type NOT NULL DEFAULT 'growth',
    is_archived BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goals;
-- +goose StatementEnd
