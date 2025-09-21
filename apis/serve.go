package apis

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/kvstore"
	"github.com/nanoteck137/watchbook/provider/myanimelist"
	"github.com/nanoteck137/watchbook/types"
)

func RegisterHandlers(app core.App, router pyrin.Router) {
	g := router.Group("/api/v1")
	InstallAuthHandlers(app, g)
	InstallSystemHandlers(app, g)
	InstallUserHandlers(app, g)

	InstallMediaHandlers(app, g)
	InstallCollectionHandlers(app, g)
	InstallProviderHandlers(app, g)

	g = router.Group("/files")
	g.Register(
		pyrin.NormalHandler{
			Name:        "GetMediaImage",
			Method:      http.MethodGet,
			Path:        "/media/:id/images/:file",
			HandlerFunc: func(c pyrin.Context) error {
				id := c.Param("id")
				file := c.Param("file")

				dir := app.WorkDir().MediaDirById(id)
				p := dir.Images()
				f := os.DirFS(p)

				return pyrin.ServeFile(c, f, file)
			},
		},

		pyrin.NormalHandler{
			Name:        "GetCollectionImage",
			Method:      http.MethodGet,
			Path:        "/collections/:id/images/:file",
			HandlerFunc: func(c pyrin.Context) error {
				id := c.Param("id")
				file := c.Param("file")

				dir := app.WorkDir().CollectionDirById(id)
				p := dir.Images()
				f := os.DirFS(p)

				return pyrin.ServeFile(c, f, file)
			},
		},
	)
}

func Server(app core.App) (*pyrin.Server, error) {
	s := pyrin.NewServer(&pyrin.ServerConfig{
		LogName: watchbook.AppName,
		RegisterHandlers: func(router pyrin.Router) {
			RegisterHandlers(app, router)
		},
	})

	app.JobProcessor().RegisterHandler("import-mal-watchlist", func(ctx context.Context, job database.Job) error {
		store, err := kvstore.Deserialize(job.Payload)
		if err != nil {
			return err
		}

		username := store["username"]
		userId := store["userId"]

		if !app.ProviderManager().IsValidProvider(myanimelist.AnimeProviderName) {
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
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return s, nil
}
