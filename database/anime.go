package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/database/adapter"
	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

// TODO(patrik): Move
type JsonColumn[T any] struct {
	Data  T
	Valid bool
}

func (j *JsonColumn[T]) Scan(src any) error {
	var res T

	if src == nil {
		j.Data = res
		j.Valid = false
		return nil
	}

	switch value := src.(type) {
	case string:
		err := json.Unmarshal([]byte(value), &j.Data)
		if err != nil {
			return err
		}

		j.Valid = true
	case []byte:
		err := json.Unmarshal(value, &j.Data)
		if err != nil {
			return err
		}

		j.Valid = true
	default:
		return fmt.Errorf("unsupported type %T", src)
	}

	return nil
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	raw, err := json.Marshal(j.Data)
	return raw, err
}

// func (j *JsonColumn[T]) Get() *T {
// 	return j.Val
// }

type AnimeUserData struct {
	List         *types.AnimeUserList `json:"list"`
	Score        *int64               `json:"score"`
	RewatchCount *int64               `json:"rewatch_count"`
	Episode      *int64               `json:"episode"`
	IsRewatching int                  `json:"is_rewatching"`
}

type AnimeImageJson struct {
	Filename string               `json:"filename"`
	Type     types.EntryImageType `json:"type"`
}

type Anime struct {
	RowId int `db:"rowid"`

	Id   string          `db:"id"`
	Type types.AnimeType `db:"type"`

	TmdbId    sql.NullString `db:"tmdb_id"`
	MalId     sql.NullString `db:"mal_id"`
	AnilistId sql.NullString `db:"anilist_id"`

	Title        string         `db:"title"`
	TitleEnglish sql.NullString `db:"title_english"`

	Description sql.NullString `db:"description"`

	Score        sql.NullFloat64   `db:"score"`
	Status       types.AnimeStatus `db:"status"`
	Rating       types.AnimeRating `db:"rating"`
	AiringSeason sql.NullString    `db:"airing_season"`

	StartDate sql.NullString `db:"start_date"`
	EndDate   sql.NullString `db:"end_date"`

	AdminStatus types.EntryAdminStatus `db:"admin_status"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	EpisodeCount int64 `db:"episode_count"`

	Studios JsonColumn[[]string]         `db:"studios"`
	Tags    JsonColumn[[]string]         `db:"tags"`
	Images  JsonColumn[[]AnimeImageJson] `db:"images"`

	UserData JsonColumn[AnimeUserData] `db:"user_data"`
}

func AnimeAiringSeasonQuery() *goqu.SelectDataset {
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

func AnimeTagQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_tags")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				tbl.Col("tag_slug"),
			).As("data"),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeStudioQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_studios")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				tbl.Col("tag_slug"),
			).As("data"),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeImageJsonQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_images")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"filename",
					tbl.Col("filename"),

					"type",
					tbl.Col("type"),
				),
			).As("data"),
		).
		Where(tbl.Col("is_primary").Gt(0)).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeUserDataQuery(userId *string) *goqu.SelectDataset {
	tbl := goqu.T("anime_user_data")

	query := dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),

			tbl.Col("list"),
			tbl.Col("episode"),
			tbl.Col("rewatch_count"),
			tbl.Col("is_rewatching"),
			tbl.Col("score"),

			goqu.Func(
				"json_object",

				"list",
				tbl.Col("list"),

				"episode",
				tbl.Col("episode"),

				"rewatch_count",
				tbl.Col("rewatch_count"),

				"is_rewatching",
				tbl.Col("is_rewatching"),

				"score",
				tbl.Col("score"),
			).As("data"),
		)

	if userId != nil {
		query = query.Where(tbl.Col("user_id").Eq(*userId))
	} else {
		query = query.Where(goqu.L("false"))
	}

	return query
}

func AnimeEpisodeCountQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_episodes")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.COUNT(tbl.Col("idx")).As("data"),
		).
		GroupBy(tbl.Col("anime_id"))
}

// TODO(patrik): Use goqu.T more
func AnimeQuery(userId *string) *goqu.SelectDataset {
	episodeCountQuery := AnimeEpisodeCountQuery()
	studiosQuery := AnimeStudioQuery()
	tagsQuery := AnimeTagQuery()
	imagesQuery := AnimeImageJsonQuery()

	userDataQuery := AnimeUserDataQuery(userId)

	query := dialect.From("animes").
		Select(
			"animes.rowid",

			"animes.id",
			"animes.type",

			"animes.tmdb_id",
			"animes.mal_id",
			"animes.anilist_id",

			"animes.title",
			"animes.description",

			"animes.score",
			"animes.status",
			"animes.rating",
			"animes.airing_season",

			"animes.start_date",
			"animes.end_date",

			"animes.admin_status",

			"animes.created",
			"animes.updated",

			goqu.I("episode_count.data").As("episode_count"),

			goqu.I("studios.data").As("studios"),
			goqu.I("tags.data").As("tags"),
			goqu.I("images.data").As("images"),

			goqu.I("user_data.data").As("user_data"),
		).
		LeftJoin(
			episodeCountQuery.As("episode_count"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("episode_count.id"))),
		).
		LeftJoin(
			studiosQuery.As("studios"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("studios.id"))),
		).
		LeftJoin(
			tagsQuery.As("tags"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("tags.id"))),
		).
		LeftJoin(
			imagesQuery.As("images"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("images.id"))),
		).
		LeftJoin(
			userDataQuery.As("user_data"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("user_data.id"))),
		)

	return query
}

type FetchOptions struct {
	PerPage int
	Page    int
}

func (db *Database) GetPagedAnimes(ctx context.Context, userId *string, filterStr, sortStr string, opts FetchOptions) ([]Anime, types.Page, error) {
	query := AnimeQuery(userId)

	var err error

	a := adapter.TrackResolverAdapter{}
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
		Select(goqu.COUNT("animes.id"))

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

	items, err := ember.Multiple[Anime](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db *Database) GetAllAnimes(ctx context.Context) ([]Anime, error) {
	query := AnimeQuery(nil)
	return ember.Multiple[Anime](db.db, ctx, query)
}

func (db *Database) GetAnimeById(ctx context.Context, userId *string, id string) (Anime, error) {
	query := AnimeQuery(userId).
		Where(goqu.I("animes.id").Eq(id))

	return ember.Single[Anime](db.db, ctx, query)
}

func (db *Database) GetAnimeByMalId(ctx context.Context, userId *string, malId string) (Anime, error) {
	query := AnimeQuery(userId).
		Where(goqu.I("animes.mal_id").Eq(malId))

	return ember.Single[Anime](db.db, ctx, query)
}

type CreateAnimeParams struct {
	Id   string
	Type types.AnimeType

	TmdbId    sql.NullString
	MalId     sql.NullString
	AnilistId sql.NullString

	Title string

	Description sql.NullString

	Score        sql.NullFloat64
	Status       types.AnimeStatus
	Rating       types.AnimeRating
	AiringSeason sql.NullString

	StartDate sql.NullString
	EndDate   sql.NullString

	AdminStatus types.EntryAdminStatus

	Created int64
	Updated int64
}

func (db *Database) CreateAnime(ctx context.Context, params CreateAnimeParams) (string, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateAnimeId()
	}

	if params.Type == "" {
		params.Type = types.AnimeTypeUnknown
	}

	if params.Status == "" {
		params.Status = types.AnimeStatusUnknown
	}

	if params.Rating == "" {
		params.Rating = types.AnimeRatingUnknown
	}

	if params.AdminStatus == "" {
		params.AdminStatus = types.EntryAdminStatusNotFixed
	}

	query := dialect.Insert("animes").Rows(goqu.Record{
		"id":   id,
		"type": params.Type,

		"tmdb_id":    params.TmdbId,
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

		"admin_status": params.AdminStatus,

		"created": created,
		"updated": updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type AnimeChanges struct {
	Type Change[types.AnimeType]

	TmdbId    Change[sql.NullString]
	MalId     Change[sql.NullString]
	AnilistId Change[sql.NullString]

	Title Change[string]

	Description Change[sql.NullString]

	Score        Change[sql.NullFloat64]
	Status       Change[types.AnimeStatus]
	Rating       Change[types.AnimeRating]
	AiringSeason Change[sql.NullString]

	StartDate Change[sql.NullString]
	EndDate   Change[sql.NullString]

	AdminStatus Change[types.EntryAdminStatus]

	Created Change[int64]
}

func (db *Database) UpdateAnime(ctx context.Context, id string, changes AnimeChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)

	addToRecord(record, "tmdb_id", changes.TmdbId)
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

	addToRecord(record, "admin_status", changes.AdminStatus)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("animes").
		Set(record).
		Where(goqu.I("animes.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAnime(ctx context.Context, id string) error {
	query := dialect.Delete("animes").
		Where(goqu.I("animes.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddTagToAnime(ctx context.Context, animeId, tagSlug string) error {
	query := dialect.Insert("anime_tags").
		Rows(goqu.Record{
			"anime_id": animeId,
			"tag_slug": tagSlug,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAllTagsFromAnime(ctx context.Context, animeId string) error {
	query := dialect.Delete("anime_tags").
		Where(goqu.I("anime_tags.anime_id").Eq(animeId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddStudioToAnime(ctx context.Context, animeId, tagSlug string) error {
	query := dialect.Insert("anime_studios").
		Rows(goqu.Record{
			"anime_id": animeId,
			"tag_slug": tagSlug,
		})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAllStudiosFromAnime(ctx context.Context, animeId string) error {
	query := dialect.Delete("anime_studios").
		Where(goqu.I("anime_studios.anime_id").Eq(animeId))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAnimeUserList(ctx context.Context, animeId, userId string) error {
	query := dialect.Delete("anime_user_list").
		Where(
			goqu.I("anime_user_list.anime_id").Eq(animeId),
			goqu.I("anime_user_list.user_id").Eq(userId),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

const DefaultAnimeUserList = types.AnimeUserListWatching

const (
	AnimeScoreMin = 1
	AnimeScoreMax = 10
)

type SetAnimeUserData struct {
	List         sql.NullString
	Episode      sql.NullInt64
	RewatchCount sql.NullInt64
	IsRewatching bool
	Score        sql.NullInt64
}

func (db *Database) SetAnimeUserData(ctx context.Context, animeId, userId string, data SetAnimeUserData) error {

	if data.List.Valid {
		switch types.AnimeUserList(data.List.String) {
		case types.AnimeUserListWatching,
			types.AnimeUserListCompleted,
			types.AnimeUserListOnHold,
			types.AnimeUserListDropped,
			types.AnimeUserListPlanToWatch:
		default:
			data.List.String = string(DefaultAnimeUserList)
		}
	}

	if data.Score.Valid {
		data.Score.Int64 = utils.Clamp(data.Score.Int64, AnimeScoreMin, AnimeScoreMax)
	}

	query := dialect.Insert("anime_user_data").
		Rows(goqu.Record{
			"anime_id": animeId,
			"user_id":  userId,

			"list":          data.List,
			"episode":       data.Episode,
			"rewatch_count": data.RewatchCount,
			"is_rewatching": data.IsRewatching,
			"score":         data.Score,
		}).
		OnConflict(
			goqu.DoUpdate("anime_id, user_id", goqu.Record{
				"list":          data.List,
				"episode":       data.Episode,
				"rewatch_count": data.RewatchCount,
				"is_rewatching": data.IsRewatching,
				"score":         data.Score,
			}),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
