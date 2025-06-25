package apis

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/downloader"
	"github.com/nanoteck137/watchbook/mal"
	"github.com/nanoteck137/watchbook/types"
	"golang.org/x/time/rate"
)

var dl = downloader.NewDownloader(
	rate.NewLimiter(rate.Every(4*time.Second), 10),
	mal.UserAgent,
)

func downloadImage(ctx context.Context, db *database.Database, workDir types.WorkDir, mediaId, url string, typ types.MediaImageType, isPrimary bool) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("downloadImage: failed http get request: %w", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("downloadImage: failed to parse Content-Type: %w", err)
	}

	ext := ""
	switch mediaType {
	case "image/png":
		ext = ".png"
	case "image/jpeg":
		ext = ".jpeg"
	default:
		return "", fmt.Errorf("downloadImage: unsupported media type (%s): %w", mediaType, err)
	}

	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return "", fmt.Errorf("downloadImage: failed io.Copy: %w", err)
	}

	mediaDir := workDir.MediaDir()
	dst := mediaDir.MediaImageDir(mediaId)
	err = os.Mkdir(dst, 0755)
	if err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("downloadImage: failed os.Mkdir: %w", err)
	}

	rawHash := sha256.Sum256(buf.Bytes())
	hash := hex.EncodeToString(rawHash[:])

	name := hash + ext

	err = db.CreateMediaImage(ctx, database.CreateMediaImageParams{
		MediaId:   mediaId,
		Hash:      hash,
		Type:      typ,
		MimeType:  mediaType,
		Filename:  name,
		IsPrimary: isPrimary,
	})
	if err != nil {
		if errors.Is(err, database.ErrItemAlreadyExists) {
			return hash, nil
		}

		return "", fmt.Errorf("downloadImage: failed to create media image: %w", err)
	}

	err = os.WriteFile(path.Join(dst, name), buf.Bytes(), 0644)
	if err != nil {
		return "", fmt.Errorf("downloadImage: failed to write image to disk: %w", err)
	}

	return hash, nil
}

