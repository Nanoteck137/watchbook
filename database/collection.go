package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Collection struct {
	RowId int `db:"rowid"`

	Id   string               `db:"id"`
	Type types.CollectionType `db:"type"`

	Name string `db:"name"`

	AdminStatus types.AdminStatus `db:"admin_status"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func CollectionQuery(userId *string) *goqu.SelectDataset {
	query := dialect.From("collections").
		Select(
			"collections.rowid",

			"collections.id",
			"collections.type",

			"collections.name",

			"collections.admin_status",

			"collections.created",
			"collections.updated",
		)

	return query
}

// func (db *Database) GetPagedCollection(ctx context.Context, userId *string, filterStr, sortStr string, opts FetchOptions) ([]Collection, types.Page, error) {
// 	query := CollectionQuery(userId)
//
// 	var err error
//
// 	a := adapter.CollectionResolverAdapter{}
// 	resolver := filter.New(&a)
//
// 	query, err = applyFilter(query, resolver, filterStr)
// 	if err != nil {
// 		return nil, types.Page{}, err
// 	}
//
// 	query, err = applySort(query, resolver, sortStr)
// 	if err != nil {
// 		return nil, types.Page{}, err
// 	}
//
// 	countQuery := query.
// 		Select(goqu.COUNT("collections.id"))
//
// 	if opts.PerPage > 0 {
// 		query = query.
// 			Limit(uint(opts.PerPage)).
// 			Offset(uint(opts.Page * opts.PerPage))
// 	}
//
// 	totalItems, err := ember.Single[int](db.db, ctx, countQuery)
// 	if err != nil {
// 		return nil, types.Page{}, err
// 	}
//
// 	totalPages := utils.TotalPages(opts.PerPage, totalItems)
// 	page := types.Page{
// 		Page:       opts.Page,
// 		PerPage:    opts.PerPage,
// 		TotalItems: totalItems,
// 		TotalPages: totalPages,
// 	}
//
// 	items, err := ember.Multiple[Collection](db.db, ctx, query)
// 	if err != nil {
// 		return nil, types.Page{}, err
// 	}
//
// 	return items, page, nil
// }

func (db *Database) GetAllCollections(ctx context.Context) ([]Collection, error) {
	query := CollectionQuery(nil)
	return ember.Multiple[Collection](db.db, ctx, query)
}

func (db *Database) GetCollectionById(ctx context.Context, userId *string, id string) (Collection, error) {
	query := CollectionQuery(userId).
		Where(goqu.I("collections.id").Eq(id))

	return ember.Single[Collection](db.db, ctx, query)
}

type CreateCollectionParams struct {
	Id   string
	Type types.CollectionType

	Name string

	AdminStatus types.AdminStatus

	Created int64
	Updated int64
}

func (db *Database) CreateCollection(ctx context.Context, params CreateCollectionParams) (string, error) {
	if params.Created == 0 && params.Updated == 0 {
		t := time.Now().UnixMilli()
		params.Created = t
		params.Updated = t
	}

	if params.Id == "" {
		params.Id = utils.CreateCollectionId()
	}

	if params.Type == "" {
		params.Type = types.CollectionTypeUnknown
	}

	if params.AdminStatus == "" {
		params.AdminStatus = types.AdminStatusNotFixed
	}

	query := dialect.Insert("collections").Rows(goqu.Record{
		"id":   params.Id,
		"type": params.Type,

		"name": params.Name,

		"admin_status": params.AdminStatus,

		"created": params.Created,
		"updated": params.Updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type CollectionChanges struct {
	Type Change[types.CollectionType]

	Name Change[string]

	AdminStatus Change[types.AdminStatus]

	Created Change[int64]
}

func (db *Database) UpdateCollection(ctx context.Context, id string, changes CollectionChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "admin_status", changes.AdminStatus)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	query := dialect.Update("collections").
		Set(record).
		Where(goqu.I("collections.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveCollection(ctx context.Context, id string) error {
	query := dialect.Delete("collections").
		Where(goqu.I("collections.id").Eq(id))

	_, err := db.db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
