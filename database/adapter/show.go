package adapter

import (
	"go/ast"

	"github.com/nanoteck137/watchbook/filter"
)

var _ filter.ResolverAdapter = (*ShowResolverAdapter)(nil)

type ShowResolverAdapter struct{}

func (a *ShowResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "shows.name", filter.SortTypeAsc
}

func (a *ShowResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "shows.id",
		}, true
	case "type":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "shows.type",
		}, true
	case "name":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "shows.name",
		}, true
	case "searchSlug":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "shows.search_slug",
		}, true
	case "created":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "shows.created",
		}, true
	case "updated":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "shows.updated",
		}, true
	}

	return filter.Name{}, false
}

func (a *ShowResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	// switch typ {
	// case "tags":
	// 	return utils.Slug(name), true
	// }

	return "", false
}

func (a *ShowResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
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

func (a *ShowResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	switch name {
	case "hasType":
		return resolver.In(name, "type", args)
	}

	return nil, filter.UnknownFunction(name)
}
