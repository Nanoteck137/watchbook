package myanimelist

import (
	"log/slog"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nanoteck137/watchbook/utils"
)

type ThemeSongType string

const (
	ThemeSongOpening ThemeSongType = "opening"
	ThemeSongEnding  ThemeSongType = "ending"
)

type ThemeSong struct {
	Index  int64         `json:"index"`
	Name   string        `json:"name"`
	Artist string        `json:"artist"`
	Type   ThemeSongType `json:"type"`

	Raw string `json:"raw"`
}

type RelatedEntry struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	Relation string `json:"relation"`
}

type Anime struct {
	Id string `json:"id"`

	Title        string `json:"title"`
	TitleEnglish string `json:"titleEnglish"`

	Description string `json:"description"`

	CoverImageUrl string `json:"coverImageUrl"`

	Type         string `json:"type"`
	Status       string `json:"status"`
	EpisodeCount *int64 `json:"episodeCount"`
	Rating       string `json:"rating"`
	Premiered    string `json:"premiered"`
	Source       string `json:"source"`
	Broadcast    string `json:"broadcast"`

	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`

	ScoreRaw string   `json:"scoreRaw"`
	Score    *float64 `json:"score"`

	Studios   []string `json:"studios"`
	Producers []string `json:"producers"`

	Genres       []string `json:"genres"`
	Themes       []string `json:"themes"`
	Demographics []string `json:"demographics"`

	ThemeSongs     []ThemeSong    `json:"themeSongs"`
	RelatedEntries []RelatedEntry `json:"relatedEntries"`

	AniDBUrl            string `json:"aniDBUrl"`
	AnimeNewsNetworkUrl string `json:"animeNewsNetworkUrl"`

	Pictures []string `json:"pictures"`
}

type SeasonalAnime struct {
	Id string `json:"id"`

	Title        string `json:"title"`
	TitleEnglish string `json:"titleEnglish"`

	Description string `json:"description"`

	CoverImageUrl string `json:"coverImageUrl"`

	Type         string `json:"type"`
	Status       string `json:"status"`
	EpisodeCount int64  `json:"episodeCount"`
	Rating       string `json:"rating"`

	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`

	Score float64 `json:"score"`

	Studios []string `json:"studios"`

	Genres       []string `json:"genres"`
	Themes       []string `json:"themes"`
	Demographics []string `json:"demographics"`
}

type EpisodeExtraInfo struct {
	Total   int64 `json:"total"`
	Current int64 `json:"current"`
}

type Episode struct {
	EnglishTitle    string  `json:"englishTitle"`
	JapaneseTitle   string  `json:"japaneseTitle"`
	Number          int64   `json:"number"`
	Aired           string  `json:"aired"`
	AverageScoreRaw string  `json:"averageScoreRaw"`
	AverageScore    float64 `json:"averageScore"`
}

type Seasonal struct {
	Animes []SeasonalAnime
}

