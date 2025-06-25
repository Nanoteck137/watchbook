package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

const DateLayout = "2006-01-02"

// TODO(patrik):
//  - Add missing in-between parts

type MediaUser struct {
	List         *types.MediaUserList `json:"list"`
	Score        *int64               `json:"score"`
	Part         *int64               `json:"part"`
	RevisitCount *int64               `json:"revisitCount"`
	IsRevisiting bool                 `json:"isRevisiting"`
}

type MediaImage struct {
	Hash    string `json:"hash"`
	Url     string `json:"url"`
	IsCover bool   `json:"isCover"`
}

type Media struct {
	Id string `json:"id"`

	Title       string  `json:"title"`
	Description *string `json:"description"`

	Type      types.MediaType   `json:"type"`
	Score     *float64          `json:"score"`
	Status    types.MediaStatus `json:"status"`
	Rating    types.MediaRating `json:"rating"`
	PartCount int64             `json:"partCount"`
	// AiringSeason *MediaTag         `json:"airingSeason"`

	// StartDate *string `json:"startDate"`
	// EndDate   *string `json:"endDate"`

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

	for _, image := range media.Images.Data {
		if image.Type == types.MediaImageTypeCover && coverUrl == nil {
			url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", media.Id, image.Filename))
			coverUrl = &url
		}

		if image.Type == types.MediaImageTypeBanner && bannerUrl == nil {
			url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", media.Id, image.Filename))
			bannerUrl = &url
		}

		if image.Type == types.MediaImageTypeLogo && logoUrl == nil {
			url := ConvertURL(c, fmt.Sprintf("/files/media/%s/%s", media.Id, image.Filename))
			logoUrl = &url
		}
	}

	var user *MediaUser
	if hasUser {
		user = &MediaUser{}

		if media.UserData.Valid {
			val := media.UserData.Data
			user.List = val.List
			user.Part = val.Part
			user.RevisitCount = val.RevisitCount
			user.Score = val.Score
			user.IsRevisiting = val.IsRevisiting > 0
		}
	}

	return Media{
		Id:          media.Id,
		Title:       media.Title,
		Description: utils.SqlNullToStringPtr(media.Description),
		Type:        media.Type,
		Score:       utils.SqlNullToFloat64Ptr(media.Score),
		Status:      media.Status,
		Rating:      media.Rating,
		PartCount:   media.PartCount,
		Studios:     utils.FixNilArrayToEmpty(media.Studios.Data),
		Tags:        utils.FixNilArrayToEmpty(media.Tags.Data),
		CoverUrl:    coverUrl,
		BannerUrl:   bannerUrl,
		LogoUrl:     logoUrl,
		User:        user,
	}
}

type SetMediaUserData struct {
	List         *types.MediaUserList `json:"list,omitempty"`
	Score        *int64               `json:"score,omitempty"`
	Part         *int64               `json:"part,omitempty"`
	RevisitCount *int64               `json:"revisitCount,omitempty"`
	IsRevisiting *bool                `json:"isRevisiting,omitempty"`
}

func (b *SetMediaUserData) Transform() {
	if b.Score != nil {
		*b.Score = utils.Clamp(*b.Score, 0, 10)
	}
}

type CreateMedia struct {
	Id string `json:"id"`
}

type CreateMediaBody struct {
	Type string `json:"type"`

	TmdbId    string `json:"tmdbId"`
	MalId     string `json:"malId"`
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
}

func (b *CreateMediaBody) Transform() {
	b.TmdbId = transform.String(b.TmdbId)
	b.MalId = transform.String(b.MalId)
	b.AnilistId = transform.String(b.AnilistId)

	b.Title = transform.String(b.Title)
	b.Description = transform.String(b.Description)

	b.Score = utils.Clamp(b.Score, 0.0, 10.0)
	b.AiringSeason = utils.TransformStringSlug(b.AiringSeason)

	b.PartCount = utils.Min(b.PartCount, 0)

	b.StartDate = transform.String(b.StartDate)
	b.EndDate = transform.String(b.EndDate)

	b.Tags = utils.TransformSlugArray(b.Tags)
	b.Studios = utils.TransformSlugArray(b.Studios)
}

