package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
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

	TmdbId    string `json:"tmdbId"`
	ImdbId    string `json:"imdbId"`
	MalId     string `json:"malId"`
	AnilistId string `json:"anilistId"`

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
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", media.Id, path.Base(media.CoverFile.String)))
		coverUrl = &url
	}

	if media.LogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", media.Id, path.Base(media.LogoFile.String)))
		logoUrl = &url
	}

	if media.BannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", media.Id, path.Base(media.BannerFile.String)))
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
		TmdbId:       media.TmdbId.String,
		ImdbId:       media.ImdbId.String,
		MalId:        media.MalId.String,
		AnilistId:    media.AnilistId.String,
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

type CreateMedia struct {
	Id string `json:"id"`
}

type CreateMediaBody struct {
	MediaType string `json:"mediaType"`

	TmdbId string `json:"tmdbId"`
	ImdbId string `json:"imdbId"`
	// TODO(patrik): Add validation for this, should start with 'anime@' or 'manga@'
	MalId string `json:"malId"`
	// TODO(patrik): Add validation for this, should start with 'anime@' or 'manga@'
	AnilistId string `json:"anilistId"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Score        float64 `json:"score"`
	Status       string  `json:"status"`
	Rating       string  `json:"rating"`
	AiringSeason string  `json:"airingSeason"`

	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`

	PartCount int `json:"partCount"`

	CoverUrl  string `json:"coverUrl"`
	BannerUrl string `json:"bannerUrl"`
	LogoUrl   string `json:"logoUrl"`

	Tags    []string `json:"tags"`
	Studios []string `json:"studios"`

	CollectionId   string `json:"collectionId,omitempty"`
	CollectionName string `json:"collectionName,omitempty"`
}

func (b *CreateMediaBody) Transform() {
	b.TmdbId = anvil.String(b.TmdbId)
	b.ImdbId = anvil.String(b.ImdbId)
	b.MalId = anvil.String(b.MalId)
	b.AnilistId = anvil.String(b.AnilistId)

	b.Title = anvil.String(b.Title)
	b.Description = anvil.String(b.Description)

	b.Score = utils.Clamp(b.Score, 0.0, 10.0)
	b.AiringSeason = utils.TransformStringSlug(b.AiringSeason)

	b.PartCount = utils.Min(b.PartCount, 0)

	b.StartDate = anvil.String(b.StartDate)
	b.EndDate = anvil.String(b.EndDate)

	b.Tags = utils.TransformSlugArray(b.Tags)
	b.Studios = utils.TransformSlugArray(b.Studios)
}

func (b CreateMediaBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.MediaType, validate.Required, validate.By(types.ValidateMediaType)),

		validate.Field(&b.Title, validate.Required),

		validate.Field(&b.Status, validate.By(types.ValidateMediaStatus)),
		validate.Field(&b.Rating, validate.By(types.ValidateMediaRating)),

		validate.Field(&b.StartDate, validate.Date(types.MediaDateLayout)),
		validate.Field(&b.EndDate, validate.Date(types.MediaDateLayout)),

		validate.Field(&b.CollectionName, validate.Required.When(b.CollectionId != "")),
	)
}

