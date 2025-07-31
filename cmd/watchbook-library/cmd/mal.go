package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/kr/pretty"
	"github.com/nanoteck137/watchbook/cmd/watchbook-library/config"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/provider/myanimelist"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

type LibraryDir string

func (d LibraryDir) String() string {
	return string(d)
}

func (d LibraryDir) MalDir() string {
	return path.Join(d.String(), "mal")
}

func (d LibraryDir) MalAnimesDir() string {
	return path.Join(d.MalDir(), "animes")
}

func (d LibraryDir) MalMangasDir() string {
	return path.Join(d.MalDir(), "mangas")
}

func (d LibraryDir) MalSeriesDir() string {
	return path.Join(d.MalDir(), "series")
}

func (d LibraryDir) MalDownloadDir() string {
	return path.Join(d.MalDir(), "download")
}

func ensureMalDirs(libraryDir LibraryDir) error {
	dirs := []string{
		libraryDir.MalDir(),
		libraryDir.MalAnimesDir(), 
		libraryDir.MalMangasDir(), 
		libraryDir.MalSeriesDir(), 
		libraryDir.MalDownloadDir(), 
	}

	for _, dir := range dirs {
		err := os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	return nil
}

func openLibrary(dir string) (*library.LibrarySearch, error) {
	search, err := library.SearchLibrary(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to search library: %w", err)
	}

	if len(search.Errors) > 0 {
		return nil, fmt.Errorf("library has errors: %v", search.Errors)
	}

	return search, nil
}

func prepareLibraryMal(search *library.LibrarySearch) map[string]library.Media {
	entries := make(map[string]library.Media)

	for _, media := range search.Media {
		id := media.Ids.MyAnimeList
		if id == "" {
			continue
		}

		entries[id] = media
	}

	return entries
}

func saveMedia(out string, media library.Media) error {
	d, err := toml.Marshal(media)
	if err != nil {
		return fmt.Errorf("failed to marshal media: %w", err)
	}

	err = os.WriteFile(path.Join(out, "media.toml"), d, 0644)
	if err != nil {
		return fmt.Errorf("failed to save marshaled media: %w", err)
	}

	return nil
}

var malCmd = &cobra.Command{
	Use: "mal",
}

var malGetCmd = &cobra.Command{
	Use:  "get <MAL_ID>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		malId := args[0]

		cfg := config.LoadedConfig

		libraryDir := LibraryDir(cfg.LibraryDir)

		err := ensureMalDirs(libraryDir)
		if err != nil {
			logger.Fatal("failed to create mal dirs", "err", err)
		}

		search, err := openLibrary(libraryDir.String())
		if err != nil {
			logger.Fatal("failed to open library", "err", err, "dir", libraryDir.String())
		}

		entries := prepareLibraryMal(search)

		if m, exists := entries[malId]; exists {
			logger.Fatal("entry with id already exists", "id", malId, "path", m.Path)
		}

		data, err := myanimelist.RawGetAnime(malId)
		if err != nil {
			logger.Fatal("failed to get anime data", "err", err)
		}

		title := data.Title
		if data.TitleEnglish != "" {
			title = data.TitleEnglish
		}

		score := 0.0
		if data.Score != nil {
			score = *data.Score
		}

		startDate := ""
		if data.StartDate != nil {
			startDate = *data.StartDate
		}

		endDate := ""
		if data.EndDate != nil {
			endDate = *data.EndDate
		}

		media := library.Media{
			Id:        utils.CreateMediaId(),
			MediaType: data.Type,
			Ids: library.MediaIds{
				MyAnimeList: malId,
			},
			General: library.MediaGeneral{
				Title:        title,
				Description:  data.Description,
				Score:        score,
				Status:       data.Status,
				Rating:       data.Rating,
				AiringSeason: data.AiringSeason,
				StartDate:    startDate,
				EndDate:      endDate,
				Studios:      data.Studios,
				Tags:         data.Tags,
			},
			Images: library.Images{},
			Parts:  []library.MediaPart{},
		}

		if media.MediaType.IsMovie() {
			media.Parts = []library.MediaPart{
				{
					Name: media.General.Title,
				},
			}
		} else if data.EpisodeCount != nil {
			media.Parts = make([]library.MediaPart, 0, *data.EpisodeCount)
			for i := range *data.EpisodeCount {
				media.Parts = append(media.Parts, library.MediaPart{
					Name: fmt.Sprintf("Episode %d", i+1),
				})
			}
		}

		out := path.Join(libraryDir.MalDownloadDir(), malId+"-"+utils.Slug(media.General.Title))
		err = os.Mkdir(out, 0755)
		if err != nil {
			logger.Fatal("failed to create dir for anime", "err", err, "title", media.General.Title)
		}

		p, err := utils.DownloadImage(data.CoverImageUrl, out, "cover")
		if err != nil {
			logger.Fatal("failed to download image", "err", err)
		}

		media.Images.Cover = path.Base(p)

		err = saveMedia(out, media)
		if err != nil {
			logger.Fatal("failed to save media", "err", err, "title", media.General.Title)
		}

		pretty.Println(media)
	},
}

var malTestCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadedConfig

		libraryDir := LibraryDir(cfg.LibraryDir)

		err := ensureMalDirs(libraryDir)
		if err != nil {
			logger.Fatal("failed to create mal dirs", "err", err)
		}

		search, err := openLibrary(libraryDir.String())
		if err != nil {
			logger.Fatal("failed to open library", "err", err, "dir", libraryDir.String())
		}

		entries := prepareLibraryMal(search)

		seasonal, err := myanimelist.FetchSeasonal("winter", 2021)
		if err != nil {
			logger.Fatal("failed", "err", err)
		}

		for _, anime := range seasonal.Animes {
			if m, exists := entries[anime.Id]; exists {
				logger.Warn("entry with id already exists", "id", anime.Id, "path", m.Path)
				continue
			}

			title := anime.Title
			if anime.TitleEnglish != "" {
				title = anime.TitleEnglish
			}

			airingSeason := types.GetAiringSeason(anime.StartDate)

			var tags []string

			for _, genre := range anime.Genres {
				tags = append(tags, genre)
			}

			for _, theme := range anime.Themes {
				tags = append(tags, theme)
			}

			for _, demographic := range anime.Demographics {
				tags = append(tags, demographic)
			}

			media := library.Media{
				Id:        utils.CreateMediaId(),
				MediaType: myanimelist.ConvertAnimeType(anime.Type),
				Ids: library.MediaIds{
					MyAnimeList: anime.Id,
				},
				General: library.MediaGeneral{
					Title:        title,
					Score:        anime.Score,
					Description:  anime.Description,
					Status:       myanimelist.ConvertAnimeStatus(anime.Status),
					Rating:       myanimelist.ConvertAnimeRating(anime.Rating),
					AiringSeason: airingSeason,
					StartDate:    anime.StartDate,
					EndDate:      anime.EndDate,
					Studios:      anime.Studios,
					Tags:         tags,
				},
				Images: library.Images{},
				Parts:  []library.MediaPart{},
			}

			if media.MediaType.IsMovie() {
				media.Parts = []library.MediaPart{
					{
						Name: media.General.Title,
					},
				}
			} else {
				media.Parts = make([]library.MediaPart, 0, anime.EpisodeCount)
				for i := range anime.EpisodeCount {
					media.Parts = append(media.Parts, library.MediaPart{
						Name: fmt.Sprintf("Episode %d", i+1),
					})
				}
			}

			out := path.Join(libraryDir.MalDownloadDir(), anime.Id+"-"+utils.Slug(media.General.Title))
			err = os.Mkdir(out, 0755)
			if err != nil {
				logger.Fatal("failed to create dir for anime", "err", err, "title", media.General.Title)
			}

			p, err := utils.DownloadImage(anime.CoverImageUrl, out, "cover")
			if err != nil {
				logger.Fatal("failed to download image", "err", err)
			}

			media.Images.Cover = path.Base(p)

			err = saveMedia(out, media)
			if err != nil {
				logger.Fatal("failed to save media", "err", err, "title", media.General.Title)
			}

			entries[anime.Id] = media
		}
	},
}

// TODO(patrik): Move to watchbook-cli
// var importMalUserCmd = &cobra.Command{
// 	Use: "import-mal-user",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		client := api.New("http://localhost:3000")
// 		client.Headers.Set("X-Api-Token", "cga632fp69xf0f1ottw4nqc55yhtziym")
//
// 		watchlist, err := myanimelist.GetUserWatchlist("nanoteck137")
// 		if err != nil {
// 			logger.Fatal("failed get mal userlist", "err", err)
// 		}
//
// 		watchlist = watchlist[1:]
//
// 		for _, entry := range watchlist {
// 			pretty.Println(entry)
//
// 			media, err := client.GetMedia(api.Options{
// 				Query: url.Values{
// 					"filter": {fmt.Sprintf("malId == \"%d\"", entry.AnimeId)},
// 				},
// 				Header: http.Header{},
// 			})
// 			if err != nil {
// 				logger.Fatal("failed get media", "err", err)
// 			}
//
// 			// pretty.Println(media)
//
// 			for _, media := range media.Media {
// 				list := "plan-to-watch"
//
// 				switch entry.Status {
// 				case myanimelist.WatchlistStatusCurrentlyWatching:
// 					list = "watching"
// 				case myanimelist.WatchlistStatusCompleted:
// 					list = "completed"
// 				case myanimelist.WatchlistStatusOnHold:
// 					list = "on-hold"
// 				case myanimelist.WatchlistStatusDropped:
// 					list = "dropped"
// 				case myanimelist.WatchlistStatusPlanToWatch:
// 					list = "plan-to-watch"
// 				}
//
// 				var score *int
// 				if media.Score != nil {
// 					s := entry.Score
// 					score = &s
// 				}
//
// 				var currentPart *int
// 				if entry.NumWatchedEpisodes > 0 {
// 					currentPart = &entry.NumWatchedEpisodes
// 				}
//
// 				var isRevisiting *bool
// 				if entry.IsRewatching > 0 {
// 					v := true
// 					isRevisiting = &v
// 				}
//
// 				_, err := client.SetMediaUserData(media.Id, api.SetMediaUserData{
// 					List:         &list,
// 					Score:        score,
// 					CurrentPart:  currentPart,
// 					RevisitCount: new(int),
// 					IsRevisiting: isRevisiting,
// 				}, api.Options{})
// 				if err != nil {
// 					logger.Fatal("failed set media user data", "err", err)
// 				}
// 			}
//
// 			// break
// 		}
//
// 	},
// }

func init() {
	malCmd.AddCommand(malGetCmd)
	malCmd.AddCommand(malTestCmd)

	rootCmd.AddCommand(malCmd)
}
