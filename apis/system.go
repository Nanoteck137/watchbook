package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/core"
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
	anime, err := db.GetAnimeById(ctx, nil, animeId)
	if err != nil {
		return err
	}

	if !anime.MalId.Valid {
		return nil
	}

	malId := anime.MalId.String

	fmt.Printf("Updating %s (%s) - %s\n", anime.Title, malId, anime.Id)

	animeData, err := mal.FetchAnimeData(dl, malId)
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

		name := "cover" + ext
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
		err := db.CreateTag(ctx, utils.Slug(theme), theme)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	// TODO(patrik): Add some sanitization
	for _, genre := range animeData.Genres {
		err := db.CreateTag(ctx, utils.Slug(genre), genre)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	// TODO(patrik): Add some sanitization
	for _, demographic := range animeData.Demographics {
		err := db.CreateTag(ctx, utils.Slug(demographic), demographic)
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

		Description: database.Change[sql.NullString]{
			Value: sql.NullString{
				String: animeData.Description,
				Valid:  animeData.Description != "",
			},
			Changed: animeData.Description != anime.Description.String,
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

		// AniDBUrl: database.Change[sql.NullString]{
		// 	Value: sql.NullString{
		// 		String: animeData.AniDBUrl,
		// 		Valid:  animeData.AniDBUrl != "",
		// 	},
		// 	Changed: animeData.AniDBUrl != anime.AniDBUrl.String,
		// },
		//
		// AnimeNewsNetworkUrl: database.Change[sql.NullString]{
		// 	Value: sql.NullString{
		// 		String: animeData.AnimeNewsNetworkUrl,
		// 		Valid:  animeData.AnimeNewsNetworkUrl != "",
		// 	},
		// 	Changed: animeData.AnimeNewsNetworkUrl != anime.AnimeNewsNetworkUrl.String,
		// },

		CoverFilename: database.Change[sql.NullString]{
			Value: sql.NullString{
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
		err := db.AddTagToAnime(ctx, animeId, utils.Slug(theme))
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	for _, genre := range animeData.Genres {
		err := db.AddTagToAnime(ctx, animeId, utils.Slug(genre))
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	for _, demographic := range animeData.Demographics {
		err := db.AddTagToAnime(ctx, animeId, utils.Slug(demographic))
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

	return nil
}

type Event struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type EventData interface {
	GetEventType() string
}

// NOTE(patrik): Based on: https://gist.github.com/Ananto30/8af841f250e89c07e122e2a838698246
type Broker struct {
	Notifier chan EventData

	newClients     chan chan EventData
	closingClients chan chan EventData
	clients        map[chan EventData]bool
}

func NewBrokerServer() (broker *Broker) {
	// Instantiate a broker
	broker = &Broker{
		Notifier:       make(chan EventData, 1),
		newClients:     make(chan chan EventData),
		closingClients: make(chan chan EventData),
		clients:        make(map[chan EventData]bool),
	}

	// Set it running - listening and broadcasting events
	go broker.listen()

	return
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			broker.clients[s] = true
			log.Debug("Client added", "numClients", len(broker.clients))
		case s := <-broker.closingClients:
			delete(broker.clients, s)
			log.Debug("Removed client", "numClients", len(broker.clients))
		case event := <-broker.Notifier:
			for clientMessageChan := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

func (broker *Broker) EmitEvent(event EventData) {
	broker.Notifier <- event
}

type GetSystemInfo struct {
	Version string `json:"version"`
}

type DownloadHandlerStatus struct {
	IsDownloading   bool   `json:"isDownloading"`
	CurrentDownload int    `json:"currentDownload"`
	TotalDownloads  int    `json:"totalDownloads"`
	LastError       string `json:"lastError"`
}

type DownloadHandler struct {
	isDownloading atomic.Bool

	mutex           sync.Mutex
	currentDownload int
	totalDownloads  int

	lastError error

	broker *Broker
}

func NewDownloadHandler() *DownloadHandler {
	return &DownloadHandler{
		isDownloading:   atomic.Bool{},
		mutex:           sync.Mutex{},
		currentDownload: 0,
		totalDownloads:  0,
		lastError:       nil,
		broker:          NewBrokerServer(),
	}
}

func (d *DownloadHandler) getStatus() DownloadHandlerStatus {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	var lastError string

	err := d.lastError
	if err != nil {
		lastError = err.Error()
	}

	return DownloadHandlerStatus{
		IsDownloading:   d.isDownloading.Load(),
		CurrentDownload: d.currentDownload,
		TotalDownloads:  d.totalDownloads,
		LastError:       lastError,
	}
}

func (d *DownloadHandler) sendStatusEvent() {
	d.broker.EmitEvent(StatusEvent{
		DownloadHandlerStatus: d.getStatus(),
	})
}

func (d *DownloadHandler) updateStatus(current, total int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.currentDownload = current
	d.totalDownloads = total
}

func (d *DownloadHandler) download(app core.App) error {
	d.isDownloading.Store(true)
	defer d.isDownloading.Store(false)

	// entries, err := mal.GetUserWatchlist(dl, "Nanoteck137")
	// if err != nil {
	// 	return err
	// }
	//
	// ctx := context.TODO()
	//
	// for _, entry := range entries {
	// 	malId := strconv.Itoa(entry.AnimeId)
	//
	// 	_, err := db.GetAnimeByMalId(ctx, malId)
	// 	if err != nil && errors.Is(err, database.ErrItemNotFound) {
	// 		_, err := db.CreateAnime(ctx, database.CreateAnimeParams{
	// 			MalId:           malId,
	// 			Title:           string(entry.AnimeTitle),
	// 			ShouldFetchData: true,
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	ctx := context.TODO()

	db := app.DB()
	workDir := app.WorkDir()

	ids, err := db.GetAnimeIdsForFetching(ctx)
	if err != nil {
		return err
	}

	d.updateStatus(0, len(ids))
	d.sendStatusEvent()

	// TODO(patrik): This is temporary
	for i, id := range ids[:4] {
	// for i, id := range ids {
		d.updateStatus(i + 1, len(ids))
		d.sendStatusEvent()

		err := fetchAndUpdateAnime(ctx, db, workDir, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DownloadHandler) startDownload(app core.App) error {
	if d.isDownloading.Load() {
		return errors.New("already downloading")
	}

	err := d.download(app)
	d.lastError = err

	d.sendStatusEvent()

	return nil
}

var _ EventData = (*StatusEvent)(nil)

type StatusEvent struct {
	DownloadHandlerStatus
}

func (e StatusEvent) GetEventType() string {
	return "status"
}

var downloadHandler = NewDownloadHandler()

func InstallSystemHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetSystemInfo",
			Path:         "/system/info",
			Method:       http.MethodGet,
			ResponseType: GetSystemInfo{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return GetSystemInfo{
					Version: watchbook.Version,
				}, nil
			},
		},

		// user, err := app.DB().GetUserByUsername(ctx, app.config.Username)
		// if err != nil {
		// 	return err
		// }
		//
		// for _, entry := range entries {
		// 	malId := strconv.Itoa(entry.AnimeId)
		//
		// 	anime, err := db.GetAnimeByMalId(ctx, malId)
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	noList := false
		// 	list := types.AnimeUserListWatching
		// 	switch entry.Status {
		// 	case mal.WatchlistStatusCurrentlyWatching:
		// 		list = types.AnimeUserListWatching
		// 	case mal.WatchlistStatusCompleted:
		// 		list = types.AnimeUserListCompleted
		// 	case mal.WatchlistStatusOnHold:
		// 		list = types.AnimeUserListOnHold
		// 	case mal.WatchlistStatusDropped:
		// 		list = types.AnimeUserListDropped
		// 	case mal.WatchlistStatusPlanToWatch:
		// 		list = types.AnimeUserListPlanToWatch
		// 	default:
		// 		noList = true
		// 	}
		//
		// 	err = db.SetAnimeUserData(ctx, anime.Id, user.Id, database.SetAnimeUserData{
		// 		List: sql.NullString{
		// 			String: string(list),
		// 			Valid:  !noList,
		// 		},
		// 		Episode: sql.NullInt64{
		// 			Int64: int64(entry.NumWatchedEpisodes),
		// 			Valid: entry.NumWatchedEpisodes != 0,
		// 		},
		// 		IsRewatching: entry.IsRewatching > 0,
		// 		Score: sql.NullInt64{
		// 			Int64: int64(entry.Score),
		// 			Valid: entry.Score != 0,
		// 		},
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	// err = db.SetAnimeUserList(ctx, anime.Id, user.Id, list)
		// 	// if err != nil {
		// 	// 	return err
		// 	// }
		// 	//
		// 	// err = db.SetAnimeUserEpisode(ctx, anime.Id, user.Id, entry.NumWatchedEpisodes)
		// 	// if err != nil {
		// 	// 	return err
		// 	// }
		// 	//
		// 	// if entry.Score > 0 {
		// 	// 	err = db.SetAnimeUserScore(ctx, anime.Id, user.Id, entry.Score)
		// 	// 	if err != nil {
		// 	// 		return err
		// 	// 	}
		// 	// } else {
		// 	// 	err = db.RemoveAnimeUserScore(ctx, anime.Id, user.Id)
		// 	// 	if err != nil {
		// 	// 		return err
		// 	// 	}
		// 	// }
		// }

		pyrin.ApiHandler{
			Name:   "StartDownload",
			Path:   "/system/download",
			Method: http.MethodGet,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				go func() {
					err := downloadHandler.startDownload(app)
					if err != nil {
						log.Error("failed to start downloader", "err", err)
					}
				}()

				return nil, nil
			},
		},

		pyrin.NormalHandler{
			Name:   "SseHandler",
			Method: http.MethodGet,
			Path:   "/system/sse",
			HandlerFunc: func(c pyrin.Context) error {
				r := c.Request()
				w := c.Response()

				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				w.Header().Set("Access-Control-Allow-Origin", "*")

				rc := http.NewResponseController(w)

				eventChan := make(chan EventData)
				downloadHandler.broker.newClients <- eventChan

				defer func() {
					downloadHandler.broker.closingClients <- eventChan
				}()

				sendEvent := func(eventData EventData) {
					fmt.Fprintf(w, "data: ")

					event := Event{
						Type: eventData.GetEventType(),
						Data: eventData,
					}

					encode := json.NewEncoder(w)
					encode.Encode(event)

					fmt.Fprintf(w, "\n\n")
					rc.Flush()
				}

				sendEvent(StatusEvent{
					DownloadHandlerStatus: downloadHandler.getStatus(),
				})

				for {
					select {
					case <-r.Context().Done():
						return nil
					case event := <-eventChan:
						sendEvent(event)
					}
				}
			},
		},
	)
}
