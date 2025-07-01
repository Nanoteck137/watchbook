package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
)

type CollectionMediaItem struct {
	RowId int `db:"rowid"`

	CollectionId string `db:"collectionId"`
	MediaId      string `db:"mediaId"`

	Name        string `db:"name"`
	OrderNumber int64  `db:"order_number"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func CollectionMediaItemQuery(userId *string) *goqu.SelectDataset {
	query := dialect.From("collection_media_items").
		Select(
			"collection_media_items.rowid",

			"collection_media_items.collection_id",
			"collection_media_items.media_id",

			"collection_media_items.name",
			"collection_media_items.order_number",

			"collection_media_items.created",
			"collection_media_items.updated",
		)

	return query
}

func (db *Database) GetAllCollectionMediaItems(ctx context.Context) ([]CollectionMediaItem, error) {
	query := CollectionMediaItemQuery(nil)
	return ember.Multiple[CollectionMediaItem](db.db, ctx, query)
}

func (db *Database) GetAllCollectionMediaItemsByCollection(ctx context.Context, collectionId string) ([]CollectionMediaItem, error) {
	query := CollectionMediaItemQuery(nil).
		Where(goqu.I("collection_media_items.collection_id").Eq(collectionId))
	return ember.Multiple[CollectionMediaItem](db.db, ctx, query)
}

func (db *Database) GetCollectionMediaItemById(ctx context.Context, collectionId, mediaId string) (CollectionMediaItem, error) {
	query := CollectionMediaItemQuery(nil).
		Where(
			goqu.I("collection_media_items.collection_id").Eq(collectionId),
			goqu.I("collection_media_items.media_id").Eq(mediaId),
		)

	return ember.Single[CollectionMediaItem](db.db, ctx, query)
}

type CreateCollectionMediaItemParams struct {
	CollectionId string
	MediaId      string

	Name        string
	OrderNumber int64

	Created int64
	Updated int64
}

func (db *Database) CreateCollectionMediaItem(ctx context.Context, params CreateCollectionMediaItemParams) error {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	query := dialect.Insert("collection_media_items").Rows(goqu.Record{
		"collection_id": params.CollectionId,
		"media_id":      params.MediaId,

		"name":         params.Name,
		"order_number": params.OrderNumber,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type CollectionMediaItemChanges struct {
	Name        Change[string]
	OrderNumber Change[int64]

	Created Change[int64]
}

func (db *Database) UpdateCollectionMediaItem(ctx context.Context, collectionId, mediaId string, changes CollectionMediaItemChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "order_number", changes.OrderNumber)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("collection_media_items").
		Set(record).
		Where(
			goqu.I("collection_media_items.collection_id").Eq(collectionId),
			goqu.I("collection_media_items.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveCollectionMediaItem(ctx context.Context, collectionId, mediaId string) error {
	query := dialect.Delete("collection_media_items").
		Where(
			goqu.I("collection_media_items.collection_id").Eq(collectionId),
			goqu.I("collection_media_items.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
