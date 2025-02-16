-- +goose Up
ALTER TABLE reports ADD COLUMN suspend_days INT;

-- +goose Down
ALTER TABLE reports DROP COLUMN suspend_days;