func ExtractAnimeData(pagePath string) (Anime, error) {
	f, err := os.Open(pagePath)
	if err != nil {
		return Anime{}, err
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return Anime{}, err
	}

	img := doc.Find(".leftside img")
	imageUrl, _ := img.Attr("data-src")

	title := doc.Find(".title-name strong").Text()
	title = strings.TrimSpace(title)

	titleEnglish := doc.Find(".title-english").Text()
	titleEnglish = strings.TrimSpace(titleEnglish)

	desc := doc.Find("p[itemprop=\"description\"]").Text()

	var score *float64
	scoreRaw := doc.Find(".score-label").First().Text()
	s, err := strconv.ParseFloat(scoreRaw, 64)
	if err == nil {
		score = &s
	}

	typ := doc.Find(".type > a").Text()
	typ = strings.TrimSpace(typ)

	leftside := doc.Find("div .leftside")

	var studios []string
	leftside.Find("span:contains(\"Studios:\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
		raw := s.Text()
		raw = strings.TrimSpace(raw)

		if raw != "add some" {
			studios = append(studios, s.Text())
		}
	})

	var producers []string
	leftside.Find("span:contains(\"Producers:\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
		raw := s.Text()
		raw = strings.TrimSpace(raw)

		if raw != "add some" {
			producers = append(producers, s.Text())
		}
	})

	status := leftside.Find("span:contains(\"Status:\")").Parent().Children().Remove().End().Text()
	status = strings.TrimSpace(status)

	episodesRaw := leftside.Find("span:contains(\"Episodes:\")").Parent().Children().Remove().End().Text()
	episodesRaw = strings.TrimSpace(episodesRaw)

	var episodeCount *int64

	if episodesRaw != "Unknown" {
		c, err := strconv.ParseInt(episodesRaw, 10, 32)
		if err != nil {
			return Anime{}, err
		}

		episodeCount = &c
	}

	var startDate *string
	var endDate *string

	aired := leftside.Find("span:contains(\"Aired:\")").Parent().Children().Remove().End().Text()
	aired = strings.TrimSpace(aired)

	if aired != "Not available" {
		dates := strings.Split(aired, "to")

		{
			start := dates[0]
			start = strings.TrimSpace(start)
			t, err := ParseDate(start)
			if err == nil {
				f := formatDate(t)
				startDate = &f
			}
		}

		if len(dates) >= 2 {
			end := dates[1]
			if end != "?" {
				end = strings.TrimSpace(end)
				t, err := ParseDate(end)
				if err == nil {
					f := formatDate(t)
					endDate = &f
				}
			}
		}
	}

	var genres []string
	leftside.Find("span:contains(\"Genres:\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
		t := s.Text()
		t = strings.TrimSpace(t)

		if t != "" {
			genres = append(genres, s.Text())
		}
	})

	leftside.Find("span:contains(\"Genre:\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
		t := s.Text()
		t = strings.TrimSpace(t)

		if t != "" {
			genres = append(genres, s.Text())
		}
	})

	var themes []string
	leftside.Find("span:contains(\"Theme:\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
		t := s.Text()
		t = strings.TrimSpace(t)

		if t != "" {
			themes = append(themes, s.Text())
		}
	})

	leftside.Find("span:contains(\"Themes:\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
		t := s.Text()
		t = strings.TrimSpace(t)

		if t != "" {
			themes = append(themes, s.Text())
		}
	})

	var demographic []string
	leftside.Find("span:contains(\"Demographic:\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
		t := s.Text()
		t = strings.TrimSpace(t)

		if t != "" {
			demographic = append(demographic, s.Text())
		}
	})

	leftside.Find("span:contains(\"Demographics:\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
		t := s.Text()
		t = strings.TrimSpace(t)

		if t != "" {
			demographic = append(demographic, s.Text())
		}
	})

	source := leftside.Find("span:contains(\"Source:\")").Parent().Children().Text() //Parent().Children().Remove().End().Text()
	source = strings.TrimPrefix(source, "Source:")
	source = strings.TrimSpace(source)

	broadcast := leftside.Find("span:contains(\"Broadcast:\")").Parent().Children().Remove().End().Text()
	broadcast = strings.TrimSpace(broadcast)

	rating := leftside.Find("span:contains(\"Rating:\")").Parent().Children().Remove().End().Text()
	rating = strings.TrimSpace(rating)

	premiered := leftside.Find("span:contains(\"Premiered:\")").Parent().Find("a").Text()
	premiered = strings.TrimSpace(premiered)

	relatedEntriesEl := doc.Find(".related-entries")

	var relatedEntries []RelatedEntry

	relatedEntriesEl.Find(".entries-tile .entry").Each(func(i int, s *goquery.Selection) {
		relation := s.Find(".relation").Text()
		relation = strings.TrimSpace(relation)
		relation = utils.FixSpaces(relation)

		titleEl := s.Find(".title a")

		href, _ := titleEl.Attr("href")
		title := titleEl.Text()
		title = strings.TrimSpace(title)

		relatedEntries = append(relatedEntries, RelatedEntry{
			Title:    title,
			Url:      href,
			Relation: relation,
		})
	})

	aniDbUrl, _ := doc.Find("a[data-ga-click-type=\"external-links-anime-pc-anidb\"]").Attr("href")
	annUrl, _ := doc.Find("a[data-ga-click-type=\"external-links-anime-pc-ann\"]").Attr("href")

	var themeSongs []ThemeSong

	doc.Find(".theme-songs > table").Each(func(i int, s *goquery.Selection) {
		ending := s.Parent().HasClass("ending")

		typ := ThemeSongOpening
		if ending {
			typ = ThemeSongEnding
		}

		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			raw := s.Text()
			raw = strings.TrimSpace(raw)

			if strings.HasPrefix(raw, "No opening themes") {
				return
			}

			if strings.HasPrefix(raw, "No ending themes") {
				return
			}

			indexRaw := s.Find(".theme-song-index").Text()
			indexRaw = strings.TrimSuffix(indexRaw, ":")
			index, _ := strconv.ParseInt(indexRaw, 10, 32)

			artist := s.Find(".theme-song-artist").Text()
			artist = strings.TrimSpace(artist)
			artist = strings.TrimPrefix(artist, "by ")

			// episode := s.Find(".theme-song-episode").Text()
			// fmt.Printf("episode: %v\n", episode)

			name := s.Find(".theme-song-artist").Parent().Children().Remove().End().Text()
			name = strings.TrimSpace(name)
			name = strings.Trim(name, "\"")

			themeSongs = append(themeSongs, ThemeSong{
				Index:  index,
				Name:   name,
				Artist: artist,
				Type:   typ,
				Raw:    raw,
			})
		})
	})

	return Anime{
		Title:               title,
		TitleEnglish:        titleEnglish,
		Description:         desc,
		CoverImageUrl:       imageUrl,
		Type:                typ,
		Status:              status,
		EpisodeCount:        episodeCount,
		Rating:              rating,
		Premiered:           premiered,
		Source:              source,
		Broadcast:           broadcast,
		StartDate:           startDate,
		EndDate:             endDate,
		ScoreRaw:            scoreRaw,
		Score:               score,
		Studios:             studios,
		Producers:           producers,
		Genres:              genres,
		Themes:              themes,
		Demographics:        demographic,
		ThemeSongs:          themeSongs,
		RelatedEntries:      relatedEntries,
		AniDBUrl:            aniDbUrl,
		AnimeNewsNetworkUrl: annUrl,
	}, nil
}