func (b CreateMediaBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Type, validate.Required, validate.By(types.ValidateMediaType)),

		validate.Field(&b.Status, validate.By(types.ValidateMediaStatus)),
		validate.Field(&b.Rating, validate.By(types.ValidateMediaRating)),

		validate.Field(&b.StartDate, validate.Date(DateLayout)),
		validate.Field(&b.EndDate, validate.Date(DateLayout)),
	)
}

type EditMediaBody struct {
	Type *string `json:"type,omitempty"`

	TmdbId    *string `json:"tmdbId,omitempty"`
	MalId     *string `json:"malId,omitempty"`
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
	b.TmdbId = transform.StringPtr(b.TmdbId)
	b.MalId = transform.StringPtr(b.MalId)
	b.AnilistId = transform.StringPtr(b.AnilistId)

	b.Title = transform.StringPtr(b.Title)
	b.Description = transform.StringPtr(b.Description)

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
		validate.Field(&b.Type, validate.Required.When(b.Type != nil), validate.By(types.ValidateMediaType)),

		validate.Field(&b.Status, validate.Required.When(b.Status != nil), validate.By(types.ValidateMediaStatus)),
		validate.Field(&b.Rating, validate.Required.When(b.Rating != nil), validate.By(types.ValidateMediaRating)),

		validate.Field(&b.StartDate, validate.Date(DateLayout)),
		validate.Field(&b.EndDate, validate.Date(DateLayout)),

		validate.Field(&b.AdminStatus, validate.Required.When(b.AdminStatus != nil), validate.By(types.ValidateMediaAdminStatus)),
	)
}

type MediaPart struct {
	Index   int64  `json:"index"`
	MediaId string `json:"mediaId"`

	Name string `json:"name"`
}

type GetMediaParts struct {
	Parts []MediaPart `json:"parts"`
}

type AddMultiplePartsBody struct {
	Count int `json:"count"`
}

func (b *AddMultiplePartsBody) Transform() {
	b.Count = utils.Min(b.Count, 0)
}

type AddPart struct {
	Index int64 `json:"index"`
}

type AddPartBody struct {
	Index int64  `json:"index"`
	Name  string `json:"name"`
}

func (b *AddPartBody) Transform() {
	b.Name = transform.String(b.Name)
	b.Index = utils.Min(b.Index, 0)
}

type EditPartBody struct {
	Name *string `json:"name"`
}

func (b *EditPartBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
}

func (b EditPartBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

type EditImageBody struct {
	Type *string `json:"type"`

	IsPriamry *bool `json:"isPrimary"`
}

func (b EditImageBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Type, validate.Required.When(b.Type != nil), validate.By(types.ValidateMediaImageType)),
	)
}

type AddImage struct {
	Hash string `json:"hash"`
}

type AddImageBody struct {
	ImageUrl string `json:"imageUrl"`
	Type     string `json:"type"`
}

func (b *AddImageBody) Transform() {
	b.ImageUrl = transform.String(b.ImageUrl)
}

