package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
)

type Anime struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	Title        string         `db:"title"`
	TitleEnglish sql.NullString `db:"title_english"`

	Description string `db:"description"`

	Type         string        `db:"type"`
	Status       string        `db:"status"`
	Rating       string        `db:"rating"`
	AiringSeason string        `db:"airing_season"`
	EpisodeCount sql.NullInt64 `db:"episode_count"`

	StartDate sql.NullString `db:"start_date"`
	EndDate   sql.NullString `db:"end_date"`

	Score sql.NullFloat64 `db:"score"`

	AniDBUrl            sql.NullString `db:"ani_db_url"`
	AnimeNewsNetworkUrl sql.NullString `db:"anime_news_network_url"`

	DownloadDate time.Time `db:"download_date"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func AnimeQuery() *goqu.SelectDataset {
	query := dialect.From("animes").
		Select(
			"animes.rowid",

			"animes.id",

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

			"animes.download_date",

			"animes.created",
			"animes.updated",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllAnimes(ctx context.Context) ([]Anime, error) {
	query := AnimeQuery()

	var items []Anime
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetAnimeById(ctx context.Context, id string) (Anime, error) {
	query := AnimeQuery().
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

type CreateAnimeParams struct {
	Id string

	Title        string
	TitleEnglish sql.NullString

	Description string

	Type         string
	Status       string
	Rating       string
	AiringSeason string
	EpisodeCount sql.NullInt64

	StartDate sql.NullString
	EndDate   sql.NullString

	Score sql.NullFloat64

	AniDBUrl            sql.NullString
	AnimeNewsNetworkUrl sql.NullString

	DownloadDate time.Time

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

	// id := params.Id
	// if id == "" {
	// 	id = utils.CreateId()
	// }

	query := dialect.Insert("animes").Rows(goqu.Record{
		"id": params.Id,

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

		"download_date": params.DownloadDate,

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
	Title        Change[string]
	TitleEnglish Change[sql.NullString]

	Description Change[string]

	Type         Change[string]
	Status       Change[string]
	Rating       Change[string]
	AiringSeason Change[string]
	EpisodeCount Change[sql.NullInt64]

	StartDate Change[sql.NullString]
	EndDate   Change[sql.NullString]

	Score Change[sql.NullFloat64]

	AniDBUrl            Change[sql.NullString]
	AnimeNewsNetworkUrl Change[sql.NullString]

	DownloadDate Change[time.Time]

	Created Change[int64]
}

func (db *Database) UpdateAnime(ctx context.Context, id string, changes AnimeChanges) error {
	record := goqu.Record{}

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

	addToRecord(record, "download_date", changes.DownloadDate)

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
			"theme_slug": themeSlug,
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

func (db *Database) AddGenreToAnime(ctx context.Context, animeId, genreSlug string) error {
	ds := dialect.Insert("anime_genres").
		Prepared(true).
		Rows(goqu.Record{
			"anime_id": animeId,
			"genre_slug": genreSlug,
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

func (db *Database) AddStudioToAnime(ctx context.Context, animeId, studioSlug string) error {
	ds := dialect.Insert("anime_studios").
		Prepared(true).
		Rows(goqu.Record{
			"anime_id": animeId,
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
			"anime_id": animeId,
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
