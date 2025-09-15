package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nanoteck137/watchbook/utils"
)

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
	VoteAverage         float64                 `json:"vote_average"`          //: 7.357,
	VoteCount           int                     `json:"vote_count"`            //: 9313
}

type TmdbTvSeasonDetailsEpisode struct {
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	EpisodeType    string  `json:"episode_type"`
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	Runtime        int     `json:"runtime"`
	SeasonNumber   int     `json:"season_number"`
	ShowId         int     `json:"show_id"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`

	//	  "crew": [
	//	    {
	//	      "department": "Writing",
	//	      "job": "Writer",
	//	      "credit_id": "52542275760ee313280006ce",
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 66633,
	//	      "known_for_department": "Writing",
	//	      "name": "Vince Gilligan",
	//	      "original_name": "Vince Gilligan",
	//	      "popularity": 0.8583,
	//	      "profile_path": "/z3E0DhBg1V1PZVEtS9vfFPzOWYB.jpg"
	//	    },
	//	    {
	//	      "department": "Directing",
	//	      "job": "Director",
	//	      "credit_id": "52542275760ee313280006e8",
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 66633,
	//	      "known_for_department": "Writing",
	//	      "name": "Vince Gilligan",
	//	      "original_name": "Vince Gilligan",
	//	      "popularity": 0.8583,
	//	      "profile_path": "/z3E0DhBg1V1PZVEtS9vfFPzOWYB.jpg"
	//	    },
	//	    {
	//	      "job": "Director of Photography",
	//	      "department": "Camera",
	//	      "credit_id": "52b7029219c29533d00dd2c1",
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 2483,
	//	      "known_for_department": "Camera",
	//	      "name": "John Toll",
	//	      "original_name": "John Toll",
	//	      "popularity": 0.1689,
	//	      "profile_path": "/cfL9A6tC7L5Ps5fq1o3WpVKGMb1.jpg"
	//	    },
	//	    {
	//	      "job": "Editor",
	//	      "department": "Editing",
	//	      "credit_id": "52b702ea19c2955402183a66",
	//	      "adult": false,
	//	      "gender": 1,
	//	      "id": 1280071,
	//	      "known_for_department": "Editing",
	//	      "name": "Lynne Willingham",
	//	      "original_name": "Lynne Willingham",
	//	      "popularity": 0.1176,
	//	      "profile_path": null
	//	    },
	//	    {
	//	      "job": "Art Direction",
	//	      "department": "Art",
	//	      "credit_id": "62feade5cf4a640080998241",
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 1018092,
	//	      "known_for_department": "Art",
	//	      "name": "James F. Oberlander",
	//	      "original_name": "James F. Oberlander",
	//	      "popularity": 0.0338,
	//	      "profile_path": null
	//	    },
	//	    {
	//	      "job": "Associate Producer",
	//	      "department": "Production",
	//	      "credit_id": "6418a04de7414600b96bf1bd",
	//	      "adult": false,
	//	      "gender": 1,
	//	      "id": 1808170,
	//	      "known_for_department": "Production",
	//	      "name": "Gina Scheerer",
	//	      "original_name": "Gina Scheerer",
	//	      "popularity": 0.0311,
	//	      "profile_path": null
	//	    }
	//	  ],
	//	  "guest_stars": [
	//	    {
	//	      "character": "Steven Gomez",
	//	      "credit_id": "5271b489760ee35b3e0881a7",
	//	      "order": 8,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 61535,
	//	      "known_for_department": "Acting",
	//	      "name": "Steven Michael Quezada",
	//	      "original_name": "Steven Michael Quezada",
	//	      "popularity": 0.4711,
	//	      "profile_path": "/pVYrDkwI6GWvCNL2kJhpDJfBFyd.jpg"
	//	    },
	//	    {
	//	      "character": "Jock",
	//	      "credit_id": "52542275760ee313280006b4",
	//	      "order": 500,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 1216132,
	//	      "known_for_department": "Acting",
	//	      "name": "Aaron Hill",
	//	      "original_name": "Aaron Hill",
	//	      "popularity": 0.8133,
	//	      "profile_path": "/rNp31SeoVqSQU6OZWxZUhGwAgyq.jpg"
	//	    },
	//	    {
	//	      "character": "Dr. Belknap",
	//	      "credit_id": "52725cb1760ee3044d0b9984",
	//	      "order": 502,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 161591,
	//	      "known_for_department": "Acting",
	//	      "name": "Gregory Chase",
	//	      "original_name": "Gregory Chase",
	//	      "popularity": 0.0327,
	//	      "profile_path": "/gNdodev00CROpXuAh5EFmkWTVOo.jpg"
	//	    },
	//	    {
	//	      "character": "Krazy-8",
	//	      "credit_id": "52725845760ee3046b09426e",
	//	      "order": 504,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 1046460,
	//	      "known_for_department": "Acting",
	//	      "name": "Max Arciniega",
	//	      "original_name": "Max Arciniega",
	//	      "popularity": 0.3614,
	//	      "profile_path": "/eqKAJKPpt41KpsLIkkBnAY4HMAL.jpg"
	//	    },
	//	    {
	//	      "character": "Bogdan Wolynetz",
	//	      "credit_id": "5272587a760ee3045009ddec",
	//	      "order": 575,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 1223197,
	//	      "known_for_department": "Acting",
	//	      "name": "Marius Stan",
	//	      "original_name": "Marius Stan",
	//	      "popularity": 0.0589,
	//	      "profile_path": "/zX8fpNkyKEtQX3zTvks1hVhrOz7.jpg"
	//	    },
	//	    {
	//	      "character": "Carmen Molina",
	//	      "credit_id": "52542273760ee31328000676",
	//	      "order": 643,
	//	      "adult": false,
	//	      "gender": 1,
	//	      "id": 115688,
	//	      "known_for_department": "Acting",
	//	      "name": "Carmen Serano",
	//	      "original_name": "Carmen Serano",
	//	      "popularity": 0.2728,
	//	      "profile_path": "/nzJEe2UqvvMIBJZP1aeFEj4EunN.jpg"
	//	    },
	//	    {
	//	      "character": "Chad's Girlfriend",
	//	      "credit_id": "56846abbc3a36836280008d4",
	//	      "order": 651,
	//	      "adult": false,
	//	      "gender": 1,
	//	      "id": 1223192,
	//	      "known_for_department": "Art",
	//	      "name": "Roberta Marquez Seret",
	//	      "original_name": "Roberta Marquez Seret",
	//	      "popularity": 0.0762,
	//	      "profile_path": null
	//	    },
	//	    {
	//	      "character": "Chad",
	//	      "credit_id": "63012a1a33a376007a442d63",
	//	      "order": 675,
	//	      "adult": false,
	//	      "gender": 0,
	//	      "id": 3670896,
	//	      "known_for_department": "Acting",
	//	      "name": "Evan Bobrick",
	//	      "original_name": "Evan Bobrick",
	//	      "popularity": 0.0239,
	//	      "profile_path": null
	//	    },
	//	    {
	//	      "character": "E.M.T",
	//	      "credit_id": "63012a3d97eab4007d00192b",
	//	      "order": 676,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 36135,
	//	      "known_for_department": "Acting",
	//	      "name": "Christopher Dempsey",
	//	      "original_name": "Christopher Dempsey",
	//	      "popularity": 0.0757,
	//	      "profile_path": "/pTngvks30p74j40TaniMkg4tbhn.jpg"
	//	    },
	//	    {
	//	      "character": "Irving",
	//	      "credit_id": "63012a5c33a376007f87247b",
	//	      "order": 677,
	//	      "adult": false,
	//	      "gender": 0,
	//	      "id": 2969089,
	//	      "known_for_department": "Production",
	//	      "name": "Allan Pacheco",
	//	      "original_name": "Allan Pacheco",
	//	      "popularity": 0.0214,
	//	      "profile_path": null
	//	    },
	//	    {
	//	      "character": "Chemistry Student",
	//	      "credit_id": "63012a655f4b73007faa4261",
	//	      "order": 678,
	//	      "adult": false,
	//	      "gender": 0,
	//	      "id": 3670897,
	//	      "known_for_department": "Acting",
	//	      "name": "Jason Byrd",
	//	      "original_name": "Jason Byrd",
	//	      "popularity": 0.0143,
	//	      "profile_path": null
	//	    },
	//	    {
	//	      "character": "Sexy Neighbor",
	//	      "credit_id": "63012a7e33a376007f872481",
	//	      "order": 679,
	//	      "adult": false,
	//	      "gender": 0,
	//	      "id": 219124,
	//	      "known_for_department": "Acting",
	//	      "name": "Linda Speciale",
	//	      "original_name": "Linda Speciale",
	//	      "popularity": 0.1767,
	//	      "profile_path": null
	//	    },
	//	    {
	//	      "character": "Jock's Friend #1",
	//	      "credit_id": "63012a8bfb5299007d660bc8",
	//	      "order": 680,
	//	      "adult": false,
	//	      "gender": 0,
	//	      "id": 3212534,
	//	      "known_for_department": "Acting",
	//	      "name": "Jesús Ramírez",
	//	      "original_name": "Jesús Ramírez",
	//	      "popularity": 0.0239,
	//	      "profile_path": "/1EfPZxdFNNi3LFLR9laLcVROAko.jpg"
	//	    },
	//	    {
	//	      "character": "Jock's Friend #2",
	//	      "credit_id": "63012ac4c2f44b007d249b54",
	//	      "order": 681,
	//	      "adult": false,
	//	      "gender": 0,
	//	      "id": 3670906,
	//	      "known_for_department": "Acting",
	//	      "name": "Joshua S. Patton",
	//	      "original_name": "Joshua S. Patton",
	//	      "popularity": 0.0289,
	//	      "profile_path": null
	//	    },
	//	    {
	//	      "character": "Emilio Koyama",
	//	      "credit_id": "631aff1f62f335007ed32cb3",
	//	      "order": 703,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 92495,
	//	      "known_for_department": "Acting",
	//	      "name": "John Koyama",
	//	      "original_name": "John Koyama",
	//	      "popularity": 0.3386,
	//	      "profile_path": "/AwtHbt8qO7D3EFonG5lqml8xgwb.jpg"
	//	    },
	//	    {
	//	      "character": "DEA Agent #1 (uncredited)",
	//	      "credit_id": "655bb4fa10923000ab494163",
	//	      "order": 848,
	//	      "adult": false,
	//	      "gender": 2,
	//	      "id": 1335375,
	//	      "known_for_department": "Crew",
	//	      "name": "Ed Duran",
	//	      "original_name": "Ed Duran",
	//	      "popularity": 0.1805,
	//	      "profile_path": "/mzPJkVKg7ve3whmvyf2TyDIuewr.jpg"
	//	    }
	//	  ]
	//	},
}

