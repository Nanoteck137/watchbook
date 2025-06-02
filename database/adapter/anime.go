package adapter

import (
	"go/ast"

	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/utils"
)

var _ filter.ResolverAdapter = (*TrackResolverAdapter)(nil)

type TrackResolverAdapter struct{}

func (a *TrackResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "animes.title", filter.SortTypeAsc
}

func (a *TrackResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "animes.id",
		}, true
	case "userList":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_data.list",
			Nullable: true,
		}, true
	case "lastDataFetch":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "animes.last_data_fetch",
			Nullable: true,
		}, true
	case "airingSeason":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "animes.airing_season",
			Nullable: true,
		}, true
	case "status":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "animes.status",
		}, true

		// mal_id TEXT,
		// ani_db_id TEXT,
		// anilist_id TEXT,
		// anime_news_network_id TEXT,
		//
		// title TEXT NOT NULL CHECK(title<>''),
		// title_english TEXT,
		//
		// description TEXT,
		//
		// type TEXT NOT NULL,
		// status TEXT NOT NULL,
		// rating TEXT NOT NULL,
		// airing_season TEXT NOT NULL,
		// episode_count INTEGER,
		//
		// start_date TEXT, 
		// end_date TEXT,
		//
		// score FLOAT,
		//
		// last_data_fetch_date INTEGER,
		//
		// created INTEGER NOT NULL,
		// updated INTEGER NOT NULL
	}

	return filter.Name{}, false
}

func (a *TrackResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	switch typ {
	case "tags":
		return utils.Slug(name), true
	}

	return "", false
}

func (a *TrackResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
	switch typ {
	case "tags":
		return filter.Table{
			Name:       "anime_tags",
			SelectName: "anime_id",
			WhereName:  "tag_slug",
		}, true
	}

	return filter.Table{}, false
}

func (a *TrackResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	switch name {
	case "hasTag":
		return resolver.InTable(name, "tags", "animes.id", args)
	}

	return nil, filter.UnknownFunction(name)
}
