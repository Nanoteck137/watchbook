package apis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/kr/pretty"
	"github.com/labstack/gommon/log"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/core"
	"github.com/nanoteck137/watchbook/database"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/pelletier/go-toml/v2"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

func slugifyArray(arr []string) []string {
	seen := map[string]bool{}
	res := make([]string, 0, len(arr))

	for _, value := range arr {
		value = utils.Slug(strings.TrimSpace(value))
		if value == "" {
			continue
		}

		if !seen[value] {
			seen[value] = true
			res = append(res, value)
		}
	}

	return res
}

// TODO(patrik): Add testing for this
func FixMedia(media *library.Media) error {
	// album := &metadata.Album
	//
	// album.Name = anvil.String(album.Name)
	//
	// if album.Year == 0 {
	// 	album.Year = metadata.General.Year
	// }
	//
	// if len(album.Artists) == 0 {
	// 	album.Artists = []string{UNKNOWN_ARTIST_NAME}
	// }
	//
	// album.Artists = fixArr(album.Artists)
	//
	// for i := range metadata.Tracks {
	// 	t := &metadata.Tracks[i]
	//
	// 	if t.Year == 0 {
	// 		t.Year = metadata.General.Year
	// 	}
	//
	// 	t.Name = anvil.String(t.Name)
	//
	// 	t.Tags = append(t.Tags, metadata.General.Tags...)
	// 	t.Tags = append(t.Tags, metadata.General.TrackTags...)
	//
	// 	if len(t.Artists) == 0 {
	// 		t.Artists = []string{UNKNOWN_ARTIST_NAME}
	// 	}
	//
	// 	t.Artists = fixArr(t.Artists)
	//
	// 	for i, tag := range t.Tags {
	// 		t.Tags[i] = utils.Slug(strings.TrimSpace(tag))
	// 	}
	// }

	media.General.Tags = slugifyArray(media.General.Tags)
	media.General.Studios = slugifyArray(media.General.Studios)

	// err := validate.ValidateStruct(&metadata.Album,
	// 	validate.Field(&metadata.Album.Name, validate.Required),
	// 	validate.Field(&metadata.Album.Artists, validate.Length(1, 0)),
	// )
	// if err != nil {
	// 	return err
	// }
	//
	// for _, track := range metadata.Tracks {
	// 	err := validate.ValidateStruct(&track,
	// 		validate.Field(&track.File, validate.Required),
	// 		validate.Field(&track.Name, validate.Required),
	// 		validate.Field(&track.Artists, validate.Length(1, 0)),
	// 	)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

type SyncHelper struct {
	media            map[string]struct{}
	mediaPathMapping map[string]string
}

// func (helper *SyncHelper) getOrCreateArtist(ctx context.Context, db *database.Database, name string) (string, error) {
// 	slug := utils.Slug(name)
//
// 	if artist, exists := helper.artists[slug]; exists {
// 		return artist, nil
// 	}
//
// 	dbArtist, err := db.GetArtistBySlug(ctx, slug)
// 	if err != nil {
// 		if errors.Is(err, database.ErrItemNotFound) {
// 			dbArtist, err = db.CreateArtist(ctx, database.CreateArtistParams{
// 				Slug: slug,
// 				Name: name,
// 			})
// 			if err != nil {
// 				return "", err
// 			}
// 		} else {
// 			return "", err
// 		}
// 	}
//
// 	helper.artists[slug] = dbArtist.Id
// 	return dbArtist.Id, nil
// }

func (helper *SyncHelper) setMediaTags(ctx context.Context, db *database.Database, mediaId string, tags []string) error {
	err := db.RemoveAllTagsFromMedia(ctx, mediaId)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		err := db.CreateTag(ctx, tag, tag)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddTagToMedia(ctx, mediaId, tag)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

func (helper *SyncHelper) setMediaStudios(ctx context.Context, db *database.Database, mediaId string, studios []string) error {
	err := db.RemoveAllStudiosFromMedia(ctx, mediaId)
	if err != nil {
		return err
	}

	for _, studio := range studios {
		err := db.CreateTag(ctx, studio, studio)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}

		err = db.AddStudioToMedia(ctx, mediaId, studio)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return err
		}
	}

	return nil
}

