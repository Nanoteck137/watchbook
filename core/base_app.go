package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"os"
	"path"
	"strconv"

	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/pyrin/trail"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/config"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/job"
	"github.com/nanoteck137/watchbook/kvstore"
	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/provider/dummy"
	"github.com/nanoteck137/watchbook/provider/myanimelist"
	"github.com/nanoteck137/watchbook/provider/tmdb"
	"github.com/nanoteck137/watchbook/tools/cache"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	logger          *trail.Logger
	db              *database.Database
	cacheDb         *ember.Database
	providerManager *provider.ProviderManager
	config          *config.Config
	jobProcessor    *job.JobProcessor
}

func (app *BaseApp) Logger() *trail.Logger {
	return app.logger
}

func (app *BaseApp) DB() *database.Database {
	return app.db
}

func (app *BaseApp) ProviderManager() *provider.ProviderManager {
	return app.providerManager
}

func (app *BaseApp) Config() *config.Config {
	return app.config
}

func (app *BaseApp) WorkDir() types.WorkDir {
	return app.config.WorkDir()
}

func ImportMedia(ctx context.Context, app App, providerName, providerId string) (string, error) {
	pm := app.ProviderManager()

	dbMedia, err := app.DB().GetMediaByProviderId(ctx, nil, providerName, providerId)
	if err == nil {
		app.Logger().Info("media already exists", "providerName", providerName, "providerId", providerId)
		return dbMedia.Id, nil
	}

	if !errors.Is(err, database.ErrItemNotFound) {
		return "", err
	}

	app.Logger().Info("media not found, creating", "providerName", providerName, "providerId", providerId)

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
		err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
			MediaId: id,
			Name:    media.Title,
			Index:   1,
		})
		if err != nil {
			return "", err
		}
	} else {
		for _, part := range media.Parts {
			err := app.DB().CreateMediaPart(ctx, database.CreateMediaPartParams{
				MediaId: id,
				Name:    part.Name,
				Index:   int64(part.Number),
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

func (app *BaseApp) Bootstrap() error {
	var err error

	workDir := app.config.WorkDir()

	dirs := []string{
		workDir.MediaDir(),
		workDir.CollectionsDir(),
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	app.db, err = database.Open(workDir.DatabaseFile())
	if err != nil {
		return err
	}

	if app.config.RunMigrations {
		err = app.db.RunMigrateUp()
		if err != nil {
			return err
		}
	}

	app.cacheDb, err = cache.OpenDatabase(workDir.CacheDatabaseFile())
	if err != nil {
		return err
	}

	cache, err := cache.NewProvider(app.cacheDb)
	if err != nil {
		return err
	}

	pm := provider.NewProviderManager(cache)
	pm.RegisterProvider(&myanimelist.MyAnimeListAnimeProvider{})
	pm.RegisterProvider(&dummy.DummyProvider{})
	pm.RegisterProvider(&tmdb.TmdbMovieProvider{})
	pm.RegisterProvider(&tmdb.TmdbTvProvider{})

	app.providerManager = pm

	app.jobProcessor = job.NewJobProcessor(app.db)
	app.jobProcessor.RegisterHandler("test", func(ctx context.Context, job database.Job) error {
		store, err := kvstore.Deserialize(job.Payload)
		if err != nil {
			return err
		}

		username := store["username"]
		userId := store["userId"]

		if !pm.IsValidProvider(myanimelist.AnimeProviderName) {
			// TODO(patrik): Better error
			return errors.New("unsupported operation")
		}

		entries, err := myanimelist.GetUserWatchlist(username)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			id := strconv.Itoa(entry.AnimeId)
			mediaId, err := ImportMedia(ctx, app, myanimelist.AnimeProviderName, id)
			if err != nil {
				return err
			}

			list := types.MediaUserListBacklog
			switch entry.Status {
			case myanimelist.WatchlistStatusCurrentlyWatching:
				list = types.MediaUserListInProgress
			case myanimelist.WatchlistStatusCompleted:
				list = types.MediaUserListCompleted
			case myanimelist.WatchlistStatusOnHold:
				list = types.MediaUserListOnHold
			case myanimelist.WatchlistStatusDropped:
				list = types.MediaUserListDropped
			case myanimelist.WatchlistStatusPlanToWatch:
				list = types.MediaUserListBacklog
			default:
				app.Logger().Error("unknown status", "status", entry.Status)
			}

			err = app.DB().SetMediaUserData(ctx, mediaId, userId, database.SetMediaUserData{
				List: list,
				Part: sql.NullInt64{
					Int64: int64(entry.NumWatchedEpisodes),
					Valid: entry.NumWatchedEpisodes != 0,
				},
				RevisitCount: sql.NullInt64{},
				IsRevisiting: false,
				Score: sql.NullInt64{
					Int64: int64(entry.Score),
					Valid: entry.Score != 0,
				},
				Created: 0,
				Updated: 0,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	_, err = os.Stat(workDir.SetupFile())
	if errors.Is(err, os.ErrNotExist) && app.config.Username != "" {
		app.logger.Info("server not setup, creating the initial user")

		ctx := context.Background()

		_, err := app.db.CreateUser(ctx, database.CreateUserParams{
			Username: app.config.Username,
			Password: app.config.InitialPassword,
			Role:     types.RoleSuperUser,
		})
		if err != nil {
			return err
		}

		f, err := os.Create(workDir.SetupFile())
		if err != nil {
			return err
		}
		f.Close()
	}

	app.jobProcessor.Start(1)

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		logger: watchbook.DefaultLogger(),
		config: config,
	}
}
