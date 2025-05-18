-- +goose Up
CREATE TABLE animes (
    id TEXT PRIMARY KEY,

    mal_id TEXT NOT NULL UNIQUE CHECK(mal_id<>''),

	title TEXT NOT NULL CHECK(title<>''),
	title_english TEXT,

    description TEXT,

	type TEXT NOT NULL,
	status TEXT NOT NULL,
    rating TEXT NOT NULL,
    airing_season TEXT NOT NULL,
	episode_count INTEGER,

	start_date TEXT, 
    end_date TEXT,
    release_date TEXT,

	score FLOAT,

    ani_db_url TEXT,
    anime_news_network_url TEXT,

    cover_filename TEXT,

    should_fetch_data BOOLEAN NOT NULL,
	last_data_fetch_date DATE NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE anime_theme_songs (
    anime_id TEXT NOT NULL REFERENCES animes(id) ON DELETE CASCADE,
    idx INTEGER NOT NULL, 

    type TEXT NOT NULL,
    raw TEXT NOT NULL,

    PRIMARY KEY(anime_id, idx)
);

CREATE TABLE anime_tags (
    anime_id TEXT REFERENCES animes(id) ON DELETE CASCADE,
    tag_slug TEXT REFERENCES tags(slug) ON DELETE CASCADE,

    PRIMARY KEY(anime_id, tag_slug)
);

CREATE TABLE anime_studios (
    anime_id TEXT REFERENCES animes(id) ON DELETE CASCADE,
    studio_slug TEXT REFERENCES studios(slug) ON DELETE CASCADE,

    PRIMARY KEY(anime_id, studio_slug)
);

CREATE TABLE anime_user_data (
    anime_id TEXT REFERENCES animes(id) ON DELETE CASCADE,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,

    list TEXT,

    episode INTEGER,
    is_rewatching BOOLEAN NOT NULL,

    score INTEGER,

    PRIMARY KEY(anime_id, user_id)
);

-- +goose Down
DROP TABLE anime_user_data;
DROP TABLE anime_studios;
DROP TABLE anime_tags;
DROP TABLE anime_theme_songs;
DROP TABLE animes;
