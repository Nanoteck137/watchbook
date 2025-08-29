package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/database/adapter"
	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type MediaPartRelease struct {
	RowId int `db:"rowid"`

	MediaId string `db:"media_id"`

	Status           types.MediaPartReleaseStatus `db:"status"`
	StartDate        time.Time                       `db:"start_date"`
	NumExpectedParts int                          `db:"num_expected_parts"`
	CurrentPart      int                          `db:"current_part"`
	NextAiring       time.Time                       `db:"next_airing"`
	IntervalDays     int                          `db:"interval_days"`
	DelayDays        int                          `db:"delay_days"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Update this to match Media
type FullMediaPartRelease struct {
	RowId int `db:"rowid"`

	MediaId string `db:"media_id"`

	Status           types.MediaPartReleaseStatus `db:"status"`
	StartDate        time.Time                    `db:"start_date"`
	NumExpectedParts int                          `db:"num_expected_parts"`
	CurrentPart      int                          `db:"current_part"`
	NextAiring       time.Time                    `db:"next_airing"`
	IntervalDays     int                          `db:"interval_days"`
	DelayDays        int                          `db:"delay_days"`

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

	MediaCreators JsonColumn[[]string] `db:"media_creators"`
	MediaTags     JsonColumn[[]string] `db:"media_tags"`

	MediaUserData JsonColumn[MediaUserData] `db:"media_user_data"`
}

// TODO(patrik): Use goqu.T more
func MediaPartReleaseQuery() *goqu.SelectDataset {
	query := dialect.From("media_part_release").
		Select(
			"media_part_release.rowid",

			"media_part_release.media_id",

			"media_part_release.status",
			"media_part_release.start_date",
			"media_part_release.num_expected_parts",
			"media_part_release.current_part",
			"media_part_release.next_airing",
			"media_part_release.interval_days",
			"media_part_release.delay_days",

			"media_part_release.created",
			"media_part_release.updated",
		)

	return query
}

// TODO(patrik): Use goqu.T more
func FullMediaPartReleaseQuery(userId *string) *goqu.SelectDataset {
	mediaQuery := MediaQuery(userId)

	query := dialect.From("media_part_release").
		Select(
			"media_part_release.rowid",

			"media_part_release.media_id",

			"media_part_release.status",
			"media_part_release.start_date",
			"media_part_release.num_expected_parts",
			"media_part_release.current_part",
			"media_part_release.next_airing",
			"media_part_release.interval_days",
			"media_part_release.delay_days",

			"media_part_release.created",
			"media_part_release.updated",

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
			goqu.On(goqu.I("media_part_release.media_id").Eq(goqu.I("media.id"))),
		)

	return query
}

func (db *Database) GetPagedFullMediaPartReleases(ctx context.Context, userId *string, filterStr, sortStr string, opts FetchOptions) ([]FullMediaPartRelease, types.Page, error) {
	query := FullMediaPartReleaseQuery(userId)

	var err error

	a := adapter.ReleaseResolverAdapter{}
	resolver := filter.New(&a)

	query, err = applyFilter(query, resolver, filterStr)
	if err != nil {
		return nil, types.Page{}, err
	}

	query, err = applySort(query, resolver, sortStr)
	if err != nil {
		return nil, types.Page{}, err
	}

	countQuery := query.
		Select(goqu.COUNT("media.id"))

	if opts.PerPage > 0 {
		query = query.
			Limit(uint(opts.PerPage)).
			Offset(uint(opts.Page * opts.PerPage))
	}

	totalItems, err := ember.Single[int](db.db, ctx, countQuery)
	if err != nil {
		return nil, types.Page{}, err
	}

	totalPages := utils.TotalPages(opts.PerPage, totalItems)
	page := types.Page{
		Page:       opts.Page,
		PerPage:    opts.PerPage,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	items, err := ember.Multiple[FullMediaPartRelease](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db *Database) GetAllMediaPartReleases(ctx context.Context) ([]MediaPartRelease, error) {
	query := MediaPartReleaseQuery()
	return ember.Multiple[MediaPartRelease](db.db, ctx, query)
}

func (db *Database) GetMediaPartReleaseById(ctx context.Context, mediaId string) (MediaPartRelease, error) {
	query := MediaPartReleaseQuery().
		Where(
			goqu.I("media_part_release.media_id").Eq(mediaId),
		)

	return ember.Single[MediaPartRelease](db.db, ctx, query)
}

func (db *Database) GetAllFullMediaPartReleases(ctx context.Context, userId *string) ([]FullMediaPartRelease, error) {
	query := FullMediaPartReleaseQuery(userId)
	return ember.Multiple[FullMediaPartRelease](db.db, ctx, query)
}

func (db *Database) GetFullMediaPartReleaseById(ctx context.Context, userId *string, mediaId string) (FullMediaPartRelease, error) {
	query := FullMediaPartReleaseQuery(userId).
		Where(
			goqu.I("media_part_release.media_id").Eq(mediaId),
		)

	return ember.Single[FullMediaPartRelease](db.db, ctx, query)
}

type CreateMediaPartReleaseParams struct {
	MediaId string

	Status           types.MediaPartReleaseStatus
	StartDate        time.Time
	NumExpectedParts int
	CurrentPart      int
	NextAiring       time.Time
	IntervalDays     int
	DelayDays        int

	Created int64
	Updated int64
}

func (db *Database) CreateMediaPartRelease(ctx context.Context, params CreateMediaPartReleaseParams) error {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if !types.IsValidMediaPartReleaseStatus(params.Status) {
		params.Status = types.MediaPartReleaseStatusUnknown
	}

	query := dialect.Insert("media_part_release").Rows(goqu.Record{
		"media_id": params.MediaId,

		"status":             params.Status,
		"start_date":         params.StartDate,
		"num_expected_parts": params.NumExpectedParts,
		"current_part":       params.CurrentPart,
		"next_airing":        params.NextAiring,
		"interval_days":      params.IntervalDays,
		"delay_days":         params.DelayDays,

		"created": params.Created,
		"updated": params.Updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type MediaPartReleaseChanges struct {
	Status           Change[types.MediaPartReleaseStatus]
	StartDate        Change[time.Time]
	NumExpectedParts Change[int]
	CurrentPart      Change[int]
	NextAiring       Change[time.Time]
	IntervalDays     Change[int]
	DelayDays        Change[int]

	Created Change[int64]
}

func (db *Database) UpdateMediaPartRelease(ctx context.Context, mediaId string, changes MediaPartReleaseChanges) error {
	record := goqu.Record{}

	addToRecord(record, "status", changes.Status)
	addToRecord(record, "start_date", changes.StartDate)
	addToRecord(record, "num_expected_parts", changes.NumExpectedParts)
	addToRecord(record, "current_part", changes.CurrentPart)
	addToRecord(record, "next_airing", changes.NextAiring)
	addToRecord(record, "interval_days", changes.IntervalDays)
	addToRecord(record, "delay_days", changes.IntervalDays)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("media_part_release").
		Set(record).
		Where(
			goqu.I("media_part_release.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveMediaPartRelease(ctx context.Context, mediaId string) error {
	query := dialect.Delete("media_part_release").
		Where(
			goqu.I("media_part_release.media_id").Eq(mediaId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