// TODO(patrik): Update the errors for album
func (helper *SyncHelper) syncMedia(ctx context.Context, media *library.Media, db *database.Database) error {
	err := FixMedia(media)
	if err != nil {
		return err
	}

	dbMedia, err := db.GetMediaById(ctx, nil, media.Id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			_, err = db.CreateMedia(ctx, database.CreateMediaParams{
				Id:    media.Id,
				Title: media.General.Title,
			})
			if err != nil {
				return fmt.Errorf("failed to create media: %w", err)
			}

			dbMedia, err = db.GetMediaById(ctx, nil, media.Id)
			if err != nil {
				return fmt.Errorf("failed to get media after creation: %w", err)
			}
		} else {
			return err
		}
	}

	helper.media[dbMedia.Id] = struct{}{}
	helper.mediaPathMapping[media.Path] = dbMedia.Id

	if media.General.AiringSeason != "" {
		err = db.CreateTag(ctx, media.General.AiringSeason, media.General.AiringSeason)
		if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
			return fmt.Errorf("failed to create airing season tag: %w", err)
		}
	}

	changes := database.MediaChanges{}

	changes.Type = database.Change[types.MediaType]{
		Value:   media.MediaType,
		Changed: media.MediaType != dbMedia.Type,
	}

	changes.TmdbId = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: media.Ids.TheMovieDB,
			Valid:  media.Ids.TheMovieDB != "",
		},
		Changed: media.Ids.TheMovieDB != dbMedia.TmdbId.String,
	}

	changes.MalId = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: media.Ids.MyAnimeList,
			Valid:  media.Ids.MyAnimeList != "",
		},
		Changed: media.Ids.MyAnimeList != dbMedia.MalId.String,
	}

	changes.AnilistId = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: media.Ids.Anilist,
			Valid:  media.Ids.Anilist != "",
		},
		Changed: media.Ids.Anilist != dbMedia.AnilistId.String,
	}

	changes.Title = database.Change[string]{
		Value:   media.General.Title,
		Changed: media.General.Title != dbMedia.Title,
	}

	// changes.Description

	changes.Score = database.Change[sql.NullFloat64]{
		Value: sql.NullFloat64{
			Float64: media.General.Score,
			Valid:   media.General.Score != 0.0,
		},
		Changed: media.General.Score != dbMedia.Score.Float64,
	}

	changes.Status = database.Change[types.MediaStatus]{
		Value:   media.General.Status,
		Changed: media.General.Status != dbMedia.Status,
	}

	changes.Rating = database.Change[types.MediaRating]{
		Value:   media.General.Rating,
		Changed: media.General.Rating != dbMedia.Rating,
	}

	changes.AiringSeason = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: media.General.AiringSeason,
			Valid:  media.General.AiringSeason != "",
		},
		Changed: media.General.AiringSeason != dbMedia.AiringSeason.String,
	}

	changes.StartDate = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: media.General.StartDate,
			Valid:  media.General.StartDate != "",
		},
		Changed: media.General.StartDate != dbMedia.StartDate.String,
	}

	changes.EndDate = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: media.General.EndDate,
			Valid:  media.General.EndDate != "",
		},
		Changed: media.General.EndDate != dbMedia.EndDate.String,
	}

	coverPath := media.GetCoverPath()
	changes.CoverFile = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: coverPath,
			Valid:  coverPath != "",
		},
		Changed: coverPath != dbMedia.CoverFile.String,
	}

	logoPath := media.GetLogoPath()
	changes.LogoFile = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: logoPath,
			Valid:  logoPath != "",
		},
		Changed: logoPath != dbMedia.LogoFile.String,
	}

	bannerPath := media.GetBannerPath()
	changes.BannerFile = database.Change[sql.NullString]{
		Value: sql.NullString{
			String: bannerPath,
			Valid:  bannerPath != "",
		},
		Changed: bannerPath != dbMedia.BannerFile.String,
	}

	err = db.UpdateMedia(ctx, dbMedia.Id, changes)
	if err != nil {
		return fmt.Errorf("failed to update media: %w", err)
	}

	err = helper.setMediaTags(ctx, db, dbMedia.Id, media.General.Tags)
	if err != nil {
		return fmt.Errorf("failed to set media tags: %w", err)
	}

	err = helper.setMediaStudios(ctx, db, dbMedia.Id, media.General.Studios)
	if err != nil {
		return fmt.Errorf("failed to set media studios: %w", err)
	}

	err = db.RemoveAllMediaParts(ctx, dbMedia.Id)
	if err != nil {
		return fmt.Errorf("failed to remove all media parts: %w", err)
	}

	for i, part := range media.Parts {
		err = db.CreateMediaPart(ctx, database.CreateMediaPartParams{
			Index:   int64(i),
			MediaId: dbMedia.Id,
			Name:    part.Name,
		})
		if err != nil {
			return fmt.Errorf("failed to add media part (%d): %w", i, err)
		}
	}

	return nil
}

