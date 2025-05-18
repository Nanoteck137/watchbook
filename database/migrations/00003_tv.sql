-- +goose Up
CREATE TABLE tv (
    id TEXT PRIMARY KEY,

	title TEXT NOT NULL CHECK(title<>''),

    description TEXT,

	type TEXT,
	score FLOAT,
	status TEXT,
    rating TEXT,
	episode_count INTEGER,

	start_date TEXT, 
    end_date TEXT,
    release_date TEXT,

    cover_filename TEXT,

    should_fetch_data BOOLEAN NOT NULL,
	last_data_fetch_date DATE NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE tv_tags (
    tv_id TEXT REFERENCES tv(id) ON DELETE CASCADE,
    tag_slug TEXT REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(tv_id, tag_slug)
);

CREATE TABLE tv_studios (
    tv_id TEXT REFERENCES tv(id) ON DELETE CASCADE,
    studio_slug TEXT REFERENCES studios(slug) ON DELETE CASCADE,

    PRIMARY KEY(tv_id, studio_slug)
);

CREATE TABLE tv_user_data (
    tv_id TEXT REFERENCES tv(id) ON DELETE CASCADE,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,

    list TEXT,

    episode INTEGER,
    is_rewatching BOOLEAN NOT NULL,

    score INTEGER,

    PRIMARY KEY(tv_id, user_id)
);

-- +goose Down
DROP TABLE tv_user_data;
DROP TABLE tv_studios;
DROP TABLE tv_tags;
DROP TABLE tv;
