package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"time"

	"maps"

	"github.com/kr/pretty"
	"github.com/maruel/natural"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/anvil"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/kvstore"
	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

type ProviderValue struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Value       string `json:"value"`
}

func createProviderValues(pm *provider.ProviderManager, providerStore kvstore.Store) []ProviderValue {
	providers := make([]ProviderValue, 0, len(providerStore))
	for name, value := range providerStore {
		info, ok := pm.GetProviderInfo(name)
		if !ok {
			continue
		}

		providers = append(providers, ProviderValue{
			Name:        info.Name,
			DisplayName: info.GetDisplayName(),
			Value:       value,
		})
	}

	sort.SliceStable(providers, func(i, j int) bool {
		return natural.Less(providers[i].Name, providers[j].Name)
	})

	return providers
}

type ProviderSupports struct {
	GetMedia    bool `json:"getMedia"`
	SearchMedia bool `json:"searchMedia"`

	GetCollection    bool `json:"getCollection"`
	SearchCollection bool `json:"searchCollection"`
}

type Provider struct {
	Name        string           `json:"name"`
	DisplayName string           `json:"displayName"`
	Supports    ProviderSupports `json:"supports"`
}

type GetProviders struct {
	Providers []Provider `json:"providers"`
}

type ProviderSearchResult struct {
	ProviderName string `json:"providerName"`
	ProviderId   string `json:"providerId"`
	Title        string `json:"title"`
	ImageUrl     string `json:"imageUrl"`
}

type GetProviderSearch struct {
	SearchResults []ProviderSearchResult `json:"searchResults"`
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

type PostProviderImportMediaBody struct {
	Ids []string `json:"ids"`
}

func (b *PostProviderImportMediaBody) Transform() {
	b.Ids = fixArr(b.Ids)
}

type PostProviderImportCollectionsBody struct {
	Ids []string `json:"ids"`
}

func (b *PostProviderImportCollectionsBody) Transform() {
	b.Ids = fixArr(b.Ids)
}

func ImportMedia(ctx context.Context, app core.App, providerName, providerId string) (string, error) {
	pm := app.ProviderManager()

	dbMedia, err := app.DB().GetMediaByProviderId(ctx, nil, providerName, providerId)
	if err == nil {
		fmt.Printf("id already exists: %v\n", providerId)
		return dbMedia.Id, nil
	}

	if !errors.Is(err, database.ErrItemNotFound) {
		return "", err
	}

	media, err := pm.GetMedia(ctx, providerName, providerId)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		// TODO(patrik): Handle err
		return "", err
	}

	id := utils.CreateMediaId()

	// TODO(patrik): Better way to do this, ensure that these directories exists
	mediaDir := app.WorkDir().MediaDirById(id)
	dirs := []string{
		mediaDir.String(),
		mediaDir.Images(),
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return "", err
		}
	}

	if media.AiringSeason != nil {
		*media.AiringSeason = utils.Slug(*media.AiringSeason)

		err := app.DB().CreateTag(ctx, *media.AiringSeason, *media.AiringSeason)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return "", err
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

	coverFilename := ""
	bannerFilename := ""
	logoFilename := ""

	if media.CoverUrl != nil {
		p, err := utils.DownloadImageHashed(*media.CoverUrl, mediaDir.Images())
		if err == nil {
			coverFilename = path.Base(p)
		} else {
			app.Logger().Error("failed to download cover image for media", "err", err)
		}
	}

	if media.BannerUrl != nil {
		p, err := utils.DownloadImageHashed(*media.BannerUrl, mediaDir.Images())
		if err == nil {
			bannerFilename = path.Base(p)
		} else {
			app.Logger().Error("failed to download banner image for media", "err", err)
		}
	}

	if media.LogoUrl != nil {
		p, err := utils.DownloadImageHashed(*media.LogoUrl, mediaDir.Images())
		if err == nil {
			logoFilename = path.Base(p)
		} else {
			app.Logger().Error("failed to download logo image for media", "err", err)
		}
	}

	_, err = app.DB().CreateMedia(ctx, database.CreateMediaParams{
		Id:           id,
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
		DefaultProvider: sql.NullString{
			String: providerName,
			Valid:  providerName != "",
		},
		Providers: providerIds,
	})
	if err != nil {
		return "", err
	}

	if media.Type.IsMovie() {
		releaseDate := ""
		if media.StartDate != nil {
			releaseDate = media.StartDate.Format(types.MediaDateLayout)
		}

		err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
			MediaId: id,
			Name:    media.Title,
			Index:   1,
			ReleaseDate: sql.NullString{
				String: releaseDate,
				Valid:  releaseDate != "",
			},
		})
		if err != nil {
			return "", err
		}
	} else {
		for _, part := range media.Parts {
			releaseDate := ""
			if part.ReleaseDate != nil {
				releaseDate = part.ReleaseDate.Format(types.MediaDateLayout)
			}

			err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
				MediaId: id,
				Name:    part.Name,
				Index:   int64(part.Number),
				ReleaseDate: sql.NullString{
					String: releaseDate,
					Valid:  releaseDate != "",
				},
			})
			if err != nil {
				return "", err
			}
		}
	}

	for _, tag := range media.Tags {
		tag = utils.Slug(tag)

		err := app.DB().CreateTag(ctx, tag, tag)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return "", err
		}

		err = app.DB().AddTagToMedia(ctx, id, tag)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return "", err
		}
	}

	for _, tag := range media.Creators {
		tag = utils.Slug(tag)

		err := app.DB().CreateTag(ctx, tag, tag)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return "", err
		}

		err = app.DB().AddCreatorToMedia(ctx, id, tag)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return "", err
		}
	}

	return id, nil
}

