-- +goose Up
CREATE TABLE animes (
    id TEXT PRIMARY KEY,

    download_type TEXT NOT NULL CHECK(download_type<>''),
    mal_id TEXT,
    ani_db_id TEXT,
    anilist_id TEXT,
    anime_news_network_id TEXT,

	title TEXT NOT NULL CHECK(title<>''),
	title_english TEXT,

    description TEXT,

	type TEXT NOT NULL,
	score FLOAT,
	status TEXT NOT NULL,
    rating TEXT NOT NULL,
	episode_count INTEGER,
    airing_season TEXT REFERENCES tags(slug) ON DELETE SET NULL,

	start_date TEXT, 
    end_date TEXT,

	last_data_fetch INTEGER,

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

CREATE TABLE anime_images (
    anime_id TEXT NOT NULL REFERENCES animes(id) ON DELETE CASCADE,
    hash TEXT NOT NULL, 

    image_type TEXT NOT NULL,
    filename TEXT NOT NULL,
    is_cover BOOLEAN NOT NULL,

    PRIMARY KEY(anime_id, hash)
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
    rewatch_count INTEGER,
    is_rewatching BOOLEAN NOT NULL,

    score INTEGER,

    PRIMARY KEY(anime_id, user_id)
);

-- +goose Down
DROP TABLE anime_user_data;
DROP TABLE anime_studios;
DROP TABLE anime_tags;
DROP TABLE anime_images;
DROP TABLE anime_theme_songs;
DROP TABLE animes;
