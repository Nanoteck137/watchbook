package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type ShowSeasonPart struct {
	RowId int `db:"rowid"`

	ShowId       string `db:"show_id"`
	SeasonNumber int    `db:"season_num"`
	Index        int    `db:"idx"`

	Name        string         `db:"name"`
	ReleaseDate sql.NullString `db:"release_date"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func ShowSeasonPartQuery() *goqu.SelectDataset {
	query := dialect.From("show_season_parts").
		Select(
			"show_season_parts.rowid",

			"show_season_parts.show_id",
			"show_season_parts.season_num",
			"show_season_parts.idx",

			"show_season_parts.name",
			"show_season_parts.release_date",

			"show_season_parts.created",
			"show_season_parts.updated",
		)

	return query
}

// func (db DB) GetAllShowSeasonParts(ctx context.Context) ([]ShowSeasonPart, error) {
// 	query := ShowSeasonPartQuery()
// 	return ember.Multiple[ShowSeasonPart](db.db, ctx, query)
// }
//
// func (db DB) GetShowSeasonPartByIndexMediaId(ctx context.Context, index int64, mediaId string) (ShowSeasonPart, error) {
// 	query := ShowSeasonPartQuery().
// 		Where(
// 			goqu.I("show_season_parts.idx").Eq(index),
// 			goqu.I("show_season_parts.media_id").Eq(mediaId),
// 		)
//
// 	return ember.Single[ShowSeasonPart](db.db, ctx, query)
// }
//
// func (db DB) GetShowSeasonPartsByMediaId(ctx context.Context, mediaId string) ([]ShowSeasonPart, error) {
// 	query := ShowSeasonPartQuery().
// 		Where(goqu.I("show_season_parts.media_id").Eq(mediaId)).
// 		Order(goqu.I("show_season_parts.idx").Asc())
//
// 	return ember.Multiple[ShowSeasonPart](db.db, ctx, query)
// }

type CreateShowSeasonPartParams struct {
	ShowId       string
	SeasonNumber int
	Index        int

	Name        string
	ReleaseDate sql.NullString

	Created int64
	Updated int64
}

func (db DB) CreateShowSeasonPart(ctx context.Context, params CreateShowSeasonPartParams) error {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	query := dialect.Insert("show_season_parts").Rows(goqu.Record{
		"show_id":    params.ShowId,
		"season_num": params.SeasonNumber,
		"idx":        params.Index,

		"name":         params.Name,
		"release_date": params.ReleaseDate,

		"created": created,
		"updated": updated,
	})

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// type ShowSeasonPartChanges struct {
// 	Name        Change[string]
// 	ReleaseDate Change[sql.NullString]
//
// 	Created Change[int64]
// }
//
// func (db DB) UpdateShowSeasonPart(ctx context.Context, index int64, mediaId string, changes ShowSeasonPartChanges) error {
// 	record := goqu.Record{}
//
// 	addToRecord(record, "name", changes.Name)
// 	addToRecord(record, "release_date", changes.ReleaseDate)
//
// 	addToRecord(record, "created", changes.Created)
//
// 	if len(record) == 0 {
// 		return nil
// 	}
//
// 	record["updated"] = time.Now().UnixMilli()
//
// 	query := dialect.Update("show_season_parts").
// 		Set(record).
// 		Where(
// 			goqu.I("show_season_parts.idx").Eq(index),
// 			goqu.I("show_season_parts.media_id").Eq(mediaId),
// 		)
//
// 	_, err := db.db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// func (db DB) RemoveShowSeasonPart(ctx context.Context, idx int64, mediaId string) error {
// 	query := dialect.Delete("show_season_parts").
// 		Where(
// 			goqu.I("show_season_parts.idx").Eq(idx),
// 			goqu.I("show_season_parts.media_id").Eq(mediaId),
// 		)
//
// 	_, err := db.db.Exec(ctx, query)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (db DB) RemoveAllShowSeasonParts(ctx context.Context, showId string, seasonNumber int) error {
	query := dialect.Delete("show_season_parts").
		Where(
			goqu.I("show_season_parts.show_id").Eq(showId),
			goqu.I("show_season_parts.season_num").Eq(seasonNumber),
		)

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