// func fetchAndUpdateMedia(ctx context.Context, db *database.Database, workDir types.WorkDir, animeId string) error {
// 	anime, err := db.GetMediaById(ctx, nil, animeId)
// 	if err != nil {
// 		return err
// 	}
//
// 	if !anime.MalId.Valid {
// 		return nil
// 	}
//
// 	malId := anime.MalId.String
//
// 	fmt.Printf("Updating %s (%s) - %s\n", anime.Title, malId, anime.Id)
//
// 	animeData, err := mal.FetchMediaData(dl, malId, true)
// 	if err != nil {
// 		return err
// 	}
//
// 	hasCover := false
// 	for _, image := range anime.Images.Data {
// 		if image.IsCover > 0 {
// 			hasCover = true
// 			break
// 		}
// 	}
//
// 	for _, url := range animeData.Pictures {
// 		err := downloadImage(ctx, db, workDir, anime.Id, url, false)
// 		if err != nil {
// 			logger.Error("failed to download image", "err", err, "animeId", anime.Id)
// 			continue
// 		}
// 	}
//
// 	if !hasCover {
// 		err := downloadImage(ctx, db, workDir, anime.Id, animeData.CoverImageUrl, true)
// 		if err != nil {
// 			logger.Error("failed to download image", "err", err, "animeId", anime.Id)
// 		}
// 	}
//
// 	// TODO(patrik): Add some sanitization
// 	for _, theme := range animeData.Themes {
// 		err := db.CreateTag(ctx, utils.Slug(theme), theme)
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	// TODO(patrik): Add some sanitization
// 	for _, genre := range animeData.Genres {
// 		err := db.CreateTag(ctx, utils.Slug(genre), genre)
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	// TODO(patrik): Add some sanitization
// 	for _, demographic := range animeData.Demographics {
// 		err := db.CreateTag(ctx, utils.Slug(demographic), demographic)
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	// TODO(patrik): Add some sanitization
// 	for _, studio := range animeData.Studios {
// 		err := db.CreateStudio(ctx, utils.Slug(studio), studio)
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	airingSeasonSlug := utils.Slug(animeData.Premiered)
// 	if airingSeasonSlug != "" {
// 		err = db.CreateTag(ctx, airingSeasonSlug, animeData.Premiered)
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	animeType := mal.ConvertMediaType(animeData.Type)
// 	animeStatus := mal.ConvertMediaStatus(animeData.Status)
// 	animeRating := mal.ConvertMediaRating(animeData.Rating)
//
// 	aniDbId := ""
// 	animeNewsNetworkId := ""
//
// 	if u, err := url.Parse(animeData.AniDBUrl); err == nil {
// 		aniDbId = u.Query().Get("aid")
// 	}
//
// 	if u, err := url.Parse(animeData.MediaNewsNetworkUrl); err == nil {
// 		animeNewsNetworkId = u.Query().Get("id")
// 	}
//
// 	err = db.UpdateMedia(ctx, animeId, database.AnimeChanges{
// 		AniDbId: database.Change[sql.NullString]{
// 			Value: sql.NullString{
// 				String: aniDbId,
// 				Valid:  aniDbId != "",
// 			},
// 			Changed: aniDbId != anime.AniDbId.String,
// 		},
//
// 		MediaNewsNetworkId: database.Change[sql.NullString]{
// 			Value: sql.NullString{
// 				String: animeNewsNetworkId,
// 				Valid:  animeNewsNetworkId != "",
// 			},
// 			Changed: animeNewsNetworkId != anime.MediaNewsNetworkId.String,
// 		},
//
// 		Title: database.Change[string]{
// 			Value:   animeData.Title,
// 			Changed: animeData.Title != anime.Title,
// 		},
//
// 		TitleEnglish: database.Change[sql.NullString]{
// 			Value: sql.NullString{
// 				String: animeData.TitleEnglish,
// 				Valid:  animeData.TitleEnglish != "",
// 			},
// 			Changed: animeData.TitleEnglish != anime.TitleEnglish.String,
// 		},
//
// 		Description: database.Change[sql.NullString]{
// 			Value: sql.NullString{
// 				String: animeData.Description,
// 				Valid:  animeData.Description != "",
// 			},
// 			Changed: animeData.Description != anime.Description.String,
// 		},
//
// 		Type: database.Change[types.MediaType]{
// 			Value:   animeType,
// 			Changed: animeType != anime.Type,
// 		},
//
// 		Status: database.Change[types.MediaStatus]{
// 			Value:   animeStatus,
// 			Changed: animeStatus != anime.Status,
// 		},
//
// 		Rating: database.Change[types.MediaRating]{
// 			Value:   animeRating,
// 			Changed: animeRating != anime.Rating,
// 		},
//
// 		AiringSeason: database.Change[sql.NullString]{
// 			Value: sql.NullString{
// 				String: airingSeasonSlug,
// 				Valid:  airingSeasonSlug != "",
// 			},
// 			Changed: true,
// 		},
//
// 		EpisodeCount: database.Change[sql.NullInt64]{
// 			Value: sql.NullInt64{
// 				Int64: utils.NullToDefault(animeData.EpisodeCount),
// 				Valid: animeData.EpisodeCount != nil,
// 			},
// 			Changed: true,
// 		},
//
// 		StartDate: database.Change[sql.NullString]{
// 			Value: sql.NullString{
// 				String: utils.NullToDefault(animeData.StartDate),
// 				Valid:  animeData.StartDate != nil,
// 			},
// 			Changed: true,
// 		},
//
// 		EndDate: database.Change[sql.NullString]{
// 			Value: sql.NullString{
// 				String: utils.NullToDefault(animeData.EndDate),
// 				Valid:  animeData.EndDate != nil,
// 			},
// 			Changed: true,
// 		},
//
// 		Score: database.Change[sql.NullFloat64]{
// 			Value: sql.NullFloat64{
// 				Float64: utils.NullToDefault(animeData.Score),
// 				Valid:   animeData.Score != nil,
// 			},
// 			Changed: true,
// 		},
//
// 		LastDataFetch: database.Change[sql.NullInt64]{
// 			Value: sql.NullInt64{
// 				Int64: time.Now().UnixMilli(),
// 				Valid: true,
// 			},
// 			Changed: true,
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	err = db.DeleteAllMediaThemeSongs(ctx, anime.Id)
// 	if err != nil {
// 		return err
// 	}
//
// 	for i, themeSong := range animeData.ThemeSongs {
// 		err := db.CreateMediaThemeSong(ctx, database.CreateAnimeThemeSongParams{
// 			MediaId: animeId,
// 			Idx:     i,
// 			Raw:     themeSong.Raw,
// 			Type:    mal.ConvertThemeSongType(themeSong.Type),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	for _, theme := range animeData.Themes {
// 		err := db.AddTagToMedia(ctx, animeId, utils.Slug(theme))
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	err = db.RemoveAllTagsFromMedia(ctx, animeId)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = db.RemoveAllStudiosFromMedia(ctx, animeId)
// 	if err != nil {
// 		return err
// 	}
//
// 	for _, genre := range animeData.Genres {
// 		err := db.AddTagToMedia(ctx, animeId, utils.Slug(genre))
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	for _, demographic := range animeData.Demographics {
// 		err := db.AddTagToMedia(ctx, animeId, utils.Slug(demographic))
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	for _, studio := range animeData.Studios {
// 		err := db.AddStudioToMedia(ctx, animeId, utils.Slug(studio))
// 		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
// 			return err
// 		}
// 	}
//
// 	return nil
// }

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
			logger.Debug("Client added", "numClients", len(broker.clients))
		case s := <-broker.closingClients:
			delete(broker.clients, s)
			logger.Debug("Removed client", "numClients", len(broker.clients))
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
	isDownloading    atomic.Bool
	downloadCanceled atomic.Bool

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

func (d *DownloadHandler) download(app core.App, mediaIds []string) error {
	d.isDownloading.Store(true)
	defer d.isDownloading.Store(false)

	// ctx := context.TODO()
	//
	// db := app.DB()
	// workDir := app.WorkDir()

	d.updateStatus(0, len(mediaIds))
	d.sendStatusEvent()

	for i, id := range mediaIds {
		_ = id

		if d.downloadCanceled.Load() {
			return fmt.Errorf("download canceled")
		}

		d.updateStatus(i+1, len(mediaIds))
		d.sendStatusEvent()

		// err := fetchAndUpdateMedia(ctx, db, workDir, id)
		// if err != nil {
		// 	return fmt.Errorf("failed to download media (%s): %w", id, err)
		// }
	}

	return nil
}

func (d *DownloadHandler) cancelDownload() {
	d.downloadCanceled.Store(true)
}

func (d *DownloadHandler) startDownload(app core.App, mediaIds []string) error {
	if d.isDownloading.Load() {
		return errors.New("already downloading")
	}

	d.downloadCanceled.Store(false)

	err := d.download(app, mediaIds)
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

type StartDownloadBody struct {
	Ids []string `json:"ids"`
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

		pyrin.ApiHandler{
			Name:     "StartDownload",
			Path:     "/system/download",
			Method:   http.MethodPost,
			BodyType: StartDownloadBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[StartDownloadBody](c)
				if err != nil {
					return nil, err
				}

				if downloadHandler.isDownloading.Load() {
					// TODO(patrik): Better error
					return nil, errors.New("already downloading")
				}

				go func() {
					err := downloadHandler.startDownload(app, body.Ids)
					if err != nil {
						logger.Error("failed to start downloader", "err", err)
					}
				}()

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "CancelDownload",
			Path:   "/system/download",
			Method: http.MethodDelete,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				if !downloadHandler.isDownloading.Load() {
					// TODO(patrik): Better error
					return nil, errors.New("not downloading")
				}

				downloadHandler.cancelDownload()

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
