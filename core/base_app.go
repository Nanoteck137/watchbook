package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nanoteck137/watchbook/config"
	"github.com/nanoteck137/watchbook/core/log"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/downloader"
	"github.com/nanoteck137/watchbook/mal"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
	"golang.org/x/time/rate"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	db     *database.Database
	config *config.Config
}

func (app *BaseApp) DB() *database.Database {
	return app.db
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

	// dirs := []string{
	// }
	//
	// for _, dir := range dirs {
	// 	err = os.Mkdir(dir, 0755)
	// 	if err != nil && !os.IsExist(err) {
	// 		return err
	// 	}
	// }

	app.db, err = database.Open(workDir)
	if err != nil {
		return err
	}

	if app.config.RunMigrations {
		err = database.RunMigrateUp(app.db)
		if err != nil {
			return err
		}
	}

	db := app.db

	dl := downloader.NewDownloader(
		rate.NewLimiter(rate.Every(4*time.Second), 10),
		mal.UserAgent,
	)

	entries, err := mal.GetUserWatchlist(dl, "Nanoteck137")
	if err != nil {
		return err
	}

	ctx := context.TODO()

	for i, entry := range entries {
		malId := strconv.Itoa(entry.AnimeId)

		id, err := db.CreateAnime(ctx, database.CreateAnimeParams{
			MalId: malId,
			Title: string(entry.AnimeTitle),
			// TitleEnglish: sql.NullString{
			// 	String: anime.TitleEnglish,
			// 	Valid:  anime.TitleEnglish != "",
			// },
			// Description:  anime.Description,
			// Type:         types.AnimeTypeUnknown,
			// Status:       types.AnimeStatusUnknown,
			// Rating:       types.AnimeRatingUnknown,
			// AiringSeason: anime.Premiered,
			// EpisodeCount: utils.Int64PtrToSqlNull(anime.EpisodeCount),
			// StartDate:    utils.StringPtrToSqlNull(anime.StartDate),
			// EndDate:      utils.StringPtrToSqlNull(anime.EndDate),
			// Score:        utils.Float64PtrToSqlNull(anime.Score),
			// DownloadDate: time.Now(),
			// AniDBUrl: sql.NullString{
			// 	String: anime.AniDBUrl,
			// 	Valid:  anime.AniDBUrl != "",
			// },
			// AnimeNewsNetworkUrl: sql.NullString{
			// 	String: anime.AnimeNewsNetworkUrl,
			// 	Valid:  anime.AnimeNewsNetworkUrl != "",
			// },
		})
		if err != nil {
			return err
		}

		continue

		fmt.Printf("Processing entry %d/%d - ", i+1, len(entries))
		fmt.Printf("Downloading %d\n", entry.AnimeId)

		anime, err := mal.FetchAnimeData(dl, workDir, malId)
		if err != nil {
			return err
		}

		// TODO(patrik): Add some sanitization
		for _, theme := range anime.Themes {
			err := db.CreateTheme(ctx, utils.Slug(theme), theme)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}

		// TODO(patrik): Add some sanitization
		for _, genre := range anime.Genres {
			err := db.CreateGenre(ctx, utils.Slug(genre), genre)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}

		// TODO(patrik): Add some sanitization
		for _, studio := range anime.Studios {
			err := db.CreateStudio(ctx, utils.Slug(studio), studio)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}

		// TODO(patrik): Add some sanitization
		for _, producer := range anime.Producers {
			err := db.CreateProducer(ctx, utils.Slug(producer), producer)
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}

		animeType := mal.ConvertAnimeType(anime.Type)
		animeStatus := mal.ConvertAnimeStatus(anime.Status)
		animeRating := mal.ConvertAnimeRating(anime.Rating)

		id, err = db.CreateAnime(ctx, database.CreateAnimeParams{
			MalId: malId,
			Title: anime.Title,
			TitleEnglish: sql.NullString{
				String: anime.TitleEnglish,
				Valid:  anime.TitleEnglish != "",
			},
			Description:  anime.Description,
			Type:         animeType,
			Status:       animeStatus,
			Rating:       animeRating,
			AiringSeason: anime.Premiered,
			EpisodeCount: utils.Int64PtrToSqlNull(anime.EpisodeCount),
			StartDate:    utils.StringPtrToSqlNull(anime.StartDate),
			EndDate:      utils.StringPtrToSqlNull(anime.EndDate),
			Score:        utils.Float64PtrToSqlNull(anime.Score),
			DownloadDate: time.Now(),
			AniDBUrl: sql.NullString{
				String: anime.AniDBUrl,
				Valid:  anime.AniDBUrl != "",
			},
			AnimeNewsNetworkUrl: sql.NullString{
				String: anime.AnimeNewsNetworkUrl,
				Valid:  anime.AnimeNewsNetworkUrl != "",
			},
		})
		if err != nil {
			return err
		}

		for i, themeSong := range anime.ThemeSongs {
			err := db.CreateAnimeThemeSong(ctx, database.CreateAnimeThemeSongParams{
				AnimeId: id,
				Idx:     i,
				Raw:     themeSong.Raw,
				Type:    mal.ConvertThemeSongType(themeSong.Type),
			})
			if err != nil {
				return err
			}
		}

		for _, theme := range anime.Themes {
			err := db.AddThemeToAnime(ctx, id, utils.Slug(theme))
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}

		for _, genre := range anime.Genres {
			err := db.AddGenreToAnime(ctx, id, utils.Slug(genre))
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}

		for _, studio := range anime.Studios {
			err := db.AddStudioToAnime(ctx, id, utils.Slug(studio))
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}

		for _, producer := range anime.Producers {
			err := db.AddProducerToAnime(ctx, id, utils.Slug(producer))
			if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
				return err
			}
		}
	}

	return errors.New("work in progress")

	_, err = os.Stat(workDir.SetupFile())
	if errors.Is(err, os.ErrNotExist) && app.config.Username != "" {
		log.Info("Server not setup, creating the initial user")

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

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		config: config,
	}
}
