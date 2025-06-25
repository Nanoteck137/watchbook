package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/types"
)

type MediaImage struct {
	RowId int `db:"rowid"`

	MediaId string `db:"media_id"`
	Hash    string `db:"hash"`

	Type      types.MediaImageType `db:"type"`
	MimeType  string               `db:"mime_type"`
	Filename  string               `db:"filename"`
	IsPrimary int                  `db:"is_primary"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func MediaImageQuery() *goqu.SelectDataset {
	query := dialect.From("media_images").
		Select(
			"media_images.rowid",

			"media_images.media_id",
			"media_images.hash",

			"media_images.type",
			"media_images.mime_type",
			"media_images.filename",
			"media_images.is_primary",

			"media_images.created",
			"media_images.updated",
		)

	return query
}

func (db *Database) GetMediaImagesByHashMediaId(ctx context.Context, mediaId, hash string) (MediaImage, error) {
	query := MediaImageQuery().
		Where(
			goqu.I("media_images.media_id").Eq(mediaId),
			goqu.I("media_images.hash").Eq(hash),
		)

	return ember.Single[MediaImage](db.db, ctx, query)
}

func (db *Database) RemoveAllImagesFromMedia(ctx context.Context, mediaId string) error {
	query := goqu.Delete("media_images").
		Where(goqu.I("media_images.media_id").Eq(mediaId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type CreateMediaImageParams struct {
	MediaId string
	Hash    string

	Type      types.MediaImageType
	MimeType  string
	Filename  string
	IsPrimary bool

	Created int64
	Updated int64
}

func (db *Database) CreateMediaImage(ctx context.Context, params CreateMediaImageParams) error {
	t := time.Now().UnixMilli()
	if params.Created == 0 && params.Updated == 0 {
		params.Created = t
		params.Updated = t
	}

	if params.Type == "" {
		params.Type = types.MediaImageTypeUnknown
	}

	query := dialect.Insert("media_images").Rows(goqu.Record{
		"media_id": params.MediaId,
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

type MediaImageChanges struct {
	Type      Change[types.MediaImageType]
	MimeType  Change[string]
	Filename  Change[string]
	IsPrimary Change[bool]

	Created Change[int64]
}

func (db *Database) UpdateMediaImage(ctx context.Context, mediaId, hash string, changes MediaImageChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)
	addToRecord(record, "mime_type", changes.MimeType)
	addToRecord(record, "filename", changes.Filename)
	addToRecord(record, "is_primary", changes.IsPrimary)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	query := dialect.Update("media_images").
		Set(record).
		Where(
			goqu.I("media_images.media_id").Eq(mediaId),
			goqu.I("media_images.hash").Eq(hash),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
