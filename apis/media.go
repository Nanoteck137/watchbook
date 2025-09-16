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
	"sort"
	"strconv"
	"time"

	"github.com/maruel/natural"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type MediaUser struct {
	HasData      bool                `json:"hasData"`
	List         types.MediaUserList `json:"list"`
	Score        *int64              `json:"score"`
	CurrentPart  *int64              `json:"currentPart"`
	RevisitCount *int64              `json:"revisitCount"`
	IsRevisiting bool                `json:"isRevisiting"`
}

type MediaRelease struct {
	ReleaseType      types.MediaPartReleaseType `json:"releaseType"`
	StartDate        string                     `json:"startDate"`
	NumExpectedParts int                        `json:"numExpectedParts"`
	PartOffset       int                        `json:"partOffset"`
	IntervalDays     int                        `json:"intervalDays"`
	DelayDays        int                        `json:"delayDays"`

	Status      types.MediaPartReleaseStatus `json:"status"`
	CurrentPart int                          `json:"currentPart"`
	NextAiring  *string                      `json:"nextAiring"`
}

type MediaProvider struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Value       string `json:"value"`
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

	Creators []string `json:"creators"`
	Tags     []string `json:"tags"`

	CoverUrl  *string `json:"coverUrl"`
	BannerUrl *string `json:"bannerUrl"`
	LogoUrl   *string `json:"logoUrl"`

	Providers []MediaProvider `json:"providers"`

	User    *MediaUser    `json:"user,omitempty"`
	Release *MediaRelease `json:"release"`
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

