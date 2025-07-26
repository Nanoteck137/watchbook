package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/kr/pretty"
	"github.com/nanoteck137/watchbook/library"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var tmdbCmd = &cobra.Command{
	Use: "tmdb",
}

type TmdbSearchRequestResult struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	Id               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
	MediaType        string  `json:"media_type"`
	OriginalLanguage string  `json:"original_language"`
	GenreIds         []int   `json:"genre_ids"`
	Popularity       float32 `json:"popularity"`
	ReleaseDate      string  `json:"release_date"`
	Video            bool    `json:"video"`
	VoteAverage      float32 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type TmdbSearchRequest struct {
	Page         int                       `json:"page"`
	TotalPages   int                       `json:"total_pages"`
	TotalResults int                       `json:"total_results"`
	Results      []TmdbSearchRequestResult `json:"results"`
}

var tmdbSearchCmd = &cobra.Command{
	Use: "search",
	Run: func(cmd *cobra.Command, args []string) {
		var searchQuery string
		i := huh.NewInput().
			Title("Search Query").
			Value(&searchQuery)
		err := i.Run()
		if err != nil {
			logger.Fatal("failed to run input", "err", err)
		}

		fmt.Printf("searchQuery: %v\n", searchQuery)

		// url := "https://api.themoviedb.org/3/search/multi?query=the%20big%20short&include_adult=true&language=en-US&page=1"
		url, err := utils.CreateUrlBase("https://api.themoviedb.org", "/3/search/multi", url.Values{
			"query":         {searchQuery},
			"include_adult": {"true"},
			"language":      {"en-US"},
			"page":          {"1"},
		})
		if err != nil {
			logger.Fatal("failed to create url", "err", err)
		}

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			logger.Fatal("failed to create request", "err", err)
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Fatal("failed to send request", "err", err)
		}
		defer res.Body.Close()

		d := json.NewDecoder(res.Body)

		var decodedRes TmdbSearchRequest
		err = d.Decode(&decodedRes)
		if err != nil {
			logger.Fatal("failed to decode response body", "err", err)
		}

		pretty.Println(decodedRes)
	},
}

