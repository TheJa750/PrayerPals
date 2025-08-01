-- +goose Up
ALTER TABLE groups
ADD COLUMN invite_code VARCHAR(9) UNIQUE NOT NULL;

-- +goose Down
ALTER TABLE groups
DROP COLUMN invite_code;