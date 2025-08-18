-- +goose Up
ALTER TABLE media_user_data ADD COLUMN updated INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE media_user_data DROP COLUMN updated;
