package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"maps"

	"github.com/kr/pretty"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/kvstore"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type Provider struct {
	Name string `json:"name"`
}

type GetProviders struct {
	Providers []Provider `json:"providers"`
}

type ProviderSearchResult struct {
	Provider   string `json:"provider"`
	ProviderId string `json:"providerId"`
	Title      string `json:"title"`
	ImageUrl   string `json:"imageUrl"`
}

type GetProviderSearch struct {
	SearchResults []ProviderSearchResult `json:"searchResults"`
}

type PostProviderImportMediaBody struct {
	Ids []string `json:"ids"`
}

func fixArr(arr []string) []string {
	if arr == nil {
		return nil
	}

	res := make([]string, 0, len(arr))

	for _, s := range arr {
		s = anvil.String(s)
		if s != "" {
			res = append(res, s)
		}
	}

	if len(res) == 0 {
		return nil
	}

	return res
}

func (b *PostProviderImportMediaBody) Transform() {
	b.Ids = fixArr(b.Ids)
}

func InstallProviderHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetProviders",
			Method:       http.MethodGet,
			Path:         "/providers",
			ResponseType: GetProviders{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providers := app.ProviderManager().GetProviders()

				res := GetProviders{
					Providers: make([]Provider, len(providers)),
				}

				for i, p := range providers {
					res.Providers[i] = Provider{
						Name: p.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "ProviderSearchMedia",
			Method:       http.MethodGet,
			Path:         "/providers/:providerName",
			ResponseType: GetProviderSearch{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")
				url := c.Request().URL
				query := url.Query().Get("query")
				_ = query

				pm := app.ProviderManager()
				_ = pm

				items, err := pm.SearchMedia(context.Background(), providerName, query)
				if err != nil {
					// TODO(patrik): Handle some of the errors the provider gives, like the provider not found
					return nil, err
				}

				pretty.Println(items)

				res := GetProviderSearch{
					SearchResults: make([]ProviderSearchResult, len(items)),
				}

				for i, item := range items {
					res.SearchResults[i] = ProviderSearchResult{
						Provider:   providerName,
						ProviderId: item.ProviderId,
						Title:      item.Title,
						ImageUrl:   item.ImageUrl,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ProviderImportMedia",
			Method: http.MethodPost,
			Path:   "/providers/:providerName/import",
			// ResponseType: PostProviderImportMedia{},
			BodyType: PostProviderImportMediaBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")

				body, err := pyrin.Body[PostProviderImportMediaBody](c)
				if err != nil {
					return nil, err
				}

				pm := app.ProviderManager()

				ctx := context.Background()

				for _, id := range body.Ids {

					_, err := app.DB().GetMediaByProviderId(ctx, nil, providerName, id)
					if err == nil {
						fmt.Printf("id already exists: %v\n", id)
						continue
					}

					if !errors.Is(err, database.ErrItemNotFound) {
						return nil, err
					}

					media, err := pm.GetMedia(ctx, providerName, id)
					if err != nil {
						fmt.Printf("err: %v\n", err)
						// TODO(patrik): Handle err
						return nil, err
					}

					pretty.Println(media)

					// id, err := app.DB().CreateMedia(ctx, database.CreateMediaParams{
					// 	Type: media.Type,
					// 	MalId: sql.NullString{
					// 		String: media.ProviderId,
					// 		Valid:  true,
					// 	},
					// 	Title:        media.Title,
					// 	Description:  utils.StringPtrToSqlNull(media.Description),
					// 	Score:        utils.Float64PtrToSqlNull(media.Score),
					// 	Status:       media.Status,
					// 	Rating:       media.Rating,
					// 	AiringSeason: utils.StringPtrToSqlNull(media.AiringSeason),
					// 	// StartDate:    sql.NullString{},
					// 	// EndDate:      sql.NullString{},
					// })
					// if err != nil {
					// 	return nil, err
					// }

					if media.AiringSeason != nil {
						err := app.DB().CreateTag(ctx, *media.AiringSeason, *media.AiringSeason)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}

					providerIds := kvstore.Store{}
					maps.Copy(providerIds, media.ExtraProviderIds)
					providerIds[providerName] = media.ProviderId

					startDate := ""
					if media.StartDate != nil {
						startDate = media.StartDate.Format(types.MediaDateLayout)
					}

					endDate := ""
					if media.EndDate != nil {
						endDate = media.EndDate.Format(types.MediaDateLayout)
					}

					id, err := app.DB().CreateMedia(ctx, database.CreateMediaParams{
						Type:         media.Type,
						Title:        media.Title,
						Description:  utils.StringPtrToSqlNull(media.Description),
						Score:        utils.Float64PtrToSqlNull(media.Score),
						Status:       media.Status,
						Rating:       media.Rating,
						AiringSeason: utils.StringPtrToSqlNull(media.AiringSeason),
						StartDate: sql.NullString{
							String: startDate,
							Valid:  startDate != "",
						},
						EndDate: sql.NullString{
							String: endDate,
							Valid:  endDate != "",
						},
						Providers: providerIds,
					})
					if err != nil {
						return nil, err
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

					if media.Type.IsMovie() {
						err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
							MediaId: id,
							Name:    media.Title,
							Index:   1,
						})
						if err != nil {
							return nil, err
						}
					} else {
						for _, part := range media.Parts {
							err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
								MediaId: id,
								Name:    part.Name,
								Index:   int64(part.Number),
							})
							if err != nil {
								return nil, err
							}
						}
					}

					changes := database.MediaChanges{}

					if media.CoverUrl != nil {
						p, err := utils.DownloadImage(*media.CoverUrl, mediaDir.Images(), "cover")
						if err != nil {
							logger.Error("failed to download cover image for media", "mediaId", id, "err", err)
						}

						n := path.Base(p)
						changes.CoverFile = database.Change[sql.NullString]{
							Value: sql.NullString{
								String: n,
								Valid:  n != "",
							},
							Changed: true,
						}
					}

					if media.BannerUrl != nil {
						p, err := utils.DownloadImage(*media.BannerUrl, mediaDir.Images(), "banner")
						if err != nil {
							logger.Error("failed to download banner image for media", "mediaId", id, "err", err)
						}

						n := path.Base(p)
						changes.BannerFile = database.Change[sql.NullString]{
							Value: sql.NullString{
								String: n,
								Valid:  n != "",
							},
							Changed: true,
						}
					}

					if media.LogoUrl != nil {
						p, err := utils.DownloadImage(*media.LogoUrl, mediaDir.Images(), "logo")
						if err != nil {
							logger.Error("failed to download logo image for media", "media", id, "err", err)
						}

						n := path.Base(p)
						changes.LogoFile = database.Change[sql.NullString]{
							Value: sql.NullString{
								String: n,
								Valid:  n != "",
							},
							Changed: true,
						}
					}

					err = app.DB().UpdateMedia(ctx, id, changes)
					if err != nil {
						return nil, err
					}

					for _, tag := range media.Tags {
						err := app.DB().CreateTag(ctx, tag, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}

						err = app.DB().AddTagToMedia(ctx, id, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}

					for _, tag := range media.Creators {
						err := app.DB().CreateTag(ctx, tag, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}

						err = app.DB().AddCreatorToMedia(ctx, id, tag)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				return nil, nil
			},
		},
	)
}
