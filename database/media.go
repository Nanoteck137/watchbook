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

type MediaUserData struct {
	List         types.MediaUserList `json:"list"`
	Score        *int64              `json:"score"`
	Part         *int64              `json:"part"`
	RevisitCount *int64              `json:"revisit_count"`
	IsRevisiting int                 `json:"is_revisiting"`
	Updated      int                 `json:"updated"`
}

type MediaRelease struct {
	Type             string `json:"type"`
	StartDate        string `json:"start_date"`
	NumExpectedParts int    `json:"num_expected_parts"`
	PartOffset       int    `json:"part_offset"`
	IntervalDays     int    `json:"interval_days"`
	DelayDays        int    `json:"delay_days"`

	Created int `json:"created"`
	Updated int `json:"updated"`
}

type Media struct {
	RowId int `db:"rowid"`

	Id   string          `db:"id"`
	Type types.MediaType `db:"type"`

	TmdbId    sql.NullString `db:"tmdb_id"`
	ImdbId    sql.NullString `db:"imdb_id"`
	MalId     sql.NullString `db:"mal_id"`
	AnilistId sql.NullString `db:"anilist_id"`

	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`

	Score        sql.NullFloat64   `db:"score"`
	Status       types.MediaStatus `db:"status"`
	Rating       types.MediaRating `db:"rating"`
	AiringSeason sql.NullString    `db:"airing_season"`

	StartDate sql.NullString `db:"start_date"`
	EndDate   sql.NullString `db:"end_date"`

	CoverFile  sql.NullString `db:"cover_file"`
	LogoFile   sql.NullString `db:"logo_file"`
	BannerFile sql.NullString `db:"banner_file"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	PartCount sql.NullInt64 `db:"part_count"`

	Creators JsonColumn[[]string] `db:"creators"`
	Tags     JsonColumn[[]string] `db:"tags"`

	UserData JsonColumn[MediaUserData] `db:"user_data"`

	Release JsonColumn[MediaRelease] `db:"release"`
}

func MediaAiringSeasonQuery() *goqu.SelectDataset {
	tbl := goqu.T("tags")

	return dialect.From(tbl).
		Select(
			tbl.Col("slug").As("slug"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					tbl.Col("slug"),
					"name",
					tbl.Col("name"),
				),
			).As("data"),
		)
}

func MediaTagQuery() *goqu.SelectDataset {
	tbl := goqu.T("media_tags")

	return dialect.From(tbl).
		Select(
			tbl.Col("media_id").As("id"),
			goqu.Func(
				"json_group_array",
				tbl.Col("tag_slug"),
			).As("data"),
		).
		GroupBy(tbl.Col("media_id"))
}

func MediaCreatorQuery() *goqu.SelectDataset {
	tbl := goqu.T("media_creators")

	return dialect.From(tbl).
		Select(
			tbl.Col("media_id").As("id"),
			goqu.Func(
				"json_group_array",
				tbl.Col("tag_slug"),
			).As("data"),
		).
		GroupBy(tbl.Col("media_id"))
}

func MediaUserDataQuery(userId *string) *goqu.SelectDataset {
	tbl := goqu.T("media_user_data")

	query := dialect.From(tbl).
		Select(
			tbl.Col("media_id").As("id"),

			tbl.Col("list"),
			tbl.Col("part"),
			tbl.Col("revisit_count"),
			tbl.Col("is_revisiting"),
			tbl.Col("score"),
			tbl.Col("updated"),

			goqu.Func(
				"json_object",

				"list",
				tbl.Col("list"),

				"part",
				tbl.Col("part"),

				"revisit_count",
				tbl.Col("revisit_count"),

				"is_revisiting",
				tbl.Col("is_revisiting"),

				"score",
				tbl.Col("score"),

				"updated",
				tbl.Col("updated"),
			).As("data"),
		)

	if userId != nil {
		query = query.Where(tbl.Col("user_id").Eq(*userId))
	} else {
		query = query.Where(goqu.L("false"))
	}

	return query
}

