package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"sort"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Collection struct {
	Id             string               `json:"id"`
	CollectionType types.CollectionType `json:"collectionType"`

	CoverUrl  *string `json:"coverUrl"`
	LogoUrl   *string `json:"logoUrl"`
	BannerUrl *string `json:"bannerUrl"`

	Name string `json:"name"`
}

type GetCollections struct {
	Page        types.Page   `json:"page"`
	Collections []Collection `json:"collections"`
}

type GetCollectionById struct {
	Collection
}

func ConvertDBCollection(c pyrin.Context, hasUser bool, collection database.Collection) Collection {
	// TODO(patrik): Add default cover
	var coverUrl *string
	var bannerUrl *string
	var logoUrl *string

	if collection.CoverFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/collections/%s/%s", collection.Id, path.Base(collection.CoverFile.String)))
		coverUrl = &url
	}

	if collection.LogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/collections/%s/%s", collection.Id, path.Base(collection.LogoFile.String)))
		logoUrl = &url
	}

	if collection.BannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/collections/%s/%s", collection.Id, path.Base(collection.BannerFile.String)))
		bannerUrl = &url
	}

	return Collection{
		Id:             collection.Id,
		CollectionType: collection.Type,
		Name:           collection.Name,
		CoverUrl:       coverUrl,
		LogoUrl:        logoUrl,
		BannerUrl:      bannerUrl,
	}
}

type CollectionItem struct {
	CollectionId string `json:"collectionId"`
	MediaId      string `json:"mediaId"`

	CollectionName string `json:"collectionName"`
	SearchSlug     string `json:"searchSlug"`
	Order          int64  `json:"order"`
	SubOrder       int64  `json:"subOrder"`

	Title       string  `json:"title"`
	Description *string `json:"description"`

	MediaType    types.MediaType   `json:"mediaType"`
	Score        *float64          `json:"score"`
	Status       types.MediaStatus `json:"status"`
	Rating       types.MediaRating `json:"rating"`
	PartCount    int64             `json:"partCount"`
	AiringSeason *string           `json:"airingSeason"`

	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`

	Studios []string `json:"studios"`
	Tags    []string `json:"tags"`

	CoverUrl  *string `json:"coverUrl"`
	LogoUrl   *string `json:"logoUrl"`
	BannerUrl *string `json:"bannerUrl"`

	User *MediaUser `json:"user,omitempty"`
}

func ConvertDBCollectionItem(c pyrin.Context, hasUser bool, item database.FullCollectionMediaItem) CollectionItem {
	// TODO(patrik): Add default cover
	var coverUrl *string
	var bannerUrl *string
	var logoUrl *string

	if item.MediaCoverFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", item.MediaId, path.Base(item.MediaCoverFile.String)))
		coverUrl = &url
	}

	if item.MediaLogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", item.MediaId, path.Base(item.MediaLogoFile.String)))
		logoUrl = &url
	}

	if item.MediaBannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", item.MediaId, path.Base(item.MediaBannerFile.String)))
		bannerUrl = &url
	}

	var user *MediaUser
	if hasUser {
		user = &MediaUser{}

		if item.MediaUserData.Valid {
			val := item.MediaUserData.Data
			user.List = val.List
			user.CurrentPart = val.Part
			user.RevisitCount = val.RevisitCount
			user.Score = val.Score
			user.IsRevisiting = val.IsRevisiting > 0
		}
	}

	return CollectionItem{
		CollectionId:   item.CollectionId,
		MediaId:        item.MediaId,
		Title:          item.MediaTitle,
		Description:    utils.SqlNullToStringPtr(item.MediaDescription),
		MediaType:      item.MediaType,
		Score:          utils.SqlNullToFloat64Ptr(item.MediaScore),
		Status:         item.MediaStatus,
		Rating:         item.MediaRating,
		PartCount:      item.MediaPartCount.Int64,
		AiringSeason:   utils.SqlNullToStringPtr(item.MediaAiringSeason),
		StartDate:      utils.SqlNullToStringPtr(item.MediaStartDate),
		EndDate:        utils.SqlNullToStringPtr(item.MediaEndDate),
		Studios:        utils.FixNilArrayToEmpty(item.MediaStudios.Data),
		Tags:           utils.FixNilArrayToEmpty(item.MediaTags.Data),
		CoverUrl:       coverUrl,
		LogoUrl:        logoUrl,
		BannerUrl:      bannerUrl,
		User:           user,
		CollectionName: item.CollectionName,
		SearchSlug:     item.SearchSlug,
		Order:          item.OrderNumber,
		SubOrder:       item.SubOrderNumber,
	}
}

type CollectionGroup struct {
	Name  string `json:"name"`
	Order int    `json:"order"`

	Entries []CollectionItem
}

type GetCollectionItems struct {
	Groups []CollectionGroup `json:"groups"`
}

func InstallCollectionHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetCollections",
			Method:       http.MethodGet,
			Path:         "/collections",
			ResponseType: GetCollections{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				ctx := context.TODO()

				filterStr := q.Get("filter")
				sortStr := q.Get("sort")
				collections, page, err := app.DB().GetPagedCollections(ctx, nil, filterStr, sortStr, opts)
				if err != nil {
					return nil, err
				}

				res := GetCollections{
					Page:        page,
					Collections: make([]Collection, len(collections)),
				}

				for i, col := range collections {
					res.Collections[i] = ConvertDBCollection(c, false, col)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetCollectionById",
			Method:       http.MethodGet,
			Path:         "/collections/:id",
			ResponseType: GetCollectionById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				var userId *string
				if user, err := User(app, c); err == nil {
					userId = &user.Id
				}

				collection, err := app.DB().GetCollectionById(c.Request().Context(), userId, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, CollectionNotFound()
					}

					return nil, err
				}

				return GetCollectionById{
					Collection: ConvertDBCollection(c, false, collection),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetCollectionItems",
			Method:       http.MethodGet,
			Path:         "/collections/:id/items",
			ResponseType: GetCollectionItems{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				var userId *string
				if user, err := User(app, c); err == nil {
					userId = &user.Id
				}

				ctx := c.Request().Context()

				collection, err := app.DB().GetCollectionById(ctx, userId, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, CollectionNotFound()
					}

					return nil, err
				}

				items, err := app.DB().GetFullAllCollectionMediaItemsByCollection(ctx, userId, collection.Id)
				if err != nil {
					return nil, err
				}

				groups := make(map[string][]database.FullCollectionMediaItem)

				for _, item := range items {
					groups[item.GroupName] = append(groups[item.GroupName], item)
				}

				res := GetCollectionItems{}

				for _, group := range groups {
					entries := make([]CollectionItem, 0, len(group))

					for _, entry := range group {
						i := ConvertDBCollectionItem(c, userId != nil, entry)
						entries = append(entries, i)
					}

					sort.SliceStable(entries, func(i, j int) bool {
						return entries[i].Order < entries[j].Order
					})

					res.Groups = append(res.Groups, CollectionGroup{
						Name:    group[0].GroupName,
						Order:   int(group[0].GroupOrder),
						Entries: entries,
					})
				}

				sort.SliceStable(res.Groups, func(i, j int) bool {
					return res.Groups[i].Order < res.Groups[j].Order
				})

				return res, nil
			},
		},
	)
}
