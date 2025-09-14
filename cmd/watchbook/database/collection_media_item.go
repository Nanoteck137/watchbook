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

	GroupName  string `db:"group_name"`
	GroupOrder int64  `db:"group_order"`

	Name           string `db:"name"`
	SearchSlug     string `db:"search_slug"`
	OrderNumber    int64  `db:"order_number"`
	SubOrderNumber int64  `db:"suborder_number"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Update this to match Media
type FullCollectionMediaItem struct {
	RowId int `db:"rowid"`

	CollectionId string `db:"collection_id"`
	MediaId      string `db:"media_id"`

	GroupName  string `db:"group_name"`
	GroupOrder int64  `db:"group_order"`

	CollectionName string `db:"collection_name"`
	SearchSlug     string `db:"search_slug"`
	OrderNumber    int64  `db:"order_number"`
	SubOrderNumber int64  `db:"suborder_number"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	MediaType types.MediaType `db:"media_type"`

	MediaTmdbId    sql.NullString `db:"media_tmdb_id"`
	MediaMalId     sql.NullString `db:"media_mal_id"`
	MediaAnilistId sql.NullString `db:"media_anilist_id"`

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

	MediaCreated int64 `db:"media_created"`
	MediaUpdated int64 `db:"media_updated"`

	MediaPartCount sql.NullInt64 `db:"media_part_count"`

	MediaCreators JsonColumn[[]string]         `db:"media_creators"`
	MediaTags     JsonColumn[[]string]         `db:"media_tags"`

	MediaUserData JsonColumn[MediaUserData] `db:"media_user_data"`
}

// TODO(patrik): Use goqu.T more
func CollectionMediaItemQuery(userId *string) *goqu.SelectDataset {
	query := dialect.From("collection_media_items").
		Select(
			"collection_media_items.rowid",

			"collection_media_items.collection_id",
			"collection_media_items.media_id",

			"collection_media_items.group_name",
			"collection_media_items.group_order",

			"collection_media_items.name",
			"collection_media_items.search_slug",
			"collection_media_items.order_number",
			"collection_media_items.suborder_number",

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

			"collection_media_items.group_name",
			"collection_media_items.group_order",

			goqu.I("collection_media_items.name").As("collection_name"),
			"collection_media_items.search_slug",
			"collection_media_items.order_number",
			"collection_media_items.suborder_number",

			"collection_media_items.created",
			"collection_media_items.updated",

			goqu.I("media.type").As("media_type"),

			goqu.I("media.tmdb_id").As("media_tmdb_id"),
			goqu.I("media.mal_id").As("media_mal_id"),
			goqu.I("media.anilist_id").As("media_anilist_id"),

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

func (db *Database) GetAllCollectionMediaItems(ctx context.Context) ([]CollectionMediaItem, error) {
	query := CollectionMediaItemQuery(nil)
	return ember.Multiple[CollectionMediaItem](db.db, ctx, query)
}

func (db *Database) GetFullAllCollectionMediaItemsByCollection(ctx context.Context, userId *string, collectionId string) ([]FullCollectionMediaItem, error) {
	query := FullCollectionMediaItemQuery(userId).
		Where(goqu.I("collection_media_items.collection_id").Eq(collectionId))
	return ember.Multiple[FullCollectionMediaItem](db.db, ctx, query)
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

	GroupName  string
	GroupOrder int64

	Name           string
	SearchSlug     string
	OrderNumber    int64
	SubOrderNumber int64

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

		"group_name":  params.GroupName,
		"group_order": params.GroupOrder,

		"name":            params.Name,
		"search_slug":     params.SearchSlug,
		"order_number":    params.OrderNumber,
		"suborder_number": params.SubOrderNumber,

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
	GroupName  Change[string]
	GroupOrder Change[int64]

	Name           Change[string]
	SearchSlug     Change[string]
	OrderNumber    Change[int64]
	SubOrderNumber Change[int64]

	Created Change[int64]
}

func (db *Database) UpdateCollectionMediaItem(ctx context.Context, collectionId, mediaId string, changes CollectionMediaItemChanges) error {
	record := goqu.Record{}

	addToRecord(record, "group_name", changes.GroupName)
	addToRecord(record, "group_order", changes.GroupOrder)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "search_slug", changes.SearchSlug)
	addToRecord(record, "order_number", changes.OrderNumber)
	addToRecord(record, "suborder_number", changes.SubOrderNumber)

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

func (db *Database) RemoveAllCollectionMediaItems(ctx context.Context, collectionId string) error {
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
