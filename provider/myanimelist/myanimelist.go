package myanimelist

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/kr/pretty"
	"github.com/nanoteck137/watchbook/downloader"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
	"golang.org/x/time/rate"
)

var dl = downloader.NewDownloader(
	rate.NewLimiter(rate.Every(4*time.Second), 10),
	UserAgent,
)

type AnimeEntry struct {
	Type types.MediaType `json:"type"`

	Title        string `json:"title"`
	TitleEnglish string `json:"titleEnglish"`

	Description string `json:"description"`

	Score        *float64          `json:"score"`
	Status       types.MediaStatus `json:"status"`
	Rating       types.MediaRating `json:"rating"`
	AiringSeason string            `json:"airingSeason"`

	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`

	Studios []string `json:"studios"`
	Tags    []string `json:"tags"`

	CoverImageUrl string `json:"coverImageUrl"`

	EpisodeCount *int64 `json:"episodeCount"`

	UsingCache bool `json:"-"`
}

func readAnimeEntry(p string) (AnimeEntry, error) {
	d, err := os.ReadFile(p)
	if err != nil {
		return AnimeEntry{}, err
	}

	var res AnimeEntry
	err = json.Unmarshal(d, &res)
	if err != nil {
		return AnimeEntry{}, err
	}

	return res, nil
}

func fetchAndCacheNewAnimeData(malId string, cacheDest string) (AnimeEntry, error) {
	data, err := FetchAnimeData(dl, malId, false)
	if err != nil {
		return AnimeEntry{}, err
	}

	pretty.Println(data)

	desc := strings.Builder{}

	if data.Description != "" {
		desc.WriteString(data.Description)
		desc.WriteString("\n\n")
	}

	fmt.Fprintf(&desc, "Type: %s\n", data.Type)
	fmt.Fprintf(&desc, "Status: %s\n", data.Status)
	if data.EpisodeCount != nil {
		fmt.Fprintf(&desc, "Episode Count: %d\n", *data.EpisodeCount)
	}
	fmt.Fprintf(&desc, "Rating: %s\n", data.Rating)
	fmt.Fprintf(&desc, "Premiered: %s\n", data.Premiered)
	fmt.Fprintf(&desc, "Source: %s\n", data.Source)
	fmt.Fprintf(&desc, "Broadcast: %s\n", data.Broadcast)

	if len(data.ThemeSongs) > 0 {
		desc.WriteString("Theme Songs:\n")
		for _, t := range data.ThemeSongs {
			ty := "UNKNOWN"

			switch t.Type {
			case ThemeSongOpening:
				ty = "OP"
			case ThemeSongEnding:
				ty = "ED"
			}

			desc.WriteString(t.Raw)
			desc.WriteString(" (")
			desc.WriteString(ty)
			desc.WriteString(")\n")
		}
	}

	var startDate *string
	var endDate *string

	if data.StartDate != nil && *data.StartDate != "" {
		startDate = data.StartDate
	}

	if data.EndDate != nil && *data.EndDate != "" {
		endDate = data.EndDate
	}

	studios := make([]string, 0, len(data.Studios))
	tags := make([]string, 0, len(data.Genres)+len(data.Themes)+len(data.Demographics))

	for _, s := range data.Studios {
		studios = append(studios, utils.Slug(s))
	}

	for _, t := range data.Genres {
		tags = append(tags, utils.Slug(t))
	}

	for _, t := range data.Themes {
		tags = append(tags, utils.Slug(t))
	}

	for _, t := range data.Demographics {
		tags = append(tags, utils.Slug(t))
	}

	res := AnimeEntry{
		Type:          ConvertAnimeType(data.Type),
		Title:         data.Title,
		TitleEnglish:  data.TitleEnglish,
		Description:   strings.TrimSpace(desc.String()),
		Score:         data.Score,
		Status:        ConvertAnimeStatus(data.Status),
		Rating:        ConvertAnimeRating(data.Rating),
		AiringSeason:  utils.Slug(data.Premiered),
		StartDate:     startDate,
		EndDate:       endDate,
		Studios:       studios,
		Tags:          tags,
		CoverImageUrl: data.CoverImageUrl,
		EpisodeCount:  data.EpisodeCount,
	}

	d, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return AnimeEntry{}, err
	}

	err = os.WriteFile(cacheDest, d, 0644)
	if err != nil {
		return AnimeEntry{}, err
	}

	return res, nil
}

func fetchAnimeData(malId string) (AnimeEntry, error) {
	data, err := FetchAnimeData(dl, malId, false)
	if err != nil {
		return AnimeEntry{}, err
	}

	desc := strings.Builder{}

	if data.Description != "" {
		desc.WriteString(data.Description)
		desc.WriteString("\n\n")
	}

	fmt.Fprintf(&desc, "Type: %s\n", data.Type)
	fmt.Fprintf(&desc, "Status: %s\n", data.Status)
	if data.EpisodeCount != nil {
		fmt.Fprintf(&desc, "Episode Count: %d\n", *data.EpisodeCount)
	}
	fmt.Fprintf(&desc, "Rating: %s\n", data.Rating)
	fmt.Fprintf(&desc, "Premiered: %s\n", data.Premiered)
	fmt.Fprintf(&desc, "Source: %s\n", data.Source)
	fmt.Fprintf(&desc, "Broadcast: %s\n", data.Broadcast)

	if len(data.ThemeSongs) > 0 {
		desc.WriteString("Theme Songs:\n")
		for _, t := range data.ThemeSongs {
			ty := "UNKNOWN"

			switch t.Type {
			case ThemeSongOpening:
				ty = "OP"
			case ThemeSongEnding:
				ty = "ED"
			}

			desc.WriteString(t.Raw)
			desc.WriteString(" (")
			desc.WriteString(ty)
			desc.WriteString(")\n")
		}
	}

	var startDate *string
	var endDate *string

	if data.StartDate != nil && *data.StartDate != "" {
		startDate = data.StartDate
	}

	if data.EndDate != nil && *data.EndDate != "" {
		endDate = data.EndDate
	}

	studios := make([]string, 0, len(data.Studios))
	tags := make([]string, 0, len(data.Genres)+len(data.Themes)+len(data.Demographics))

	for _, s := range data.Studios {
		studios = append(studios, utils.Slug(s))
	}

	for _, t := range data.Genres {
		tags = append(tags, utils.Slug(t))
	}

	for _, t := range data.Themes {
		tags = append(tags, utils.Slug(t))
	}

	for _, t := range data.Demographics {
		tags = append(tags, utils.Slug(t))
	}

	res := AnimeEntry{
		Type:          ConvertAnimeType(data.Type),
		Title:         data.Title,
		TitleEnglish:  data.TitleEnglish,
		Description:   strings.TrimSpace(desc.String()),
		Score:         data.Score,
		Status:        ConvertAnimeStatus(data.Status),
		Rating:        ConvertAnimeRating(data.Rating),
		AiringSeason:  utils.Slug(data.Premiered),
		StartDate:     startDate,
		EndDate:       endDate,
		Studios:       studios,
		Tags:          tags,
		CoverImageUrl: data.CoverImageUrl,
		EpisodeCount:  data.EpisodeCount,
	}

	return res, nil
}

func RawGetAnime(malId string) (AnimeEntry, error) {
	res, err := fetchAnimeData(malId)
	if err != nil {
		return AnimeEntry{}, err
	}

	return res, nil
}

// TODO(patrik): Remove
func GetAnime(workDir types.WorkDir, malId string, useCache bool) (AnimeEntry, error) {
	panic("REMOVE")
	cacheDir := "" 

	err := os.Mkdir(cacheDir, 0755)
	if err != nil && !os.IsExist(err) {
		return AnimeEntry{}, err
	}

	cacheName := fmt.Sprintf("mal-%s-anime", malId)
	fullCacheName := cacheName + ".json"
	cachePath := path.Join(cacheDir, fullCacheName)

	if !useCache {
		res, err := fetchAndCacheNewAnimeData(malId, cachePath)
		if err != nil {
			return AnimeEntry{}, err
		}

		return res, nil
	}

	cacheEntry, err := readAnimeEntry(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			res, err := fetchAndCacheNewAnimeData(malId, cachePath)
			if err != nil {
				return AnimeEntry{}, err
			}

			return res, nil
		}

		return AnimeEntry{}, err
	}

	cacheEntry.UsingCache = true
	return cacheEntry, nil
}

type AnimeEpisodes struct{}

type AnimePictures struct{}

type MangaEntry struct{}

func ConvertAnimeType(typ string) types.MediaType {
	switch typ {
	case "TV":
		return types.MediaTypeAnimeSeason
	case "OVA":
		return types.MediaTypeAnimeSeason
	case "Movie":
		return types.MediaTypeAnimeMovie
	case "Special":
		return types.MediaTypeAnimeSeason
	case "ONA":
		return types.MediaTypeAnimeSeason
	case "Music":
		return types.MediaTypeUnknown
	case "CM":
		return types.MediaTypeAnimeSeason
	case "PV":
		return types.MediaTypeAnimeSeason
	case "TV Special":
		return types.MediaTypeAnimeSeason
	case "":
		return types.MediaTypeUnknown
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown anime type \"%s\"\n", typ)
	}

	return types.MediaTypeUnknown
}

func ConvertAnimeStatus(status string) types.MediaStatus {
	switch status {
	case "Currently Airing":
		return types.MediaStatusAiring
	case "Finished Airing":
		return types.MediaStatusFinished
	case "Not yet aired":
		return types.MediaStatusNotAired
	case "":
		return types.MediaStatusUnknown
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown anime status \"%s\"\n", status)
	}

	return types.MediaStatusUnknown
}

func ConvertAnimeRating(rating string) types.MediaRating {
	switch rating {
	case "G - All Ages":
		return types.MediaRatingAllAges
	case "PG - Children":
		return types.MediaRatingPG
	case "PG-13 - Teens 13 or older":
		return types.MediaRatingPG13
	case "R - 17+ (violence & profanity)":
		return types.MediaRatingR17
	case "R+ - Mild Nudity":
		return types.MediaRatingRMildNudity
	case "Rx - Hentai":
		return types.MediaRatingRHentai
	case "":
		return types.MediaRatingUnknown
	default:
		// TODO(patrik): Better logging
		fmt.Printf("WARN: Unknown anime rating \"%s\"\n", rating)
	}

	return types.MediaRatingUnknown
}
