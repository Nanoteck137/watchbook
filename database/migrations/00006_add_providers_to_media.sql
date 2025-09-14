-- +goose Up
ALTER TABLE media ADD COLUMN providers TEXT NOT NULL DEFAULT "{}";

-- +goose Down
ALTER TABLE media DROP COLUMN providers;
