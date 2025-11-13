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
	"strconv"
	"time"

	"github.com/kr/pretty"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/provider/myanimelist"
	"github.com/nanoteck137/watchbook/provider/sonarr"
	"github.com/nanoteck137/watchbook/provider/tmdb"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Show struct {
	Id   string         `json:"id"`
	Type types.ShowType `json:"type"`

	Name string `json:"name"`

	CoverUrl  *string `json:"coverUrl"`
	LogoUrl   *string `json:"logoUrl"`
	BannerUrl *string `json:"bannerUrl"`

	DefaultProvider *string         `json:"defaultProvider"`
	Providers       []ProviderValue `json:"providers"`
}

type GetShows struct {
	Page  types.Page `json:"page"`
	Shows []Show     `json:"shows"`
}

type GetShowById struct {
	Show
}

func ConvertDBShow(c pyrin.Context, pm *provider.ProviderManager, hasUser bool, show database.Show) Show {
	// TODO(patrik): Add default cover
	var coverUrl *string
	var bannerUrl *string
	var logoUrl *string

	if show.CoverFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/shows/%s/images/%s", show.Id, path.Base(show.CoverFile.String)))
		coverUrl = &url
	}

	if show.LogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/shows/%s/images/%s", show.Id, path.Base(show.LogoFile.String)))
		logoUrl = &url
	}

	if show.BannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/shows/%s/images/%s", show.Id, path.Base(show.BannerFile.String)))
		bannerUrl = &url
	}

	return Show{
		Id:              show.Id,
		Type:            show.Type,
		Name:            show.Name,
		CoverUrl:        coverUrl,
		LogoUrl:         logoUrl,
		BannerUrl:       bannerUrl,
		DefaultProvider: utils.SqlNullToStringPtr(show.DefaultProvider),
		Providers:       createProviderValues(pm, show.Providers),
	}
}

type ShowSeason struct {
	Num    int    `json:"num"`
	ShowId string `json:"showId"`

	Name       string `json:"name"`
	SearchSlug string `json:"searchSlug"`

	Items []ShowSeasonItem `json:"items"`
}

func ConvertDBShowSeason(c pyrin.Context, item database.ShowSeason) ShowSeason {
	return ShowSeason{
		Num:        item.Num,
		ShowId:     item.ShowId,
		Name:       item.Name,
		SearchSlug: item.SearchSlug,
		Items:      []ShowSeasonItem{},
	}
}

type ShowSeasonItem struct {
	ShowSeasonNum int    `json:"showSeasonNum"`
	ShowId        string `json:"showId"`
	MediaId       string `json:"mediaId"`

	Position int `json:"position"`

	Title       string  `json:"title"`
	Description *string `json:"description"`

	Type         types.MediaType   `json:"type"`
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

	DefaultProvider *string         `json:"defaultProvider"`
	Providers       []ProviderValue `json:"providers"`

	User    *MediaUser    `json:"user,omitempty"`
	Release *MediaRelease `json:"release"`
}

