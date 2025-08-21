package adapter

import (
	"go/ast"

	"github.com/nanoteck137/watchbook/filter"
	"github.com/nanoteck137/watchbook/utils"
)

var _ filter.ResolverAdapter = (*NotificationResolverAdapter)(nil)

type NotificationResolverAdapter struct{}

func (a *NotificationResolverAdapter) DefaultSort() (string, filter.SortType) {
	return "notifications.created", filter.SortTypeAsc
}

func (a *NotificationResolverAdapter) ResolveVariableName(name string) (filter.Name, bool) {
	switch name {
	case "id":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "notifications.id",
		}, true
	case "type":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "notifications.type",
		}, true
	case "title":
		return filter.Name{
			Kind: filter.NameKindString,
			Name: "notifications.title",
		}, true
	case "isRead":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "notifications.is_read",
		}, true
	case "created":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "notifications.created",
		}, true
	case "updated":
		return filter.Name{
			Kind: filter.NameKindNumber,
			Name: "notifications.updated",
		}, true
	}

	return filter.Name{}, false
}

func (a *NotificationResolverAdapter) ResolveNameToId(typ, name string) (string, bool) {
	switch typ {
	case "tags":
		return utils.Slug(name), true
	}

	return "", false
}

func (a *NotificationResolverAdapter) ResolveTable(typ string) (filter.Table, bool) {
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

func (a *NotificationResolverAdapter) ResolveFunctionCall(resolver *filter.Resolver, name string, args []ast.Expr) (filter.FilterExpr, error) {
	// switch name {
	// case "hasTag":
	// 	return resolver.InTable(name, "tags", "media.id", args)
	// }

	return nil, filter.UnknownFunction(name)
}
