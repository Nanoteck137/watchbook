package tmdb

type TmdbMovieSearchRequestResult struct {
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

type TmdbTvSearchRequestResult struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	Id               int     `json:"id"`
	Name             string  `json:"name"`
	OriginalName     string  `json:"original_name"`
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

type TmdbSearchRequest[T any] struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
	Results      []T `json:"results"`
}

// var tmdbSearchCmd = &cobra.Command{
// 	Use: "search",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		var searchQuery string
// 		i := huh.NewInput().
// 			Title("Search Query").
// 			Value(&searchQuery)
// 		err := i.Run()
// 		if err != nil {
// 			logger.Fatal("failed to run input", "err", err)
// 		}
//
// 		fmt.Printf("searchQuery: %v\n", searchQuery)
//
// 		// url := "https://api.themoviedb.org/3/search/multi?query=the%20big%20short&include_adult=true&language=en-US&page=1"
// 		url, err := utils.CreateUrlBase("https://api.themoviedb.org", "/3/search/multi", url.Values{
// 			"query":         {searchQuery},
// 			"include_adult": {"true"},
// 			"language":      {"en-US"},
// 			"page":          {"1"},
// 		})
// 		if err != nil {
// 			logger.Fatal("failed to create url", "err", err)
// 		}
//
// 		req, err := http.NewRequest("GET", url.String(), nil)
// 		if err != nil {
// 			logger.Fatal("failed to create request", "err", err)
// 		}
//
// 		req.Header.Add("accept", "application/json")
// 		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")
//
// 		res, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			logger.Fatal("failed to send request", "err", err)
// 		}
// 		defer res.Body.Close()
//
// 		d := json.NewDecoder(res.Body)
//
// 		var decodedRes TmdbSearchRequest
// 		err = d.Decode(&decodedRes)
// 		if err != nil {
// 			logger.Fatal("failed to decode response body", "err", err)
// 		}
//
// 		pretty.Println(decodedRes)
// 	},
// }

// var tmdbMovieCmd = &cobra.Command{
// 	Use:  "movie <ID>",
// 	Args: cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		movieId := args[0]
//
// 		apiAddress, _ := cmd.Flags().GetString("api-address")
// 		client := api.New(apiAddress)
//
// 		tempDir, err := os.MkdirTemp("", "watchbook-cli-*")
// 		if err != nil {
// 			logger.Fatal("failed to create temp dir", "err", err)
// 		}
// 		defer os.RemoveAll(tempDir)
//
// 		search, err := client.GetMedia(api.Options{
// 			Query: url.Values{
// 				"filter": {fmt.Sprintf(`tmdbId=="movie@%s"`, movieId)},
// 			},
// 		})
// 		if err != nil {
// 			logger.Fatal("failed to get media", "err", err)
// 		}
//
// 		if len(search.Media) > 0 {
// 			logger.Warn("entry with id already exists", "id", movieId)
// 			return
// 		}
//
// 		// url := "https://api.themoviedb.org/3/movie/318846?language=en-US"
// 		url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/movie/%s", movieId), url.Values{
// 			"language": {"en-US"},
// 		})
// 		if err != nil {
// 			logger.Fatal("failed to create url", "err", err)
// 		}
//
// 		req, err := http.NewRequest("GET", url.String(), nil)
// 		if err != nil {
// 			logger.Fatal("failed to create request", "err", err)
// 		}
//
// 		req.Header.Add("accept", "application/json")
// 		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")
//
// 		res, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			logger.Fatal("failed to send request", "err", err)
// 		}
// 		defer res.Body.Close()
//
// 		d := json.NewDecoder(res.Body)
//
// 		var decodedRes TmdbMovieDetails
// 		err = d.Decode(&decodedRes)
// 		if err != nil {
// 			logger.Fatal("failed to decode response body", "err", err)
// 		}
//
// 		pretty.Println(decodedRes)
//
// 		status := types.MediaStatusUpcoming
// 		switch decodedRes.Status {
// 		case "Released":
// 			status = types.MediaStatusCompleted
// 		}
//
// 		studios := make([]string, len(decodedRes.ProductionCompanies))
//
// 		for i, company := range decodedRes.ProductionCompanies {
// 			studios[i] = company.Name
// 		}
//
// 		tags := make([]string, 0, len(decodedRes.Genres))
//
// 		for _, genre := range decodedRes.Genres {
// 			tags = append(tags, genre.Name)
// 		}
//
// 		airingSeason := types.GetAiringSeason(decodedRes.ReleaseDate)
//
// 		{
// 			res, err := client.CreateMedia(api.CreateMediaBody{
// 				MediaType:    string(types.MediaTypeMovie),
// 				TmdbId:       fmt.Sprintf("movie@%s", movieId),
// 				ImdbId:       decodedRes.ImdbId,
// 				Title:        decodedRes.Title,
// 				Description:  decodedRes.Overview,
// 				Score:        float32(utils.RoundFloat(float64(decodedRes.VoteAverage), 2)),
// 				Status:       string(status),
// 				Rating:       string(types.MediaRatingUnknown),
// 				AiringSeason: airingSeason,
// 				StartDate:    decodedRes.ReleaseDate,
// 				EndDate:      decodedRes.ReleaseDate,
// 				Creators:     studios,
// 				Tags:         tags,
// 			}, api.Options{})
// 			if err != nil {
// 				logger.Fatal("failed to create media", "err", err, "id", movieId)
// 			}
//
// 			_, err = client.SetParts(res.Id, api.SetPartsBody{
// 				Parts: []api.PartBody{
// 					{
// 						Name: decodedRes.Title,
// 					},
// 				},
// 			}, api.Options{})
// 			if err != nil {
// 				logger.Fatal("failed to set parts", "err", err, "id", movieId)
// 			}
//
// 			coverUrl := "http://image.tmdb.org/t/p/original" + decodedRes.PosterPath
// 			coverPath, err := utils.DownloadImage(coverUrl, tempDir, "cover")
// 			if err != nil {
// 				logger.Fatal("failed to download image", "err", err, "id", movieId)
// 			}
//
// 			bannerUrl := "http://image.tmdb.org/t/p/original" + decodedRes.BackdropPath
// 			bannerPath, err := utils.DownloadImage(bannerUrl, tempDir, "banner")
// 			if err != nil {
// 				logger.Fatal("failed to download image", "err", err, "id", movieId)
// 			}
//
// 			images, err := createImageForm(coverPath, "", bannerPath)
// 			if err != nil {
// 				logger.Fatal("failed to create image form", "err", err, "id", movieId)
// 			}
//
// 			_, err = client.ChangeMediaImages(res.Id, images.Boundary, &images.Buf, api.Options{})
// 			if err != nil {
// 				logger.Fatal("failed to set images", "err", err, "id", movieId)
// 			}
// 		}
// 	},
// }

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

