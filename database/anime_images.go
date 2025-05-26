package database

import (
	"context"

	"github.com/doug-martin/goqu/v9"
)

type AnimeImage struct {
	RowId int `db:"rowid"`

	AnimeId string `db:"anime_id"`
	Hash    string `db:"hash"`

	ImageType string `db:"image_type"`
	Filename  string `db:"filename"`
	IsCover   int    `db:"is_cover"`
}

// TODO(patrik): Use goqu.T more
func AnimeImageQuery() *goqu.SelectDataset {
	query := dialect.From("anime_images").
		Select(
			"anime_images.rowid",

			"anime_images.anime_id",
			"anime_images.hash",

			"anime_images.image_type",
			"anime_images.filename",
			"anime_images.is_cover",
		)

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

func (db *Database) RemoveAllImagesFromAnime(ctx context.Context, animeId string) error {
	query := goqu.Delete("anime_images").
		Where(goqu.I("anime_images.anime_id").Eq(animeId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type CreateAnimeImageParams struct {
	AnimeId string
	Hash    string

	ImageType string
	Filename  string
	IsCover   bool
}

func (db *Database) CreateAnimeImage(ctx context.Context, params CreateAnimeImageParams) error {
	query := dialect.Insert("anime_images").Rows(goqu.Record{
		"anime_id": params.AnimeId,
		"hash":     params.Hash,

		"image_type": params.ImageType,
		"filename":   params.Filename,
		"is_cover":   params.IsCover,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type AnimeImageChanges struct {
	ImageType Change[string]
	Filename  Change[string]
	IsCover   Change[bool]
}

func (db *Database) UpdateAnimeImage(ctx context.Context, animeId, hash string, changes AnimeImageChanges) error {
	record := goqu.Record{}

	addToRecord(record, "image_type", changes.ImageType)
	addToRecord(record, "filename", changes.Filename)
	addToRecord(record, "is_cover", changes.IsCover)

	if len(record) == 0 {
		return nil
	}

	query := dialect.Update("anime_images").
		Set(record).
		Where(
			goqu.I("anime_images.anime_id").Eq(animeId),
			goqu.I("anime_images.hash").Eq(hash),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAnimeCover(ctx context.Context, animeId string) error {
	query := dialect.Update("anime_images").
		Set(goqu.Record{
			"is_cover": false,
		}).
		Where(goqu.I("anime_images.anime_id").Eq(animeId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
