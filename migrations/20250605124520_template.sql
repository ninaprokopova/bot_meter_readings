-- +goose Up
CREATE TABLE template (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE REFERENCES users(user_id),
    template VARCHAR(1000) NOT NULL DEFAULT '$показания$' 
);

CREATE INDEX idx_template ON template(user_id);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS template;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
