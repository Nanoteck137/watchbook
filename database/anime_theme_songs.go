package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/watchbook/types"
)

type AnimeThemeSong struct {
	RowId int `db:"rowid"`

	AnimeId string `db:"anime_id"`
	Index   int    `db:"idx"`

	Type types.ThemeSongType `db:"type"`
	Raw  string              `db:"raw"`
}

// TODO(patrik): Use goqu.T more
func AnimeThemeSongQuery() *goqu.SelectDataset {
	query := dialect.From("anime_theme_songs").
		Select(
			"anime_theme_songs.rowid",

			"anime_theme_songs.anime_id",
			"anime_theme_songs.idx",

			"anime_theme_songs.type",
			"anime_theme_songs.raw",
		).
		Prepared(true)

	return query
}

// func (db *Database) GetAllAnimes(ctx context.Context) ([]Anime, error) {
// 	query := AnimeQuery()
//
// 	var items []Anime
// 	err := db.Select(&items, query)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return items, nil
// }

// func (db *Database) GetAnimeById(ctx context.Context, id string) (Anime, error) {
// 	query := AnimeQuery().
// 		Where(goqu.I("animes.id").Eq(id))
//
// 	var item Anime
// 	err := db.Get(&item, query)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return Anime{}, ErrItemNotFound
// 		}
//
// 		return Anime{}, err
// 	}
//
// 	return item, nil
// }

func (db *Database) RemoveAllThemeSongsFromAnime(ctx context.Context, animeId string) error {
	query := goqu.Delete("anime_theme_songs").
		Where(goqu.I("anime_theme_songs.anime_id").Eq(animeId)).
		Prepared(true)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type CreateAnimeThemeSongParams struct {
	AnimeId string
	Idx     int

	Type types.ThemeSongType
	Raw  string
}

func (db *Database) CreateAnimeThemeSong(ctx context.Context, params CreateAnimeThemeSongParams) error {
	query := dialect.Insert("anime_theme_songs").Rows(goqu.Record{
		"anime_id": params.AnimeId,
		"idx":      params.Idx,

		"type": params.Type,
		"raw":  params.Raw,
	}).
		Prepared(true)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
