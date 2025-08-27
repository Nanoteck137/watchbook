package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type ReleaseUser struct {
	HasData      bool                `json:"hasData"`
	List         types.MediaUserList `json:"list"`
	Score        *int64              `json:"score"`
	CurrentPart  *int64              `json:"currentPart"`
	RevisitCount *int64              `json:"revisitCount"`
	IsRevisiting bool                `json:"isRevisiting"`
}

type Release struct {
	MediaId string `json:"mediaId"`

	NumExpectedParts int    `json:"numExpectedParts"`
	CurrentPart      int    `json:"currentPart"`
	NextAiring       string `json:"nextAiring"`
	IntervalDays     int    `json:"intervalDays"`
	IsActive         int    `json:"isActive"`

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

func ConvertDBRelease(c pyrin.Context, hasUser bool, media database.FullMediaPartRelease) Release {
	// TODO(patrik): Add default cover
	var coverUrl *string
	var bannerUrl *string
	var logoUrl *string

	if media.MediaCoverFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", media.MediaId, path.Base(media.MediaCoverFile.String)))
		coverUrl = &url
	}

	if media.MediaLogoFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", media.MediaId, path.Base(media.MediaLogoFile.String)))
		logoUrl = &url
	}

	if media.MediaBannerFile.Valid {
		url := ConvertURL(c, fmt.Sprintf("/files/media/%s/images/%s", media.MediaId, path.Base(media.MediaBannerFile.String)))
		bannerUrl = &url
	}

	var user *MediaUser
	if hasUser {
		user = &MediaUser{}

		if media.MediaUserData.Valid {
			val := media.MediaUserData.Data
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


	return Release{
		MediaId:     media.MediaId,
		Title:       media.MediaTitle,
		Description: utils.SqlNullToStringPtr(media.MediaDescription),
		TmdbId:      media.MediaTmdbId.String,
		// TODO(patrik): Fix
		// ImdbId:           media.MediaImdbId.String,
		MalId:            media.MediaMalId.String,
		AnilistId:        media.MediaAnilistId.String,
		MediaType:        media.MediaType,
		Score:            utils.SqlNullToFloat64Ptr(media.MediaScore),
		Status:           media.MediaStatus,
		Rating:           media.MediaRating,
		PartCount:        media.MediaPartCount.Int64,
		Creators:         utils.FixNilArrayToEmpty(media.MediaCreators.Data),
		Tags:             utils.FixNilArrayToEmpty(media.MediaTags.Data),
		AiringSeason:     utils.SqlNullToStringPtr(media.MediaAiringSeason),
		StartDate:        utils.SqlNullToStringPtr(media.MediaStartDate),
		EndDate:          utils.SqlNullToStringPtr(media.MediaEndDate),
		CoverUrl:         coverUrl,
		BannerUrl:        bannerUrl,
		LogoUrl:          logoUrl,
		User:             user,
		NumExpectedParts: media.NumExpectedParts,
		CurrentPart:      media.CurrentPart,
		NextAiring:       media.NextAiring,
		IntervalDays:     media.IntervalDays,
		IsActive:         media.IsActive,
	}
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
	)
}