func ExtractEpisodeData(pagePath string) ([]Episode, EpisodeExtraInfo, error) {
	f, err := os.Open(pagePath)
	if err != nil {
		return nil, EpisodeExtraInfo{}, err
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return nil, EpisodeExtraInfo{}, err
	}

	r := regexp.MustCompile(`\((\d+)\/(\d+)\)`)

	episodeCount := doc.Find("tbody > tr > td > div > .di-ib").Text()
	captures := r.FindStringSubmatch(episodeCount)

	total := int64(0)
	current := int64(0)

	if captures != nil {
		currentRaw := captures[1]
		totalRaw := captures[2]

		current, _ = strconv.ParseInt(currentRaw, 10, 32)
		total, _ = strconv.ParseInt(totalRaw, 10, 32)
	}

	var episodes []Episode

	doc.Find(".episode-list-data").Each(func(i int, s *goquery.Selection) {
		// fmt.Printf("s.Text(): %v\n", s.Text())

		titles := s.Find(".episode-title")

		english := titles.Find(".fl-l").Text()
		english = strings.TrimSpace(english)

		japanese := titles.Find(".di-ib").Text()
		japanese = strings.TrimSpace(japanese)

		episodeNumberRaw := s.Find(".episode-number").Text()
		number, _ := strconv.ParseInt(episodeNumberRaw, 10, 32)

		aired := s.Find(".episode-aired").Text()
		aired = strings.TrimSpace(aired)

		t, _ := ParseDate(aired)
		aired = formatDate(t)

		averageScoreRaw := s.Find(".average > .value").Text()
		averageScore, _ := strconv.ParseFloat(averageScoreRaw, 64)

		episodes = append(episodes, Episode{
			EnglishTitle:    english,
			JapaneseTitle:   japanese,
			Number:          number,
			Aired:           aired,
			AverageScoreRaw: averageScoreRaw,
			AverageScore:    averageScore,
		})
	})

	return episodes, EpisodeExtraInfo{
		Total:   total,
		Current: current,
	}, nil
}

func ExtractPictures(pagePath string) ([]string, error) {
	f, err := os.Open(pagePath)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return nil, err
	}

	var images []string

	doc.Find(".js-picture-gallery img").Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("data-src")
		val = strings.TrimSpace(val)

		if exists && val != "" {
			images = append(images, val)
		}
	})

	return images, nil
}

var rxClassTrim = regexp.MustCompile("[\t\r\n]")

func getClassesSlice(classes string) []string {
	return strings.Split(rxClassTrim.ReplaceAllString(" "+classes+" ", " "), " ")
}