func (helper *SyncHelper) syncCollection(ctx context.Context, collection *library.Collection, db *database.Database) error {
	// err := FixMedia(media)
	// if err != nil {
	// 	return err
	// }

	dbCollection, err := db.GetCollectionById(ctx, nil, collection.Id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			_, err = db.CreateCollection(ctx, database.CreateCollectionParams{
				Id:          collection.Id,
				Type:        types.CollectionTypeAnime,
				Name:        collection.General.Name,
			})
			if err != nil {
				return fmt.Errorf("failed to create media: %w", err)
			}

			dbCollection, err = db.GetCollectionById(ctx, nil, collection.Id)
			if err != nil {
				return fmt.Errorf("failed to get media after creation: %w", err)
			}
		} else {
			return err
		}
	}

	// TODO(patrik): Fill out
	changes := database.CollectionChanges{}

	changes.Name = database.Change[string]{
		Value:   collection.General.Name,
		Changed: collection.General.Name != dbCollection.Name,
	}

	err = db.UpdateCollection(ctx, dbCollection.Id, changes)
	if err != nil {
		return fmt.Errorf("failed to update collection: %w", err)
	}

	err = db.RemoveAllCollectionMediaItems(ctx, dbCollection.Id)
	if err != nil {
		return fmt.Errorf("failed to remove all media items from collection: %w", err)
	}

	for _, entry := range collection.Entries {
		p := path.Join(collection.Path, entry.Path)
		mediaId, ok := helper.mediaPathMapping[p]
		if !ok {
			return fmt.Errorf("failed to map path to media: %v", p)
		}

		err := db.CreateCollectionMediaItem(ctx, database.CreateCollectionMediaItemParams{
			CollectionId:   dbCollection.Id,
			MediaId:        mediaId,
			Name:           entry.SearchSlug,
			OrderNumber:    int64(entry.Order),
			SubOrderNumber: int64(entry.SubOrder),
		})
		if err != nil {
			return fmt.Errorf("failed to add media to collection: %w", err)
		}
	}

	return nil
}

type ReportType string

const (
	ReportTypeSearch ReportType = "search"
	ReportTypeSync   ReportType = "sync"
)

type SyncError struct {
	Type        ReportType `json:"type"`
	Message     string     `json:"message"`
	FullMessage *string    `json:"fullMessage,omitempty"`
}

type MissingMedia struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type SyncHandler struct {
	broker *Broker

	mutex sync.RWMutex

	isSyncing atomic.Bool

	syncErrors   []SyncError
	missingMedia []MissingMedia

	SyncChannel chan bool
}

type Report struct {
	SyncErrors   []SyncError    `json:"syncErrors"`
	MissingMedia []MissingMedia `json:"missingMedia"`
}

func (s *SyncHandler) GetReport() Report {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return Report{
		SyncErrors:   s.syncErrors,
		MissingMedia: s.missingMedia,
	}
}

