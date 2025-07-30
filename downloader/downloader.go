package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/time/rate"
)

type Downloader struct {
	client    *RLHTTPClient
	userAgent string
}

var NotFound = errors.New("page not found")

func NewDownloader(rateLimiter *rate.Limiter, userAgent string) *Downloader {
	return &Downloader{
		client: &RLHTTPClient{
			Client: &http.Client{
				Transport: &http.Transport{},
			},
			Ratelimiter: rateLimiter,
		},
		userAgent: userAgent,
	}
}

func (scraper *Downloader) DownloadToFile(url, dest string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if scraper.userAgent != "" {
		req.Header.Set("User-Agent", scraper.userAgent)
	}

	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := scraper.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return NotFound
		default:
			// TODO(patrik): Better error
			return fmt.Errorf("not successful %d", resp.StatusCode)
		}
	}

	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (scraper *Downloader) DownloadJson(url string, dest any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if scraper.userAgent != "" {
		req.Header.Set("User-Agent", scraper.userAgent)
	}

	resp, err := scraper.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return NotFound
		default:
			// TODO(patrik): Better error
			return fmt.Errorf("not successful %d", resp.StatusCode)
		}
	}

	// TODO(patrik): Check Content-Type

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&dest)
	if err != nil {
		return err
	}

	return nil
}
