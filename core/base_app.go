package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
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

var dl = downloader.NewDownloader(
	rate.NewLimiter(rate.Every(4*time.Second), 10),
	mal.UserAgent,
)

func fetchAndUpdateAnime(ctx context.Context, db *database.Database, workDir types.WorkDir, animeId string) error {
	anime, err := db.GetAnimeById(ctx, animeId)
	if err != nil {
		return err
	}

	fmt.Printf("Updating %s (%s) - %s\n", anime.Title, anime.MalId, anime.Id)

	animeData, err := mal.FetchAnimeData(dl, anime.MalId)
	if err != nil {
		return err
	}

	coverFilename := ""
	if animeData.CoverImageUrl != "" && !anime.CoverFilename.Valid {
		// dl.DownloadToFile()
		resp, err := http.Get(animeData.CoverImageUrl)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		contentType := resp.Header.Get("Content-Type")
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return err
		}

		ext := ""
		switch mediaType {
		case "image/png":
			ext = ".png"
		case "image/jpeg":
			ext = ".jpeg"
		default:
			return fmt.Errorf("Unsupported media type for cover: %s", mediaType)
		}

		dst := path.Join(workDir.ImagesEntriesDir(), animeId)
		err = os.Mkdir(dst, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}

		name := "cover"+ext
		p := path.Join(dst, name)
		f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return err
		}

		coverFilename = name
	}

	// TODO(patrik): Add some sanitization
	for _, theme := range animeData.Themes {
		err := db.CreateTheme(ctx, utils.Slug(theme), theme)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	// TODO(patrik): Add some sanitization
	for _, genre := range animeData.Genres {
		err := db.CreateGenre(ctx, utils.Slug(genre), genre)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	// TODO(patrik): Add some sanitization
	for _, studio := range animeData.Studios {
		err := db.CreateStudio(ctx, utils.Slug(studio), studio)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	// TODO(patrik): Add some sanitization
	for _, producer := range animeData.Producers {
		err := db.CreateProducer(ctx, utils.Slug(producer), producer)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	animeType := mal.ConvertAnimeType(animeData.Type)
	animeStatus := mal.ConvertAnimeStatus(animeData.Status)
	animeRating := mal.ConvertAnimeRating(animeData.Rating)

	err = db.UpdateAnime(ctx, animeId, database.AnimeChanges{
		Title: database.Change[string]{
			Value:   animeData.Title,
			Changed: animeData.Title != anime.Title,
		},

		TitleEnglish: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: animeData.TitleEnglish,
				Valid:  animeData.TitleEnglish != "",
			},
			Changed: animeData.TitleEnglish != anime.TitleEnglish.String,
		},

		Description: database.Change[string]{
			Value:   animeData.Description,
			Changed: animeData.Description != anime.Description,
		},

		Type: database.Change[types.AnimeType]{
			Value:   animeType,
			Changed: animeType != anime.Type,
		},

		Status: database.Change[types.AnimeStatus]{
			Value:   animeStatus,
			Changed: animeStatus != anime.Status,
		},

		Rating: database.Change[types.AnimeRating]{
			Value:   animeRating,
			Changed: animeRating != anime.Rating,
		},

		AiringSeason: database.Change[string]{
			Value:   animeData.Premiered,
			Changed: animeData.Premiered != anime.AiringSeason,
		},

		EpisodeCount: database.Change[sql.NullInt64]{
			Value: sql.NullInt64{
				Int64: utils.NullToDefault(animeData.EpisodeCount),
				Valid: animeData.EpisodeCount != nil,
			},
			Changed: true,
		},

		StartDate: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: utils.NullToDefault(animeData.StartDate),
				Valid:  animeData.StartDate != nil,
			},
			Changed: true,
		},

		EndDate: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: utils.NullToDefault(animeData.EndDate),
				Valid:  animeData.EndDate != nil,
			},
			Changed: true,
		},

		Score: database.Change[sql.NullFloat64]{
			Value: sql.NullFloat64{
				Float64: utils.NullToDefault(animeData.Score),
				Valid:   animeData.Score != nil,
			},
			Changed: true,
		},

		AniDBUrl: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: animeData.AniDBUrl,
				Valid:  animeData.AniDBUrl != "",
			},
			Changed: animeData.AniDBUrl != anime.AniDBUrl.String,
		},

		AnimeNewsNetworkUrl: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: animeData.AnimeNewsNetworkUrl,
				Valid:  animeData.AnimeNewsNetworkUrl != "",
			},
			Changed: animeData.AnimeNewsNetworkUrl != anime.AnimeNewsNetworkUrl.String,
		},

		CoverFilename: database.Change[sql.NullString]{
			Value:   sql.NullString{
				String: coverFilename,
				Valid:  coverFilename != "",
			},
			Changed: coverFilename != anime.CoverFilename.String,
		},

		ShouldFetchData: database.Change[bool]{
			Value:   false,
			Changed: true,
		},

		LastDataFetchDate: database.Change[time.Time]{
			Value:   time.Now(),
			Changed: true,
		},
	})
	if err != nil {
		return err
	}

	for i, themeSong := range animeData.ThemeSongs {
		err := db.CreateAnimeThemeSong(ctx, database.CreateAnimeThemeSongParams{
			AnimeId: animeId,
			Idx:     i,
			Raw:     themeSong.Raw,
			Type:    mal.ConvertThemeSongType(themeSong.Type),
		})
		if err != nil {
			return err
		}
	}

	for _, theme := range animeData.Themes {
		err := db.AddThemeToAnime(ctx, animeId, utils.Slug(theme))
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	for _, genre := range animeData.Genres {
		err := db.AddGenreToAnime(ctx, animeId, utils.Slug(genre))
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	for _, studio := range animeData.Studios {
		err := db.AddStudioToAnime(ctx, animeId, utils.Slug(studio))
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	for _, producer := range animeData.Producers {
		err := db.AddProducerToAnime(ctx, animeId, utils.Slug(producer))
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

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

	dirs := []string{
		workDir.ImagesDir(),
		workDir.ImagesEntriesDir(),
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

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

	entries, err := mal.GetUserWatchlist(dl, "Nanoteck137")
	if err != nil {
		return err
	}

	ctx := context.TODO()

	for _, entry := range entries {
		malId := strconv.Itoa(entry.AnimeId)

		_, err := db.GetAnimeByMalId(ctx, malId)
		if err != nil && errors.Is(err, database.ErrItemNotFound) {
			_, err := db.CreateAnime(ctx, database.CreateAnimeParams{
				MalId:           malId,
				Title:           string(entry.AnimeTitle),
				ShouldFetchData: true,
			})
			if err != nil {
				return err
			}
		}
	}

	// ids, err := db.GetAnimeIdsForFetching(ctx)
	// if err != nil {
	// 	return err
	// }

	// for _, id := range ids[:4] {
	// 	err := fetchAndUpdateAnime(ctx, db, workDir, id)
	// 	if err != nil {
	// 		fmt.Printf("err: %v\n", err)
	// 		continue
	// 	}
	// }

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
