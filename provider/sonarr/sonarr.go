package sonarr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Sonarr API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new Sonarr API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Series represents a TV series in Sonarr
type Series struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	TvdbID        int       `json:"tvdbId"`
	ImdbID        string    `json:"imdbId"`
	Overview      string    `json:"overview"`
	Year          int       `json:"year"`
	Status        string    `json:"status"`
	Monitored     bool      `json:"monitored"`
	Path          string    `json:"path"`
	Seasons       []Season  `json:"seasons"`
	Images        []Image   `json:"images"`
	Added         time.Time `json:"added"`
	FirstAired    time.Time `json:"firstAired"`
	Runtime       int       `json:"runtime"`
	Network       string    `json:"network"`
	Genres        []string  `json:"genres"`
	Certification string    `json:"certification"`
}

// Season represents a season of a TV series
type Season struct {
	SeasonNumber int             `json:"seasonNumber"`
	Monitored    bool            `json:"monitored"`
	Statistics   SeasonStatistics `json:"statistics"`
}

// SeasonStatistics contains statistics about a season
type SeasonStatistics struct {
	EpisodeCount      int     `json:"episodeCount"`
	EpisodeFileCount  int     `json:"episodeFileCount"`
	TotalEpisodeCount int     `json:"totalEpisodeCount"`
	SizeOnDisk        int64   `json:"sizeOnDisk"`
	PercentOfEpisodes float64 `json:"percentOfEpisodes"`
}

// Image represents an image (poster, banner, etc.)
type Image struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
	RemoteURL string `json:"remoteUrl"`
}

// Episode represents an episode in Sonarr
type Episode struct {
	ID                int       `json:"id"`
	SeriesID          int       `json:"seriesId"`
	EpisodeFileID     int       `json:"episodeFileId"`
	SeasonNumber      int       `json:"seasonNumber"`
	EpisodeNumber     int       `json:"episodeNumber"`
	Title             string    `json:"title"`
	AirDate           string    `json:"airDate"`
	AirDateUTC        time.Time `json:"airDateUtc"`
	Overview          string    `json:"overview"`
	HasFile           bool      `json:"hasFile"`
	Monitored         bool      `json:"monitored"`
	AbsoluteEpisodeNumber int   `json:"absoluteEpisodeNumber"`
}

// makeRequest makes an HTTP request to the Sonarr API
func (c *Client) makeRequest(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v3/%s", c.BaseURL, endpoint)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("Accept", "application/json")
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	return body, nil
}

// GetAllSeries fetches all series from Sonarr
func (c *Client) GetAllSeries() ([]Series, error) {
	body, err := c.makeRequest("series")
	if err != nil {
		return nil, err
	}
	
	var series []Series
	if err := json.Unmarshal(body, &series); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return series, nil
}

// GetSeriesByID fetches a specific series by its Sonarr ID
func (c *Client) GetSeriesByID(id int) (*Series, error) {
	body, err := c.makeRequest(fmt.Sprintf("series/%d", id))
	if err != nil {
		return nil, err
	}
	
	var series Series
	if err := json.Unmarshal(body, &series); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &series, nil
}

// GetEpisodesBySeries fetches all episodes for a specific series
func (c *Client) GetEpisodesBySeries(seriesID int) ([]Episode, error) {
	body, err := c.makeRequest(fmt.Sprintf("episode?seriesId=%d", seriesID))
	if err != nil {
		return nil, err
	}
	
	var episodes []Episode
	if err := json.Unmarshal(body, &episodes); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return episodes, nil
}

// TestConnection tests the connection to Sonarr
func (c *Client) TestConnection() error {
	body, err := c.makeRequest("system/status")
	if err != nil {
		return err
	}
	
	var status map[string]any
	if err := json.Unmarshal(body, &status); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	
	return nil
}
