package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/utils"
)

type AnimeEpisode struct {
	RowId int `db:"rowid"`

	Id      string `db:"id"`
	AnimeId string `db:"anime_id"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func AnimeEpisodeQuery() *goqu.SelectDataset {
	query := dialect.From("anime_episodes").
		Select(
			"animes.rowid",

			"animes.created",
			"animes.updated",
		)

	return query
}

func (db *Database) GetAllAnimeEpisodes(ctx context.Context) ([]AnimeEpisode, error) {
	query := AnimeEpisodeQuery()
	return ember.Multiple[AnimeEpisode](db.db, ctx, query)
}

func (db *Database) GetAnimeEpisodeById(ctx context.Context, id string) (AnimeEpisode, error) {
	query := AnimeEpisodeQuery().
		Where(goqu.I("anime_episodes.id").Eq(id))

	return ember.Single[AnimeEpisode](db.db, ctx, query)
}

func (db *Database) GetAnimeEpisodesByAnimeId(ctx context.Context, animeId string) ([]AnimeEpisode, error) {
	query := AnimeEpisodeQuery().
		Where(goqu.I("anime_episodes.anime_id").Eq(animeId))

	return ember.Multiple[AnimeEpisode](db.db, ctx, query)
}

type CreateAnimeEpisodeParams struct {
	Id      string
	AnimeId string

	Name  string
	Index int64

	Created int64
	Updated int64
}

func (db *Database) CreateAnimeEpisode(ctx context.Context, params CreateAnimeEpisodeParams) (string, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateAnimeEpisodeId()
	}

	query := dialect.Insert("anime_episodes").Rows(goqu.Record{
		"id":       id,
		"anime_id": params.AnimeId,

		"name": params.Name,
		"idx":  params.Index,

		"created": created,
		"updated": updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type AnimeEpisodeChanges struct {
	AnimeId Change[string]

	Name  Change[string]
	Index Change[int64]

	Created Change[int64]
}

func (db *Database) UpdateAnimeEpisode(ctx context.Context, id string, changes AnimeEpisodeChanges) error {
	record := goqu.Record{}

	addToRecord(record, "anime_id", changes.AnimeId)

	addToRecord(record, "name", changes.AnimeId)
	addToRecord(record, "index", changes.Index)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("anime_episodes").
		Set(record).
		Where(goqu.I("animes.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveAnimeEpisode(ctx context.Context, id string) error {
	query := dialect.Delete("anime_episodes").
		Where(goqu.I("anime_episodes.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
