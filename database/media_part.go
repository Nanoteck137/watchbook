package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
)

type MediaPart struct {
	RowId int `db:"rowid"`

	Index   int64  `db:"idx"`
	MediaId string `db:"media_id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func MediaPartQuery() *goqu.SelectDataset {
	query := dialect.From("media_parts").
		Select(
			"media_parts.rowid",

			"media_parts.idx",
			"media_parts.media_id",

			"media_parts.name",

			"media_parts.created",
			"media_parts.updated",
		)

	return query
}

func (db *Database) GetAllMediaParts(ctx context.Context) ([]MediaPart, error) {
	query := MediaPartQuery()
	return ember.Multiple[MediaPart](db.db, ctx, query)
}

func (db *Database) GetMediaPartByIndexMediaId(ctx context.Context, index int64, mediaId string) (MediaPart, error) {
	query := MediaPartQuery().
		Where(
			goqu.I("media_parts.idx").Eq(index),
			goqu.I("media_parts.media_id").Eq(mediaId),
		)

	return ember.Single[MediaPart](db.db, ctx, query)
}

func (db *Database) GetMediaPartsByMediaId(ctx context.Context, mediaId string) ([]MediaPart, error) {
	query := MediaPartQuery().
		Where(goqu.I("media_parts.media_id").Eq(mediaId)).
		Order(goqu.I("media_parts.idx").Asc())

	return ember.Multiple[MediaPart](db.db, ctx, query)
}

type CreateMediaPartParams struct {
	Index int64
	MediaId string

	Name  string

	Created int64
	Updated int64
}

func (db *Database) CreateMediaPart(ctx context.Context, params CreateMediaPartParams) error {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	query := dialect.Insert("media_parts").Rows(goqu.Record{
		"idx":  params.Index,
		"media_id": params.MediaId,

		"name": params.Name,

		"created": created,
		"updated": updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type MediaPartChanges struct {
	Name  Change[string]

	Created Change[int64]
}

func (db *Database) UpdateMediaPart(ctx context.Context, index int64, mediaId string, changes MediaPartChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("media_parts").
		Set(record).
		Where(
			goqu.I("media_parts.idx").Eq(index),
			goqu.I("media_parts.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveMediaPart(ctx context.Context, idx int64, mediaId string) error {
	query := dialect.Delete("media_parts").
		Where(
			goqu.I("media_parts.idx").Eq(idx),
			goqu.I("media_parts.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