func (b AddImageBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.ImageUrl, validate.Required),
		validate.Field(&b.Type, validate.Required, validate.By(types.ValidateMediaImageType)),
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

				ty := types.MediaType(body.Type)

				id, err := app.DB().CreateMedia(ctx, database.CreateMediaParams{
					Type: ty,
					TmdbId: sql.NullString{
						String: body.TmdbId,
						Valid:  body.TmdbId != "",
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
							Name:    fmt.Sprintf("Part %d", i+1),
							Index:   int64(i + 1),
						})
						if err != nil {
							return nil, err
						}
					}
				}

				if body.CoverUrl != "" {
					_, err := downloadImage(ctx, app.DB(), app.WorkDir(), id, body.CoverUrl, types.MediaImageTypeCover, true)
					if err != nil {
						logger.Error("failed to download cover image for media", "mediaId", id, "err", err)
					}
				}

				if body.BannerUrl != "" {
					_, err := downloadImage(ctx, app.DB(), app.WorkDir(), id, body.BannerUrl, types.MediaImageTypeBanner, true)
					if err != nil {
						logger.Error("failed to download banner image for media", "mediaId", id, "err", err)
					}
				}

				if body.LogoUrl != "" {
					_, err := downloadImage(ctx, app.DB(), app.WorkDir(), id, body.LogoUrl, types.MediaImageTypeLogo, true)
					if err != nil {
						logger.Error("failed to download logo image for media", "media", id, "err", err)
					}
				}

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

				if body.Type != nil {
					t := types.MediaType(*body.Type)

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

				if body.AdminStatus != nil {
					s := types.MediaAdminStatus(*body.AdminStatus)
					changes.AdminStatus = database.Change[types.MediaAdminStatus]{
						Value:   s,
						Changed: s != dbMedia.AdminStatus,
					}
				}

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

				name := body.Name
				if name == "" {
					name = fmt.Sprintf("Part %d", index)
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
			Name:         "AddMultipleParts",
			Method:       http.MethodPost,
			Path:         "/media/:id/multiple/parts",
			ResponseType: nil,
			BodyType:     AddMultiplePartsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				id := c.Param("id")

				body, err := pyrin.Body[AddMultiplePartsBody](c)
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

				// TODO(patrik): A better implementation would be getting
				// the last part from the database
				parts, err := app.DB().GetMediaPartsByMediaId(ctx, dbMedia.Id)
				if err != nil {
					return nil, err
				}

				lastIndex := int64(0)
				if len(parts) > 0 {
					part := parts[len(parts)-1]
					lastIndex = part.Index
				}

				for i := range body.Count {
					idx := lastIndex + int64(i) + 1

					err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
						Index:   idx,
						MediaId: dbMedia.Id,
						Name:    fmt.Sprintf("Part %d", idx),
					})
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
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
			Name:         "AddImage",
			Method:       http.MethodPost,
			Path:         "/media/:id/images",
			ResponseType: AddImage{},
			BodyType:     AddImageBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				id := c.Param("id")

				body, err := pyrin.Body[AddImageBody](c)
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

				t := types.MediaImageType(body.Type)
				hash, err := downloadImage(ctx, app.DB(), app.WorkDir(), dbMedia.Id, body.ImageUrl, t, false)
				if err != nil {
					logger.Error("failed to download image for media", "mediaId", dbMedia.Id, "err", err)
					return nil, err
				}

				return AddImage{
					Hash: hash,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditImage",
			Method:       http.MethodPatch,
			Path:         "/media/:id/images/:hash",
			ResponseType: nil,
			BodyType:     EditImageBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				// TODO(patrik): Add admin check

				id := c.Param("id")
				hash := c.Param("hash")

				body, err := pyrin.Body[EditImageBody](c)
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

				dbImage, err := app.DB().GetMediaImagesByHashMediaId(ctx, dbMedia.Id, hash)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ImageNotFound()
					}

					return nil, err
				}

				changes := database.MediaImageChanges{}

				if body.Type != nil {
					t := types.MediaImageType(*body.Type)
					changes.Type = database.Change[types.MediaImageType]{
						Value:   t,
						Changed: t != dbImage.Type,
					}
				}

				if body.IsPriamry != nil {
					changes.IsPrimary = database.Change[bool]{
						Value:   *body.IsPriamry,
						Changed: *body.IsPriamry != (dbImage.IsPrimary > 0),
					}
				}

				err = app.DB().UpdateMediaImage(ctx, dbImage.MediaId, dbImage.Hash, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
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

				if body.Part != nil {
					data.Part = sql.NullInt64{
						Int64: *body.Part,
						Valid: *body.Part != 0,
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
