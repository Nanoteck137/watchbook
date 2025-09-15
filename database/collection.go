package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/watchbook/database/adapter"
	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/kvstore"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Collection struct {
	RowId int `db:"rowid"`

	Id   string               `db:"id"`
	Type types.CollectionType `db:"type"`

	Name string `db:"name"`

	CoverFile  sql.NullString `db:"cover_file"`
	LogoFile   sql.NullString `db:"logo_file"`
	BannerFile sql.NullString `db:"banner_file"`

	Providers kvstore.Store `db:"providers"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

// TODO(patrik): Use goqu.T more
func CollectionQuery() *goqu.SelectDataset {
	query := dialect.From("collections").
		Select(
			"collections.rowid",

			"collections.id",
			"collections.type",

			"collections.name",

			"collections.cover_file",
			"collections.logo_file",
			"collections.banner_file",

			"collections.providers",

			"collections.created",
			"collections.updated",
		)

	return query
}

func (db *Database) GetPagedCollections(ctx context.Context, filterStr, sortStr string, opts FetchOptions) ([]Collection, types.Page, error) {
	query := CollectionQuery()

	var err error

	a := adapter.CollectionResolverAdapter{}
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
		Select(goqu.COUNT("collections.id"))

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

	items, err := ember.Multiple[Collection](db.db, ctx, query)
	if err != nil {
		return nil, types.Page{}, err
	}

	return items, page, nil
}

func (db *Database) GetAllCollections(ctx context.Context) ([]Collection, error) {
	query := CollectionQuery()
	return ember.Multiple[Collection](db.db, ctx, query)
}

func (db *Database) GetCollectionById(ctx context.Context, id string) (Collection, error) {
	query := CollectionQuery().
		Where(goqu.I("collections.id").Eq(id))

	return ember.Single[Collection](db.db, ctx, query)
}

func (db *Database) GetCollectionByProviderId(ctx context.Context, providerName, value string) (Collection, error) {
	query := CollectionQuery().
		Where(
			goqu.Func("json_extract", goqu.I("collections.providers"), "$."+providerName).Eq(value),
		)

	return ember.Single[Collection](db.db, ctx, query)
}

type CreateCollectionParams struct {
	Id   string
	Type types.CollectionType

	Name string

	CoverFile  sql.NullString
	LogoFile   sql.NullString
	BannerFile sql.NullString

	Providers kvstore.Store

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

	query := dialect.Insert("collections").Rows(goqu.Record{
		"id":   params.Id,
		"type": params.Type,

		"name": params.Name,

		"cover_file":  params.CoverFile,
		"logo_file":   params.LogoFile,
		"banner_file": params.BannerFile,

		"providers": params.Providers,

		"created": params.Created,
		"updated": params.Updated,
	}).
		Returning("id")

	return ember.Single[string](db.db, ctx, query)
}

type CollectionChanges struct {
	Type Change[types.CollectionType]

	Name Change[string]

	CoverFile  Change[sql.NullString]
	LogoFile   Change[sql.NullString]
	BannerFile Change[sql.NullString]

	Providers kvstore.Store

	Created Change[int64]
}

func (db *Database) UpdateCollection(ctx context.Context, id string, changes CollectionChanges) error {
	record := goqu.Record{}

	addToRecord(record, "type", changes.Type)

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "cover_file", changes.CoverFile)
	addToRecord(record, "logo_file", changes.LogoFile)
	addToRecord(record, "banner_file", changes.BannerFile)

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