func MediaReleaseQuery() *goqu.SelectDataset {
	tbl := goqu.T("media_part_release")

	query := dialect.From(tbl).
		Select(
			tbl.Col("media_id").As("id"),

			tbl.Col("type"),
			tbl.Col("start_date"),
			tbl.Col("num_expected_parts"),
			tbl.Col("part_offset"),
			tbl.Col("interval_days"),
			tbl.Col("delay_days"),

			tbl.Col("created"),
			tbl.Col("updated"),

			goqu.Func(
				"json_object",

				"type",
				tbl.Col("type"),

				"start_date",
				tbl.Col("start_date"),

				"num_expected_parts",
				tbl.Col("num_expected_parts"),

				"part_offset",
				tbl.Col("part_offset"),

				"interval_days",
				tbl.Col("interval_days"),

				"delay_days",
				tbl.Col("delay_days"),

				"created",
				tbl.Col("created"),

				"updated",
				tbl.Col("updated"),
			).As("data"),
		)

	return query
}

func MediaPartCountQuery() *goqu.SelectDataset {
	tbl := goqu.T("media_parts")

	return dialect.From(tbl).
		Select(
			tbl.Col("media_id").As("id"),
			goqu.COUNT(tbl.Col("idx")).As("data"),
		).
		GroupBy(tbl.Col("media_id"))
}

// TODO(patrik): Use goqu.T more
func MediaQuery(userId *string) *goqu.SelectDataset {
	partCountQuery := MediaPartCountQuery()
	creatorsQuery := MediaCreatorQuery()
	tagsQuery := MediaTagQuery()

	userDataQuery := MediaUserDataQuery(userId)
	releaseQuery := MediaReleaseQuery()

	query := dialect.From("media").
		Select(
			"media.rowid",

			"media.id",
			"media.type",

			"media.tmdb_id",
			"media.imdb_id",
			"media.mal_id",
			"media.anilist_id",

			"media.title",
			"media.description",

			"media.score",
			"media.status",
			"media.rating",
			"media.airing_season",

			"media.start_date",
			"media.end_date",

			"media.cover_file",
			"media.logo_file",
			"media.banner_file",

			"media.created",
			"media.updated",

			goqu.I("part_count.data").As("part_count"),

			goqu.I("creators.data").As("creators"),
			goqu.I("tags.data").As("tags"),

			goqu.I("user_data.data").As("user_data"),

			goqu.I("release.data").As("release"),
		).
		LeftJoin(
			partCountQuery.As("part_count"),
			goqu.On(goqu.I("media.id").Eq(goqu.I("part_count.id"))),
		).
		LeftJoin(
			creatorsQuery.As("creators"),
			goqu.On(goqu.I("media.id").Eq(goqu.I("creators.id"))),
		).
		LeftJoin(
			tagsQuery.As("tags"),
			goqu.On(goqu.I("media.id").Eq(goqu.I("tags.id"))),
		).
		LeftJoin(
			userDataQuery.As("user_data"),
			goqu.On(goqu.I("media.id").Eq(goqu.I("user_data.id"))),
		).
		LeftJoin(
			releaseQuery.As("release"),
			goqu.On(goqu.I("media.id").Eq(goqu.I("release.id"))),
		)

	return query
}

func (db *Database) GetAllMediaIds(ctx context.Context) ([]string, error) {
	query := dialect.From("media").
		Select("media.id")

	return ember.Multiple[string](db.db, ctx, query)
}

type FetchOptions struct {
	PerPage int
	Page    int
}

