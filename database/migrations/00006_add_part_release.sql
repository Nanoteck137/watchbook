-- +goose Up
CREATE TABLE media_part_release (
    media_id TEXT NOT NULL PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,

    num_expected_parts INTEGER NOT NULL,
    current_part INTEGER NOT NULL,
    next_airing TEXT NOT NULL,
    interval_days INTEGER NOT NULL,
    is_active INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

-- +goose Down
DROP TABLE media_part_release;