func ConvertDBMedia(c pyrin.Context, pm *provider.ProviderManager, hasUser bool, media database.Media) Media {
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

	providers := make([]MediaProvider, 0, len(media.Providers))
	for name, value := range media.Providers {
		info, ok := pm.GetProviderInfo(name)
		if !ok {
			continue
		}

		providers = append(providers, MediaProvider{
			Name:        info.Name,
			DisplayName: info.GetDisplayName(),
			Value:       value,
		})
	}

	sort.SliceStable(providers, func(i, j int) bool {
		return natural.Less(providers[i].Name, providers[j].Name)
	})

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
			user.HasData = true
		}
	}

	var release *MediaRelease
	if media.Release.Valid {
		data := media.Release.Data

		status := types.MediaPartReleaseStatusUnknown
		startDate, _ := time.Parse(time.RFC3339, data.StartDate)
		startDate = startDate.UTC()

		t := time.Now()

		currentPart := utils.CurrentPart(startDate, data.DelayDays, data.IntervalDays)
		currentPart += data.PartOffset

		var nextAiring *time.Time

		if t.Before(startDate) {
			status = types.MediaPartReleaseStatusWaiting
			nextAiring = &startDate
		} else {
			status = types.MediaPartReleaseStatusRunning

			t := utils.NextAiringDate(startDate, data.DelayDays, data.IntervalDays)
			nextAiring = &t

			if data.NumExpectedParts > 0 && currentPart >= data.NumExpectedParts {
				currentPart = data.NumExpectedParts
				status = types.MediaPartReleaseStatusCompleted
				nextAiring = nil
			}
		}

		formatNextTime := func() *string {
			if nextAiring != nil {
				s := nextAiring.Format(time.RFC3339)
				return &s
			}

			return nil
		}

		typ := types.MediaPartReleaseType(data.Type)
		if !types.IsValidMediaPartReleaseType(typ) {
			typ = types.MediaPartReleaseTypeNotConfirmed
		}

		release = &MediaRelease{
			ReleaseType:      typ,
			StartDate:        startDate.Format(time.RFC3339),
			NumExpectedParts: data.NumExpectedParts,
			PartOffset:       data.PartOffset,
			IntervalDays:     data.IntervalDays,
			DelayDays:        data.DelayDays,
			Status:           status,
			CurrentPart:      currentPart,
			NextAiring:       formatNextTime(),
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
		Creators:     utils.FixNilArrayToEmpty(media.Creators.Data),
		Tags:         utils.FixNilArrayToEmpty(media.Tags.Data),
		AiringSeason: utils.SqlNullToStringPtr(media.AiringSeason),
		StartDate:    utils.SqlNullToStringPtr(media.StartDate),
		EndDate:      utils.SqlNullToStringPtr(media.EndDate),
		CoverUrl:     coverUrl,
		BannerUrl:    bannerUrl,
		LogoUrl:      logoUrl,
		Providers:    providers,
		User:         user,
		Release:      release,
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

	Tags     []string `json:"tags"`
	Creators []string `json:"creators"`

	CollectionId   string `json:"collectionId,omitempty"`
	CollectionName string `json:"collectionName,omitempty"`
}

func (b *CreateMediaBody) Transform() {
	b.Title = anvil.String(b.Title)
	b.Description = anvil.String(b.Description)

	b.Score = utils.Clamp(b.Score, 0.0, 10.0)
	b.AiringSeason = utils.TransformStringSlug(b.AiringSeason)

	b.PartCount = utils.Min(b.PartCount, 0)

	b.StartDate = anvil.String(b.StartDate)
	b.EndDate = anvil.String(b.EndDate)

	b.Tags = utils.TransformSlugArray(b.Tags)
	b.Creators = utils.TransformSlugArray(b.Creators)
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

	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`

	Score        *float64 `json:"score,omitempty"`
	Status       *string  `json:"status,omitempty"`
	Rating       *string  `json:"rating,omitempty"`
	AiringSeason *string  `json:"airingSeason,omitempty"`

	StartDate *string `json:"startDate,omitempty"`
	EndDate   *string `json:"endDate,omitempty"`

	AdminStatus *string `json:"adminStatus,omitempty"`

	Tags     *[]string `json:"tags,omitempty"`
	Creators *[]string `json:"creators,omitempty"`
}

func (b *EditMediaBody) Transform() {
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

	if b.Creators != nil {
		*b.Creators = utils.TransformSlugArray(*b.Creators)
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

type SetMediaReleaseBody struct {
	MediaId          string                     `json:"mediaId"`
	ReleaseType      types.MediaPartReleaseType `json:"releaseType"`
	StartDate        string                     `json:"startDate"`
	NumExpectedParts int                        `json:"numExpectedParts"`
	IntervalDays     int                        `json:"intervalDays"`
	DelayDays        int                        `json:"delayDays"`
}

func (b *SetMediaReleaseBody) Transform() {
	b.NumExpectedParts = utils.Min(b.NumExpectedParts, 0)
	b.IntervalDays = utils.Min(b.IntervalDays, 0)
	b.DelayDays = utils.Min(b.DelayDays, 0)
}

func (b SetMediaReleaseBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.MediaId, validate.Required),
		validate.Field(&b.ReleaseType, validate.Required, validate.By(types.ValidateMediaPartReleaseType)),
		validate.Field(&b.StartDate, validate.Required, validate.Date(time.RFC3339)),
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
				pm := app.ProviderManager()

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
					res.Media[i] = ConvertDBMedia(c, pm, userId != nil, m)
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
				pm := app.ProviderManager()

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
					Media: ConvertDBMedia(c, pm, userId != nil, media),
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
			Name:         "CreateMedia",
			Method:       http.MethodPost,
			Path:         "/media",
			ResponseType: CreateMedia{},
			BodyType:     CreateMediaBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

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
					Type:  ty,
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

				for _, tag := range body.Creators {
					err := app.DB().CreateTag(ctx, tag, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				for _, tag := range body.Creators {
					err := app.DB().AddCreatorToMedia(ctx, id, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				if body.CollectionId != "" {
					col, err := app.DB().GetCollectionById(ctx, body.CollectionId)
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
						Position:     0,
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

				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

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

				if body.Creators != nil {
					err := app.DB().RemoveAllCreatorsFromMedia(ctx, dbMedia.Id)
					if err != nil {
						return nil, err
					}

					creators := *body.Creators

					for _, tag := range creators {
						err := app.DB().CreateTag(ctx, tag, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}

					for _, tag := range creators {
						err := app.DB().AddCreatorToMedia(ctx, id, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				return nil, nil
			},
		},

		pyrin.FormApiHandler{
			Name:         "ChangeMediaImages",
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

				_, err := User(app, c, HasEditPrivilege)
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
							p := path.Join(mediaDir.Images(), old.String)
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

						outFile, err := os.OpenFile(path.Join(mediaDir.Images(), name+ext), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

				changes.CoverFile, err = test(dbMedia.CoverFile, "cover")
				if err != nil {
					return nil, err
				}

				changes.LogoFile, err = test(dbMedia.LogoFile, "logo")
				if err != nil {
					return nil, err
				}

				changes.BannerFile, err = test(dbMedia.BannerFile, "banner")
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
				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

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
				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

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
				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

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
			HandlerFunc: func(c pyrin.Context) (any, error) {
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
					List:         val.List,
					Part:         utils.Int64PtrToSqlNull(val.Part),
					IsRevisiting: val.IsRevisiting > 0,
					Score:        utils.Int64PtrToSqlNull(val.Score),
				}

				if body.List != nil {
					data.List = *body.List
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
			Name:         "DeleteMediaUserData",
			Method:       http.MethodDelete,
			Path:         "/media/:id/user",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

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

				err = app.DB().DeleteMediaUserData(ctx, media.Id, user.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SetMediaRelease",
			Method:       http.MethodPost,
			Path:         "/media/:id/release",
			ResponseType: nil,
			BodyType:     SetMediaReleaseBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				body, err := pyrin.Body[SetMediaReleaseBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				media, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				t, err := time.Parse(time.RFC3339, body.StartDate)
				if err != nil {
					return nil, err
				}

				err = app.DB().SetMediaPartRelease(ctx, media.Id, database.SetMediaPartRelease{
					Type:             body.ReleaseType,
					StartDate:        t.Format(time.RFC3339),
					NumExpectedParts: body.NumExpectedParts,
					IntervalDays:     body.IntervalDays,
					DelayDays:        body.DelayDays,
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "DeleteMediaRelease",
			Method:       http.MethodDelete,
			Path:         "/media/:id/release",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				media, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveMediaPartRelease(ctx, media.Id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
