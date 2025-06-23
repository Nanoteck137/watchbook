package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/types"
)

type AnimeImage struct {
	RowId int `db:"rowid"`

	AnimeId string `db:"anime_id"`
	Hash    string `db:"hash"`

	Type      types.EntryImageType `db:"type"`
	MimeType  string               `db:"mime_type"`
	Filename  string               `db:"filename"`
	IsPrimary int                  `db:"is_primary"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func AnimeImageQuery() *goqu.SelectDataset {
	query := dialect.From("anime_images").
		Select(
			"anime_images.rowid",

			"anime_images.anime_id",
			"anime_images.hash",

			"anime_images.type",
			"anime_images.mime_type",
			"anime_images.filename",
			"anime_images.is_primary",

			"anime_images.created",
			"anime_images.updated",
		)

	return query
}

func (db *Database) GetAnimeImagesByHashAnimeId(ctx context.Context, animeId, hash string) (AnimeImage, error) {
	query := AnimeImageQuery().
		Where(
			goqu.I("anime_images.anime_id").Eq(animeId),
			goqu.I("anime_images.hash").Eq(hash),
		)

	return ember.Single[AnimeImage](db.db, ctx, query)
}

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

	Type      types.EntryImageType
	MimeType  string
	Filename  string
	IsPrimary bool

	Created int64
	Updated int64
}

func (db *Database) CreateAnimeImage(ctx context.Context, params CreateAnimeImageParams) error {
	t := time.Now().UnixMilli()
	if params.Created == 0 && params.Updated == 0 {
		params.Created = t
		params.Updated = t
	}

	if params.Type == "" {
		params.Type = types.EntryImageTypeUnknown
	}

	query := dialect.Insert("anime_images").Rows(goqu.Record{
		"anime_id": params.AnimeId,
		"hash":     params.Hash,

		"type":       params.Type,
		"mime_type":  params.MimeType,
		"filename":   params.Filename,
		"is_primary": params.IsPrimary,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type AnimeImageChanges struct {
	Type      Change[types.EntryImageType]
	MimeType  Change[string]
	Filename  Change[string]
	IsPrimary Change[bool]

	Created Change[int64]
}

func (db *Database) UpdateAnimeImage(ctx context.Context, animeId, hash string, changes AnimeImageChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)
	addToRecord(record, "mime_type", changes.MimeType)
	addToRecord(record, "filename", changes.Filename)
	addToRecord(record, "is_primary", changes.IsPrimary)

	addToRecord(record, "created", changes.Created)

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
