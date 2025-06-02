-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL COLLATE NOCASE CHECK(username<>'') UNIQUE,
    password TEXT NOT NULL CHECK(password<>''),
    role TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE api_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    name TEXT NOT NULL CHECK(name<>''),

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

CREATE TABLE users_settings (
    id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    display_name TEXT
);

-- +goose Down
DROP TABLE users_settings;

DROP TABLE studios;
DROP TABLE tags;

DROP TABLE api_tokens;
DROP TABLE users;
