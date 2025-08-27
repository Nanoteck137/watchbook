package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
)

type MediaPartRelease struct {
	RowId int `db:"rowid"`

	MediaId string `db:"media_id"`

	NumExpectedParts int    `db:"num_expected_parts"`
	CurrentPart      int    `db:"current_part"`
	NextAiring       string `db:"next_airing"`
	IntervalDays     int    `db:"interval_days"`
	IsActive         int    `db:"is_active"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func MediaPartReleaseQuery(userId *string) *goqu.SelectDataset {
	query := dialect.From("media_part_release").
		Select(
			"media_part_release.rowid",

			"media_part_release.media_id",

			"media_part_release.num_expected_parts",
			"media_part_release.current_part",
			"media_part_release.next_airing",
			"media_part_release.interval_days",
			"media_part_release.is_active",

			"media_part_release.created",
			"media_part_release.updated",
		)

	return query
}

func (db *Database) GetAllMediaPartReleases(ctx context.Context) ([]MediaPartRelease, error) {
	query := MediaPartReleaseQuery(nil)
	return ember.Multiple[MediaPartRelease](db.db, ctx, query)
}

func (db *Database) GetMediaPartReleaseById(ctx context.Context, mediaId string) (MediaPartRelease, error) {
	query := MediaPartReleaseQuery(nil).
		Where(
			goqu.I("media_part_release.media_id").Eq(mediaId),
		)

	return ember.Single[MediaPartRelease](db.db, ctx, query)
}

type CreateMediaPartReleaseParams struct {
	MediaId string

	NumExpectedParts int
	CurrentPart      int
	NextAiring       string
	IntervalDays     int
	IsActive         int

	Created int64
	Updated int64
}

func (db *Database) CreateMediaPartRelease(ctx context.Context, params CreateMediaPartReleaseParams) error {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	query := dialect.Insert("media_part_release").Rows(goqu.Record{
		"media_id": params.MediaId,

		"num_expected_parts": params.NumExpectedParts,
		"current_part":       params.CurrentPart,
		"next_airing":        params.NextAiring,
		"interval_days":      params.IntervalDays,
		"is_active":          params.IsActive,

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
	NumExpectedParts Change[int]
	CurrentPart      Change[int]
	NextAiring       Change[string]
	IntervalDays     Change[int]
	IsActive         Change[int]

	Created Change[int64]
}

func (db *Database) UpdateMediaPartRelease(ctx context.Context, mediaId string, changes MediaPartReleaseChanges) error {
	record := goqu.Record{}

	addToRecord(record, "num_expected_parts", changes.NumExpectedParts)
	addToRecord(record, "current_part", changes.CurrentPart)
	addToRecord(record, "next_airing", changes.NextAiring)
	addToRecord(record, "interval_days", changes.IntervalDays)
	addToRecord(record, "is_active", changes.IsActive)

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
