package apis

import (
	"context"
	"errors"
	"net/http"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
)

type Collection struct {
	Id   string               `json:"id"`
	Type types.CollectionType `json:"type"`

	Name string `json:"name"`
}

type GetCollections struct {
	// Page       types.Page `json:"page"`
	Collections []Collection `json:"collections"`
}

type GetCollectionById struct {
	Collection
}

func ConvertDBCollection(c pyrin.Context, hasUser bool, collection database.Collection) Collection {
	return Collection{
		Id:   collection.Id,
		Type: collection.Type,
		Name: collection.Name,
	}
}

type CreateCollection struct {
	Id string `json:"id"`
}

type CreateCollectionBody struct {
	Type string `json:"type"`

	Name string `json:"name"`
}

func (b *CreateCollectionBody) Transform() {
	b.Name = transform.String(b.Name)
}

func (b CreateCollectionBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Type, validate.Required, validate.By(types.ValidateCollectionType)),
		validate.Field(&b.Name, validate.Required),
	)
}

type EditCollectionBody struct {
	Type *string `json:"type,omitempty"`

	Name *string `json:"name,omitempty"`

	AdminStatus *string `json:"adminStatus,omitempty"`
}

func (b *EditCollectionBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
}

func (b EditCollectionBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Type, validate.Required.When(b.Type != nil), validate.By(types.ValidateCollectionType)),

		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),

		validate.Field(&b.AdminStatus, validate.Required.When(b.AdminStatus != nil), validate.By(types.ValidateAdminStatus)),
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
				// q := c.Request().URL.Query()
				// opts := getPageOptions(q)

				ctx := context.TODO()

				// filterStr := q.Get("filter")
				// sortStr := q.Get("sort")
				collections, err := app.DB().GetAllCollections(ctx)
				if err != nil {
					return nil, err
				}

				res := GetCollections{
					// Page:       types.Page{},
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
			Name:         "CreateCollection",
			Method:       http.MethodPost,
			Path:         "/collections",
			ResponseType: CreateCollection{},
			BodyType:     CreateCollectionBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				body, err := pyrin.Body[CreateCollectionBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				ty := types.CollectionType(body.Type)

				id, err := app.DB().CreateCollection(ctx, database.CreateCollectionParams{
					Type:        ty,
					Name:        body.Name,
					AdminStatus: "",
				})
				if err != nil {
					return nil, err
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

				// TODO(patrik): Add admin check

				body, err := pyrin.Body[EditCollectionBody](c)
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

				if body.Type != nil {
					t := types.CollectionType(*body.Type)

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

				if body.AdminStatus != nil {
					s := types.AdminStatus(*body.AdminStatus)
					changes.AdminStatus = database.Change[types.AdminStatus]{
						Value:   s,
						Changed: s != dbCollection.AdminStatus,
					}
				}

				err = app.DB().UpdateCollection(ctx, dbCollection.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
