package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/types"
)

type ShowSeasonItem struct {
	RowId int `db:"rowid"`

	ShowSeasonNum int    `db:"show_season_num"`
	ShowId        string `db:"show_id"`
	MediaId       string `db:"media_id"`

	Position int `db:"position"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Update this to match Media
type FullShowSeasonItem struct {
	RowId int `db:"rowid"`

	ShowSeasonNum int    `db:"show_season_num"`
	ShowId        string `db:"show_id"`
	MediaId       string `db:"media_id"`

	Position int `db:"position"`

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

	MediaDefaultProvider sql.NullString `db:"media_default_provider"`
	MediaProviders       ember.KVStore  `db:"media_providers"`

	MediaCreated int64 `db:"media_created"`
	MediaUpdated int64 `db:"media_updated"`

	MediaPartCount sql.NullInt64 `db:"media_part_count"`

	MediaCreators ember.JsonColumn[[]string] `db:"media_creators"`
	MediaTags     ember.JsonColumn[[]string] `db:"media_tags"`

	MediaUserData ember.JsonColumn[MediaUserData] `db:"media_user_data"`

	MediaRelease ember.JsonColumn[MediaRelease] `db:"media_release"`
}

// TODO(patrik): Use goqu.T more
func ShowSeasonItemQuery() *goqu.SelectDataset {
	query := dialect.From("show_season_items").
		Select(
			"show_season_items.rowid",

			"show_season_items.show_season_num",
			"show_season_items.show_id",
			"show_season_items.media_id",

			"show_season_items.position",

			"show_season_items.created",
			"show_season_items.updated",
		)

	return query
}

// TODO(patrik): Use goqu.T more
func FullShowSeasonItemQuery(userId *string) *goqu.SelectDataset {
	mediaQuery := MediaQuery(userId)

	query := dialect.From("show_season_items").
		Select(
			"show_season_items.rowid",

			"show_season_items.show_season_num",
			"show_season_items.show_id",
			"show_season_items.media_id",

			"show_season_items.position",

			"show_season_items.created",
			"show_season_items.updated",

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

			goqu.I("media.default_provider").As("media_default_provider"),
			goqu.I("media.providers").As("media_providers"),

			goqu.I("media.created").As("media_created"),
			goqu.I("media.updated").As("media_updated"),

			goqu.I("media.part_count").As("media_part_count"),

			goqu.I("media.creators").As("media_creators"),
			goqu.I("media.tags").As("media_tags"),

			goqu.I("media.user_data").As("media_user_data"),

			goqu.I("media.release").As("media_release"),
		).
		Join(
			mediaQuery.As("media"),
			goqu.On(goqu.I("show_season_items.media_id").Eq(goqu.I("media.id"))),
		)

	return query
}

func (db DB) GetAllShowSeasonItems(ctx context.Context) ([]ShowSeasonItem, error) {
	query := ShowSeasonItemQuery()
	return ember.Multiple[ShowSeasonItem](db.db, ctx, query)
}

func (db DB) GetAllShowSeasonItemsByShowSeason(ctx context.Context, showSeasonNum int, showId string) ([]ShowSeasonItem, error) {
	query := ShowSeasonItemQuery().
		Where(
			goqu.I("show_season_items.show_season_num").Eq(showSeasonNum),
			goqu.I("show_season_items.show_id").Eq(showId),
		).
		Order(goqu.I("show_season_items.position").Asc())
	return ember.Multiple[ShowSeasonItem](db.db, ctx, query)
}

func (db DB) GetFullAllShowSeasonItemsByShowSeason(ctx context.Context, userId *string, showSeasonNum int, showId string) ([]FullShowSeasonItem, error) {
	query := FullShowSeasonItemQuery(userId).
		Where(
			goqu.I("show_season_items.show_season_num").Eq(showSeasonNum),
			goqu.I("show_season_items.show_id").Eq(showId),
		).
		Order(goqu.I("show_season_items.position").Asc())
	return ember.Multiple[FullShowSeasonItem](db.db, ctx, query)
}

func (db DB) GetShowSeasonItemById(ctx context.Context, showSeasonNum int, showId, mediaId string) (ShowSeasonItem, error) {
	query := ShowSeasonItemQuery().
		Where(
			goqu.I("show_season_items.show_season_num").Eq(showSeasonNum),
			goqu.I("show_season_items.show_id").Eq(showId),
			goqu.I("show_season_items.media_id").Eq(mediaId),
		)

	return ember.Single[ShowSeasonItem](db.db, ctx, query)
}

type CreateShowSeasonItemParams struct {
	ShowSeasonNum int
	ShowId        string
	MediaId       string

	Position int

	Created int64
	Updated int64
}

func (db DB) CreateShowSeasonItem(ctx context.Context, params CreateShowSeasonItemParams) error {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	query := dialect.Insert("show_season_items").Rows(goqu.Record{
		"show_season_num": params.ShowSeasonNum,
		"show_id":         params.ShowId,
		"media_id":        params.MediaId,

		"position": params.Position,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type ShowSeasonItemChanges struct {
	Position Change[int]

	Created Change[int64]
}

func (db DB) UpdateShowSeasonItem(ctx context.Context, showSeasonNum int, showId, mediaId string, changes ShowSeasonItemChanges) error {
	record := goqu.Record{}

	addToRecord(record, "position", changes.Position)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("show_season_items").
		Set(record).
		Where(
			goqu.I("show_season_items.show_season_num").Eq(showSeasonNum),
			goqu.I("show_season_items.show_id").Eq(showId),
			goqu.I("show_season_items.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveShowSeasonItem(ctx context.Context, showSeasonNum int, showId, mediaId string) error {
	query := dialect.Delete("show_season_items").
		Where(
			goqu.I("show_season_items.show_season_num").Eq(showSeasonNum),
			goqu.I("show_season_items.show_id").Eq(showId),
			goqu.I("show_season_items.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveAllShowSeasonItems(ctx context.Context, showSeasonNum int, showId string) error {
	query := dialect.Delete("show_season_items").
		Where(
			goqu.I("show_season_items.show_season_num").Eq(showSeasonNum),
			goqu.I("show_season_items.show_id").Eq(showId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
