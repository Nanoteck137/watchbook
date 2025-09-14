package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Collection struct {
	Id             string               `json:"id"`
	CollectionType types.CollectionType `json:"collectionType"`

	Name string `json:"name"`

	CoverUrl  *string `json:"coverUrl"`
	LogoUrl   *string `json:"logoUrl"`
	BannerUrl *string `json:"bannerUrl"`
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
		url := ConvertURL(c, fmt.Sprintf("/files/collections/%s/images/%s", collection.Id, path.Base(collection.CoverFile.String)))
		coverUrl = &url
	}

	if collection.LogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/collections/%s/images/%s", collection.Id, path.Base(collection.LogoFile.String)))
		logoUrl = &url
	}

	if collection.BannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/collections/%s/images/%s", collection.Id, path.Base(collection.BannerFile.String)))
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
	Position       int    `json:"position"`

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

	Creators []string `json:"creators"`
	Tags     []string `json:"tags"`

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
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", item.MediaId, path.Base(item.MediaCoverFile.String)))
		coverUrl = &url
	}

	if item.MediaLogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", item.MediaId, path.Base(item.MediaLogoFile.String)))
		logoUrl = &url
	}

	if item.MediaBannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", item.MediaId, path.Base(item.MediaBannerFile.String)))
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
		Creators:       utils.FixNilArrayToEmpty(item.MediaCreators.Data),
		Tags:           utils.FixNilArrayToEmpty(item.MediaTags.Data),
		CoverUrl:       coverUrl,
		LogoUrl:        logoUrl,
		BannerUrl:      bannerUrl,
		User:           user,
		CollectionName: item.CollectionName,
		SearchSlug:     item.SearchSlug,
		Position:       item.Position,
	}
}

type GetCollectionItems struct {
	Items []CollectionItem `json:"items"`
}

type CreateCollection struct {
	Id string `json:"id"`
}

type CreateCollectionBody struct {
	CollectionType string `json:"collectionType"`

	Name string `json:"name"`
}

func (b *CreateCollectionBody) Transform() {
	b.Name = anvil.String(b.Name)
}

func (b CreateCollectionBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.CollectionType, validate.Required, validate.By(types.ValidateCollectionType)),
		validate.Field(&b.Name, validate.Required),
	)
}

type EditCollectionBody struct {
	CollectionType *string `json:"collectionType,omitempty"`

	Name *string `json:"name,omitempty"`

	AdminStatus *string `json:"adminStatus,omitempty"`
}

func (b *EditCollectionBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)
}

func (b EditCollectionBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.CollectionType, validate.Required.When(b.CollectionType != nil), validate.By(types.ValidateCollectionType)),

		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),

		validate.Field(&b.AdminStatus, validate.Required.When(b.AdminStatus != nil), validate.By(types.ValidateAdminStatus)),
	)
}

type AddCollectionItemBody struct {
	MediaId string `json:"mediaId"`

	Name       string `json:"name"`
	SearchSlug string `json:"searchSlug"`
	Position   int    `json:"position"`
}

func (b *AddCollectionItemBody) Transform() {
	b.Name = anvil.String(b.Name)
	b.SearchSlug = utils.Slug(b.SearchSlug)
}

func (b AddCollectionItemBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type EditCollectionItemBody struct {
	Name       *string `json:"name,omitempty"`
	SearchSlug *string `json:"searchSlug,omitempty"`
	Position   *int    `json:"position,omitempty"`
}

func (b *EditCollectionItemBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)
	if b.SearchSlug != nil {
		*b.SearchSlug = utils.Slug(*b.SearchSlug)
	}
}

