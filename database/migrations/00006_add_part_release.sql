-- +goose Up
CREATE TABLE media_part_release (
    media_id TEXT NOT NULL PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,

    start_date DATETIME NOT NULL,
    num_expected_parts INTEGER NOT NULL,
    part_offset INTEGER NOT NULL,
    interval_days INTEGER NOT NULL,
    delay_days INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

-- +goose Down
DROP TABLE media_part_release;