func (db *Database) GetPagedMedia(ctx context.Context, userId *string, filterStr, sortStr string, opts FetchOptions) ([]Media, types.Page, error) {
	query := MediaQuery(userId)

	var err error

	a := adapter.MediaResolverAdapter{}
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

	items, err := ember.Multiple[Media](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db *Database) GetAllMedia(ctx context.Context) ([]Media, error) {
	query := MediaQuery(nil)
	return ember.Multiple[Media](db.db, ctx, query)
}

func (db *Database) GetAllMediaForPartPredict(ctx context.Context) ([]Media, error) {
	query := MediaQuery(nil).
		Where(goqu.I("media.release_interval").Gt(0))
	return ember.Multiple[Media](db.db, ctx, query)
}

func (db *Database) GetMediaById(ctx context.Context, userId *string, id string) (Media, error) {
	query := MediaQuery(userId).
		Where(goqu.I("media.id").Eq(id))

	return ember.Single[Media](db.db, ctx, query)
}

func (db *Database) GetMediaByMalId(ctx context.Context, userId *string, malId string) (Media, error) {
	query := MediaQuery(userId).
		Where(goqu.I("media.mal_id").Eq(malId))

	return ember.Single[Media](db.db, ctx, query)
}

type CreateMediaParams struct {
	Id   string
	Type types.MediaType

	TmdbId    sql.NullString
	ImdbId    sql.NullString
	MalId     sql.NullString
	AnilistId sql.NullString

	Title string

	Description sql.NullString

	Score        sql.NullFloat64
	Status       types.MediaStatus
	Rating       types.MediaRating
	AiringSeason sql.NullString

	StartDate sql.NullString
	EndDate   sql.NullString

	CoverFile  sql.NullString
	LogoFile   sql.NullString
	BannerFile sql.NullString

	Created int64
	Updated int64
}

func (db *Database) CreateMedia(ctx context.Context, params CreateMediaParams) (string, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateMediaId()
	}

	if params.Type == "" {
		params.Type = types.MediaTypeUnknown
	}

	if params.Status == "" {
		params.Status = types.MediaStatusUnknown
	}

	if params.Rating == "" {
		params.Rating = types.MediaRatingUnknown
	}

	query := dialect.Insert("media").Rows(goqu.Record{
		"id":   id,
		"type": params.Type,

		"tmdb_id":    params.TmdbId,
		"imdb_id":    params.ImdbId,
		"mal_id":     params.MalId,
		"anilist_id": params.AnilistId,

		"title": params.Title,

		"description": params.Description,

		"score":         params.Score,
		"status":        params.Status,
		"rating":        params.Rating,
		"airing_season": params.AiringSeason,

		"start_date": params.StartDate,
		"end_date":   params.EndDate,

		"cover_file":  params.CoverFile,
		"logo_file":   params.LogoFile,
		"banner_file": params.BannerFile,

		"created": created,
		"updated": updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type MediaChanges struct {
	Type Change[types.MediaType]

	TmdbId    Change[sql.NullString]
	ImdbId    Change[sql.NullString]
	MalId     Change[sql.NullString]
	AnilistId Change[sql.NullString]

	Title Change[string]

	Description Change[sql.NullString]

	Score        Change[sql.NullFloat64]
	Status       Change[types.MediaStatus]
	Rating       Change[types.MediaRating]
	AiringSeason Change[sql.NullString]

	StartDate Change[sql.NullString]
	EndDate   Change[sql.NullString]

	CoverFile  Change[sql.NullString]
	LogoFile   Change[sql.NullString]
	BannerFile Change[sql.NullString]

	Created Change[int64]
}

func (db *Database) UpdateMedia(ctx context.Context, id string, changes MediaChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)

	addToRecord(record, "tmdb_id", changes.TmdbId)
	addToRecord(record, "imdb_id", changes.ImdbId)
	addToRecord(record, "mal_id", changes.MalId)
	addToRecord(record, "anilist_id", changes.AnilistId)

	addToRecord(record, "title", changes.Title)

	addToRecord(record, "description", changes.Description)

	addToRecord(record, "score", changes.Score)
	addToRecord(record, "status", changes.Status)
	addToRecord(record, "rating", changes.Rating)
	addToRecord(record, "airing_season", changes.AiringSeason)

	addToRecord(record, "start_date", changes.StartDate)
	addToRecord(record, "end_date", changes.EndDate)

	addToRecord(record, "cover_file", changes.CoverFile)
	addToRecord(record, "logo_file", changes.LogoFile)
	addToRecord(record, "banner_file", changes.BannerFile)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("media").
		Set(record).
		Where(goqu.I("media.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveMedia(ctx context.Context, id string) error {
	query := dialect.Delete("media").
		Where(goqu.I("media.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddTagToMedia(ctx context.Context, mediaId, tagSlug string) error {
	query := dialect.Insert("media_tags").
		Rows(goqu.Record{
			"media_id": mediaId,
			"tag_slug": tagSlug,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAllTagsFromMedia(ctx context.Context, mediaId string) error {
	query := dialect.Delete("media_tags").
		Where(goqu.I("media_tags.media_id").Eq(mediaId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddCreatorToMedia(ctx context.Context, mediaId, tagSlug string) error {
	query := dialect.Insert("media_creators").
		Rows(goqu.Record{
			"media_id": mediaId,
			"tag_slug": tagSlug,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAllCreatorsFromMedia(ctx context.Context, mediaId string) error {
	query := dialect.Delete("media_creators").
		Where(goqu.I("media_creators.media_id").Eq(mediaId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveMediaUserList(ctx context.Context, mediaId, userId string) error {
	query := dialect.Delete("media_user_list").
		Where(
			goqu.I("media_user_list.media_id").Eq(mediaId),
			goqu.I("media_user_list.user_id").Eq(userId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

const DefaultMediaUserList = types.MediaUserListBacklog

const (
	MediaScoreMin = 1
	MediaScoreMax = 10
)

type SetMediaUserData struct {
	List         types.MediaUserList
	Part         sql.NullInt64
	RevisitCount sql.NullInt64
	IsRevisiting bool
	Score        sql.NullInt64
}

func (db *Database) SetMediaUserData(ctx context.Context, mediaId, userId string, data SetMediaUserData) error {

	if data.List == "" {
		data.List = DefaultMediaUserList
	}

	if data.Score.Valid {
		data.Score.Int64 = utils.Clamp(data.Score.Int64, MediaScoreMin, MediaScoreMax)
	}

	updated := time.Now().UnixMilli()

	query := dialect.Insert("media_user_data").
		Rows(goqu.Record{
			"media_id": mediaId,
			"user_id":  userId,

			"list":          data.List,
			"part":          data.Part,
			"revisit_count": data.RevisitCount,
			"is_revisiting": data.IsRevisiting,
			"score":         data.Score,

			"updated": updated,
		}).
		OnConflict(
			goqu.DoUpdate("media_id, user_id", goqu.Record{
				"list":          data.List,
				"part":          data.Part,
				"revisit_count": data.RevisitCount,
				"is_revisiting": data.IsRevisiting,
				"score":         data.Score,

				"updated": updated,
			}),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteMediaUserData(ctx context.Context, mediaId, userId string) error {
	query := dialect.Delete("media_user_data").
		Where(
			goqu.I("media_user_data.media_id").Eq(mediaId),
			goqu.I("media_user_data.user_id").Eq(userId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

type SetMediaPartRelease struct {
	Type             types.MediaPartReleaseType
	StartDate        string
	NumExpectedParts int
	PartOffset       int
	IntervalDays     int
	DelayDays        int

	Created int64
	Updated int64
}

func (db *Database) SetMediaPartRelease(ctx context.Context, mediaId string, data SetMediaPartRelease) error {
	if data.Created == 0 && data.Updated == 0 {
		t := time.Now().UnixMilli()
		data.Created = t
		data.Updated = t
	}

	if !types.IsValidMediaPartReleaseType(data.Type) {
		data.Type = types.MediaPartReleaseTypeNotConfirmed
	}

	query := dialect.Insert("media_part_release").
		Rows(goqu.Record{
			"media_id": mediaId,

			"type":               data.Type,
			"start_date":         data.StartDate,
			"num_expected_parts": data.NumExpectedParts,
			"part_offset":        data.PartOffset,
			"interval_days":      data.IntervalDays,
			"delay_days":         data.DelayDays,

			"created": data.Created,
			"updated": data.Updated,
		}).
		OnConflict(
			goqu.DoUpdate("media_id", goqu.Record{
				"type":               data.Type,
				"start_date":         data.StartDate,
				"num_expected_parts": data.NumExpectedParts,
				"part_offset":        data.PartOffset,
				"interval_days":      data.IntervalDays,
				"delay_days":         data.DelayDays,

				"updated": data.Updated,
			}),
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