func ConvertDBShowSeasonItem(c pyrin.Context, pm *provider.ProviderManager, hasUser bool, item database.FullShowSeasonItem) ShowSeasonItem {
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

	var release *MediaRelease
	if item.MediaRelease.Valid {
		data := item.MediaRelease.Data

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

	return ShowSeasonItem{
		ShowSeasonNum: item.ShowSeasonNum,
		ShowId:        item.ShowId,
		MediaId:       item.MediaId,

		Position: item.Position,

		Title:       item.MediaTitle,
		Description: utils.SqlNullToStringPtr(item.MediaDescription),

		Type:         item.MediaType,
		Score:        utils.SqlNullToFloat64Ptr(item.MediaScore),
		Status:       item.MediaStatus,
		Rating:       item.MediaRating,
		PartCount:    item.MediaPartCount.Int64,
		AiringSeason: utils.SqlNullToStringPtr(item.MediaAiringSeason),

		StartDate: utils.SqlNullToStringPtr(item.MediaStartDate),
		EndDate:   utils.SqlNullToStringPtr(item.MediaEndDate),

		Creators: utils.FixNilArrayToEmpty(item.MediaCreators.Data),
		Tags:     utils.FixNilArrayToEmpty(item.MediaTags.Data),

		CoverUrl:  coverUrl,
		BannerUrl: bannerUrl,
		LogoUrl:   logoUrl,

		DefaultProvider: utils.SqlNullToStringPtr(item.MediaDefaultProvider),
		Providers:       createProviderValues(pm, item.MediaProviders),

		User:    user,
		Release: release,
	}
}

type GetShowSeasons struct {
	Seasons []ShowSeason `json:"seasons"`
}

type CreateShow struct {
	Id string `json:"id"`
}

type CreateShowBody struct {
	Type string `json:"type"`

	Name string `json:"name"`

	CoverUrl  string `json:"coverUrl"`
	BannerUrl string `json:"bannerUrl"`
	LogoUrl   string `json:"logoUrl"`
}

func (b *CreateShowBody) Transform() {
	b.Name = anvil.String(b.Name)

	b.CoverUrl = anvil.String(b.CoverUrl)
	b.BannerUrl = anvil.String(b.BannerUrl)
	b.LogoUrl = anvil.String(b.LogoUrl)
}

func (b CreateShowBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Type, validate.Required, validate.By(types.ValidateShowType)),
		validate.Field(&b.Name, validate.Required),
	)
}

type EditShowBody struct {
	Type *string `json:"type,omitempty"`

	Name *string `json:"name,omitempty"`

	CoverUrl  *string `json:"coverUrl,omitempty"`
	BannerUrl *string `json:"bannerUrl,omitempty"`
	LogoUrl   *string `json:"logoUrl,omitempty"`
}

func (b *EditShowBody) Transform() {
	b.Name = anvil.StringPtr(b.Name)

	b.CoverUrl = anvil.StringPtr(b.CoverUrl)
	b.BannerUrl = anvil.StringPtr(b.BannerUrl)
	b.LogoUrl = anvil.StringPtr(b.LogoUrl)
}

func (b EditShowBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Type, validate.Required.When(b.Type != nil), validate.By(types.ValidateShowType)),

		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),

		validate.Field(&b.CoverUrl, validate.Required.When(b.CoverUrl != nil)),
		validate.Field(&b.BannerUrl, validate.Required.When(b.BannerUrl != nil)),
		validate.Field(&b.LogoUrl, validate.Required.When(b.LogoUrl != nil)),
	)
}

type AddShowSeasonBody struct {
	Num        int    `json:"num"`
	Name       string `json:"name"`
	SearchSlug string `json:"searchSlug"`
}

func (b *AddShowSeasonBody) Transform() {
	b.Num = utils.Min(b.Num, 0)
	b.Name = anvil.String(b.Name)
	b.SearchSlug = utils.Slug(b.SearchSlug)
}

