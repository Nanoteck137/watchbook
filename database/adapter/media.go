package adapter

import (
	"go/ast"

	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/utils"
)

var _ filter.ResolverAdapter = (*MediaResolverAdapter)(nil)

type MediaResolverAdapter struct{}

func (a *MediaResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "media.title", filter.SortTypeAsc
}

func (a *MediaResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "media.id",
		}, true
	case "malId":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "media.mal_id",
			Nullable: true,
		}, true
	case "title":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "media.title",
		}, true
	case "userList":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "user_data.list",
			Nullable: true,
		}, true
	case "airingSeason":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "media.airing_season",
			Nullable: true,
		}, true
	case "status":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "media.status",
		}, true
	}

	return filter.Name{}, false
}

func (a *MediaResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	switch typ {
	case "tags":
		return utils.Slug(name), true
	}

	return "", false
}

func (a *MediaResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
	switch typ {
	case "tags":
		return filter.Table{
			Name:       "media_tags",
			SelectName: "media_id",
			WhereName:  "tag_slug",
		}, true
	}

	return filter.Table{}, false
}

func (a *MediaResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	switch name {
	case "hasTag":
		return resolver.InTable(name, "tags", "media.id", args)
	}

	return nil, filter.UnknownFunction(name)
}
