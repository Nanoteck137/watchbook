package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/validate"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Release struct {
	MediaId string `json:"mediaId"`

	ReleaseStatus    types.MediaPartReleaseStatus `json:"releaseStatus"`
	StartDate        string                       `json:"startDate"`
	NumExpectedParts int                          `json:"numExpectedParts"`
	CurrentPart      int                          `json:"currentPart"`
	NextAiring       string                       `json:"nextAiring"`
	IntervalDays     int                          `json:"intervalDays"`
	DelayDays        int                          `json:"delayDays"`

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

	// StartDate *string `json:"startDate"`
	// EndDate   *string `json:"endDate"`

	Creators []string `json:"creators"`
	Tags     []string `json:"tags"`

	CoverUrl  *string `json:"coverUrl"`
	BannerUrl *string `json:"bannerUrl"`
	LogoUrl   *string `json:"logoUrl"`

	User *MediaUser `json:"user,omitempty"`
}

type GetReleases struct {
	Page     types.Page `json:"page"`
	Releases []Release  `json:"releases"`
}

type GetReleaseById struct {
	Release
}

func NextAiringDate(start time.Time, delayDays, intervalDays int) time.Time {
	effectiveStart := start.Add(time.Duration(delayDays) * 24 * time.Hour)
	now := time.Now().UTC()

	// If the show hasn't started yet, return the effective start date
	if now.Before(effectiveStart) {
		return effectiveStart
	}

	// Calculate how many intervals have passed
	diff := now.Sub(effectiveStart)
	intervalsPassed := int(diff.Hours() / (24 * float64(intervalDays)))

	// Next airing date = effectiveStart + (intervalsPassed + 1) * interval
	nextAiring := effectiveStart.Add(time.Duration(intervalsPassed+1) * time.Duration(intervalDays) * 24 * time.Hour)
	return nextAiring
}

func CurrentEpisode(start time.Time, delayDays, intervalDays int) int {
	effectiveStart := start.Add(time.Duration(delayDays) * 24 * time.Hour)
	now := time.Now().UTC()

	// If current time is before start, episode = 0
	if now.Before(effectiveStart) {
		return 0
	}

	// Calculate elapsed time since effective start
	elapsed := now.Sub(effectiveStart)

	// Calculate how many full intervals have passed (including the first episode at start)
	episodesPassed := int(elapsed.Hours() / (24 * float64(intervalDays)))

	// Current episode = episodesPassed + 1 (because first episode is at start)
	return episodesPassed + 1
}

func ConvertDBRelease(c pyrin.Context, hasUser bool, release database.FullMediaPartRelease) Release {
	// TODO(patrik): Add default cover
	var coverUrl *string
	var bannerUrl *string
	var logoUrl *string

	if release.MediaCoverFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", release.MediaId, path.Base(release.MediaCoverFile.String)))
		coverUrl = &url
	}

	if release.MediaLogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", release.MediaId, path.Base(release.MediaLogoFile.String)))
		logoUrl = &url
	}

	if release.MediaBannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", release.MediaId, path.Base(release.MediaBannerFile.String)))
		bannerUrl = &url
	}

	var user *MediaUser
	if hasUser {
		user = &MediaUser{}

		if release.MediaUserData.Valid {
			val := release.MediaUserData.Data
			user.List = val.List
			user.CurrentPart = val.Part
			user.RevisitCount = val.RevisitCount
			user.Score = val.Score
			user.IsRevisiting = val.IsRevisiting > 0
			user.HasData = true
		}
	}

	// nextAiring, _ := time.Parse(types.MediaDateLayout, media.NextAiring)
	//
	// // d, _ := nextAiring.MarshalText()
	//
	// // d := media.NextAiring
	// d := nextAiring.Format(types.MediaDateLayout)

	t := time.Now()

	status := types.MediaPartReleaseStatusUnknown

	var nextAiring time.Time
	currentPart := CurrentEpisode(release.StartDate, release.DelayDays, release.IntervalDays)

	if t.Before(release.StartDate) {
		status = types.MediaPartReleaseStatusWaiting
		nextAiring = release.StartDate
	} else {
		status = types.MediaPartReleaseStatusRunning

		if release.NumExpectedParts > 0 && currentPart >= release.NumExpectedParts {
			currentPart = release.NumExpectedParts
			status = types.MediaPartReleaseStatusCompleted
		}

		nextAiring = NextAiringDate(release.StartDate, release.DelayDays, release.IntervalDays)
	}

	return Release{
		MediaId:     release.MediaId,
		Title:       release.MediaTitle,
		Description: utils.SqlNullToStringPtr(release.MediaDescription),
		TmdbId:      release.MediaTmdbId.String,
		// TODO(patrik): Fix
		// ImdbId:           media.MediaImdbId.String,
		MalId:        release.MediaMalId.String,
		AnilistId:    release.MediaAnilistId.String,
		MediaType:    release.MediaType,
		Score:        utils.SqlNullToFloat64Ptr(release.MediaScore),
		Status:       release.MediaStatus,
		Rating:       release.MediaRating,
		PartCount:    release.MediaPartCount.Int64,
		Creators:     utils.FixNilArrayToEmpty(release.MediaCreators.Data),
		Tags:         utils.FixNilArrayToEmpty(release.MediaTags.Data),
		AiringSeason: utils.SqlNullToStringPtr(release.MediaAiringSeason),
		// StartDate:        utils.SqlNullToStringPtr(release.MediaStartDate),
		// EndDate:          utils.SqlNullToStringPtr(release.MediaEndDate),
		CoverUrl:         coverUrl,
		BannerUrl:        bannerUrl,
		LogoUrl:          logoUrl,
		User:             user,
		NumExpectedParts: release.NumExpectedParts,
		CurrentPart:      currentPart,
		NextAiring:       nextAiring.Format(time.RFC3339),
		IntervalDays:     release.IntervalDays,
		// IsActive:         release.IsActive,
		ReleaseStatus: status,
		// TODO(patrik): Add
		ImdbId:    "",
		StartDate: release.StartDate.Format(time.RFC3339),
		DelayDays: release.DelayDays,
	}
}