type TmdbTvSeasonDetails struct {
	Id           string                       `json:"_id"`
	AirDate      string                       `json:"air_date"`
	Name         string                       `json:"name"`
	Episodes     []TmdbTvSeasonDetailsEpisode `json:"episodes"`
	Overview     string                       `json:"overview"`
	SerieId      int                          `json:"id"`
	PosterPath   string                       `json:"poster_path"`
	SeasonNumber int                          `json:"season_number"`
	VoteAverage  float64                      `json:"vote_average"`
}

type TmdbImage struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	AspectRatio float64 `json:"aspect_ratio"`
	Iso6391     *string `json:"iso_639_1"`
	FilePath    string  `json:"file_path"`
	VoteAverge  float64     `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

type TmdbImages struct {
	Id        int         `json:"id"`
	Backdrops []TmdbImage `json:"backdrops"`
	Logos     []TmdbImage `json:"logos"`
	Posters   []TmdbImage `json:"posters"`
}

func movieSearch(query string) (TmdbSearchRequest[TmdbMovieSearchRequestResult], error) {
	url, err := utils.CreateUrlBase("https://api.themoviedb.org", "/3/search/movie", url.Values{
		"query":         {query},
		"include_adult": {"true"},
		"language":      {"en-US"},
		"page":          {"1"},
	})

	if err != nil {
		return TmdbSearchRequest[TmdbMovieSearchRequestResult]{}, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return TmdbSearchRequest[TmdbMovieSearchRequestResult]{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TmdbSearchRequest[TmdbMovieSearchRequestResult]{}, err
	}

	defer res.Body.Close()

	d := json.NewDecoder(res.Body)

	var decodedRes TmdbSearchRequest[TmdbMovieSearchRequestResult]
	err = d.Decode(&decodedRes)
	if err != nil {
		return TmdbSearchRequest[TmdbMovieSearchRequestResult]{}, err
	}

	return decodedRes, nil
}

func tvSearch(query string) (TmdbSearchRequest[TmdbTvSearchRequestResult], error) {
	url, err := utils.CreateUrlBase("https://api.themoviedb.org", "/3/search/tv", url.Values{
		"query":         {query},
		"include_adult": {"true"},
		"language":      {"en-US"},
		"page":          {"1"},
	})

	if err != nil {
		return TmdbSearchRequest[TmdbTvSearchRequestResult]{}, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return TmdbSearchRequest[TmdbTvSearchRequestResult]{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TmdbSearchRequest[TmdbTvSearchRequestResult]{}, err
	}

	defer res.Body.Close()

	d := json.NewDecoder(res.Body)

	var decodedRes TmdbSearchRequest[TmdbTvSearchRequestResult]
	err = d.Decode(&decodedRes)
	if err != nil {
		return TmdbSearchRequest[TmdbTvSearchRequestResult]{}, err
	}

	return decodedRes, nil
}

func getMovieDetails(id string) (TmdbMovieDetails, error) {
	// url := "https://api.themoviedb.org/3/movie/318846?language=en-US"
	url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/movie/%s", id), url.Values{
		"language": {"en-US"},
	})
	if err != nil {
		return TmdbMovieDetails{}, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return TmdbMovieDetails{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TmdbMovieDetails{}, err
	}
	defer res.Body.Close()

	d := json.NewDecoder(res.Body)

	var decodedRes TmdbMovieDetails
	err = d.Decode(&decodedRes)
	if err != nil {
		return TmdbMovieDetails{}, err
	}

	return decodedRes, nil
}

func getTvDetails(id string) (TmdbTvDetails, error) {
	url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/tv/%s", id), url.Values{
		"language": {"en-US"},
	})
	if err != nil {
		return TmdbTvDetails{}, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return TmdbTvDetails{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TmdbTvDetails{}, err
	}
	defer res.Body.Close()

	// if res.StatusCode != 200 {
	// 	return TmdbTvDetails{}, err
	// }

	d := json.NewDecoder(res.Body)

	var decodedRes TmdbTvDetails
	err = d.Decode(&decodedRes)
	if err != nil {
		return TmdbTvDetails{}, err
	}

	return decodedRes, nil
}

func getTvSeasonDetails(id, season string) (TmdbTvSeasonDetails, error) {
	url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/tv/%s/season/%s", id, season), url.Values{
		"language": {"en-US"},
	})
	if err != nil {
		return TmdbTvSeasonDetails{}, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return TmdbTvSeasonDetails{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TmdbTvSeasonDetails{}, err
	}
	defer res.Body.Close()

	// if res.StatusCode != 200 {
	// 	return TmdbTvDetails{}, err
	// }

	d := json.NewDecoder(res.Body)

	var decodedRes TmdbTvSeasonDetails
	err = d.Decode(&decodedRes)
	if err != nil {
		return TmdbTvSeasonDetails{}, err
	}

	return decodedRes, nil
}

func getMovieImages(id string) (TmdbImages, error) {
	url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/movie/%s/images", id), url.Values{
		"include_image_language": {"en"},
	})
	if err != nil {
		return TmdbImages{}, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return TmdbImages{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TmdbImages{}, err
	}
	defer res.Body.Close()

	// if res.StatusCode != 200 {
	// 	return TmdbTvDetails{}, err
	// }

	d := json.NewDecoder(res.Body)

	var decodedRes TmdbImages
	err = d.Decode(&decodedRes)
	if err != nil {
		return TmdbImages{}, err
	}

	return decodedRes, nil
}

func getTvImages(id string) (TmdbImages, error) {
	url, err := utils.CreateUrlBase("https://api.themoviedb.org", fmt.Sprintf("/3/tv/%s/images", id), url.Values{
		"include_image_language": {"en"},
	})
	if err != nil {
		return TmdbImages{}, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return TmdbImages{}, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMWVlNGNjODJiM2MxNDI3MTI4MjQ0Zjg2MmRmYzdmNCIsIm5iZiI6MTc0MDUzMDU2MC4xNzEsInN1YiI6IjY3YmU2MzgwMGZmNjY0M2U3MTNkM2FiYiIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.95qulxC-Skebj9avAqPtxgEuKZwlT_Bqp36iuqOF9wc")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TmdbImages{}, err
	}
	defer res.Body.Close()

	// if res.StatusCode != 200 {
	// 	return TmdbTvDetails{}, err
	// }

	d := json.NewDecoder(res.Body)

	var decodedRes TmdbImages
	err = d.Decode(&decodedRes)
	if err != nil {
		return TmdbImages{}, err
	}

	return decodedRes, nil
}