// var tmdbTvCmd = &cobra.Command{
// 	Use:  "tv <ID>",
// 	Args: cobra.ExactArgs(1),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		tvId := args[0]
//
// 		apiAddress, _ := cmd.Flags().GetString("api-address")
// 		client := api.New(apiAddress)
//
// 		tempDir, err := os.MkdirTemp("", "watchbook-cli-*")
// 		if err != nil {
// 			logger.Fatal("failed to create temp dir", "err", err)
// 		}
// 		defer os.RemoveAll(tempDir)
//
// 		// url := "https://api.themoviedb.org/3/movie/318846?language=en-US"
// 		// url := "https://api.themoviedb.org/3/tv/66732?language=en-US"
// 		url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/tv/%s", tvId), url.Values{
// 			"language": {"en-US"},
// 		})
// 		if err != nil {
// 			logger.Fatal("failed to create url", "err", err)
// 		}
//
// 		req, err := http.NewRequest("GET", url.String(), nil)
// 		if err != nil {
// 			logger.Fatal("failed to create request", "err", err)
// 		}
//
// 		req.Header.Add("accept", "application/json")
// 		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")
//
// 		res, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			logger.Fatal("failed to send request", "err", err)
// 		}
// 		defer res.Body.Close()
//
// 		if res.StatusCode != 200 {
// 			logger.Fatal("failed to get response", "status", res.Status)
// 		}
//
// 		d := json.NewDecoder(res.Body)
//
// 		var decodedRes TmdbTvDetails
// 		err = d.Decode(&decodedRes)
// 		if err != nil {
// 			logger.Fatal("failed to decode response body", "err", err)
// 		}
//
// 		pretty.Println(decodedRes)
//
// 		col, err := client.CreateCollection(api.CreateCollectionBody{
// 			CollectionType: "series",
// 			Name:           decodedRes.Name,
// 		}, api.Options{})
// 		if err != nil {
// 			pretty.Println(err)
// 			logger.Fatal("failed to create tv collection", "err", err, "id", tvId)
// 		}
//
// 		coverUrl := "http://image.tmdb.org/t/p/original" + decodedRes.PosterPath
// 		coverPath, err := utils.DownloadImage(coverUrl, tempDir, "cover")
// 		if err != nil {
// 			logger.Fatal("failed to download image", "err", err, "id", tvId)
// 		}
//
// 		bannerUrl := "http://image.tmdb.org/t/p/original" + decodedRes.BackdropPath
// 		bannerPath, err := utils.DownloadImage(bannerUrl, tempDir, "banner")
// 		if err != nil {
// 			logger.Fatal("failed to download image", "err", err, "id", tvId)
// 		}
//
// 		images, err := createImageForm(coverPath, "", bannerPath)
// 		if err != nil {
// 			logger.Fatal("failed to create image form", "err", err, "id", tvId)
// 		}
//
// 		_, err = client.ChangeCollectionImages(col.Id, images.Boundary, &images.Buf, api.Options{})
// 		if err != nil {
// 			logger.Fatal("failed to set images", "err", err, "id", tvId)
// 		}
//
// 		for _, season := range decodedRes.Seasons {
// 			id := fmt.Sprintf("tv@%s/%d", tvId, season.SeasonNumber)
//
// 			name := ""
// 			entryName := ""
// 			if season.SeasonNumber == 0 {
// 				name = "Specials"
// 				name = fmt.Sprintf("%s (Specials)", decodedRes.Name)
// 				entryName = "Specials"
// 			} else {
// 				name = fmt.Sprintf("%s (Season %d)", decodedRes.Name, season.SeasonNumber)
// 				entryName = fmt.Sprintf("Season %d", season.SeasonNumber)
// 			}
//
// 			status := types.MediaStatusUpcoming
// 			if types.IsReleased(season.AirDate) {
// 				status = types.MediaStatusCompleted
// 			}
//
// 			studios := make([]string, 0, len(decodedRes.ProductionCompanies)+len(decodedRes.Networks))
//
// 			for _, company := range decodedRes.ProductionCompanies {
// 				studios = append(studios, company.Name)
// 			}
//
// 			for _, network := range decodedRes.Networks {
// 				studios = append(studios, network.Name)
// 			}
//
// 			tags := make([]string, 0, len(decodedRes.Genres))
//
// 			for _, genre := range decodedRes.Genres {
// 				tags = append(tags, genre.Name)
// 			}
//
// 			airingSeason := types.GetAiringSeason(season.AirDate)
//
// 			res, err := client.CreateMedia(api.CreateMediaBody{
// 				MediaType:    string(types.MediaTypeTV),
// 				TmdbId:       id,
// 				Title:        name,
// 				Description:  season.Overview,
// 				Score:        float32(utils.RoundFloat(float64(season.VoteAverage), 2)),
// 				Status:       string(status),
// 				Rating:       string(types.MediaRatingUnknown),
// 				AiringSeason: airingSeason,
// 				StartDate:    season.AirDate,
// 				EndDate:      season.AirDate,
// 				Creators:     studios,
// 				Tags:         tags,
// 			}, api.Options{})
// 			if err != nil {
// 				logger.Fatal("failed to create media", "err", err, "id", tvId)
// 			}
//
// 			parts := make([]api.PartBody, 0, season.EpisodeCount)
// 			for i := range season.EpisodeCount {
// 				parts = append(parts, api.PartBody{
// 					Name: fmt.Sprintf("Episode %d", i+1),
// 				})
// 			}
//
// 			_, err = client.SetParts(res.Id, api.SetPartsBody{
// 				Parts: parts,
// 			}, api.Options{})
// 			if err != nil {
// 				logger.Fatal("failed to set parts", "err", err, "id", id)
// 			}
//
// 			coverUrl := "http://image.tmdb.org/t/p/original" + season.PosterPath
// 			coverPath, err := utils.DownloadImage(coverUrl, tempDir, "cover")
// 			if err != nil {
// 				logger.Fatal("failed to download image", "err", err, "id", id)
// 			}
//
// 			// bannerUrl := "http://image.tmdb.org/t/p/original" + season.BackdropPath
// 			// bannerPath, err := utils.DownloadImage(bannerUrl, tempDir, "banner")
// 			// if err != nil {
// 			// 	logger.Fatal("failed to download image", "err", err, "id", id)
// 			// }
//
// 			images, err := createImageForm(coverPath, "", "")
// 			if err != nil {
// 				logger.Fatal("failed to create image form", "err", err, "id", id)
// 			}
//
// 			_, err = client.ChangeMediaImages(res.Id, images.Boundary, &images.Buf, api.Options{})
// 			if err != nil {
// 				logger.Fatal("failed to set images", "err", err, "id", id)
// 			}
//
// 			_, err = client.AddCollectionItem(col.Id, api.AddCollectionItemBody{
// 				MediaId:    res.Id,
// 				Name:       entryName,
// 				SearchSlug: utils.Slug(entryName),
// 				Order:      season.SeasonNumber,
// 			}, api.Options{})
// 			if err != nil {
// 				logger.Fatal("failed to add media to collection", "err", err, "id", id, "colId", col.Id, "mediaId", res.Id)
// 			}
// 		}
// 	},
// }
