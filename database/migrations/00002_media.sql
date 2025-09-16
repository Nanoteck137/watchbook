-- +goose Up
CREATE TABLE media (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,

	title TEXT NOT NULL CHECK(title<>''),
    description TEXT,

	score FLOAT,
	status TEXT NOT NULL,
    rating TEXT NOT NULL,
    airing_season TEXT REFERENCES tags(slug) ON DELETE SET NULL,

	start_date TEXT, 
    end_date TEXT,

    cover_file TEXT,
    logo_file TEXT,
    banner_file TEXT,

    default_provider TEXT,
    providers TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE media_tags (
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    tag_slug TEXT NOT NULL REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(media_id, tag_slug)
);

CREATE TABLE media_creators (
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    tag_slug TEXT NOT NULL REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(media_id, tag_slug)
);

CREATE TABLE media_user_data (
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    list TEXT NOT NULL,

    part INTEGER,
    revisit_count INTEGER,
    is_revisiting BOOLEAN NOT NULL,

    score INTEGER,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(media_id, user_id)
);

CREATE TABLE media_parts (
    idx INTEGER NOT NULL,
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,

    name TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(idx, media_id)
);

CREATE TABLE media_part_release (
    media_id TEXT NOT NULL PRIMARY KEY REFERENCES media(id) ON DELETE CASCADE,

    type TEXT NOT NULL,
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
DROP TABLE media_parts;
DROP TABLE media_user_data;
DROP TABLE media_creators;
DROP TABLE media_tags;
DROP TABLE media;