type ProviderMediaUpdateBody struct {
	ReplaceImages bool `json:"replaceImages,omitempty"`
	OverrideParts bool `json:"overrideParts,omitempty"`
	SetRelease    bool `json:"setRelease,omitempty"`
}

type ProviderCollectionUpdateBody struct {
	ReplaceImages bool `json:"replaceImages,omitempty"`
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
					displayName := p.DisplayName
					if displayName == "" {
						displayName = p.Name
					}

					res.Providers[i] = Provider{
						Name:        p.Name,
						DisplayName: displayName,
						Supports: ProviderSupports{
							GetMedia:         p.SupportGetMedia,
							SearchMedia:      p.SupportSearchMedia,
							GetCollection:    p.SupportGetCollection,
							SearchCollection: p.SupportSearchCollection,
						},
					}
				}

				sort.SliceStable(res.Providers, func(i, j int) bool {
					return natural.Less(res.Providers[i].Name, res.Providers[j].Name)
				})

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "ProviderSearchMedia",
			Method:       http.MethodGet,
			Path:         "/providers/:providerName/media",
			ResponseType: GetProviderSearch{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")
				url := c.Request().URL
				query := url.Query().Get("query")

				pm := app.ProviderManager()

				items, err := pm.SearchMedia(context.Background(), providerName, query)
				if err != nil {
					// TODO(patrik): Handle some of the errors the provider gives, like the provider not found
					return nil, err
				}

				res := GetProviderSearch{
					SearchResults: make([]ProviderSearchResult, len(items)),
				}

				for i, item := range items {
					res.SearchResults[i] = ProviderSearchResult{
						ProviderName: providerName,
						ProviderId:   item.ProviderId,
						Title:        item.Title,
						ImageUrl:     item.ImageUrl,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "ProviderSearchCollections",
			Method:       http.MethodGet,
			Path:         "/providers/:providerName/collections",
			ResponseType: GetProviderSearch{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")
				url := c.Request().URL
				query := url.Query().Get("query")

				pm := app.ProviderManager()

				items, err := pm.SearchCollection(context.Background(), providerName, query)
				if err != nil {
					// TODO(patrik): Handle some of the errors the provider gives, like the provider not found
					return nil, err
				}

				res := GetProviderSearch{
					SearchResults: make([]ProviderSearchResult, len(items)),
				}

				for i, item := range items {
					res.SearchResults[i] = ProviderSearchResult{
						ProviderName: providerName,
						ProviderId:   item.ProviderId,
						Title:        item.Title,
						ImageUrl:     item.ImageUrl,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ProviderImportMedia",
			Method: http.MethodPost,
			Path:   "/providers/:providerName/media/import",
			// ResponseType: PostProviderImportMedia{},
			BodyType: PostProviderImportMediaBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")

				body, err := pyrin.Body[PostProviderImportMediaBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.Background()

				for _, id := range body.Ids {
					_, err := ImportMedia(ctx, app, providerName, id)
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "ProviderImportCollections",
			Method: http.MethodPost,
			Path:   "/providers/:providerName/collections/import",
			// ResponseType: PostProviderImportMedia{},
			BodyType: PostProviderImportCollectionsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")

				body, err := pyrin.Body[PostProviderImportCollectionsBody](c)
				if err != nil {
					return nil, err
				}

				pm := app.ProviderManager()

				ctx := context.Background()

				for _, providerId := range body.Ids {
					_, err := app.DB().GetCollectionByProviderId(ctx, providerName, providerId)
					if err == nil {
						fmt.Printf("id already exists: %v\n", providerId)
						continue
					}

					if !errors.Is(err, database.ErrItemNotFound) {
						return nil, err
					}

					data, err := pm.GetCollection(ctx, providerName, providerId)
					if err != nil {
						// TODO(patrik): Handle err
						return nil, err
					}

					id := utils.CreateCollectionId()

					collectionDir := app.WorkDir().CollectionDirById(id)
					dirs := []string{
						collectionDir.String(),
						collectionDir.Images(),
					}

					for _, dir := range dirs {
						err = os.Mkdir(dir, 0755)
						if err != nil && !os.IsExist(err) {
							return nil, err
						}
					}

					providerIds := kvstore.Store{}
					maps.Copy(providerIds, data.ExtraProviderIds)
					providerIds[providerName] = data.ProviderId

					coverFilename := ""
					bannerFilename := ""
					logoFilename := ""

					if data.CoverUrl != nil {
						p, err := utils.DownloadImageHashed(*data.CoverUrl, collectionDir.Images())
						if err == nil {
							coverFilename = path.Base(p)
						} else {
							app.Logger().Error("failed to download cover image for collection", "err", err)
						}
					}

					if data.BannerUrl != nil {
						p, err := utils.DownloadImageHashed(*data.BannerUrl, collectionDir.Images())
						if err == nil {
							bannerFilename = path.Base(p)
						} else {
							app.Logger().Error("failed to download banner image for collection", "err", err)
						}
					}

					if data.LogoUrl != nil {
						p, err := utils.DownloadImageHashed(*data.LogoUrl, collectionDir.Images())
						if err == nil {
							logoFilename = path.Base(p)
						} else {
							app.Logger().Error("failed to download logo image for collection", "err", err)
						}
					}

					_, err = app.DB().CreateCollection(ctx, database.CreateCollectionParams{
						Id:   id,
						Type: data.Type,
						Name: data.Name,
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
						DefaultProvider: sql.NullString{
							String: providerName,
							Valid:  providerName != "",
						},
						Providers: providerIds,
					})
					if err != nil {
						return nil, err
					}

					for i, item := range data.Items {
						mediaId, err := ImportMedia(ctx, app, providerName, item.Id)
						if err != nil {
							return nil, err
						}

						err = app.DB().CreateCollectionMediaItem(ctx, database.CreateCollectionMediaItemParams{
							CollectionId: id,
							MediaId:      mediaId,
							Name:         item.Name,
							SearchSlug:   utils.Slug(item.Name),
							Position:     i,
						})
						if err != nil {
							return nil, err
						}
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "ProviderUpdateMedia",
			Method:   http.MethodPatch,
			Path:     "/providers/:providerName/media/:mediaId",
			BodyType: ProviderMediaUpdateBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")
				mediaId := c.Param("mediaId")

				body, err := pyrin.Body[ProviderMediaUpdateBody](c)
				if err != nil {
					return nil, err
				}

				pm := app.ProviderManager()

				ctx := context.Background()

				dbMedia, err := app.DB().GetMediaById(ctx, nil, mediaId)
				if err != nil {
					// TODO(patrik): Handle err
					return nil, err
				}

				providerId, ok := dbMedia.Providers[providerName]
				if !ok {
					return nil, errors.New("provider not found on media")
				}

				data, err := pm.GetMedia(ctx, providerName, providerId)
				if err != nil {
					fmt.Printf("err: %v\n", err)
					// TODO(patrik): Handle err
					return "", err
				}

				pretty.Println(data)

				if data.AiringSeason != nil {
					err := app.DB().CreateTag(ctx, *data.AiringSeason, *data.AiringSeason)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return "", err
					}
				}

				changes := database.MediaChanges{}

				changes.Title = database.Change[string]{
					Value:   data.Title,
					Changed: data.Title != dbMedia.Title,
				}

				desc := utils.NullToDefault(data.Description)
				changes.Description = database.Change[sql.NullString]{
					Value: sql.NullString{
						String: desc,
						Valid:  desc != "",
					},
					Changed: desc != dbMedia.Description.String,
				}

				score := utils.NullToDefault(data.Score)
				changes.Score = database.Change[sql.NullFloat64]{
					Value: sql.NullFloat64{
						Float64: score,
						Valid:   score != 0.0,
					},
					Changed: score != dbMedia.Score.Float64,
				}

				changes.Status = database.Change[types.MediaStatus]{
					Value:   data.Status,
					Changed: data.Status != dbMedia.Status,
				}

				changes.Rating = database.Change[types.MediaRating]{
					Value:   data.Rating,
					Changed: data.Rating != dbMedia.Rating,
				}

				changes.Rating = database.Change[types.MediaRating]{
					Value:   data.Rating,
					Changed: data.Rating != dbMedia.Rating,
				}

				airingSeason := utils.NullToDefault(data.AiringSeason)
				changes.AiringSeason = database.Change[sql.NullString]{
					Value: sql.NullString{
						String: airingSeason,
						Valid:  airingSeason != "",
					},
					Changed: airingSeason != dbMedia.AiringSeason.String,
				}

				formatNullTime := func(t *time.Time) string {
					if t == nil {
						return ""
					}

					return t.Format(types.MediaDateLayout)
				}

				startDate := formatNullTime(data.StartDate)
				changes.StartDate = database.Change[sql.NullString]{
					Value: sql.NullString{
						String: startDate,
						Valid:  startDate != "",
					},
					Changed: startDate != dbMedia.StartDate.String,
				}

				endDate := formatNullTime(data.EndDate)
				changes.EndDate = database.Change[sql.NullString]{
					Value: sql.NullString{
						String: endDate,
						Valid:  endDate != "",
					},
					Changed: endDate != dbMedia.EndDate.String,
				}

				if body.ReplaceImages {
					// TODO(patrik): Cover
					// TODO(patrik): Banner
					// TODO(patrik): Logo
					mediaDir := app.WorkDir().MediaDirById(dbMedia.Id)

					if data.CoverUrl != nil {
						p, err := utils.DownloadImageHashed(*data.CoverUrl, mediaDir.Images())
						if err != nil {
							return "", fmt.Errorf("failed to download cover image for media: %w", err)
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

					if data.BannerUrl != nil {
						p, err := utils.DownloadImageHashed(*data.BannerUrl, mediaDir.Images())
						if err != nil {
							return "", fmt.Errorf("failed to download banner image for media: %w", err)
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

					if data.LogoUrl != nil {
						p, err := utils.DownloadImageHashed(*data.LogoUrl, mediaDir.Images())
						if err != nil {
							return "", fmt.Errorf("failed to download logo image for media: %w", err)
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
				}

				// TODO(patrik): Extra Provider Ids?

				err = app.DB().UpdateMedia(ctx, dbMedia.Id, changes)
				if err != nil {
					return nil, err
				}

				if body.OverrideParts {
					err = app.DB().RemoveAllMediaParts(ctx, dbMedia.Id)
					if err != nil {
						return nil, err
					}

					if data.Type.IsMovie() {
						releaseDate := ""
						if data.StartDate != nil {
							releaseDate = data.StartDate.Format(types.MediaDateLayout)
						}

						err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
							MediaId: dbMedia.Id,
							Name:    data.Title,
							Index:   1,
							ReleaseDate: sql.NullString{
								String: releaseDate,
								Valid:  releaseDate != "",
							},
						})
						if err != nil {
							return nil, err
						}
					} else {
						for _, part := range data.Parts {
							releaseDate := ""
							if part.ReleaseDate != nil {
								releaseDate = part.ReleaseDate.Format(types.MediaDateLayout)
							}

							// TODO(patrik): Check for duplicated numbers
							err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
								MediaId: dbMedia.Id,
								Name:    part.Name,
								Index:   int64(part.Number),
								ReleaseDate: sql.NullString{
									String: releaseDate,
									Valid:  releaseDate != "",
								},
							})
							if err != nil {
								return nil, err
							}
						}
					}
				}

				for _, tag := range data.Tags {
					tag = utils.Slug(tag)

					err := app.DB().CreateTag(ctx, tag, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}

					err = app.DB().AddTagToMedia(ctx, dbMedia.Id, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				for _, tag := range data.Creators {
					tag = utils.Slug(tag)

					err := app.DB().CreateTag(ctx, tag, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}

					err = app.DB().AddCreatorToMedia(ctx, dbMedia.Id, tag)
					if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
						return nil, err
					}
				}

				if body.SetRelease {
					if data.Release != nil {
						err := app.DB().SetMediaPartRelease(ctx, dbMedia.Id, database.SetMediaPartRelease{
							Type:             types.MediaPartReleaseTypeNotConfirmed,
							StartDate:        data.Release.Format(time.RFC3339),
							NumExpectedParts: len(data.Parts),
							PartOffset:       0,
							IntervalDays:     7,
							DelayDays:        0,
						})
						if err != nil {
							return nil, err
						}
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "ProviderUpdateCollection",
			Method:   http.MethodPatch,
			Path:     "/providers/:providerName/collections/:collectionId",
			BodyType: ProviderCollectionUpdateBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				providerName := c.Param("providerName")
				collectionId := c.Param("collectionId")

				body, err := pyrin.Body[ProviderCollectionUpdateBody](c)
				if err != nil {
					return nil, err
				}

				pm := app.ProviderManager()

				ctx := context.Background()

				dbCollection, err := app.DB().GetCollectionById(ctx, collectionId)
				if err != nil {
					// TODO(patrik): Handle err
					return nil, err
				}

				providerId, ok := dbCollection.Providers[providerName]
				if !ok {
					// TODO(patrik): Better error
					return nil, errors.New("provider not found on media")
				}

				data, err := pm.GetCollection(ctx, providerName, providerId)
				if err != nil {
					fmt.Printf("err: %v\n", err)
					// TODO(patrik): Handle err
					return "", err
				}

				changes := database.CollectionChanges{}

				changes.Type = database.Change[types.CollectionType]{
					Value:   data.Type,
					Changed: data.Type != dbCollection.Type,
				}

				changes.Name = database.Change[string]{
					Value:   data.Name,
					Changed: data.Name != dbCollection.Name,
				}

				if body.ReplaceImages {
					collectionDir := app.WorkDir().CollectionDirById(dbCollection.Id)

					if data.CoverUrl != nil {
						p, err := utils.DownloadImageHashed(*data.CoverUrl, collectionDir.Images())
						if err != nil {
							// TODO(patrik): Better error
							return "", fmt.Errorf("failed to download cover image for collection: %w", err)
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

					if data.BannerUrl != nil {
						p, err := utils.DownloadImageHashed(*data.BannerUrl, collectionDir.Images())
						if err != nil {
							// TODO(patrik): Better error
							return "", fmt.Errorf("failed to download banner image for collection: %w", err)
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

					if data.LogoUrl != nil {
						p, err := utils.DownloadImageHashed(*data.LogoUrl, collectionDir.Images())
						if err != nil {
							// TODO(patrik): Better error
							return "", fmt.Errorf("failed to download logo image for collection: %w", err)
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
				}

				err = app.DB().UpdateCollection(ctx, dbCollection.Id, changes)
				if err != nil {
					return nil, err
				}

				items, err := app.DB().GetFullAllCollectionMediaItemsByCollection(ctx, nil, dbCollection.Id)
				if err != nil {
					return nil, err
				}

				pretty.Println(items)

				itemToMedia := map[string]string{}
				for _, item := range data.Items {
					mediaId, err := ImportMedia(ctx, app, providerName, item.Id)
					if err != nil {
						return nil, err
					}

					itemToMedia[item.Id] = mediaId
				}

				err = app.DB().RemoveAllCollectionMediaItems(ctx, dbCollection.Id)
				if err != nil {
					return nil, err
				}

				for _, item := range data.Items {
					mediaId, ok := itemToMedia[item.Id]
					if !ok {
						continue
					}

					err := app.DB().CreateCollectionMediaItem(ctx, database.CreateCollectionMediaItemParams{
						CollectionId: dbCollection.Id,
						MediaId:      mediaId,
						Name:         item.Name,
						SearchSlug:   utils.Slug(item.Name),
						Position:     item.Position,
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