type EditMediaBody struct {
	MediaType *string `json:"mediaType,omitempty"`

	TmdbId *string `json:"tmdbId,omitempty"`
	ImdbId *string `json:"imdbId,omitempty"`
	// TODO(patrik): Add validation for this, should start with 'anime@' or 'manga@'
	MalId *string `json:"malId,omitempty"`
	// TODO(patrik): Add validation for this, should start with 'anime@' or 'manga@'
	AnilistId *string `json:"anilistId,omitempty"`

	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`

	Score        *float64 `json:"score,omitempty"`
	Status       *string  `json:"status,omitempty"`
	Rating       *string  `json:"rating,omitempty"`
	AiringSeason *string  `json:"airingSeason,omitempty"`

	StartDate *string `json:"startDate,omitempty"`
	EndDate   *string `json:"endDate,omitempty"`

	AdminStatus *string `json:"adminStatus,omitempty"`

	Tags    *[]string `json:"tags,omitempty"`
	Studios *[]string `json:"studios,omitempty"`
}

func (b *EditMediaBody) Transform() {
	b.TmdbId = anvil.StringPtr(b.TmdbId)
	b.ImdbId = anvil.StringPtr(b.ImdbId)
	b.MalId = anvil.StringPtr(b.MalId)
	b.AnilistId = anvil.StringPtr(b.AnilistId)

	b.Title = anvil.StringPtr(b.Title)
	b.Description = anvil.StringPtr(b.Description)

	if b.Score != nil {
		*b.Score = utils.Clamp(*b.Score, 0.0, 10.0)
	}

	if b.AiringSeason != nil {
		*b.AiringSeason = utils.TransformStringSlug(*b.AiringSeason)
	}

	if b.Tags != nil {
		*b.Tags = utils.TransformSlugArray(*b.Tags)
	}

	if b.Studios != nil {
		*b.Studios = utils.TransformSlugArray(*b.Studios)
	}
}

func (b EditMediaBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.MediaType, validate.Required.When(b.MediaType != nil), validate.By(types.ValidateMediaType)),

		validate.Field(&b.Title, validate.Required.When(b.Title != nil)),

		validate.Field(&b.Status, validate.Required.When(b.Status != nil), validate.By(types.ValidateMediaStatus)),
		validate.Field(&b.Rating, validate.Required.When(b.Rating != nil), validate.By(types.ValidateMediaRating)),

		validate.Field(&b.StartDate, validate.Date(types.MediaDateLayout)),
		validate.Field(&b.EndDate, validate.Date(types.MediaDateLayout)),

		validate.Field(&b.AdminStatus, validate.Required.When(b.AdminStatus != nil), validate.By(types.ValidateAdminStatus)),
	)
}

type AddPart struct {
	Index int64 `json:"index"`
}

type AddPartBody struct {
	Index int64  `json:"index"`
	Name  string `json:"name"`
}

func (b *AddPartBody) Transform() {
	b.Name = anvil.String(b.Name)
	b.Index = utils.Min(b.Index, 0)
}

type EditPartBody struct {
	Name *string `json:"name"`
}

func (b *EditPartBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)
}

func (b EditPartBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

type PartBody struct {
	Name string `json:"name"`
}

func (b *PartBody) Transform() {
	b.Name = anvil.String(b.Name)
}

func (b PartBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type SetPartsBody struct {
	Parts []PartBody `json:"parts"`
}

func (b *SetPartsBody) Transform() {
	for i := range b.Parts {
		b.Parts[i].Transform()
	}
}

func (b SetPartsBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Parts),
	)
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

		pyrin.ApiHandler{
			Name:         "CreateMedia",
			Method:       http.MethodPost,
			Path:         "/media",
			ResponseType: CreateMedia{},
			BodyType:     CreateMediaBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				body, err := pyrin.Body[CreateMediaBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				if body.AiringSeason != "" {
					err := app.DB().CreateTag(ctx, body.AiringSeason, body.AiringSeason)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				ty := types.MediaType(body.MediaType)

				id, err := app.DB().CreateMedia(ctx, database.CreateMediaParams{
					Type: ty,
					TmdbId: sql.NullString{
						String: body.TmdbId,
						Valid:  body.TmdbId != "",
					},
					ImdbId: sql.NullString{
						String: body.ImdbId,
						Valid:  body.ImdbId != "",
					},
					MalId: sql.NullString{
						String: body.MalId,
						Valid:  body.MalId != "",
					},
					AnilistId: sql.NullString{
						String: body.AnilistId,
						Valid:  body.AnilistId != "",
					},
					Title: body.Title,
					Description: sql.NullString{
						String: body.Description,
						Valid:  body.Description != "",
					},
					Score: sql.NullFloat64{
						Float64: body.Score,
						Valid:   body.Score != 0.0,
					},
					Status: types.MediaStatus(body.Status),
					Rating: types.MediaRating(body.Rating),
					AiringSeason: sql.NullString{
						String: body.AiringSeason,
						Valid:  body.AiringSeason != "",
					},
					StartDate: sql.NullString{
						String: body.StartDate,
						Valid:  body.StartDate != "",
					},
					EndDate: sql.NullString{
						String: body.EndDate,
						Valid:  body.EndDate != "",
					},
				})
				if err != nil {
					return nil, err
				}

				if ty.IsMovie() {
					err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
						MediaId: id,
						Name:    body.Title,
						Index:   1,
					})
					if err != nil {
						return nil, err
					}
				} else {
					for i := range body.PartCount {
						err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
							MediaId: id,
							Name:    fmt.Sprintf("Episode %d", i+1),
							Index:   int64(i + 1),
						})
						if err != nil {
							return nil, err
						}
					}
				}

				// if body.CoverUrl != "" {
				// 	_, err := downloadImage(ctx, app.DB(), app.WorkDir(), id, body.CoverUrl, types.MediaImageTypeCover, true)
				// 	if err != nil {
				// 		logger.Error("failed to download cover image for media", "mediaId", id, "err", err)
				// 	}
				// }
				//
				// if body.BannerUrl != "" {
				// 	_, err := downloadImage(ctx, app.DB(), app.WorkDir(), id, body.BannerUrl, types.MediaImageTypeBanner, true)
				// 	if err != nil {
				// 		logger.Error("failed to download banner image for media", "mediaId", id, "err", err)
				// 	}
				// }
				//
				// if body.LogoUrl != "" {
				// 	_, err := downloadImage(ctx, app.DB(), app.WorkDir(), id, body.LogoUrl, types.MediaImageTypeLogo, true)
				// 	if err != nil {
				// 		logger.Error("failed to download logo image for media", "media", id, "err", err)
				// 	}
				// }

				for _, tag := range body.Tags {
					err := app.DB().CreateTag(ctx, tag, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				for _, tag := range body.Tags {
					err := app.DB().AddTagToMedia(ctx, id, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				for _, tag := range body.Studios {
					err := app.DB().CreateTag(ctx, tag, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				for _, tag := range body.Studios {
					err := app.DB().AddStudioToMedia(ctx, id, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				if body.CollectionId != "" {
					col, err := app.DB().GetCollectionById(ctx, nil, body.CollectionId)
					if err != nil {
						// TODO(patrik): Better handling of error
						if !errors.Is(err, database.ErrItemNotFound) {
							return nil, CollectionNotFound()
						}

						return nil, err
					}

					err = app.DB().CreateCollectionMediaItem(ctx, database.CreateCollectionMediaItemParams{
						CollectionId: col.Id,
						MediaId:      id,
						Name:         body.CollectionName,
						OrderNumber:  0,
					})
					if err != nil {
						return nil, err
					}
				}

				mediaDir := app.WorkDir().MediaDirById(id)
				dirs := []string{
					mediaDir.String(),
					mediaDir.Images(),
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						return nil, err
					}
				}

				return CreateMedia{
					Id: id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditMedia",
			Method:       http.MethodPatch,
			Path:         "/media/:id",
			ResponseType: nil,
			BodyType:     EditMediaBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				// TODO(patrik): Add admin check

				body, err := pyrin.Body[EditMediaBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbMedia, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				changes := database.MediaChanges{}

				if body.MediaType != nil {
					t := types.MediaType(*body.MediaType)

					changes.Type = database.Change[types.MediaType]{
						Value:   t,
						Changed: t != dbMedia.Type,
					}
				}

				if body.TmdbId != nil {
					changes.TmdbId = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.TmdbId,
							Valid:  *body.TmdbId != "",
						},
						Changed: *body.TmdbId != dbMedia.TmdbId.String,
					}
				}

				if body.ImdbId != nil {
					changes.ImdbId = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.ImdbId,
							Valid:  *body.ImdbId != "",
						},
						Changed: *body.ImdbId != dbMedia.ImdbId.String,
					}
				}

				if body.MalId != nil {
					changes.MalId = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.MalId,
							Valid:  *body.MalId != "",
						},
						Changed: *body.MalId != dbMedia.MalId.String,
					}
				}

				if body.AnilistId != nil {
					changes.AnilistId = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.AnilistId,
							Valid:  *body.AnilistId != "",
						},
						Changed: *body.AnilistId != dbMedia.AnilistId.String,
					}
				}

				if body.Title != nil {
					changes.Title = database.Change[string]{
						Value:   *body.Title,
						Changed: *body.Title != dbMedia.Title,
					}
				}

				if body.Description != nil {
					changes.Description = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.Description,
							Valid:  *body.Description != "",
						},
						Changed: *body.Description != dbMedia.Description.String,
					}
				}

				if body.Score != nil {
					changes.Score = database.Change[sql.NullFloat64]{
						Value: sql.NullFloat64{
							Float64: *body.Score,
							Valid:   *body.Score != 0.0,
						},
						Changed: *body.Score != dbMedia.Score.Float64,
					}
				}

				if body.Status != nil {
					s := types.MediaStatus(*body.Status)
					changes.Status = database.Change[types.MediaStatus]{
						Value:   s,
						Changed: s != dbMedia.Status,
					}
				}

				if body.Rating != nil {
					r := types.MediaRating(*body.Rating)
					changes.Rating = database.Change[types.MediaRating]{
						Value:   r,
						Changed: r != dbMedia.Rating,
					}
				}

				if body.StartDate != nil {
					changes.StartDate = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.StartDate,
							Valid:  *body.StartDate != "",
						},
						Changed: *body.StartDate != dbMedia.StartDate.String,
					}
				}

				if body.EndDate != nil {
					changes.EndDate = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: *body.EndDate,
							Valid:  *body.EndDate != "",
						},
						Changed: *body.EndDate != dbMedia.EndDate.String,
					}
				}

				if body.AiringSeason != nil {
					airingSeason := *body.AiringSeason
					if airingSeason != "" {
						err := app.DB().CreateTag(ctx, airingSeason, airingSeason)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}

					changes.AiringSeason = database.Change[sql.NullString]{
						Value: sql.NullString{
							String: airingSeason,
							Valid:  airingSeason != "",
						},
						Changed: airingSeason != dbMedia.AiringSeason.String,
					}
				}

				// if body.AdminStatus != nil {
				// 	s := types.AdminStatus(*body.AdminStatus)
				// 	changes.AdminStatus = database.Change[types.AdminStatus]{
				// 		Value:   s,
				// 		Changed: s != dbMedia.AdminStatus,
				// 	}
				// }

				err = app.DB().UpdateMedia(ctx, dbMedia.Id, changes)
				if err != nil {
					return nil, err
				}

				if body.Tags != nil {
					err := app.DB().RemoveAllTagsFromMedia(ctx, dbMedia.Id)
					if err != nil {
						return nil, err
					}

					tags := *body.Tags

					for _, tag := range tags {
						err := app.DB().CreateTag(ctx, tag, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}

					for _, tag := range tags {
						err := app.DB().AddTagToMedia(ctx, id, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				if body.Studios != nil {
					err := app.DB().RemoveAllStudiosFromMedia(ctx, dbMedia.Id)
					if err != nil {
						return nil, err
					}

					studios := *body.Studios

					for _, tag := range studios {
						err := app.DB().CreateTag(ctx, tag, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}

					for _, tag := range studios {
						err := app.DB().AddStudioToMedia(ctx, id, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				return nil, nil
			},
		},

		pyrin.FormApiHandler{
			Name:         "ChangeImages",
			Method:       http.MethodPatch,
			Path:         "/media/:id/images",
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

				// TODO(patrik): Add admin check

				ctx := context.Background()

				dbMedia, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				mediaDir := app.WorkDir().MediaDirById(dbMedia.Id)
				dirs := []string{
					mediaDir.String(),
					mediaDir.Images(),
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						return nil, err
					}
				}

				changes := database.MediaChanges{}

				// TODO(patrik): Change name
				test := func(name string) (database.Change[sql.NullString], error) {
					coverFiles, err := pyrin.FormFiles(c, name)
					if err != nil {
						return database.Change[sql.NullString]{}, err
					}

					if len(coverFiles) > 0 {
						file := coverFiles[0]

						// TODO(patrik): Add better size limiting
						if file.Size > 25*1024*1024 {
							return database.Change[sql.NullString]{}, errors.New("file too big")
						}

						contentType := file.Header.Get("Content-Type")
						ext, err := utils.GetImageExtFromContentType(contentType)
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, errors.New("file too big")
						}

						if dbMedia.CoverFile.Valid {
							p := path.Join(mediaDir.Images(), dbMedia.CoverFile.String)
							err = os.Remove(p)
							if err != nil {
								return database.Change[sql.NullString]{}, errors.New("file too big")
							}
						}

						f, err := file.Open()
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, errors.New("file too big")
						}
						defer f.Close()

						outFile, err := os.OpenFile(path.Join(mediaDir.Images(), name+ext), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, errors.New("file too big")
						}
						defer outFile.Close()

						_, err = io.Copy(outFile, f)
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, errors.New("file too big")
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

				changes.CoverFile, err = test("cover")
				if err != nil {
					return nil, err
				}

				changes.LogoFile, err = test("logo")
				if err != nil {
					return nil, err
				}

				changes.BannerFile, err = test("banner")
				if err != nil {
					return nil, err
				}

				err = app.DB().UpdateMedia(ctx, dbMedia.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AddPart",
			Method:       http.MethodPost,
			Path:         "/media/:id/single/parts",
			ResponseType: AddPart{},
			BodyType:     AddPartBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				id := c.Param("id")

				body, err := pyrin.Body[AddPartBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbMedia, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				index := body.Index

				if index == 0 {
					// TODO(patrik): A better implementation would be getting
					// the last part from the database
					parts, err := app.DB().GetMediaPartsByMediaId(ctx, dbMedia.Id)
					if err != nil {
						return nil, err
					}

					if len(parts) > 0 {
						part := parts[len(parts)-1]
						index = part.Index + 1
					}
				}

				// TODO(patrik): Change this
				name := body.Name
				if name == "" {
					name = fmt.Sprintf("Episode %d", index)
				}

				err = app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
					Index:   index,
					MediaId: dbMedia.Id,
					Name:    name,
				})
				if err != nil {
					if errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, PartAlreadyExists()
					}

					return nil, err
				}

				return AddPart{
					Index: index,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditPart",
			Method:       http.MethodPatch,
			Path:         "/media/:id/parts/:index",
			ResponseType: nil,
			BodyType:     EditPartBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				id := c.Param("id")
				index, err := strconv.ParseInt(c.Param("index"), 10, 64)
				if err != nil {
					// TODO(patrik): Handle error better
					return nil, errors.New("failed to parse 'index' path param as integer")
				}

				body, err := pyrin.Body[EditPartBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbMedia, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				dbPart, err := app.DB().GetMediaPartByIndexMediaId(ctx, index, dbMedia.Id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PartNotFound()
					}

					return nil, err
				}

				changes := database.MediaPartChanges{}

				if body.Name != nil {
					changes.Name = database.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != dbPart.Name,
					}
				}

				err = app.DB().UpdateMediaPart(ctx, dbPart.Index, dbPart.MediaId, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "RemovePart",
			Method:       http.MethodDelete,
			Path:         "/media/:id/parts/:index",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				id := c.Param("id")
				index, err := strconv.ParseInt(c.Param("index"), 10, 64)
				if err != nil {
					// TODO(patrik): Handle error better
					return nil, errors.New("failed to parse 'index' path param as integer")
				}

				ctx := context.Background()

				dbMedia, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveMediaPart(ctx, index, dbMedia.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SetParts",
			Method:       http.MethodPost,
			Path:         "/media/:id/parts",
			ResponseType: nil,
			BodyType:     SetPartsBody{},
			HandlerFunc:  func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[SetPartsBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbMedia, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveAllMediaParts(ctx, dbMedia.Id)
				if err != nil {
					return nil, err
				}

				for i, part := range body.Parts {
					err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
						Index:   int64(i),
						MediaId: id,
						Name:    part.Name,
					})
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},
	)
}
