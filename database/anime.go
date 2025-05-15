package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

// TODO(patrik): Move
type JsonColumn[T any] struct {
	Has bool
	Val T
}

func (j *JsonColumn[T]) Scan(src any) error {
	var res T

	if src == nil {
		j.Val = res
		j.Has = false
		return nil
	}

	switch value := src.(type) {
	case string:
		err := json.Unmarshal([]byte(value), &j.Val)
		if err != nil {
			return err
		}

		j.Has = true
	case []byte:
		err := json.Unmarshal(value, &j.Val)
		if err != nil {
			return err
		}

		j.Has = true
	default:
		return fmt.Errorf("unsupported type %T", src)
	}

	return nil
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	raw, err := json.Marshal(j.Val)
	return raw, err
}

// func (j *JsonColumn[T]) Get() *T {
// 	return j.Val
// }

type AnimeStudio struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeProducer struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeTag struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type AnimeUserData struct {
	List         *types.AnimeUserList `json:"list"`
	Score        *int64               `json:"score"`
	Episode      *int64               `json:"episode"`
	IsRewatching int                  `json:"is_rewatching"`
}

type Anime struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	MalId string `db:"mal_id"`

	Title        string         `db:"title"`
	TitleEnglish sql.NullString `db:"title_english"`

	Description sql.NullString `db:"description"`

	Type         types.AnimeType   `db:"type"`
	Status       types.AnimeStatus `db:"status"`
	Rating       types.AnimeRating `db:"rating"`
	AiringSeason string            `db:"airing_season"`
	EpisodeCount sql.NullInt64     `db:"episode_count"`

	StartDate sql.NullString `db:"start_date"`
	EndDate   sql.NullString `db:"end_date"`

	Score sql.NullFloat64 `db:"score"`

	AniDBUrl            sql.NullString `db:"ani_db_url"`
	AnimeNewsNetworkUrl sql.NullString `db:"anime_news_network_url"`

	CoverFilename sql.NullString `db:"cover_filename"`

	ShouldFetchData   bool      `db:"should_fetch_data"`
	LastDataFetchDate time.Time `db:"last_data_fetch_date"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Studios      JsonColumn[[]AnimeStudio]   `db:"studios"`
	Producers    JsonColumn[[]AnimeProducer] `db:"producers"`
	Themes       JsonColumn[[]AnimeTag]      `db:"themes"`
	Genres       JsonColumn[[]AnimeTag]      `db:"genres"`
	Demographics JsonColumn[[]AnimeTag]      `db:"demographics"`

	UserData JsonColumn[AnimeUserData] `db:"user_data"`
}

func AnimeStudioQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_studios")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					goqu.I("studios.slug"),
					"name",
					goqu.I("studios.name"),
				),
			).As("studios"),
		).
		Join(
			goqu.I("studios"),
			goqu.On(tbl.Col("studio_slug").Eq(goqu.I("studios.slug"))),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeProducerQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_producers")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					goqu.I("producers.slug"),
					"name",
					goqu.I("producers.name"),
				),
			).As("producers"),
		).
		Join(
			goqu.I("producers"),
			goqu.On(tbl.Col("producer_slug").Eq(goqu.I("producers.slug"))),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeThemeQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_themes")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					goqu.I("tags.slug"),
					"name",
					goqu.I("tags.name"),
				),
			).As("data"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(tbl.Col("tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeGenreQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_genres")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					goqu.I("tags.slug"),
					"name",
					goqu.I("tags.name"),
				),
			).As("data"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(tbl.Col("tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeDemographicQuery() *goqu.SelectDataset {
	tbl := goqu.T("anime_demographics")

	return dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"slug",
					goqu.I("tags.slug"),
					"name",
					goqu.I("tags.name"),
				),
			).As("data"),
		).
		Join(
			goqu.I("tags"),
			goqu.On(tbl.Col("tag_slug").Eq(goqu.I("tags.slug"))),
		).
		GroupBy(tbl.Col("anime_id"))
}

func AnimeUserDataQuery(userId *string) *goqu.SelectDataset {
	tbl := goqu.T("anime_user_data")

	query := dialect.From(tbl).
		Select(
			tbl.Col("anime_id").As("id"),
			goqu.Func(
				"json_object",

				"list",
				tbl.Col("list"),

				"episode",
				tbl.Col("episode"),

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

// TODO(patrik): Use goqu.T more
func AnimeQuery(userId *string) *goqu.SelectDataset {
	studiosQuery := AnimeStudioQuery()
	producersQuery := AnimeProducerQuery()
	themesQuery := AnimeThemeQuery()
	genresQuery := AnimeGenreQuery()
	demographicsQuery := AnimeDemographicQuery()

	userDataQuery := AnimeUserDataQuery(userId)

	query := dialect.From("animes").
		Select(
			"animes.rowid",

			"animes.id",

			"animes.mal_id",

			"animes.title",
			"animes.title_english",

			"animes.description",

			"animes.type",
			"animes.status",
			"animes.rating",
			"animes.airing_season",
			"animes.episode_count",

			"animes.start_date",
			"animes.end_date",

			"animes.score",

			"animes.ani_db_url",
			"animes.anime_news_network_url",

			"animes.cover_filename",

			"animes.should_fetch_data",
			"animes.last_data_fetch_date",

			"animes.created",
			"animes.updated",

			goqu.I("studios.studios").As("studios"),
			goqu.I("producers.producers").As("producers"),
			goqu.I("themes.data").As("themes"),
			goqu.I("genres.data").As("genres"),
			goqu.I("demographics.data").As("demographics"),

			goqu.I("user_data.data").As("user_data"),
		).
		Prepared(true).
		LeftJoin(
			studiosQuery.As("studios"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("studios.id"))),
		).
		LeftJoin(
			producersQuery.As("producers"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("producers.id"))),
		).
		LeftJoin(
			themesQuery.As("themes"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("themes.id"))),
		).
		LeftJoin(
			genresQuery.As("genres"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("genres.id"))),
		).
		LeftJoin(
			demographicsQuery.As("demographics"),
			goqu.On(goqu.I("animes.id").Eq(goqu.I("demographics.id"))),
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

func (db *Database) GetPagedAnimes(ctx context.Context, userId *string, opts FetchOptions) ([]Anime, types.Page, error) {
	query := AnimeQuery(userId)

	var err error

	countQuery := query.
		Select(goqu.COUNT("animes.id"))

	if opts.PerPage > 0 {
		query = query.
			Limit(uint(opts.PerPage)).
			Offset(uint(opts.Page * opts.PerPage))
	}

	var totalItems int
	err = db.Get(&totalItems, countQuery)
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

	var items []Anime
	err = db.Select(&items, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db *Database) GetAllAnimes(ctx context.Context) ([]Anime, error) {
	query := AnimeQuery(nil)

	var items []Anime
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetAnimeIdsForFetching(ctx context.Context) ([]string, error) {
	query := AnimeQuery(nil).
		Select("animes.id").
		Where(goqu.I("animes.should_fetch_data").Eq(true))

	var items []string
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetAnimeById(ctx context.Context, userId *string, id string) (Anime, error) {
	query := AnimeQuery(userId).
		Where(goqu.I("animes.id").Eq(id))

	var item Anime
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Anime{}, ErrItemNotFound
		}

		return Anime{}, err
	}

	return item, nil
}

func (db *Database) GetAnimeByMalId(ctx context.Context, malId string) (Anime, error) {
	query := AnimeQuery(nil).
		Where(goqu.I("animes.mal_id").Eq(malId))

	var item Anime
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Anime{}, ErrItemNotFound
		}

		return Anime{}, err
	}

	return item, nil
}

type CreateAnimeParams struct {
	Id string

	MalId string

	Title        string
	TitleEnglish sql.NullString

	Description sql.NullString

	Type         types.AnimeType
	Status       types.AnimeStatus
	Rating       types.AnimeRating
	AiringSeason string
	EpisodeCount sql.NullInt64

	StartDate sql.NullString
	EndDate   sql.NullString

	Score sql.NullFloat64

	AniDBUrl            sql.NullString
	AnimeNewsNetworkUrl sql.NullString

	CoverFilename sql.NullString

	ShouldFetchData   bool
	LastDataFetchDate time.Time

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

	query := dialect.Insert("animes").Rows(goqu.Record{
		"id": id,

		"mal_id": params.MalId,

		"title":         params.Title,
		"title_english": params.TitleEnglish,

		"description": params.Description,

		"type":          params.Type,
		"status":        params.Status,
		"rating":        params.Rating,
		"airing_season": params.AiringSeason,
		"episode_count": params.EpisodeCount,

		"start_date": params.StartDate,
		"end_date":   params.EndDate,

		"score": params.Score,

		"ani_db_url":             params.AniDBUrl,
		"anime_news_network_url": params.AnimeNewsNetworkUrl,

		"should_fetch_data":    params.ShouldFetchData,
		"last_data_fetch_date": params.LastDataFetchDate,

		"cover_filename": params.CoverFilename,

		"created": created,
		"updated": updated,
	}).
		Prepared(true).
		Returning("id")

	var item string
	err := db.Get(&item, query)
	if err != nil {
		return "", err
	}

	return item, nil
}

type AnimeChanges struct {
	MalId Change[string]

	Title        Change[string]
	TitleEnglish Change[sql.NullString]

	Description Change[sql.NullString]

	Type         Change[types.AnimeType]
	Status       Change[types.AnimeStatus]
	Rating       Change[types.AnimeRating]
	AiringSeason Change[string]
	EpisodeCount Change[sql.NullInt64]

	StartDate Change[sql.NullString]
	EndDate   Change[sql.NullString]

	Score Change[sql.NullFloat64]

	AniDBUrl            Change[sql.NullString]
	AnimeNewsNetworkUrl Change[sql.NullString]

	CoverFilename Change[sql.NullString]

	ShouldFetchData   Change[bool]
	LastDataFetchDate Change[time.Time]

	Created Change[int64]
}

func (db *Database) UpdateAnime(ctx context.Context, id string, changes AnimeChanges) error {
	record := goqu.Record{}

	addToRecord(record, "mal_id", changes.MalId)

	addToRecord(record, "title", changes.Title)
	addToRecord(record, "title_english", changes.TitleEnglish)

	addToRecord(record, "description", changes.Description)

	addToRecord(record, "type", changes.Type)
	addToRecord(record, "status", changes.Status)
	addToRecord(record, "rating", changes.Rating)
	addToRecord(record, "airing_season", changes.AiringSeason)
	addToRecord(record, "episode_count", changes.EpisodeCount)

	addToRecord(record, "start_date", changes.StartDate)
	addToRecord(record, "end_date", changes.EndDate)

	addToRecord(record, "score", changes.Score)

	addToRecord(record, "ani_db_url", changes.AniDBUrl)
	addToRecord(record, "anime_news_network_url", changes.AnimeNewsNetworkUrl)

	addToRecord(record, "cover_filename", changes.CoverFilename)

	addToRecord(record, "should_fetch_data", changes.ShouldFetchData)
	addToRecord(record, "last_data_fetch_date", changes.LastDataFetchDate)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("animes").
		Prepared(true).
		Set(record).
		Where(goqu.I("animes.id").Eq(id))

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddThemeToAnime(ctx context.Context, animeId, themeSlug string) error {
	ds := dialect.Insert("anime_themes").
		Prepared(true).
		Rows(goqu.Record{
			"anime_id": animeId,
			"tag_slug": themeSlug,
		})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (db *Database) RemoveAllThemesFromAnime(ctx context.Context, animeId string) error {
	query := dialect.Delete("anime_themes").
		Where(goqu.I("anime_themes.anime_id").Eq(animeId)).
		Prepared(true)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddGenreToAnime(ctx context.Context, animeId, tagSlug string) error {
	ds := dialect.Insert("anime_genres").
		Prepared(true).
		Rows(goqu.Record{
			"anime_id": animeId,
			"tag_slug": tagSlug,
		})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (db *Database) RemoveAllGenresFromAnime(ctx context.Context, animeId string) error {
	query := dialect.Delete("anime_genres").
		Where(goqu.I("anime_genres.anime_id").Eq(animeId)).
		Prepared(true)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddDemographicToAnime(ctx context.Context, animeId, tagSlug string) error {
	ds := dialect.Insert("anime_demographics").
		Prepared(true).
		Rows(goqu.Record{
			"anime_id": animeId,
			"tag_slug": tagSlug,
		})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (db *Database) RemoveAllDemographicsFromAnime(ctx context.Context, animeId string) error {
	query := dialect.Delete("anime_demographics").
		Where(goqu.I("anime_demographics.anime_id").Eq(animeId)).
		Prepared(true)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddStudioToAnime(ctx context.Context, animeId, studioSlug string) error {
	ds := dialect.Insert("anime_studios").
		Prepared(true).
		Rows(goqu.Record{
			"anime_id":    animeId,
			"studio_slug": studioSlug,
		})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (db *Database) RemoveAllStudiosFromAnime(ctx context.Context, animeId string) error {
	query := dialect.Delete("anime_studios").
		Where(goqu.I("anime_studios.anime_id").Eq(animeId)).
		Prepared(true)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddProducerToAnime(ctx context.Context, animeId, producerSlug string) error {
	ds := dialect.Insert("anime_producers").
		Prepared(true).
		Rows(goqu.Record{
			"anime_id":      animeId,
			"producer_slug": producerSlug,
		})

	_, err := db.Exec(ctx, ds)
	if err != nil {
		var e sqlite3.Error
		if errors.As(err, &e) {
			if e.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (db *Database) RemoveAllProducersFromAnime(ctx context.Context, animeId string) error {
	query := dialect.Delete("anime_producers").
		Where(goqu.I("anime_producers.anime_id").Eq(animeId)).
		Prepared(true)

	_, err := db.Exec(ctx, query)
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

	_, err := db.Exec(ctx, query)
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
			"is_rewatching": data.IsRewatching,
			"score":         data.Score,
		}).
		OnConflict(
			goqu.DoUpdate("anime_id, user_id", goqu.Record{
				"list":          data.List,
				"episode":       data.Episode,
				"is_rewatching": data.IsRewatching,
				"score":         data.Score,
			}),
		)

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// const (
// 	AnimeScoreMin = 1
// 	AnimeScoreMax = 10
// )
//
// func (db *Database) RemoveAnimeUserScore(ctx context.Context, animeId, userId string) error {
// 	query := dialect.Delete("anime_user_score").
// 		Where(
// 			goqu.I("anime_user_score.anime_id").Eq(animeId),
// 			goqu.I("anime_user_score.user_id").Eq(userId),
// 		)
//
// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (db *Database) SetAnimeUserScore(ctx context.Context, animeId, userId string, score int) error {
// 	score = utils.Clamp(score, AnimeScoreMin, AnimeScoreMax)
//
// 	query := dialect.Insert("anime_user_score").
// 		Rows(goqu.Record{
// 			"anime_id": animeId,
// 			"user_id":  userId,
//
// 			"score": score,
// 		}).
// 		OnConflict(
// 			goqu.DoUpdate("anime_id, user_id", goqu.Record{
// 				"score": score,
// 			}),
// 		)
//
// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (db *Database) RemoveAnimeUserEpisode(ctx context.Context, animeId, userId string) error {
// 	query := dialect.Delete("anime_user_episode").
// 		Where(
// 			goqu.I("anime_user_episode.anime_id").Eq(animeId),
// 			goqu.I("anime_user_episode.user_id").Eq(userId),
// 		)
//
// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (db *Database) SetAnimeUserEpisode(ctx context.Context, animeId, userId string, episode int) error {
// 	query := dialect.Insert("anime_user_episode").
// 		Rows(goqu.Record{
// 			"anime_id": animeId,
// 			"user_id":  userId,
//
// 			"episode": episode,
// 		}).
// 		OnConflict(
// 			goqu.DoUpdate("anime_id, user_id", goqu.Record{
// 				"episode": episode,
// 			}),
// 		)
//
// 	_, err := db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
