-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL COLLATE NOCASE CHECK(username<>'') UNIQUE,
    password TEXT NOT NULL CHECK(password<>''),
    role TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE tags (
    slug TEXT PRIMARY KEY CHECK(slug<>''),
    name TEXT NOT NULL CHECK(name<>'')
);

CREATE TABLE studios (
    slug TEXT PRIMARY KEY CHECK(slug<>''),
    name TEXT NOT NULL CHECK(name<>'')
);

CREATE TABLE producers (
    slug TEXT PRIMARY KEY CHECK(slug<>''),
    name TEXT NOT NULL CHECK(name<>'')
);

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

CREATE TABLE anime_producers (
    anime_id TEXT REFERENCES animes(id) ON DELETE CASCADE,
    producer_slug TEXT REFERENCES producers(slug) ON DELETE CASCADE,

    PRIMARY KEY(anime_id, producer_slug)
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

CREATE TABLE users_settings (
    id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT
);

-- +goose Down
DROP TABLE users_settings;

DROP TABLE anime_user_data;
DROP TABLE anime_producers;
DROP TABLE anime_studios;
DROP TABLE anime_demographics;
DROP TABLE anime_genres;
DROP TABLE anime_themes;
DROP TABLE anime_theme_songs;
DROP TABLE animes;

DROP TABLE producers;
DROP TABLE studios;
DROP TABLE tags;

DROP TABLE users;
