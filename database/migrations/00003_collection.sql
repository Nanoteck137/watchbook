-- +goose Up
CREATE TABLE collections (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,

	name TEXT NOT NULL CHECK(name<>''),

    cover_file TEXT,
    logo_file TEXT,
    banner_file TEXT,

    providers TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE collection_media_items (
    collection_id TEXT NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,

    name TEXT NOT NULL,
    search_slug TEXT NOT NULL,
    position INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(collection_id, media_id)
);

-- +goose Down
DROP TABLE collection_media_items;
DROP TABLE collections;
