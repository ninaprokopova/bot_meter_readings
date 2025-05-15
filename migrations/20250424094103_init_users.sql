-- +goose Up
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY,
    is_subscribed BOOLEAN NOT NULL DEFAULT TRUE,
    has_submitted BOOLEAN NOT NULL DEFAULT FALSE,
    last_reminder_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_subscribed ON users(is_subscribed) 
WHERE is_subscribed = TRUE;
-- +goose StatementBegin
COMMENT ON TABLE users IS 'Хранит статусы подписки пользователей бота';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS users;