func (b EditCollectionItemBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
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

				res := GetCollectionItems{
					Items: make([]CollectionItem, 0, len(items)),
				}

				for _, item := range items {
					res.Items = append(res.Items, ConvertDBCollectionItem(c, userId != nil, item))
				}

				sort.SliceStable(res.Items, func(i, j int) bool {
					return res.Items[i].Position < res.Items[j].Position
				})

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateCollection",
			Method:       http.MethodPost,
			Path:         "/collections",
			ResponseType: CreateCollection{},
			BodyType:     CreateCollectionBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreateCollectionBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				ty := types.CollectionType(body.CollectionType)

				id, err := app.DB().CreateCollection(ctx, database.CreateCollectionParams{
					Type: ty,
					Name: body.Name,
				})
				if err != nil {
					return nil, err
				}

				collectionDir := app.WorkDir().CollectionDirById(id)
				dirs := []string{
					collectionDir.String(),
					collectionDir.Images(),
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						return nil, err
					}
				}

				return CreateCollection{
					Id: id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditCollection",
			Method:       http.MethodPatch,
			Path:         "/collections/:id",
			ResponseType: nil,
			BodyType:     EditCollectionBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[EditCollectionBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbCollection, err := app.DB().GetCollectionById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, CollectionNotFound()
					}

					return nil, err
				}

				changes := database.CollectionChanges{}

				if body.CollectionType != nil {
					t := types.CollectionType(*body.CollectionType)

					changes.Type = database.Change[types.CollectionType]{
						Value:   t,
						Changed: t != dbCollection.Type,
					}
				}

				if body.Name != nil {
					changes.Name = database.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != dbCollection.Name,
					}
				}

				err = app.DB().UpdateCollection(ctx, dbCollection.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "DeleteCollection",
			Method:       http.MethodDelete,
			Path:         "/collections/:id",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbCollection, err := app.DB().GetCollectionById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, CollectionNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveCollection(ctx, dbCollection.Id)
				if err != nil {
					return nil, err
				}

				dir := app.WorkDir().CollectionDirById(dbCollection.Id)
				err = os.RemoveAll(dir.String())
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.FormApiHandler{
			Name:         "ChangeCollectionImages",
			Method:       http.MethodPatch,
			Path:         "/collections/:id/images",
			ResponseType: nil,
			Spec: pyrin.FormSpec{
				Files: map[string]pyrin.FormFileSpec{
					"cover": {
						NumExpected: 0,
					},
					"logo": {
						NumExpected: 0,
					},
					"banner": {
						NumExpected: 0,
					},
				},
			},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbCollection, err := app.DB().GetCollectionById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, CollectionNotFound()
					}

					return nil, err
				}

				collectionDir := app.WorkDir().CollectionDirById(id)
				dirs := []string{
					collectionDir.String(),
					collectionDir.Images(),
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						return nil, err
					}
				}

				changes := database.CollectionChanges{}

				// TODO(patrik): Change name
				test := func(old sql.NullString, name string) (database.Change[sql.NullString], error) {
					files, err := pyrin.FormFiles(c, name)
					if err != nil {
						return database.Change[sql.NullString]{}, err
					}

					if len(files) > 0 {
						file := files[0]

						// TODO(patrik): Add better size limiting
						if file.Size > 25*1024*1024 {
							return database.Change[sql.NullString]{}, errors.New("file too big")
						}

						contentType := file.Header.Get("Content-Type")
						ext, err := utils.GetImageExtFromContentType(contentType)
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, err
						}

						if old.Valid {
							p := path.Join(collectionDir.Images(), old.String)
							err = os.Remove(p)
							if err != nil {
								return database.Change[sql.NullString]{}, err
							}
						}

						f, err := file.Open()
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, err
						}
						defer f.Close()

						outFile, err := os.OpenFile(path.Join(collectionDir.Images(), name+ext), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, err
						}
						defer outFile.Close()

						_, err = io.Copy(outFile, f)
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, err
						}

						return database.Change[sql.NullString]{
							Value: sql.NullString{
								String: name + ext,
								Valid:  true,
							},
							Changed: true,
						}, nil
					}

					return database.Change[sql.NullString]{}, nil
				}

				changes.CoverFile, err = test(dbCollection.CoverFile, "cover")
				if err != nil {
					return nil, err
				}

				changes.LogoFile, err = test(dbCollection.LogoFile, "logo")
				if err != nil {
					return nil, err
				}

				changes.BannerFile, err = test(dbCollection.BannerFile, "banner")
				if err != nil {
					return nil, err
				}

				err = app.DB().UpdateCollection(ctx, dbCollection.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AddCollectionItem",
			Method:       http.MethodPost,
			Path:         "/collections/:id/items",
			ResponseType: nil,
			BodyType:     AddCollectionItemBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[AddCollectionItemBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbCollection, err := app.DB().GetCollectionById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, CollectionNotFound()
					}

					return nil, err
				}

				dbMedia, err := app.DB().GetMediaById(ctx, nil, body.MediaId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				searchSlug := body.SearchSlug
				if searchSlug == "" {
					searchSlug = utils.Slug(body.Name)
				}

				err = app.DB().CreateCollectionMediaItem(ctx, database.CreateCollectionMediaItemParams{
					CollectionId: dbCollection.Id,
					MediaId:      dbMedia.Id,

					Name:       body.Name,
					SearchSlug: searchSlug,
					Position:   body.Position,
				})
				if err != nil {
					// TODO(patrik): Better handling of error
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "RemoveCollectionItem",
			Method:       http.MethodDelete,
			Path:         "/collections/:id/items/:mediaId",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				mediaId := c.Param("mediaId")

				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				item, err := app.DB().GetCollectionMediaItemById(ctx, id, mediaId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, CollectionItemNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveCollectionMediaItem(ctx, item.CollectionId, item.MediaId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditCollectionItem",
			Method:       http.MethodPatch,
			Path:         "/collections/:id/items/:mediaId",
			ResponseType: nil,
			BodyType:     EditCollectionItemBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				mediaId := c.Param("mediaId")

				body, err := pyrin.Body[EditCollectionItemBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				item, err := app.DB().GetCollectionMediaItemById(ctx, id, mediaId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, CollectionItemNotFound()
					}

					return nil, err
				}

				changes := database.CollectionMediaItemChanges{}

				if body.Name != nil {
					changes.Name = database.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != item.Name,
					}
				}

				if body.SearchSlug != nil {
					changes.SearchSlug = database.Change[string]{
						Value:   *body.SearchSlug,
						Changed: *body.SearchSlug != item.SearchSlug,
					}
				}

				if body.Position != nil {
					changes.Position = database.Change[int]{
						Value:   *body.Position,
						Changed: *body.Position != item.Position,
					}
				}

				err = app.DB().UpdateCollectionMediaItem(ctx, id, mediaId, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
