-- +goose Up
CREATE TABLE folders (
    id TEXT PRIMARY KEY,

    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

	name TEXT NOT NULL CHECK(name<>''),

    cover_file TEXT,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE folder_items (
    folder_id TEXT NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,

    position INTEGER NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(folder_id, media_id)
);

-- +goose Down
DROP TABLE folder_items;
DROP TABLE folders;
