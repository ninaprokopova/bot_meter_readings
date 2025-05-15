-- +goose Up
CREATE TABLE meter_readings (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(user_id),
    cold_water INTEGER NOT NULL,
    hot_water INTEGER NOT NULL,
    electricity_day INTEGER NOT NULL,
    electricity_night INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_meter_readings_user ON meter_readings(user_id);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS meter_readings;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
