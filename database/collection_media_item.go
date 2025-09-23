package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/types"
)

type CollectionMediaItem struct {
	RowId int `db:"rowid"`

	CollectionId string `db:"collection_id"`
	MediaId      string `db:"media_id"`

	Name       string `db:"name"`
	SearchSlug string `db:"search_slug"`
	Position   int    `db:"position"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Update this to match Media
type FullCollectionMediaItem struct {
	RowId int `db:"rowid"`

	CollectionId string `db:"collection_id"`
	MediaId      string `db:"media_id"`

	CollectionName string `db:"collection_name"`
	SearchSlug     string `db:"search_slug"`
	Position       int    `db:"position"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	MediaType types.MediaType `db:"media_type"`

	MediaTitle       string         `db:"media_title"`
	MediaDescription sql.NullString `db:"media_description"`

	MediaScore        sql.NullFloat64   `db:"media_score"`
	MediaStatus       types.MediaStatus `db:"media_status"`
	MediaRating       types.MediaRating `db:"media_rating"`
	MediaAiringSeason sql.NullString    `db:"media_airing_season"`

	MediaStartDate sql.NullString `db:"media_start_date"`
	MediaEndDate   sql.NullString `db:"media_end_date"`

	MediaCoverFile  sql.NullString `db:"media_cover_file"`
	MediaLogoFile   sql.NullString `db:"media_logo_file"`
	MediaBannerFile sql.NullString `db:"media_banner_file"`

	MediaProviders ember.KVStore `db:"media_providers"`

	MediaCreated int64 `db:"media_created"`
	MediaUpdated int64 `db:"media_updated"`

	MediaPartCount sql.NullInt64 `db:"media_part_count"`

	MediaCreators ember.JsonColumn[[]string] `db:"media_creators"`
	MediaTags     ember.JsonColumn[[]string] `db:"media_tags"`

	MediaUserData ember.JsonColumn[MediaUserData] `db:"media_user_data"`
}

// TODO(patrik): Use goqu.T more
func CollectionMediaItemQuery(userId *string) *goqu.SelectDataset {
	query := dialect.From("collection_media_items").
		Select(
			"collection_media_items.rowid",

			"collection_media_items.collection_id",
			"collection_media_items.media_id",

			"collection_media_items.name",
			"collection_media_items.search_slug",
			"collection_media_items.position",

			"collection_media_items.created",
			"collection_media_items.updated",
		)

	return query
}

// TODO(patrik): Use goqu.T more
func FullCollectionMediaItemQuery(userId *string) *goqu.SelectDataset {
	mediaQuery := MediaQuery(userId)

	query := dialect.From("collection_media_items").
		Select(
			"collection_media_items.rowid",

			"collection_media_items.collection_id",
			"collection_media_items.media_id",

			goqu.I("collection_media_items.name").As("collection_name"),
			"collection_media_items.search_slug",
			"collection_media_items.position",

			"collection_media_items.created",
			"collection_media_items.updated",

			goqu.I("media.type").As("media_type"),

			goqu.I("media.title").As("media_title"),
			goqu.I("media.description").As("media_description"),

			goqu.I("media.score").As("media_score"),
			goqu.I("media.status").As("media_status"),
			goqu.I("media.rating").As("media_rating"),
			goqu.I("media.airing_season").As("media_airing_season"),

			goqu.I("media.start_date").As("media_start_date"),
			goqu.I("media.end_date").As("media_end_date"),

			goqu.I("media.cover_file").As("media_cover_file"),
			goqu.I("media.logo_file").As("media_logo_file"),
			goqu.I("media.banner_file").As("media_banner_file"),

			goqu.I("media.providers").As("media_providers"),

			goqu.I("media.created").As("media_created"),
			goqu.I("media.updated").As("media_updated"),

			goqu.I("media.part_count").As("media_part_count"),

			goqu.I("media.creators").As("media_creators"),
			goqu.I("media.tags").As("media_tags"),

			goqu.I("media.user_data").As("media_user_data"),
		).
		Join(
			mediaQuery.As("media"),
			goqu.On(goqu.I("collection_media_items.media_id").Eq(goqu.I("media.id"))),
		)

	return query
}

func (db DB) GetAllCollectionMediaItems(ctx context.Context) ([]CollectionMediaItem, error) {
	query := CollectionMediaItemQuery(nil)
	return ember.Multiple[CollectionMediaItem](db.db, ctx, query)
}

func (db DB) GetFullAllCollectionMediaItemsByCollection(ctx context.Context, userId *string, collectionId string) ([]FullCollectionMediaItem, error) {
	query := FullCollectionMediaItemQuery(userId).
		Where(goqu.I("collection_media_items.collection_id").Eq(collectionId))
	return ember.Multiple[FullCollectionMediaItem](db.db, ctx, query)
}

func (db DB) GetCollectionMediaItemById(ctx context.Context, collectionId, mediaId string) (CollectionMediaItem, error) {
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

	Name       string
	SearchSlug string
	Position   int

	Created int64
	Updated int64
}

func (db DB) CreateCollectionMediaItem(ctx context.Context, params CreateCollectionMediaItemParams) error {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	query := dialect.Insert("collection_media_items").Rows(goqu.Record{
		"collection_id": params.CollectionId,
		"media_id":      params.MediaId,

		"name":        params.Name,
		"search_slug": params.SearchSlug,
		"position":    params.Position,

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
	Name       Change[string]
	SearchSlug Change[string]
	Position   Change[int]

	Created Change[int64]
}

func (db DB) UpdateCollectionMediaItem(ctx context.Context, collectionId, mediaId string, changes CollectionMediaItemChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "search_slug", changes.SearchSlug)
	addToRecord(record, "position", changes.Position)

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

func (db DB) RemoveCollectionMediaItem(ctx context.Context, collectionId, mediaId string) error {
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

func (db DB) RemoveAllCollectionMediaItems(ctx context.Context, collectionId string) error {
	query := dialect.Delete("collection_media_items").
		Where(
			goqu.I("collection_media_items.collection_id").Eq(collectionId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
