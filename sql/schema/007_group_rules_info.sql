-- +goose Up
ALTER TABLE groups
ADD COLUMN rules_info TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE groups
DROP COLUMN rules_info;