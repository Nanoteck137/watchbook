package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"path"
	"sort"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Folder struct {
	Id string `json:"id"`

	UserId string `json:"userId"`

	Name string `json:"name"`

	ItemCount int `json:"itemCount"`

	CoverUrl *string `json:"coverUrl"`
}

type GetFolders struct {
	Folders []Folder `json:"folders"`
}

type GetFolderById struct {
	Folder
}

func ConvertDBFolder(c pyrin.Context, hasUser bool, folder database.Folder) Folder {
	// TODO(patrik): Add default cover
	var coverUrl *string

	if folder.CoverFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/folders/%s/images/%s", folder.Id, path.Base(folder.CoverFile.String)))
		coverUrl = &url
	}

	return Folder{
		Id:        folder.Id,
		UserId:    folder.UserId,
		Name:      folder.Name,
		ItemCount: int(folder.ItemCount.Int64),
		CoverUrl:  coverUrl,
	}
}

type FolderItem struct {
	FolderId string `json:"folderId"`
	MediaId  string `json:"mediaId"`

	Position int `json:"position"`

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

func ConvertDBFolderItem(c pyrin.Context, hasUser bool, item database.FullFolderItem) FolderItem {
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

	return FolderItem{
		FolderId:     item.FolderId,
		MediaId:      item.MediaId,
		Title:        item.MediaTitle,
		Description:  utils.SqlNullToStringPtr(item.MediaDescription),
		MediaType:    item.MediaType,
		Score:        utils.SqlNullToFloat64Ptr(item.MediaScore),
		Status:       item.MediaStatus,
		Rating:       item.MediaRating,
		PartCount:    item.MediaPartCount.Int64,
		AiringSeason: utils.SqlNullToStringPtr(item.MediaAiringSeason),
		StartDate:    utils.SqlNullToStringPtr(item.MediaStartDate),
		EndDate:      utils.SqlNullToStringPtr(item.MediaEndDate),
		Creators:     utils.FixNilArrayToEmpty(item.MediaCreators.Data),
		Tags:         utils.FixNilArrayToEmpty(item.MediaTags.Data),
		CoverUrl:     coverUrl,
		LogoUrl:      logoUrl,
		BannerUrl:    bannerUrl,
		User:         user,
		Position:     item.Position,
	}
}

type GetFolderItems struct {
	Items []FolderItem `json:"items"`
}

type CreateFolder struct {
	Id string `json:"id"`
}

type CreateFolderBody struct {
	Name string `json:"name"`

	CoverUrl string `json:"coverUrl"`
}

func (b *CreateFolderBody) Transform() {
	b.Name = anvil.String(b.Name)

	b.CoverUrl = anvil.String(b.CoverUrl)
}

func (b CreateFolderBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type EditFolderBody struct {
	Name *string `json:"name,omitempty"`

	CoverUrl *string `json:"coverUrl,omitempty"`
}

func (b *EditFolderBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)

	b.CoverUrl = anvil.StringPtr(b.CoverUrl)
}

func (b EditFolderBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),

		validate.Field(&b.CoverUrl, validate.Required.When(b.CoverUrl != nil)),
	)
}

// type AddFolderItemBody struct {
// 	MediaId string `json:"mediaId"`
//
// 	Position int `json:"position"`
// }
//
// func (b *AddFolderItemBody) Transform() {
// 	b.Name = anvil.String(b.Name)
// 	b.SearchSlug = utils.Slug(b.SearchSlug)
// }
//
// func (b AddFolderItemBody) Validate() error {
// 	return validate.ValidateStruct(&b,
// 		validate.Field(&b.Name, validate.Required),
// 	)
// }

// type EditFolderItemBody struct {
// 	Name       *string `json:"name,omitempty"`
// 	SearchSlug *string `json:"searchSlug,omitempty"`
// 	Position   *int    `json:"position,omitempty"`
// }
//
// func (b *EditFolderItemBody) Transform() {
// 	b.Name = anvil.StringPtr(b.Name)
// 	if b.SearchSlug != nil {
// 		*b.SearchSlug = utils.Slug(*b.SearchSlug)
// 	}
// }
//
// func (b EditFolderItemBody) Validate() error {
// 	return validate.ValidateStruct(&b,
// 		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
// 	)
// }

func InstallFolderHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetFolders",
			Method:       http.MethodGet,
			Path:         "/folders",
			ResponseType: GetFolders{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				folders, err := app.DB().GetAllFoldersByUserId(ctx, user.Id)
				if err != nil {
					return nil, err
				}

				res := GetFolders{
					Folders: make([]Folder, len(folders)),
				}

				for i, col := range folders {
					res.Folders[i] = ConvertDBFolder(c, true, col)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetFolderById",
			Method:       http.MethodGet,
			Path:         "/folders/:id",
			ResponseType: GetFolderById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				folder, err := app.DB().GetFolderById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderNotFound()
					}

					return nil, err
				}

				if folder.UserId != user.Id {
					return nil, FolderNotFound()
				}

				return GetFolderById{
					Folder: ConvertDBFolder(c, true, folder),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetFolderItems",
			Method:       http.MethodGet,
			Path:         "/folders/:id/items",
			ResponseType: GetFolderItems{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				folder, err := app.DB().GetFolderById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderNotFound()
					}

					return nil, err
				}

				if folder.UserId != user.Id {
					return nil, FolderNotFound()
				}

				items, err := app.DB().GetFullAllFolderItemsByFolder(ctx, &user.Id, folder.Id)
				if err != nil {
					return nil, err
				}

				res := GetFolderItems{
					Items: make([]FolderItem, 0, len(items)),
				}

				for _, item := range items {
					res.Items = append(res.Items, ConvertDBFolderItem(c, true, item))
				}

				// TODO(patrik): Just use sql order
				sort.SliceStable(res.Items, func(i, j int) bool {
					return res.Items[i].Position < res.Items[j].Position
				})

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateFolder",
			Method:       http.MethodPost,
			Path:         "/folders",
			ResponseType: CreateFolder{},
			BodyType:     CreateFolderBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[CreateFolderBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				id := utils.CreateFolderId()

				// folderDir := app.WorkDir().FolderDirById(id)
				// dirs := []string{
				// 	folderDir.String(),
				// 	folderDir.Images(),
				// }
				//
				// for _, dir := range dirs {
				// 	err = os.Mkdir(dir, 0755)
				// 	if err != nil && !os.IsExist(err) {
				// 		return nil, err
				// 	}
				// }

				coverFile := ""
				//
				// if body.CoverUrl != "" {
				// 	p, err := utils.DownloadImageHashed(body.CoverUrl, folderDir.Images())
				// 	if err == nil {
				// 		coverFile = path.Base(p)
				// 	} else {
				// 		app.Logger().Error("failed to download cover image for folder", "err", err)
				// 	}
				// }

				_, err = app.DB().CreateFolder(ctx, database.CreateFolderParams{
					Id:     id,
					UserId: user.Id,
					Name:   body.Name,
					CoverFile: sql.NullString{
						String: coverFile,
						Valid:  coverFile != "",
					},
				})
				if err != nil {
					return nil, err
				}

				return CreateFolder{
					Id: id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditFolder",
			Method:       http.MethodPatch,
			Path:         "/folders/:id",
			ResponseType: nil,
			BodyType:     EditFolderBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[EditFolderBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbFolder, err := app.DB().GetFolderById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderNotFound()
					}

					return nil, err
				}

				if dbFolder.UserId != user.Id {
					return nil, FolderNotFound()
				}

				// folderDir := app.WorkDir().FolderDirById(dbFolder.Id)

				changes := database.FolderChanges{}

				if body.Name != nil {
					changes.Name = database.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != dbFolder.Name,
					}
				}

				// if body.CoverUrl != nil {
				// 	p, err := utils.DownloadImageHashed(*body.CoverUrl, folderDir.Images())
				// 	if err == nil {
				// 		n := path.Base(p)
				// 		changes.CoverFile = database.Change[sql.NullString]{
				// 			Value: sql.NullString{
				// 				String: n,
				// 				Valid:  n != "",
				// 			},
				// 			Changed: n != dbFolder.CoverFile.String,
				// 		}
				// 	} else {
				// 		app.Logger().Error("failed to download cover image for folder", "err", err)
				// 	}
				// }

				err = app.DB().UpdateFolder(ctx, dbFolder.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "DeleteFolder",
			Method:       http.MethodDelete,
			Path:         "/folders/:id",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbFolder, err := app.DB().GetFolderById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderNotFound()
					}

					return nil, err
				}

				if dbFolder.UserId != user.Id {
					return nil, FolderNotFound()
				}

				err = app.DB().RemoveFolder(ctx, dbFolder.Id)
				if err != nil {
					return nil, err
				}

				// dir := app.WorkDir().FolderDirById(dbFolder.Id)
				// err = os.RemoveAll(dir.String())
				// if err != nil {
				// 	return nil, err
				// }

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "AddFolderItem",
			Method: http.MethodPost,
			Path:   "/folders/:id/items/:mediaId",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				mediaId := c.Param("mediaId")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbFolder, err := app.DB().GetFolderById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderNotFound()
					}

					return nil, err
				}

				if dbFolder.UserId != user.Id {
					return nil, FolderNotFound()
				}

				dbMedia, err := app.DB().GetMediaById(ctx, nil, mediaId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				pos, _ := app.DB().GetLastFolderItemPosition(ctx, dbFolder.Id)

				err = app.DB().CreateFolderItem(ctx, database.CreateFolderItemParams{
					FolderId: dbFolder.Id,
					MediaId:  dbMedia.Id,
					Position: pos + 1,
				})
				if err != nil {
					// TODO(patrik): Better handling of error
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "RemoveFolderItem",
			Method:       http.MethodDelete,
			Path:         "/folders/:id/items/:mediaId",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				mediaId := c.Param("mediaId")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbFolder, err := app.DB().GetFolderById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderNotFound()
					}

					return nil, err
				}

				if dbFolder.UserId != user.Id {
					return nil, FolderNotFound()
				}

				item, err := app.DB().GetFolderItemById(ctx, id, mediaId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderItemNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveFolderItem(ctx, item.FolderId, item.MediaId)
				if err != nil {
					return nil, err
				}

				err = app.DB().RepackFolderItems(ctx, dbFolder.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "MoveFolderItem",
			Method: http.MethodPost,
			Path:   "/folders/:id/items/:mediaId/move/:pos",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				mediaId := c.Param("mediaId")
				// TODO(patrik): Error Handling?
				pos, _ := strconv.Atoi(c.Param("pos"))

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbFolder, err := app.DB().GetFolderById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderNotFound()
					}

					return nil, err
				}

				if dbFolder.UserId != user.Id {
					return nil, FolderNotFound()
				}

				// dbMedia, err := app.DB().GetMediaById(ctx, nil, mediaId)
				// if err != nil {
				// 	if errors.Is(err, database.ErrItemNotFound) {
				// 		return nil, MediaNotFound()
				// 	}
				//
				// 	return nil, err
				// }

				_, err = app.DB().GetFolderItemById(ctx, dbFolder.Id, mediaId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, FolderItemNotFound()
					}

					return nil, err
				}

				err = app.DB().MoveFolderItem(ctx, dbFolder.Id, mediaId, pos)
				if err != nil {
					return nil, err
				}

				err = app.DB().RepackFolderItems(ctx, dbFolder.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		// pyrin.ApiHandler{
		// 	Name:         "EditFolderItem",
		// 	Method:       http.MethodPatch,
		// 	Path:         "/folders/:id/items/:mediaId",
		// 	ResponseType: nil,
		// 	BodyType:     EditFolderItemBody{},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		id := c.Param("id")
		// 		mediaId := c.Param("mediaId")
		//
		// 		body, err := pyrin.Body[EditFolderItemBody](c)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		_, err = User(app, c, HasEditPrivilege)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		ctx := context.Background()
		//
		// 		item, err := app.DB().GetFolderMediaItemById(ctx, id, mediaId)
		// 		if err != nil {
		// 			if errors.Is(err, database.ErrItemNotFound) {
		// 				return nil, FolderItemNotFound()
		// 			}
		//
		// 			return nil, err
		// 		}
		//
		// 		changes := database.FolderMediaItemChanges{}
		//
		// 		if body.Name != nil {
		// 			changes.Name = database.Change[string]{
		// 				Value:   *body.Name,
		// 				Changed: *body.Name != item.Name,
		// 			}
		// 		}
		//
		// 		if body.SearchSlug != nil {
		// 			changes.SearchSlug = database.Change[string]{
		// 				Value:   *body.SearchSlug,
		// 				Changed: *body.SearchSlug != item.SearchSlug,
		// 			}
		// 		}
		//
		// 		if body.Position != nil {
		// 			changes.Position = database.Change[int]{
		// 				Value:   *body.Position,
		// 				Changed: *body.Position != item.Position,
		// 			}
		// 		}
		//
		// 		err = app.DB().UpdateFolderMediaItem(ctx, id, mediaId, changes)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		return nil, nil
		// 	},
		// },
	)
}