type TmdbProductionCompany struct {
	Id            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type TmdbGenre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TmdbMovieDetails struct {
	Adult               bool                    `json:"adult"`                 //: false,
	BackdropPath        string                  `json:"backdrop_path"`         //: "/i7UCf0ysjbYYaqcSKUox9BJz4Kp.jpg",
	BelongsToCollection any                     `json:"belongs_to_collection"` //: null,
	Budget              int                     `json:"budget"`                //: 28000000,
	Genres              []TmdbGenre             `json:"genres"`                //: [
	Homepage            string                  `json:"homepage"`              //: "http://www.thebigshortmovie.com",
	Id                  int                     `json:"id"`                    //: 318846,
	ImdbId              string                  `json:"imdb_id"`               //: "tt1596363",
	OriginCountry       any                     `json:"origin_country"`        //: [
	OriginalLanguage    string                  `json:"original_language"`     //: "en",
	OriginalTitle       string                  `json:"original_title"`        //: "The Big Short",
	Overview            string                  `json:"overview"`              //: "The men who made millions from a global economic meltdown.",
	Popularity          float32                 `json:"popularity"`            //: 7.7052,
	PosterPath          string                  `json:"poster_path"`           //: "/scVEaJEwP8zUix8vgmMoJJ9Nq0w.jpg",
	ProductionCompanies []TmdbProductionCompany `json:"production_companies"`  //: [
	ProductionCountries any                     `json:"production_countries"`  //: [
	ReleaseDate         string                  `json:"release_date"`          //: "2015-12-11",
	Revenue             int                     `json:"revenue"`               //: 133346506,
	Runtime             int                     `json:"runtime"`               //: 131,
	SpokenLanguages     any                     `json:"spoken_languages"`      //: [
	Status              string                  `json:"status"`                //: "Released",
	Tagline             string                  `json:"tagline"`               //: "This is a true story.",
	Title               string                  `json:"title"`                 //: "The Big Short",
	Video               bool                    `json:"video"`                 //: false,
	VoteAverage         float32                 `json:"vote_average"`          //: 7.357,
	VoteCount           int                     `json:"vote_count"`            //: 9313
}

func downloadImage(url, outDir, name string) (string, error) {
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

	out := path.Join(outDir, name+ext)

	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("downloadImage: failed to open output file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", fmt.Errorf("downloadImage: failed io.Copy: %w", err)
	}

	return out, nil
}

var tmdbMovieCmd = &cobra.Command{
	Use:  "movie <ID>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := "./work/library-tmdb"
		movieId := 318846

		// url := "https://api.themoviedb.org/3/movie/318846?language=en-US"
		url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/movie/%d", movieId), url.Values{
			"language": {"en-US"},
		})
		if err != nil {
			logger.Fatal("failed to create url", "err", err)
		}

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			logger.Fatal("failed to create request", "err", err)
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Fatal("failed to send request", "err", err)
		}
		defer res.Body.Close()

		d := json.NewDecoder(res.Body)

		var decodedRes TmdbMovieDetails
		err = d.Decode(&decodedRes)
		if err != nil {
			logger.Fatal("failed to decode response body", "err", err)
		}

		pretty.Println(decodedRes)

		status := types.MediaStatusNotAired
		switch decodedRes.Status {
		case "Released":
			status = types.MediaStatusFinished
		}

		studios := make([]string, len(decodedRes.ProductionCompanies))

		for i, company := range decodedRes.ProductionCompanies {
			studios[i] = company.Name
		}

		tags := make([]string, 0, len(decodedRes.Genres))

		for _, genre := range decodedRes.Genres {
			tags = append(tags, genre.Name)
		}

		airingSeason := types.GetAiringSeason(decodedRes.ReleaseDate)

		media := library.Media{
			Id:        utils.CreateMediaId(),
			MediaType: types.MediaTypeMovie,
			Ids: library.MediaIds{
				TheMovieDB:  fmt.Sprintf("movie@%d", movieId),
				Imdb:        decodedRes.ImdbId,
				MyAnimeList: "",
				Anilist:     "",
			},
			General: library.MediaGeneral{
				Title:        decodedRes.Title,
				Score:        utils.RoundFloat(float64(decodedRes.VoteAverage), 2),
				Status:       status,
				Rating:       types.MediaRatingUnknown,
				AiringSeason: airingSeason,
				StartDate:    decodedRes.ReleaseDate,
				EndDate:      decodedRes.ReleaseDate,
				Studios:      studios,
				Tags:         tags,
			},
			Parts: []library.MediaPart{
				{
					Name: decodedRes.Title,
				},
			},
		}

		out := path.Join(dir, strconv.Itoa(movieId)+"-"+utils.Slug(media.General.Title))
		err = os.Mkdir(out, 0755)
		if err != nil {
			logger.Fatal("failed to create dir for movie", "err", err, "title", media.General.Title)
		}

		{
			url := "http://image.tmdb.org/t/p/original" + decodedRes.PosterPath
			out, err := downloadImage(url, out, "cover")
			if err != nil {
				logger.Fatal("failed to download cover image", "err", err, "title", media.General.Title)
			}

			media.Images.Cover = path.Base(out)
		}

		{
			url := "http://image.tmdb.org/t/p/original" + decodedRes.BackdropPath
			out, err := downloadImage(url, out, "banner")
			if err != nil {
				logger.Fatal("failed to download banner image", "err", err, "title", media.General.Title)
			}

			media.Images.Banner = path.Base(out)
		}

		{
			d, err := toml.Marshal(media)
			if err != nil {
				logger.Fatal("failed to marshal media", "err", err, "title", media.General.Title)
			}

			dst := path.Join(out, "media.toml")
			err = os.WriteFile(dst, d, 0644)
			if err != nil {
				logger.Fatal("failed to write media for movie", "err", err, "title", media.General.Title, "dstPath", dst)
			}
		}
	},
}

type TmdbTvDetailsSeason struct {
	AirDate      string  `json:"air_date"`      //:      "2016-07-15",
	EpisodeCount int     `json:"episode_count"` //: float64(8),
	Id           int     `json:"id"`            //:            float64(77680),
	Name         string  `json:"name"`          //:          "Season 1",
	Overview     string  `json:"overview"`      //:      "Strange things are afoot in Hawkins, Indiana, where a young boy's sudden
	PosterPath   string  `json:"poster_path"`   //:   "/hb4KmdA4F1XMnf0vVjevPKKBEjV.jpg",
	SeasonNumber int     `json:"season_number"` //: float64(1),
	VoteAverage  float64 `json:"vote_average"`  //:  float64(8.4),
}

