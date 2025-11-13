-- +goose Up
CREATE TABLE shows (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,

	name TEXT NOT NULL CHECK(name<>''),
    search_slug TEXT NOT NULL,

    cover_file TEXT,
    logo_file TEXT,
    banner_file TEXT,

    default_provider TEXT,
    providers TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE show_seasons (
    num INTEGER NOT NULL,
    show_id TEXT NOT NULL REFERENCES shows(id) ON DELETE CASCADE,

    name TEXT NOT NULL,
    search_slug TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    PRIMARY KEY(num, show_id)
);

CREATE TABLE show_season_parts (
    show_id TEXT NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    season_num INTEGER NOT NULL,
    idx INTEGER NOT NULL,

    name TEXT NOT NULL,
	release_date TEXT, 

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    FOREIGN KEY(season_num, show_id) REFERENCES show_seasons(num, show_id),
    PRIMARY KEY(idx, season_num, show_id)
);

-- CREATE TABLE show_season_items (
--     show_season_num INTEGER NOT NULL,
--     show_id TEXT NOT NULL,
--     media_id TEXT NOT NULL REFERENCES media(id) ON DELETE CASCADE,
--
--     position INTEGER NOT NULL,
--
--     created INTEGER NOT NULL,
--     updated INTEGER NOT NULL,
--
--     FOREIGN KEY(show_season_num, show_id) REFERENCES show_seasons(num, show_id),
--     PRIMARY KEY(show_season_num, show_id, media_id)
-- );

-- +goose Down
-- DROP TABLE show_season_items;
DROP TABLE show_season_parts;
DROP TABLE show_seasons;
DROP TABLE shows;
