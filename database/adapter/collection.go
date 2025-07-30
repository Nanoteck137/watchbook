package adapter

import (
	"go/ast"

	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/utils"
)

var _ filter.ResolverAdapter = (*CollectionResolverAdapter)(nil)

type CollectionResolverAdapter struct{}

func (a *CollectionResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "collections.name", filter.SortTypeAsc
}

func (a *CollectionResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "collections.id",
		}, true
	case "name":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "collections.name",
		}, true
	// case "userList":
	// 	return filter.Name{
	// 		Kind: filter.NameKindString,
	// 		Name: "user_data.list",
	// 		Nullable: true,
	// 	}, true
	// case "lastDataFetch":
	// 	return filter.Name{
	// 		Kind: filter.NameKindNumber,
	// 		Name: "media.last_data_fetch",
	// 		Nullable: true,
	// 	}, true
	// case "airingSeason":
	// 	return filter.Name{
	// 		Kind: filter.NameKindString,
	// 		Name: "media.airing_season",
	// 		Nullable: true,
	// 	}, true
	// case "status":
	// 	return filter.Name{
	// 		Kind: filter.NameKindString,
	// 		Name: "media.status",
	// 	}, true
	}

	return filter.Name{}, false
}

func (a *CollectionResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	switch typ {
	case "tags":
		return utils.Slug(name), true
	}

	return "", false
}

func (a *CollectionResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
	// switch typ {
	// case "tags":
	// 	return filter.Table{
	// 		Name:       "media_tags",
	// 		SelectName: "media_id",
	// 		WhereName:  "tag_slug",
	// 	}, true
	// }

	return filter.Table{}, false
}

func (a *CollectionResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	// switch name {
	// case "hasTag":
	// 	return resolver.InTable(name, "tags", "media.id", args)
	// }

	return nil, filter.UnknownFunction(name)
}