type TmdbTvDetails struct {
	Adult               bool                    `json:"adult"`                //: false,
	BackdropPath        string                  `json:"backdrop_path"`        //: "/56v2KjBlU4XaOv9rVYEQypROD7P.jpg",
	CreatedBy           any                     `json:"created_by"`           //: [
	EpisodeRunTime      any                     `json:"episode_run_time"`     //: [],
	FirstAirDate        string                  `json:"first_air_date"`       //: "2016-07-15",
	Genres              []TmdbGenre             `json:"genres"`               //: [
	Homepage            string                  `json:"homepage"`             //: "https://www.netflix.com/title/80057281",
	Id                  int                     `json:"id"`                   //: 66732,
	InProduction        bool                    `json:"in_production"`        //: true,
	Languages           any                     `json:"languages"`            //: [
	LastAirDate         string                  `json:"last_air_date"`        //: "2022-07-01",
	LastEpisodeToAir    any                     `json:"last_episode_to_air"`  //: {
	Name                string                  `json:"name"`                 //: "Stranger Things",
	NextEpisodeToAir    any                     `json:"next_episode_to_air"`  //: {
	Networks            []TmdbProductionCompany `json:"networks"`             //: [
	NumberOfEpisodes    int                     `json:"number_of_episodes"`   //: 42,
	NumberOfSeasons     int                     `json:"number_of_seasons"`    //: 5,
	OriginCountry       any                     `json:"origin_country"`       //: [
	OriginalLanguage    string                  `json:"original_language"`    //: "en",
	OriginalName        string                  `json:"original_name"`        //: "Stranger Things",
	Overview            string                  `json:"overview"`             //: "When a young boy vanishes, a small town uncovers a mystery involving secret experiments, terrifying supernatural forces, and one strange little girl.",
	Popularity          float64                 `json:"popularity"`           //: 167.3198,
	PosterPath          string                  `json:"poster_path"`          //: "/uOOtwVbSr4QDjAGIifLDwpb2Pdl.jpg",
	ProductionCompanies []TmdbProductionCompany `json:"production_companies"` //: [
	ProductionCountries any                     `json:"production_countries"` //: [
	Seasons             []TmdbTvDetailsSeason   `json:"seasons"`              //: [
	SpokenLanguages     any                     `json:"spoken_languages"`     //: [
	Status              string                  `json:"status"`               //: "Returning Series",
	Tagline             string                  `json:"tagline"`              //: "One last adventure.",
	Type                string                  `json:"type"`                 //: "Scripted",
	VoteAverage         float64                 `json:"vote_average"`         //: 8.6,
	VoteCount           int                     `json:"vote_count"`           //: 18534
}