func ExtractSeasonalAnimes(pagePath string) (Seasonal, error) {
	f, err := os.Open(pagePath)
	if err != nil {
		return Seasonal{}, err
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return Seasonal{}, err
	}

	var animes []SeasonalAnime

	doc.Find(".seasonal-anime").Each(func(i int, s *goquery.Selection) {
		titleEl := s.Find(".title")

		linkTitle := titleEl.Find(".link-title")
		link, _ := linkTitle.Attr("href")
		title := strings.TrimSpace(linkTitle.Text())

		englishTitle := titleEl.Find(".h3_anime_subtitle").Text()
		englishTitle = strings.TrimSpace(englishTitle)

		class, _ := s.Attr("class")

		cs := getClassesSlice(class)

		animeType := "Unknown"
		rating := ""
		status := ""

		isForKids := false

		for _, clx := range cs {
			if clx == "r18" {
				rating = "Rx - Hentai"
			}

			if clx == "kids" {
				isForKids = true
			}

			if strings.Contains(clx, "js-anime-type-") {
				ty := strings.TrimPrefix(clx, "js-anime-type-")
				switch ty {
				case "all":
				case "1":
					animeType = "TV"
				case "2":
					animeType = "OVA"
				case "3":
					animeType = "Movie"
				case "4":
					animeType = "Special"
				case "5":
					animeType = "ONA"
				case "9":
					animeType = "TV Special"
				default:
					slog.Warn("Unknown anime type", "type", ty, "title", title)
				}
			}
		}

		u, _ := url.Parse(link)
		splits := strings.Split(u.Path, "/")
		id := splits[2]

		infoEl := s.Find(".info")
		startDateStr := infoEl.Children().Eq(0).Text()
		startDateStr = strings.TrimSpace(startDateStr)
		startDate, _ := ParseDate(startDateStr)

		img := s.Find(".image img")
		coverImageUrl, _ := img.Attr("data-src")
		if coverImageUrl == "" {
			coverImageUrl, _ = img.Attr("src")
		}

		scoreStr := s.Find(".js-score").Text()
		score, _ := strconv.ParseFloat(scoreStr, 64)

		epsStr := infoEl.Children().Eq(1).Children().Eq(0).Text()
		epsStr = strings.TrimSpace(epsStr)
		eps := int64(utils.ExtractNumber(epsStr))

		var endDate *time.Time
		if eps != 0 {
			t := startDate.AddDate(0, 0, 7*(int(eps)-1))
			sub := time.Now().Sub(t)

			if sub < 0 {
				status = "Currently Airing"
			} else {
				status = "Finished Airing"
			}

			endDate = &t
		}

		var genres []string

		genresEl := s.Find(".genres > .genres-inner > .genre")
		genresEl.Each(func(i int, s *goquery.Selection) {
			genre := s.Find("a").Text()
			genres = append(genres, genre)
		})

		leftside := s.Find("div .synopsis")

		preline := leftside.Find(".preline")
		description := strings.TrimSpace(preline.Text())

		var studios []string
		leftside.Find(".caption:contains(\"Studio\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
			raw := s.Text()
			raw = strings.TrimSpace(raw)

			if raw != "Unknown" {
				studios = append(studios, s.Text())
			}
		})

		leftside.Find(".caption:contains(\"Studios\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
			raw := s.Text()
			raw = strings.TrimSpace(raw)

			if raw != "Unknown" {
				studios = append(studios, s.Text())
			}
		})

		var themes []string
		leftside.Find("span:contains(\"Theme\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
			t := s.Text()
			t = strings.TrimSpace(t)

			if t != "" {
				themes = append(themes, s.Text())
			}
		})

		leftside.Find("span:contains(\"Themes\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
			t := s.Text()
			t = strings.TrimSpace(t)

			if t != "" {
				themes = append(themes, s.Text())
			}
		})

		var demographic []string
		leftside.Find("span:contains(\"Demographic\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
			t := s.Text()
			t = strings.TrimSpace(t)

			if t != "" {
				demographic = append(demographic, s.Text())
			}
		})

		leftside.Find("span:contains(\"Demographics\")").Parent().Find("a").Each(func(i int, s *goquery.Selection) {
			t := s.Text()
			t = strings.TrimSpace(t)

			if t != "" {
				demographic = append(demographic, s.Text())
			}
		})

		if isForKids {
			return
		}

		endDateS := ""
		if endDate != nil {
			endDateS = formatDate(*endDate)
		}

		animes = append(animes, SeasonalAnime{
			Id:            id,
			Title:         title,
			TitleEnglish:  englishTitle,
			Description:   description,
			CoverImageUrl: coverImageUrl,
			Type:          animeType,
			Status:        status,
			EpisodeCount:  eps,
			Rating:        rating,
			StartDate:     formatDate(startDate),
			EndDate:       endDateS,
			Score:         score,
			Studios:       studios,
			Genres:        genres,
			Themes:        themes,
			Demographics:  demographic,
		})
	})

	return Seasonal{
		Animes: animes,
	}, nil
}
