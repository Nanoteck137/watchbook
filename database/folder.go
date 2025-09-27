package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/utils"
)

type Folder struct {
	RowId int `db:"rowid"`

	Id string `db:"id"`

	UserId string `db:"user_id"`

	Name string `db:"name"`

	CoverFile sql.NullString `db:"cover_file"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func FolderQuery() *goqu.SelectDataset {
	query := dialect.From("folders").
		Select(
			"folders.rowid",

			"folders.id",

			"folders.user_id",

			"folders.name",

			"folders.cover_file",

			"folders.created",
			"folders.updated",
		)

	return query
}

func (db DB) GetAllFolders(ctx context.Context) ([]Folder, error) {
	query := FolderQuery()
	return ember.Multiple[Folder](db.db, ctx, query)
}

func (db DB) GetAllFoldersByUserId(ctx context.Context, userId string) ([]Folder, error) {
	query := FolderQuery().Where(goqu.I("folders.user_id").Eq(userId))
	return ember.Multiple[Folder](db.db, ctx, query)
}

func (db DB) GetFolderById(ctx context.Context, id string) (Folder, error) {
	query := FolderQuery().
		Where(goqu.I("folders.id").Eq(id))

	return ember.Single[Folder](db.db, ctx, query)
}

type CreateFolderParams struct {
	Id   string

	UserId string

	Name string

	CoverFile  sql.NullString

	Created int64
	Updated int64
}

func (db DB) CreateFolder(ctx context.Context, params CreateFolderParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateFolderId()
	}

	query := dialect.Insert("folders").Rows(goqu.Record{
		"id":   params.Id,
		
		"user_id": params.UserId,

		"name": params.Name,

		"cover_file":  params.CoverFile,

		"created": params.Created,
		"updated": params.Updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type FolderChanges struct {
	Name Change[string]

	CoverFile  Change[sql.NullString]

	Created Change[int64]
}

func (db DB) UpdateFolder(ctx context.Context, id string, changes FolderChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "cover_file", changes.CoverFile)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("folders").
		Set(record).
		Where(goqu.I("folders.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) RemoveFolder(ctx context.Context, id string) error {
	query := dialect.Delete("folders").
		Where(goqu.I("folders.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