var tmdbTvCmd = &cobra.Command{
	Use:  "tv <ID>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := "./work/library-tmdb"
		_ = dir
		// tvId := 1396
		tvId := 66732

		// url := "https://api.themoviedb.org/3/movie/318846?language=en-US"
		// url := "https://api.themoviedb.org/3/tv/66732?language=en-US"
		url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/tv/%d", tvId), url.Values{
			"language": {"en-US"},
		})
		if err != nil {
			logger.Fatal("failed to create url", "err", err)
		}

		req, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			logger.Fatal("failed to create request", "err", err)
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Fatal("failed to send request", "err", err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			logger.Fatal("failed to get response", "status", res.Status)
		}

		d := json.NewDecoder(res.Body)

		var decodedRes TmdbTvDetails
		err = d.Decode(&decodedRes)
		if err != nil {
			logger.Fatal("failed to decode response body", "err", err)
		}

		pretty.Println(decodedRes)

		out := path.Join(dir, strconv.Itoa(tvId)+"-"+utils.Slug(decodedRes.Name))
		err = os.Mkdir(out, 0755)
		if err != nil {
			logger.Fatal("failed to create dir for tv", "err", err)
		}

		collectionEntries := make([]library.CollectionEntry, 0, len(decodedRes.Seasons))

		for _, season := range decodedRes.Seasons {
			name := ""
			if season.SeasonNumber == 0 {
				name = "Specials"
			} else {
				name = fmt.Sprintf("%s (Season %d)", decodedRes.Name, season.SeasonNumber)
			}

			out := path.Join(out, utils.Slug(name))
			err = os.Mkdir(out, 0755)
			if err != nil {
				logger.Fatal("failed to create dir for tv season", "err", err)
			}

			status := types.MediaStatusNotAired
			if types.IsReleased(season.AirDate) {
				status = types.MediaStatusFinished
			}

			studios := make([]string, 0, len(decodedRes.ProductionCompanies)+len(decodedRes.Networks))

			for _, company := range decodedRes.ProductionCompanies {
				studios = append(studios, company.Name)
			}

			for _, network := range decodedRes.Networks {
				studios = append(studios, network.Name)
			}

			tags := make([]string, 0, len(decodedRes.Genres))

			for _, genre := range decodedRes.Genres {
				tags = append(tags, genre.Name)
			}

			airingSeason := types.GetAiringSeason(season.AirDate)

			media := library.Media{
				Id:        utils.CreateMediaId(),
				MediaType: types.MediaTypeSeason,
				Ids: library.MediaIds{
					TheMovieDB:  fmt.Sprintf("tv@%d/%d", tvId, season.SeasonNumber),
					Imdb:        "",
					MyAnimeList: "",
					Anilist:     "",
				},
				General: library.MediaGeneral{
					Title:        name,
					Score:        utils.RoundFloat(float64(season.VoteAverage), 2),
					Status:       status,
					Rating:       types.MediaRatingUnknown,
					AiringSeason: airingSeason,
					StartDate:    season.AirDate,
					EndDate:      season.AirDate,
					Studios:      studios,
					Tags:         tags,
				},
				Parts: []library.MediaPart{},
			}

			media.Parts = make([]library.MediaPart, season.EpisodeCount)

			for i := range season.EpisodeCount {
				media.Parts[i] = library.MediaPart{
					Name: fmt.Sprintf("Episode %d", i+1),
				}
			}

			{
				url := "http://image.tmdb.org/t/p/original" + season.PosterPath
				out, err := downloadImage(url, out, "cover")
				if err != nil {
					logger.Fatal("failed to download cover image", "err", err, "title", media.General.Title)
				}

				media.Images.Cover = path.Base(out)
			}

			// {
			// 	url := "http://image.tmdb.org/t/p/original" + decodedRes.BackdropPath
			// 	out, err := downloadImage(url, out, "banner")
			// 	if err != nil {
			// 		logger.Fatal("failed to download banner image", "err", err, "title", media.General.Title)
			// 	}
			//
			// 	media.Images.Banner = path.Base(out)
			// }

			{
				d, err := toml.Marshal(media)
				if err != nil {
					logger.Fatal("failed to marshal media", "err", err, "title", media.General.Title)
				}

				dst := path.Join(out, "media.toml")
				err = os.WriteFile(dst, d, 0644)
				if err != nil {
					logger.Fatal("failed to write media for tv season", "err", err, "title", media.General.Title, "dstPath", dst)
				}
			}

			collectionEntries = append(collectionEntries, library.CollectionEntry{
				Path:       path.Base(out),
				SearchSlug: fmt.Sprintf("season-%d", season.SeasonNumber),
				Order:      season.SeasonNumber,
			})
		}

		collection := library.Collection{
			Id: utils.CreateCollectionId(),
			General: library.CollectionGeneral{
				Name: decodedRes.Name,
			},
			Images:  library.Images{},
			Entries: collectionEntries,
		}

		{
			url := "http://image.tmdb.org/t/p/original" + decodedRes.PosterPath
			out, err := downloadImage(url, out, "cover")
			if err != nil {
				logger.Fatal("failed to download cover image", "err", err, "name", collection.General.Name)
			}

			collection.Images.Cover = path.Base(out)
		}

		{
			url := "http://image.tmdb.org/t/p/original" + decodedRes.BackdropPath
			out, err := downloadImage(url, out, "banner")
			if err != nil {
				logger.Fatal("failed to download banner image", "err", err, "name", collection.General.Name)
			}

			collection.Images.Cover = path.Base(out)
		}

		{
			d, err := toml.Marshal(collection)
			if err != nil {
				logger.Fatal("failed to marshal collection", "err", err, "name", collection.General.Name)
			}

			dst := path.Join(out, "collection.toml")
			err = os.WriteFile(dst, d, 0644)
			if err != nil {
				logger.Fatal("failed to write collection for tv", "err", err, "name", collection.General.Name, "dstPath", dst)
			}
		}
	},
}

func init() {
	tmdbCmd.AddCommand(tmdbSearchCmd)
	tmdbCmd.AddCommand(tmdbMovieCmd)
	tmdbCmd.AddCommand(tmdbTvCmd)

	rootCmd.AddCommand(tmdbCmd)
}