func (b AddShowSeasonBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type EditShowSeasonBody struct {
	Num        *int    `json:"num,omitempty"`
	Name       *string `json:"name,omitempty"`
	SearchSlug *string `json:"searchSlug,omitempty"`
}

func (b *EditShowSeasonBody) Transform() {
	if b.Num != nil {
		*b.Num = utils.Min(*b.Num, 0)
	}

	b.Name = anvil.StringPtr(b.Name)

	if b.SearchSlug != nil {
		*b.SearchSlug = utils.Slug(*b.SearchSlug)
	}
}

func (b EditShowSeasonBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

type AddShowSeasonItemBody struct {
	MediaId  string `json:"mediaId"`
	Position int    `json:"position"`
}

type EditShowSeasonItemBody struct {
	Position *int `json:"position,omitempty"`
}

type GetShowSeasonEpisodes struct {
	Episodes []MediaPart `json:"episodes"`
}

type GetShowSeason struct {
	ShowSeason
}

func InstallShowHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetShows",
			Method:       http.MethodGet,
			Path:         "/shows",
			ResponseType: GetShows{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				pm := app.ProviderManager()

				q := c.Request().URL.Query()
				opts := getPageOptions(q)

				ctx := context.TODO()

				filterStr := q.Get("filter")
				sortStr := q.Get("sort")
				shows, page, err := app.DB().GetPagedShows(ctx, filterStr, sortStr, opts)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetShows{
					Page:  page,
					Shows: make([]Show, len(shows)),
				}

				for i, col := range shows {
					res.Shows[i] = ConvertDBShow(c, pm, false, col)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetShowById",
			Method:       http.MethodGet,
			Path:         "/shows/:id",
			ResponseType: GetShowById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				pm := app.ProviderManager()

				show, err := app.DB().GetShowById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowNotFound()
					}

					return nil, err
				}

				return GetShowById{
					Show: ConvertDBShow(c, pm, false, show),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetShowSeasons",
			Method:       http.MethodGet,
			Path:         "/shows/:id/seasons",
			ResponseType: GetShowSeasons{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				pm := app.ProviderManager()

				var userId *string
				if user, err := User(app, c); err == nil {
					userId = &user.Id
				}

				ctx := c.Request().Context()

				show, err := app.DB().GetShowById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowNotFound()
					}

					return nil, err
				}

				seasons, err := app.DB().GetAllShowSeasonsByShowId(ctx, show.Id)
				if err != nil {
					return nil, err
				}

				res := GetShowSeasons{
					Seasons: make([]ShowSeason, 0, len(seasons)),
				}

				_ = pm
				_ = userId

				for _, season := range seasons {
					s := ConvertDBShowSeason(c, season)

					// items, err := app.DB().GetFullAllShowSeasonItemsByShowSeason(ctx, userId, season.Num, season.ShowId)
					// if err != nil {
					// 	return nil, err
					// }
					//
					// for _, item := range items {
					// 	s.Items = append(s.Items, ConvertDBShowSeasonItem(c, pm, userId != nil, item))
					// }

					res.Seasons = append(res.Seasons, s)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateShow",
			Method:       http.MethodPost,
			Path:         "/shows",
			ResponseType: CreateShow{},
			BodyType:     CreateShowBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreateShowBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				ty := types.ShowType(body.Type)

				id := utils.CreateShowId()

				showDir := app.WorkDir().ShowDirById(id)

				err = showDir.Create()
				if err != nil {
					return nil, err
				}

				coverFile := ""
				bannerFile := ""
				logoFile := ""

				if body.CoverUrl != "" {
					p, err := utils.DownloadImageHashed(body.CoverUrl, showDir.Images())
					if err == nil {
						coverFile = path.Base(p)
					} else {
						app.Logger().Error("failed to download cover image for show", "err", err)
					}
				}

				if body.BannerUrl != "" {
					p, err := utils.DownloadImageHashed(body.BannerUrl, showDir.Images())
					if err == nil {
						bannerFile = path.Base(p)
					} else {
						app.Logger().Error("failed to download banner image for show", "err", err)
					}
				}

				if body.LogoUrl != "" {
					p, err := utils.DownloadImageHashed(body.LogoUrl, showDir.Images())
					if err == nil {
						logoFile = path.Base(p)
					} else {
						app.Logger().Error("failed to download logo image for show", "err", err)
					}
				}

				_, err = app.DB().CreateShow(ctx, database.CreateShowParams{
					Id:   id,
					Type: ty,
					Name: body.Name,
					CoverFile: sql.NullString{
						String: coverFile,
						Valid:  coverFile != "",
					},
					BannerFile: sql.NullString{
						String: bannerFile,
						Valid:  bannerFile != "",
					},
					LogoFile: sql.NullString{
						String: logoFile,
						Valid:  logoFile != "",
					},
				})
				if err != nil {
					return nil, err
				}

				return CreateShow{
					Id: id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditShow",
			Method:       http.MethodPatch,
			Path:         "/shows/:id",
			ResponseType: nil,
			BodyType:     EditShowBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[EditShowBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbShow, err := app.DB().GetShowById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowNotFound()
					}

					return nil, err
				}

				showDir := app.WorkDir().ShowDirById(dbShow.Id)

				changes := database.ShowChanges{}

				if body.Type != nil {
					t := types.ShowType(*body.Type)

					changes.Type = database.Change[types.ShowType]{
						Value:   t,
						Changed: t != dbShow.Type,
					}
				}

				if body.Name != nil {
					changes.Name = database.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != dbShow.Name,
					}
				}

				if body.CoverUrl != nil {
					p, err := utils.DownloadImageHashed(*body.CoverUrl, showDir.Images())
					if err == nil {
						n := path.Base(p)
						changes.CoverFile = database.Change[sql.NullString]{
							Value: sql.NullString{
								String: n,
								Valid:  n != "",
							},
							Changed: n != dbShow.CoverFile.String,
						}
					} else {
						app.Logger().Error("failed to download cover image for show", "err", err)
					}
				}

				if body.BannerUrl != nil {
					p, err := utils.DownloadImageHashed(*body.BannerUrl, showDir.Images())
					if err == nil {
						n := path.Base(p)
						changes.BannerFile = database.Change[sql.NullString]{
							Value: sql.NullString{
								String: n,
								Valid:  n != "",
							},
							Changed: n != dbShow.BannerFile.String,
						}
					} else {
						app.Logger().Error("failed to download banner image for show", "err", err)
					}
				}

				if body.LogoUrl != nil {
					p, err := utils.DownloadImageHashed(*body.LogoUrl, showDir.Images())
					if err == nil {
						n := path.Base(p)
						changes.LogoFile = database.Change[sql.NullString]{
							Value: sql.NullString{
								String: n,
								Valid:  n != "",
							},
							Changed: n != dbShow.LogoFile.String,
						}
					} else {
						app.Logger().Error("failed to download logo image for show", "err", err)
					}
				}

				err = app.DB().UpdateShow(ctx, dbShow.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "DeleteShow",
			Method:       http.MethodDelete,
			Path:         "/shows/:id",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				_, err := User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbShow, err := app.DB().GetShowById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveShow(ctx, dbShow.Id)
				if err != nil {
					return nil, err
				}

				dir := app.WorkDir().ShowDirById(dbShow.Id)
				err = os.RemoveAll(dir.String())
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		// TODO(patrik): Change this, hash and don't remove old images
		pyrin.FormApiHandler{
			Name:         "ChangeShowImages",
			Method:       http.MethodPatch,
			Path:         "/shows/:id/images",
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

				dbShow, err := app.DB().GetShowById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowNotFound()
					}

					return nil, err
				}

				showDir := app.WorkDir().ShowDirById(id)
				err = showDir.Create()
				if err != nil {
					return nil, err
				}

				// TODO(patrik): Move to helper.go
				upload := func(name string) (database.Change[sql.NullString], error) {
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

						f, err := file.Open()
						// TODO(patrik): Better error
						if err != nil {
							return database.Change[sql.NullString]{}, err
						}
						defer f.Close()

						data, err := io.ReadAll(f)
						if err != nil {
							return database.Change[sql.NullString]{}, fmt.Errorf("failed to read body: %w", err)
						}

						filename, err := utils.WriteHashedFile(data, showDir.Images(), ext)
						if err != nil {
							return database.Change[sql.NullString]{}, err
						}

						return database.Change[sql.NullString]{
							Value: sql.NullString{
								String: filename,
								Valid:  true,
							},
							Changed: true,
						}, nil
					}

					return database.Change[sql.NullString]{}, nil
				}

				changes := database.ShowChanges{}

				changes.CoverFile, err = upload("cover")
				if err != nil {
					return nil, err
				}

				changes.LogoFile, err = upload("logo")
				if err != nil {
					return nil, err
				}

				changes.BannerFile, err = upload("banner")
				if err != nil {
					return nil, err
				}

				err = app.DB().UpdateShow(ctx, dbShow.Id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AddShowSeason",
			Method:       http.MethodPost,
			Path:         "/shows/:id/seasons",
			ResponseType: nil,
			BodyType:     AddShowSeasonBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[AddShowSeasonBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbShow, err := app.DB().GetShowById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowNotFound()
					}

					return nil, err
				}

				searchSlug := body.SearchSlug
				if searchSlug == "" {
					searchSlug = utils.Slug(body.Name)
				}

				err = app.DB().CreateShowSeason(ctx, database.CreateShowSeasonParams{
					Num:        body.Num,
					ShowId:     dbShow.Id,
					Name:       body.Name,
					SearchSlug: searchSlug,
				})
				if err != nil {
					// TODO(patrik): Better handling of error
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetShowSeason",
			Method:       http.MethodGet,
			Path:         "/shows/:id/seasons/:seasonNum",
			ResponseType: GetShowSeason{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				num, err := strconv.Atoi(c.Param("seasonNum"))
				if err != nil {
					// TODO(patrik): Better error?
					return nil, ShowSeasonNotFound()
				}

				pm := app.ProviderManager()

				var userId *string
				if user, err := User(app, c); err == nil {
					userId = &user.Id
				}

				ctx := context.Background()

				dbShowSeason, err := app.DB().GetShowSeasonById(ctx, num, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowSeasonNotFound()
					}

					return nil, err
				}

				res := GetShowSeason{
					ShowSeason: ConvertDBShowSeason(c, dbShowSeason),
				}

				items, err := app.DB().GetFullAllShowSeasonItemsByShowSeason(ctx, userId, res.ShowSeason.Num, res.ShowSeason.ShowId)
				if err != nil {
					return nil, err
				}

				for _, item := range items {
					res.ShowSeason.Items = append(res.ShowSeason.Items, ConvertDBShowSeasonItem(c, pm, userId != nil, item))
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "RemoveShowSeason",
			Method:       http.MethodDelete,
			Path:         "/shows/:id/seasons/:seasonNum",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				num, err := strconv.Atoi(c.Param("seasonNum"))
				if err != nil {
					// TODO(patrik): Better error?
					return nil, ShowSeasonNotFound()
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				item, err := app.DB().GetShowSeasonById(ctx, num, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowSeasonNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveShowSeason(ctx, item.Num, item.ShowId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditShowSeason",
			Method:       http.MethodPatch,
			Path:         "/shows/:id/seasons/:seasonNum",
			ResponseType: nil,
			BodyType:     EditShowSeasonBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				num, err := strconv.Atoi(c.Param("seasonNum"))
				if err != nil {
					// TODO(patrik): Better error?
					return nil, ShowSeasonNotFound()
				}

				body, err := pyrin.Body[EditShowSeasonBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				item, err := app.DB().GetShowSeasonById(ctx, num, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowSeasonNotFound()
					}

					return nil, err
				}

				changes := database.ShowSeasonChanges{}

				if body.Num != nil {
					changes.Num = database.Change[int]{
						Value:   *body.Num,
						Changed: *body.Num != item.Num,
					}
				}

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

				err = app.DB().UpdateShowSeason(ctx, item.Num, item.ShowId, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "AddShowSeasonItem",
			Method:       http.MethodPost,
			Path:         "/shows/:id/seasons/:seasonNum/items",
			ResponseType: nil,
			BodyType:     AddShowSeasonItemBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				num, err := strconv.Atoi(c.Param("seasonNum"))
				if err != nil {
					// TODO(patrik): Better error?
					return nil, ShowSeasonNotFound()
				}

				body, err := pyrin.Body[AddShowSeasonItemBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				dbShowSeason, err := app.DB().GetShowSeasonById(ctx, num, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowSeasonNotFound()
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

				err = app.DB().CreateShowSeasonItem(ctx, database.CreateShowSeasonItemParams{
					ShowSeasonNum: dbShowSeason.Num,
					ShowId:        dbShowSeason.ShowId,
					MediaId:       dbMedia.Id,
					Position:      body.Position,
				})
				if err != nil {
					// TODO(patrik): Better handling of error
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "RemoveShowSeasonItem",
			Method:       http.MethodDelete,
			Path:         "/shows/:id/seasons/:seasonNum/items/:mediaId",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				num, err := strconv.Atoi(c.Param("seasonNum"))
				if err != nil {
					// TODO(patrik): Better error?
					return nil, ShowSeasonNotFound()
				}
				mediaId := c.Param("mediaId")

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				item, err := app.DB().GetShowSeasonItemById(ctx, num, id, mediaId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowSeasonItemNotFound()
					}

					return nil, err
				}

				err = app.DB().RemoveShowSeasonItem(ctx, item.ShowSeasonNum, item.ShowId, item.MediaId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditShowSeasonItem",
			Method:       http.MethodPatch,
			Path:         "/shows/:id/seasons/:seasonNum/items/:mediaId",
			ResponseType: nil,
			BodyType:     EditShowSeasonItemBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				num, err := strconv.Atoi(c.Param("seasonNum"))
				if err != nil {
					// TODO(patrik): Better error?
					return nil, ShowSeasonNotFound()
				}
				mediaId := c.Param("mediaId")

				body, err := pyrin.Body[EditShowSeasonItemBody](c)
				if err != nil {
					return nil, err
				}

				_, err = User(app, c, HasEditPrivilege)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				item, err := app.DB().GetShowSeasonItemById(ctx, num, id, mediaId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowSeasonItemNotFound()
					}

					return nil, err
				}

				changes := database.ShowSeasonItemChanges{}

				if body.Position != nil {
					changes.Position = database.Change[int]{
						Value:   *body.Position,
						Changed: *body.Position != item.Position,
					}
				}

				err = app.DB().UpdateShowSeasonItem(ctx, item.ShowSeasonNum, item.ShowId, item.MediaId, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetShowSeasonEpisodes",
			Method:       http.MethodGet,
			Path:         "/shows/:id/seasons/:seasonNum/episodes",
			ResponseType: GetShowSeasonEpisodes{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")
				num, err := strconv.Atoi(c.Param("seasonNum"))
				if err != nil {
					// TODO(patrik): Better error?
					return nil, ShowSeasonNotFound()
				}

				ctx := context.Background()

				dbShowSeason, err := app.DB().GetShowSeasonById(ctx, num, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ShowSeasonNotFound()
					}

					return nil, err
				}

				items, err := app.DB().GetAllShowSeasonItemsByShowSeason(ctx, dbShowSeason.Num, dbShowSeason.ShowId)
				if err != nil {
					return nil, err
				}

				res := GetShowSeasonEpisodes{
					Episodes: []MediaPart{},
				}

				index := int64(1)

				for _, item := range items {
					parts, err := app.DB().GetMediaPartsByMediaId(ctx, item.MediaId)
					if err != nil {
						return nil, err
					}

					for _, part := range parts {
						res.Episodes = append(res.Episodes, MediaPart{
							Index:       index,
							MediaId:     part.MediaId,
							Name:        part.Name,
							ReleaseDate: utils.SqlNullToStringPtr(part.ReleaseDate),
						})

						index += 1
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "Test",
			Method: http.MethodPost,
			Path:   "/test",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				sonarrApiKey := os.Getenv("TEMP_SONARR_API_KEY")
				sc := sonarr.NewClient("https://sonarr.nanoteck137.net", sonarrApiKey)

				ctx := context.Background()

				series, err := sc.GetAllSeries()
				if err != nil {
					return nil, err
				}

				idx := 0
				for i, s := range series {
					if s.ID == 16 {
						idx = i 
						break
					}
				}

				serie := series[idx]


				id := utils.CreateShowId()

				showDir := app.WorkDir().ShowDirById(id)
				dirs := []string{
					showDir.String(),
					showDir.Images(),
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						return nil, err
					}
				}

				// providerIds := ember.KVStore{}
				// maps.Copy(providerIds, data.ExtraProviderIds)
				// providerIds[providerName] = data.ProviderId

				fmt.Printf("serie.TmdbId: %v\n", serie.TmdbId)

				m, err := app.ProviderManager().GetCollection(ctx, tmdb.TvProviderName, strconv.Itoa(serie.TmdbId))
				if err != nil {
					fmt.Printf("err: %v\n", err)
					return nil, err
				}

				pretty.Println(m)

				coverFilename := ""
				bannerFilename := ""
				logoFilename := ""

				coverUrl := utils.NullToDefault(m.CoverUrl)
				bannerUrl := utils.NullToDefault(m.BannerUrl)
				logoUrl := utils.NullToDefault(m.LogoUrl)

				for _, img := range serie.Images {
					if img.CoverType == "poster" && coverUrl == "" {
						coverUrl = img.RemoteURL
					}

					if img.CoverType == "banner" && bannerUrl == "" {
						bannerUrl = img.RemoteURL
					}

					if img.CoverType == "clearlogo" && logoUrl == "" {
						logoUrl = img.RemoteURL
					}
				}

				if coverUrl != "" {
					p, err := utils.DownloadImageHashed(coverUrl, showDir.Images())
					if err == nil {
						coverFilename = path.Base(p)
					} else {
						app.Logger().Error("failed to download cover image for collection", "err", err)
					}
				}

				if bannerUrl != "" {
					p, err := utils.DownloadImageHashed(bannerUrl, showDir.Images())
					if err == nil {
						bannerFilename = path.Base(p)
					} else {
						app.Logger().Error("failed to download banner image for collection", "err", err)
					}
				}

				if logoUrl != "" {
					p, err := utils.DownloadImageHashed(logoUrl, showDir.Images())
					if err == nil {
						logoFilename = path.Base(p)
					} else {
						app.Logger().Error("failed to download logo image for collection", "err", err)
					}
				}

				ty := types.ShowTypeUnknown
				// switch data.Type {
				// case types.CollectionTypeAnime:
				// 	ty = types.ShowTypeAnime
				// case types.CollectionTypeSeries:
				// 	ty = types.ShowTypeTVSeries
				// }

				_, err = app.DB().CreateShow(ctx, database.CreateShowParams{
					Id:         id,
					Type:       ty,
					Name:       serie.Title,
					SearchSlug: utils.Slug(serie.Title),
					CoverFile: sql.NullString{
						String: coverFilename,
						Valid:  coverFilename != "",
					},
					BannerFile: sql.NullString{
						String: bannerFilename,
						Valid:  bannerFilename != "",
					},
					LogoFile: sql.NullString{
						String: logoFilename,
						Valid:  logoFilename != "",
					},
					// DefaultProvider: sql.NullString{
					// 	String: providerName,
					// 	Valid:  providerName != "",
					// },
					// Providers: providerIds,
				})
				if err != nil {
					fmt.Printf("err: %v\n", err)
					return nil, err
				}

				for _, item := range serie.Seasons {
					// mediaId, err := ImportMedia(ctx, app, providerName, item.Id)
					// if err != nil {
					// 	return nil, err
					// }

					name := fmt.Sprintf("Season %d", item.SeasonNumber)

					err = app.DB().CreateShowSeason(ctx, database.CreateShowSeasonParams{
						Num:        item.SeasonNumber,
						ShowId:     id,
						Name:       name,
						SearchSlug: utils.Slug(name),
					})
					if err != nil {
						return nil, err
					}

					episodes, err := sc.GetEpisodesBySeason(serie.ID, item.SeasonNumber)
					if err != nil {
						return nil, err
					}

					for _, episode := range episodes {
						err := app.DB().CreateShowSeasonPart(ctx, database.CreateShowSeasonPartParams{
							ShowId:       id,
							SeasonNumber: item.SeasonNumber,
							Index:        episode.EpisodeNumber,
							Name:         episode.Title,
							ReleaseDate:  sql.NullString{},
						})
						if err != nil {
							return nil, err
						}
					}

					// err = app.DB().CreateShowSeasonItem(ctx, database.CreateShowSeasonItemParams{
					// 	ShowSeasonNum: item.Position,
					// 	ShowId:        id,
					// 	MediaId:       mediaId,
					// 	Position:      0,
					// })
					// if err != nil {
					// 	return nil, err
					// }
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "Test2",
			Method: http.MethodPost,
			Path:   "/test2",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.Background()

				id, err := ImportMedia(ctx, app, myanimelist.AnimeProviderName, "52807")
				if err != nil {
					return nil, err
				}

				fmt.Printf("id: %v\n", id)

				// id, err := NewImportMedia(ctx, app, myanimelist.AnimeProviderName, "35760", []string{"38524"})
				// if err != nil {
				// 	fmt.Printf("err: %v\n", err)
				// 	return nil, err
				// }
				//
				// fmt.Printf("id: %v\n", id)

				return nil, nil
			},
		},
	)
}
