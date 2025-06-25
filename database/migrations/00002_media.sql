-- +goose Up
CREATE TABLE media (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,

    tmdb_id TEXT,
    mal_id TEXT,
    anilist_id TEXT,

	title TEXT NOT NULL CHECK(title<>''),
    description TEXT,

	score FLOAT,
	status TEXT NOT NULL,
    rating TEXT NOT NULL,
    airing_season TEXT REFERENCES tags(slug) ON DELETE SET NULL,

	start_date TEXT, 
    end_date TEXT,

    admin_status TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE media_images (
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    hash TEXT NOT NULL, 

    type TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    filename TEXT NOT NULL,
    is_primary BOOL NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(media_id, hash)
);

CREATE TABLE media_tags (
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    tag_slug TEXT NOT NULL REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(media_id, tag_slug)
);

CREATE TABLE media_studios (
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    tag_slug TEXT NOT NULL REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(media_id, tag_slug)
);

CREATE TABLE media_user_data (
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    list TEXT,

    part INTEGER,
    revisit_count INTEGER,
    is_revisiting BOOLEAN NOT NULL,

    score INTEGER,

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

-- +goose Down
DROP TABLE media_parts;
DROP TABLE media_user_data;
DROP TABLE media_studios;
DROP TABLE media_tags;
DROP TABLE media_images;
DROP TABLE media;
