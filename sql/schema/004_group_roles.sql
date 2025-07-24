-- +goose Up
ALTER TABLE users_groups
ADD COLUMN role TEXT NOT NULL DEFAULT 'member';

-- +goose Down
ALTER TABLE users_groups
DROP COLUMN role;