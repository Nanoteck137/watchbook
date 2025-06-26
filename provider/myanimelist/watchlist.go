package myanimelist

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nanoteck137/watchbook/downloader"
)

type WatchlistStatus int

const (
	WatchlistStatusCurrentlyWatching WatchlistStatus = 1
	WatchlistStatusCompleted         WatchlistStatus = 2
	WatchlistStatusOnHold            WatchlistStatus = 3
	WatchlistStatusDropped           WatchlistStatus = 4
	WatchlistStatusPlanToWatch       WatchlistStatus = 6
)

type WatchlistGenre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Title string

func (t *Title) UnmarshalJSON(data []byte) error {

	var n any
	err := json.Unmarshal(data, &n)
	if err != nil {
		return err
	}

	switch n := n.(type) {
	case string:
		*t = Title(n)
	case int:
		*t = Title(strconv.FormatInt(int64(n), 10))
	case float64:
		*t = Title(strconv.FormatInt(int64(n), 10))
	}

	return nil
}

type WatchlistEntry struct {
	Status             WatchlistStatus `json:"status"`
	Score              int             `json:"score"`
	AnimeId            int             `json:"anime_id"`
	Tags               string          `json:"tags"`
	IsRewatching       int             `json:"is_rewatching"`
	NumWatchedEpisodes int             `json:"num_watched_episodes"`
	CreatedAt          int             `json:"created_at"`
	UpdatedAt          int             `json:"updated_at"`
	Storage            string          `json:"storage_string"`
	Priority           string          `json:"priority_string"`
	Notes              string          `json:"notes"`
	EditableNotes      string          `json:"editable_notes"`

	AnimeTitle         Title           `json:"anime_title"`
	AnimeTitleEnglish  Title           `json:"anime_title_eng"`

	// AnimeStudios       string           `json:"anime_studios"`        //null,
	// AnimeLicensors     string           `json:"anime_licensors"`      //null,
	// AnimeSeason        string           `json:"anime_season"`         //null,
	// Demographics       string           `json:"demographics"`             //[],
	// TitleLocalized     string           `json:"title_localized"`          //null,
	// Days               string           `json:"days_string"`              //null,
	// StartDate          string           `json:"start_date_string"`        //null,
	// FinishedDate       string           `json:"finish_date_string"`       //null,
	// AnimeTitle         string           `json:"anime_title"`
	// AnimeTitleEnglish  string           `json:"anime_title_eng"`
	// AnimeNumEpisodes   int              `json:"anime_num_episodes"`
	// AnimeAiringStatus  int              `json:"anime_airing_status"`
	// AnimeTotalMembers  int              `json:"anime_total_members"`
	// AnimeTotalScore    int              `json:"anime_total_scores"`
	// AnimeScoreVal      float32          `json:"anime_score_val"`
	// AnimeScoreDiff     float32          `json:"anime_score_diff"`
	// AnimePopularity    int              `json:"anime_popularity"`
	// HasEpisodeVideo    bool             `json:"has_episode_video"`
	// HasPromotionVideo  bool             `json:"has_promotion_video"`
	// HasVideo           bool             `json:"has_video"`
	// VideoUrl           string           `json:"video_url"`
	// Genres             []WatchlistGenre `json:"genres"`
	// AneimUrl           string           `json:"anime_url"`
	// AnimeImagePath     string           `json:"anime_image_path"`
	// IsAddedToList      bool             `json:"is_added_to_list"`
	// AnimeMediaType     string           `json:"anime_media_type_string"`
	// AnimeMpaaRating    string           `json:"anime_mpaa_rating_string"`
	// AnimeStartDate     string           `json:"anime_start_date_string"`
	// AnimeEndDate       string           `json:"anime_end_date_string"`
}

const watchlistPerPage = 300

func GetUserWatchlistPage(dl *downloader.Downloader, page int, username string) ([]WatchlistEntry, error) {
	url := fmt.Sprintf("https://myanimelist.net/animelist/%s/load.json?offset=%d", username, page*watchlistPerPage)

	var dest []WatchlistEntry
	err := dl.DownloadJson(url, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func GetUserWatchlist(dl *downloader.Downloader, username string) ([]WatchlistEntry, error) {
	total, err := GetUserWatchlistPage(dl, 0, username)
	if err != nil {
		return nil, err
	}

	done := len(total) < watchlistPerPage
	page := 1

	for !done {
		next, err := GetUserWatchlistPage(dl, page, username)
		if err != nil {
			return nil, err
		}

		total = append(total, next...)

		done = len(next) < watchlistPerPage
		page++
	}

	return total, nil
}
