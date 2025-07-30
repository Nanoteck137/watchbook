package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type MediaUser struct {
	List         *types.MediaUserList `json:"list"`
	Score        *int64               `json:"score"`
	CurrentPart  *int64               `json:"currentPart"`
	RevisitCount *int64               `json:"revisitCount"`
	IsRevisiting bool                 `json:"isRevisiting"`
}

type Media struct {
	Id string `json:"id"`

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
	BannerUrl *string `json:"bannerUrl"`
	LogoUrl   *string `json:"logoUrl"`

	User *MediaUser `json:"user,omitempty"`
}

type GetMedia struct {
	Page  types.Page `json:"page"`
	Media []Media    `json:"media"`
}

type GetMediaById struct {
	Media
}

// TODO(patrik): Move
func getPageOptions(q url.Values) database.FetchOptions {
	perPage := 100
	page := 0

	if s := q.Get("perPage"); s != "" {
		i, _ := strconv.Atoi(s)
		if i > 0 {
			perPage = i
		}
	}

	if s := q.Get("page"); s != "" {
		i, _ := strconv.Atoi(s)
		page = i
	}

	return database.FetchOptions{
		PerPage: perPage,
		Page:    page,
	}
}

func ConvertDBMedia(c pyrin.Context, hasUser bool, media database.Media) Media {
	// TODO(patrik): Add default cover
	var coverUrl *string
	var bannerUrl *string
	var logoUrl *string

	if media.CoverFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", media.Id, path.Base(media.CoverFile.String)))
		coverUrl = &url
	}

	if media.LogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", media.Id, path.Base(media.LogoFile.String)))
		logoUrl = &url
	}

	if media.BannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", media.Id, path.Base(media.BannerFile.String)))
		bannerUrl = &url
	}

	var user *MediaUser
	if hasUser {
		user = &MediaUser{}

		if media.UserData.Valid {
			val := media.UserData.Data
			user.List = val.List
			user.CurrentPart = val.Part
			user.RevisitCount = val.RevisitCount
			user.Score = val.Score
			user.IsRevisiting = val.IsRevisiting > 0
		}
	}

	return Media{
		Id:           media.Id,
		Title:        media.Title,
		Description:  utils.SqlNullToStringPtr(media.Description),
		MediaType:    media.Type,
		Score:        utils.SqlNullToFloat64Ptr(media.Score),
		Status:       media.Status,
		Rating:       media.Rating,
		PartCount:    media.PartCount.Int64,
		Studios:      utils.FixNilArrayToEmpty(media.Studios.Data),
		Tags:         utils.FixNilArrayToEmpty(media.Tags.Data),
		AiringSeason: utils.SqlNullToStringPtr(media.AiringSeason),
		StartDate:    utils.SqlNullToStringPtr(media.StartDate),
		EndDate:      utils.SqlNullToStringPtr(media.EndDate),
		CoverUrl:     coverUrl,
		BannerUrl:    bannerUrl,
		LogoUrl:      logoUrl,
		User:         user,
	}
}

type SetMediaUserData struct {
	List         *types.MediaUserList `json:"list,omitempty"`
	Score        *int64               `json:"score,omitempty"`
	CurrentPart  *int64               `json:"currentPart,omitempty"`
	RevisitCount *int64               `json:"revisitCount,omitempty"`
	IsRevisiting *bool                `json:"isRevisiting,omitempty"`
}

func (b *SetMediaUserData) Transform() {
	if b.Score != nil {
		*b.Score = utils.Clamp(*b.Score, 0, 10)
	}
}

type MediaPart struct {
	Index   int64  `json:"index"`
	MediaId string `json:"mediaId"`

	Name string `json:"name"`
}

type GetMediaParts struct {
	Parts []MediaPart `json:"parts"`
}

func InstallMediaHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetMedia",
			Method:       http.MethodGet,
			Path:         "/media",
			ResponseType: GetMedia{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				ctx := context.TODO()

				var userId *string

				if q.Has("userId") {
					id := q.Get("userId")

					user, err := app.DB().GetUserById(ctx, id)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							return nil, UserNotFound()
						}

						return nil, err
					}

					userId = &user.Id
				} else {
					if user, err := User(app, c); err == nil {
						userId = &user.Id
					}
				}

				filterStr := q.Get("filter")
				sortStr := q.Get("sort")
				media, p, err := app.DB().GetPagedMedia(ctx, userId, filterStr, sortStr, opts)
				if err != nil {
					return nil, err
				}

				res := GetMedia{
					Page:  p,
					Media: make([]Media, len(media)),
				}

				for i, m := range media {
					res.Media[i] = ConvertDBMedia(c, userId != nil, m)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaById",
			Method:       http.MethodGet,
			Path:         "/media/:id",
			ResponseType: GetMediaById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				var userId *string
				if user, err := User(app, c); err == nil {
					userId = &user.Id
				}

				media, err := app.DB().GetMediaById(c.Request().Context(), userId, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				return GetMediaById{
					Media: ConvertDBMedia(c, userId != nil, media),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaParts",
			Method:       http.MethodGet,
			Path:         "/media/:id/parts",
			ResponseType: GetMediaParts{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.Background()

				dbMedia, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				parts, err := app.DB().GetMediaPartsByMediaId(ctx, dbMedia.Id)
				if err != nil {
					return nil, err
				}

				res := make([]MediaPart, len(parts))

				for i, part := range parts {
					res[i] = MediaPart{
						Index:   part.Index,
						MediaId: part.MediaId,
						Name:    part.Name,
					}
				}

				return GetMediaParts{
					Parts: res,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SetMediaUserData",
			Method:       http.MethodPost,
			Path:         "/media/:id/user",
			ResponseType: nil,
			BodyType:     SetMediaUserData{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				body, err := pyrin.Body[SetMediaUserData](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				media, err := app.DB().GetMediaById(ctx, &user.Id, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				val := media.UserData.Data

				data := database.SetMediaUserData{
					List:         utils.MediaUserListPtrToSqlNull(val.List),
					Part:         utils.Int64PtrToSqlNull(val.Part),
					IsRevisiting: val.IsRevisiting > 0,
					Score:        utils.Int64PtrToSqlNull(val.Score),
				}

				if body.List != nil {
					data.List = sql.NullString{
						String: string(*body.List),
						Valid:  *body.List != "",
					}
				}

				if body.CurrentPart != nil {
					data.Part = sql.NullInt64{
						Int64: *body.CurrentPart,
						Valid: *body.CurrentPart != 0,
					}
				}

				if body.RevisitCount != nil {
					data.RevisitCount = sql.NullInt64{
						Int64: *body.RevisitCount,
						Valid: *body.RevisitCount != 0,
					}
				}

				if body.IsRevisiting != nil {
					data.IsRevisiting = *body.IsRevisiting
				}

				if body.Score != nil {
					data.Score = sql.NullInt64{
						Int64: *body.Score,
						Valid: *body.Score != 0,
					}
				}

				err = app.DB().SetMediaUserData(ctx, media.Id, user.Id, data)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
