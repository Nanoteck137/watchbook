-- +goose Up
CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL, 

    status TEXT NOT NULL,      -- queued | running | success | failed
    priority INTEGER NOT NULL, -- higher runs first
    run_at INTEGER NOT NULL, -- schedule

    attempts INTEGER NOT NULL,
    max_attempts INTEGER NOT NULL,

    payload TEXT NOT NULL,
    error TEXT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE INDEX idx_jobs_status_priority_runat ON jobs (status, priority DESC, run_at ASC);

-- +goose Down
DROP INDEX idx_jobs_status_priority_runat;
DROP TABLE jobs;