func (s *SyncHandler) Cleanup(app core.App) error {
	ctx := context.TODO()

	for _, track := range s.missingMedia {
		err := app.DB().RemoveMedia(ctx, track.Id)
		if err != nil {
			return err
		}

		log.Info("Deleted media", "media", track)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.missingMedia = []MissingMedia{}

	return nil
}

func (s *SyncHandler) RunSync(app core.App) error {
	s.isSyncing.Store(true)
	defer s.isSyncing.Store(false)

	// TODO(patrik): Check for duplicated ids
	librarySearch, err := library.SearchLibrary(app.Config().LibraryDir)
	if err != nil {
		return err
	}

	log.Debug("Done searching for media")

	ctx := context.TODO()

	helper := SyncHelper{
		media:            map[string]struct{}{},
		mediaPathMapping: map[string]string{},
	}

	var syncErrors []error

	pretty.Println(librarySearch.Collections)

	for _, media := range librarySearch.Media {
		log.Debug("Syncing media", "path", media.Path)

		err := helper.syncMedia(ctx, &media, app.DB())
		if err != nil {
			syncErrors = append(syncErrors, err)
		}
	}

	for _, collection := range librarySearch.Collections {
		log.Debug("Syncing collection", "path", collection.Path)

		err := helper.syncCollection(ctx, &collection, app.DB())
		if err != nil {
			syncErrors = append(syncErrors, err)
		}
	}


	var missingMedia []MissingMedia

	{
		ids, err := app.DB().GetAllMediaIds(ctx)
		if err != nil {
			return err
		}

		for _, id := range ids {
			_, exists := helper.media[id]
			if !exists {
				media, err := app.DB().GetMediaById(ctx, nil, id)
				if err != nil {
					// TODO(patrik): How should we handle the error?
					continue
				}

				missingMedia = append(missingMedia, MissingMedia{
					Id:    id,
					Title: media.Title,
				})
			}
		}
	}

	errs := make([]SyncError, 0, len(librarySearch.Errors)+len(syncErrors))
	for _, err := range librarySearch.Errors {
		var fullMessage *string

		var tomlError *toml.DecodeError
		if errors.As(err, &tomlError) {
			m := tomlError.String()
			fullMessage = &m
		}

		errs = append(errs, SyncError{
			Type:        ReportTypeSearch,
			Message:     err.Error(),
			FullMessage: fullMessage,
		})
	}

	for _, err := range syncErrors {
		errs = append(errs, SyncError{
			Type:    ReportTypeSync,
			Message: err.Error(),
		})
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.syncErrors = errs
	s.missingMedia = missingMedia

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

func NewServer() (broker *Broker) {
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

var syncHandler = SyncHandler{
	SyncChannel: make(chan bool),
	broker:      NewServer(),
}

const (
	EventSyncing string = "syncing"
	EventReport  string = "report"
)

type SyncEvent struct {
	Syncing bool `json:"syncing"`
}

func (s SyncEvent) GetEventType() string {
	return EventSyncing
}

type ReportEvent struct {
	Report
}

func (s ReportEvent) GetEventType() string {
	return EventReport
}

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
			Name:         "SyncLibrary",
			Method:       http.MethodPost,
			Path:         "/system/library",
			ResponseType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				go func() {
					if syncHandler.isSyncing.Load() {
						log.Info("Syncing already")
						return
					}

					log.Info("Started library sync")

					syncHandler.broker.EmitEvent(SyncEvent{
						Syncing: true,
					})

					err := syncHandler.RunSync(app)
					if err != nil {
						log.Error("Failed to run sync", "err", err)
					}

					syncHandler.broker.EmitEvent(SyncEvent{
						Syncing: false,
					})

					syncHandler.broker.EmitEvent(ReportEvent{
						Report: syncHandler.GetReport(),
					})

					log.Info("Library sync done")
				}()

				return nil, nil
			},
		},

		// TODO(patrik): Better name?
		pyrin.ApiHandler{
			Name:   "CleanupLibrary",
			Method: http.MethodPost,
			Path:   "/system/library/cleanup",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				if syncHandler.isSyncing.Load() {
					return nil, errors.New("library is syncing")
				}

				err := syncHandler.Cleanup(app)
				if err != nil {
					return nil, err
				}

				syncHandler.broker.EmitEvent(ReportEvent{
					Report: syncHandler.GetReport(),
				})

				return nil, nil
			},
		},

		pyrin.NormalHandler{
			Name:   "SseHandler",
			Method: http.MethodGet,
			Path:   "/system/library/sse",
			HandlerFunc: func(c pyrin.Context) error {
				r := c.Request()
				w := c.Response()

				w.Header().Set("Content-Type", "text/event-stream")
				w.Header().Set("Cache-Control", "no-cache")
				w.Header().Set("Connection", "keep-alive")

				w.Header().Set("Access-Control-Allow-Origin", "*")

				rc := http.NewResponseController(w)

				eventChan := make(chan EventData)
				syncHandler.broker.newClients <- eventChan

				defer func() {
					syncHandler.broker.closingClients <- eventChan
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

				sendEvent(SyncEvent{
					Syncing: syncHandler.isSyncing.Load(),
				})

				sendEvent(ReportEvent{
					Report: syncHandler.GetReport(),
				})

				for {
					select {
					case <-r.Context().Done():
						syncHandler.broker.closingClients <- eventChan
						return nil

					case event := <-eventChan:
						sendEvent(event)
					}
				}
			},
		},
	)
}