type CreateReleaseBody struct {
	MediaId          string `json:"mediaId"`
	StartDate        string `json:"startDate"`
	NumExpectedParts int    `json:"numExpectedParts"`
	IntervalDays     int    `json:"intervalDays"`
	DelayDays        int    `json:"delayDays"`
}

func (b *CreateReleaseBody) Transform() {
	b.NumExpectedParts = utils.Min(b.NumExpectedParts, 0)
	b.IntervalDays = utils.Min(b.IntervalDays, 0)
	b.DelayDays = utils.Min(b.DelayDays, 0)
}

func (b CreateReleaseBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.MediaId, validate.Required),
		validate.Field(&b.StartDate, validate.Required, validate.Date(time.RFC3339)),
	)
}

type EditReleaseBody struct {
	StartDate        *string `json:"startDate,omitempty"`
	NumExpectedParts *int    `json:"numExpectedParts,omitempty"`
	IntervalDays     *int    `json:"intervalDays,omitempty"`
	DelayDays        *int    `json:"delayDays,omitempty"`
}

func (b *EditReleaseBody) Transform() {
	if b.NumExpectedParts != nil {
		*b.NumExpectedParts = utils.Min(*b.NumExpectedParts, 0)
	}

	if b.IntervalDays != nil {
		*b.IntervalDays = utils.Min(*b.IntervalDays, 0)
	}

	if b.DelayDays != nil {
		*b.DelayDays = utils.Min(*b.DelayDays, 0)
	}
}

func (b EditReleaseBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.StartDate, validate.Required.When(b.StartDate != nil), validate.Date(time.RFC3339)),
	)
}

func InstallReleaseHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetReleases",
			Method:       http.MethodGet,
			Path:         "/releases",
			ResponseType: GetReleases{},
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
				releases, p, err := app.DB().GetPagedFullMediaPartReleases(ctx, userId, filterStr, sortStr, opts)
				if err != nil {
					return nil, err
				}

				res := GetReleases{
					Page:     p,
					Releases: make([]Release, len(releases)),
				}

				for i, r := range releases {
					res.Releases[i] = ConvertDBRelease(c, userId != nil, r)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetReleaseById",
			Method:       http.MethodGet,
			Path:         "/releases/:id",
			ResponseType: GetReleaseById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return nil, errors.New("not implemented")

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
			Name:         "CreateRelease",
			Method:       http.MethodPost,
			Path:         "/releases",
			ResponseType: nil,
			BodyType:     CreateReleaseBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[CreateReleaseBody](c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				media, err := app.DB().GetMediaById(ctx, nil, body.MediaId)
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

				err = app.DB().CreateMediaPartRelease(ctx, database.CreateMediaPartReleaseParams{
					MediaId:          media.Id,
					StartDate:        t,
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
			Name:         "EditRelease",
			Method:       http.MethodPatch,
			Path:         "/releases/:id",
			ResponseType: nil,
			BodyType:     EditReleaseBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[EditReleaseBody](c)
				if err != nil {
					return nil, err
				}

				ctx := c.Request().Context()

				release, err := app.DB().GetMediaPartReleaseById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, MediaPartReleaseNotFound()
					}

					return nil, err
				}

				changes := database.MediaPartReleaseChanges{}

				if body.StartDate != nil {
					t, err := time.Parse(time.RFC3339, *body.StartDate)
					if err != nil {
						return nil, err
					}

					changes.StartDate = database.Change[time.Time]{
						Value:   t,
						Changed: t != release.StartDate,
					}
				}

				if body.NumExpectedParts != nil {
					changes.NumExpectedParts = database.Change[int]{
						Value:   *body.NumExpectedParts,
						Changed: *body.NumExpectedParts != release.NumExpectedParts,
					}
				}

				if body.IntervalDays != nil {
					changes.IntervalDays = database.Change[int]{
						Value:   *body.IntervalDays,
						Changed: *body.IntervalDays != release.IntervalDays,
					}
				}

				if body.DelayDays != nil {
					changes.DelayDays = database.Change[int]{
						Value:   *body.DelayDays,
						Changed: *body.DelayDays != release.DelayDays,
					}
				}

				err = app.DB().UpdateMediaPartRelease(ctx, release.MediaId, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
