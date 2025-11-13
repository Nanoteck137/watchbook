package myanimelist

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nanoteck137/watchbook/provider"
	"github.com/nanoteck137/watchbook/provider/downloader"
	"github.com/nanoteck137/watchbook/types"
	"github.com/nanoteck137/watchbook/utils"
	"golang.org/x/time/rate"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:136.0) Gecko/20100101 Firefox/136.0"

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

	StartDate *string    `json:"startDate"`
	EndDate   *string    `json:"endDate"`
	Release   *time.Time `json:"release"`

	Studios []string `json:"studios"`
	Tags    []string `json:"tags"`

	CoverImageUrl string `json:"coverImageUrl"`

	EpisodeCount *int64 `json:"episodeCount"`
}

func parseDateTimeUTC(dateStr, schedule string) (time.Time, error) {
	// Parse base date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date: %w", err)
	}

	// Split schedule, e.g. "Saturdays at 23:00 (JST)"
	parts := strings.Split(schedule, " at ")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid schedule format")
	}

	// Extract weekday
	weekdayStr := strings.TrimSuffix(parts[0], "s") // "Saturdays" → "Saturday"
	weekdayMap := map[string]time.Weekday{
		"Sunday":    time.Sunday,
		"Monday":    time.Monday,
		"Tuesday":   time.Tuesday,
		"Wednesday": time.Wednesday,
		"Thursday":  time.Thursday,
		"Friday":    time.Friday,
		"Saturday":  time.Saturday,
	}
	weekday, ok := weekdayMap[weekdayStr]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid weekday: %s", weekdayStr)
	}

	// Extract time and tz, e.g. "23:00 (JST)"
	timeAndTZ := parts[1]
	timeParts := strings.SplitN(timeAndTZ, " ", 2)
	if len(timeParts) < 2 {
		return time.Time{}, fmt.Errorf("invalid time format")
	}

	hm := timeParts[0] // "23:00"
	hmSplit := strings.Split(hm, ":")
	if len(hmSplit) != 2 {
		return time.Time{}, fmt.Errorf("invalid hour:minute")
	}
	hour, _ := strconv.Atoi(hmSplit[0])
	min, _ := strconv.Atoi(hmSplit[1])

	// Timezone
	tz := strings.Trim(timeParts[1], "()")
	var loc *time.Location
	switch tz {
	case "JST":
		loc = time.FixedZone("JST", 9*60*60)
	default:
		return time.Time{}, fmt.Errorf("unsupported timezone: %s", tz)
	}

	// Adjust to the correct weekday
	for date.Weekday() != weekday {
		date = date.AddDate(0, 0, 1)
	}

	// Build full datetime in given TZ
	localTime := time.Date(date.Year(), date.Month(), date.Day(), hour, min, 0, 0, loc)

	// Convert to UTC
	return localTime.UTC(), nil
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

	var release *time.Time
	if startDate != nil && data.Broadcast != "" {
		t, err := parseDateTimeUTC(*startDate, data.Broadcast)
		if err == nil {
			release = &t
		}
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
		Release:       release,
	}

	return res, nil
}

type AnimeEpisode struct{}

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
		return types.MediaStatusOngoing
	case "Finished Airing":
		return types.MediaStatusCompleted
	case "Not yet aired":
		return types.MediaStatusUpcoming
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

const AnimeProviderName string = "myanimelist-anime"

var _ (provider.Provider) = (*MyAnimeListAnimeProvider)(nil)

type MyAnimeListAnimeProvider struct {
}

func (m *MyAnimeListAnimeProvider) Info() provider.Info {
	return provider.Info{
		Name:                    AnimeProviderName,
		DisplayName:             "MyAnimeList Anime",
		SupportGetMedia:         true,
		SupportSearchMedia:      true,
		SupportGetCollection:    false,
		SupportSearchCollection: false,
	}
}

func (m *MyAnimeListAnimeProvider) GetCollection(c provider.Context, id string) (provider.Collection, error) {
	panic("unsupported")
}

func (m *MyAnimeListAnimeProvider) GetMedia(c provider.Context, id string) (provider.Media, error) {
	anime, err := fetchAnimeData(id)
	if err != nil {
		return provider.Media{}, err
	}

	title := anime.Title
	if anime.TitleEnglish != "" {
		title = anime.TitleEnglish
	}

	var description *string
	if anime.Description != "" {
		description = &anime.Description
	}

	var airingSeason *string
	if anime.AiringSeason != "" {
		airingSeason = &anime.AiringSeason
	}

	var startDate *time.Time
	if anime.StartDate != nil {
		d, err := time.Parse(types.MediaDateLayout, *anime.StartDate)
		if err == nil {
			startDate = &d
		}
	}

	var endDate *time.Time
	if anime.EndDate != nil {
		d, err := time.Parse(types.MediaDateLayout, *anime.EndDate)
		if err == nil {
			endDate = &d
		}
	}

	var coverUrl *string
	if anime.CoverImageUrl != "" {
		coverUrl = &anime.CoverImageUrl
	}

	episodeCount := 0
	if anime.EpisodeCount != nil {
		episodeCount = int(*anime.EpisodeCount)
	}

	parts := []provider.MediaPart{}

	episodes, _ := FetchAnimeEpisodes(dl, id)

	if episodeCount > 0 && len(episodes) > episodeCount {
		episodes = episodes[:episodeCount]
	}

	numEpisodesFound := len(episodes)
	missingEpisodes := max(episodeCount-numEpisodesFound, 0)

	lastEpisodeNumber := 0
	for _, episode := range episodes {
		var releaseDate *time.Time
		d, err := time.Parse(types.MediaDateLayout, episode.Aired)
		if err == nil {
			releaseDate = &d
		}

		n := int(episode.Number)
		parts = append(parts, provider.MediaPart{
			Name:        episode.EnglishTitle,
			Number:      n,
			ReleaseDate: releaseDate,
		})

		if n > lastEpisodeNumber {
			lastEpisodeNumber = n
		}
	}

	for i := range missingEpisodes {
		n := lastEpisodeNumber + 1 + i
		parts = append(parts, provider.MediaPart{
			Name:   fmt.Sprintf("Episode %d", n),
			Number: n,
		})
	}

	return provider.Media{
		ProviderId:       id,
		Type:             anime.Type,
		Title:            title,
		Description:      description,
		Score:            anime.Score,
		Status:           anime.Status,
		Rating:           anime.Rating,
		AiringSeason:     airingSeason,
		StartDate:        startDate,
		EndDate:          endDate,
		Release:          anime.Release,
		CoverUrl:         coverUrl,
		LogoUrl:          nil,
		BannerUrl:        nil,
		Creators:         anime.Studios,
		Tags:             anime.Tags,
		Parts:            parts,
		ExtraProviderIds: map[string]string{},
	}, nil
}

func (m *MyAnimeListAnimeProvider) SearchCollection(c provider.Context, query string) ([]provider.SearchResult, error) {
	panic("unsupported")
}

func (m *MyAnimeListAnimeProvider) SearchMedia(c provider.Context, query string) ([]provider.SearchResult, error) {
	items, err := FetchSearch(query)
	if err != nil {
		return nil, err
	}

	res := make([]provider.SearchResult, len(items))

	for i, item := range items {
		res[i] = provider.SearchResult{
			SearchType: provider.SearchResultTypeMedia,
			ProviderId: item.Id,
			Title:      item.Title,
			MediaType:  item.Type,
			ImageUrl:   item.ImageUrl,
		}
	}

	return res, nil
}
