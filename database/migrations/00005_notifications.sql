-- +goose Up
CREATE TABLE notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    type TEXT NOT NULL,
    title TEXT NOT NULL CHECK(title<>''),
    message TEXT NOT NULL,
    metadata TEXT,
    is_read INTEGER NOT NULL,

    dedup_key TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_notifications_dedup
    ON notifications(user_id, type, dedup_key);

CREATE INDEX IF NOT EXISTS idx_notifications_user_created
    ON notifications(user_id, created DESC);

-- +goose Down
DROP INDEX idx_notifications_user_created;
DROP INDEX idx_notifications_dedup;
DROP TABLE notifications;
