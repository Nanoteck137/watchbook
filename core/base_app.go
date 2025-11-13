package core

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/nanoteck137/pyrin/ember"
	"github.com/nanoteck137/pyrin/trail"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/config"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/job"
	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/provider/dummy"
	"github.com/nanoteck137/watchbook/provider/myanimelist"
	"github.com/nanoteck137/watchbook/provider/sonarr"
	"github.com/nanoteck137/watchbook/provider/tmdb"
	"github.com/nanoteck137/watchbook/tools/cache"
	"github.com/nanoteck137/watchbook/types"
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

func (app *BaseApp) JobProcessor() *job.JobProcessor {
	return app.jobProcessor
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

func (app *BaseApp) Bootstrap() error {
	var err error

	workDir := app.config.WorkDir()

	dirs := []string{
		workDir.MediaDir(),
		workDir.CollectionsDir(),
		workDir.ShowsDir(),
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

	c, err := cache.NewProvider(app.cacheDb)
	if err != nil {
		return err
	}

	pm := provider.NewProviderManager(c)
	pm.RegisterProvider(&myanimelist.MyAnimeListAnimeProvider{})
	pm.RegisterProvider(&dummy.DummyProvider{})
	pm.RegisterProvider(&tmdb.TmdbMovieProvider{})
	pm.RegisterProvider(&tmdb.TmdbTvProvider{})

	app.providerManager = pm
	app.jobProcessor = job.NewJobProcessor(app.db)

	ac := sonarr.NewApiClient(&cache.NoOpCache{}, app.config.SonarrUrl, app.config.SonarrApiKey)
	serie, err := ac.GetSerieById(context.TODO(), "2")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	pretty.Println(serie)

	series, err := ac.GetSeries(context.TODO())
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	pretty.Println(series)

	episodes, err := ac.GetSeasonEpisodes(context.TODO(), "2", 1)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	pretty.Println(episodes)

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
