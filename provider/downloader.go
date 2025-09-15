package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/nanoteck137/watchbook/utils"
	"golang.org/x/time/rate"
)

var NotFound = errors.New("page not found")

type HTTPClient struct {
	BaseUrl      string

	client    *http.Client
	limiter   *rate.Limiter
	userAgent string
	timeout   time.Duration
	cacheTtl  time.Duration
}

type ClientOption func(*HTTPClient)

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *HTTPClient) {
		c.timeout = timeout
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *HTTPClient) {
		c.userAgent = userAgent
	}
}

func WithCacheTTL(ttl time.Duration) ClientOption {
	return func(c *HTTPClient) {
		c.cacheTtl = ttl
	}
}

func WithRate(rps float64, burst int) ClientOption {
	return func(c *HTTPClient) {
		c.limiter = rate.NewLimiter(rate.Limit(rps), burst)
	}
}

func NewHttpClient(baseUrl string, opts ...ClientOption) *HTTPClient {
	downloader := &HTTPClient{
		BaseUrl:      baseUrl,
		client:       &http.Client{Timeout: 30 * time.Second},
		limiter:      rate.NewLimiter(1, 1),
		userAgent:    "",
		timeout:      30 * time.Second,
		cacheTtl:     1 * time.Hour,
	}

	for _, opt := range opts {
		opt(downloader)
	}

	return downloader
}

type RequestOptions struct {
	Headers http.Header
	Query   url.Values
}

func (c *HTTPClient) Get(ctx context.Context, key, path string, opts RequestOptions) ([]byte, error) {
	url, err := utils.CreateUrlBase(c.BaseUrl, path, opts.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to create url: %w", err)
	}

	err = c.limiter.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("rate limiter wait failed: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	u := url.String()
	req, err := http.NewRequestWithContext(reqCtx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	headers := opts.Headers.Clone()

	if c.userAgent != "" {
		headers.Set("User-Agent", c.userAgent)
	}

	req.Header = headers

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response data: %w", err)
	}

	return data, nil
}

// func (d *Downloader) DownloadToFile(url, dest string) error {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	if scraper.userAgent != "" {
// 		req.Header.Set("User-Agent", scraper.userAgent)
// 	}
//
// 	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
//
// 	resp, err := scraper.client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
//
// 	if resp.StatusCode != http.StatusOK {
// 		switch resp.StatusCode {
// 		case http.StatusNotFound:
// 			return NotFound
// 		default:
// 			// TODO(patrik): Better error
// 			return fmt.Errorf("not successful %d", resp.StatusCode)
// 		}
// 	}
//
// 	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
//
// 	_, err = io.Copy(f, resp.Body)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
//
// func (scraper *Downloader) DownloadJson(url string, dest any) error {
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	if scraper.userAgent != "" {
// 		req.Header.Set("User-Agent", scraper.userAgent)
// 	}
//
// 	resp, err := scraper.client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
//
// 	if resp.StatusCode != http.StatusOK {
// 		switch resp.StatusCode {
// 		case http.StatusNotFound:
// 			return NotFound
// 		default:
// 			// TODO(patrik): Better error
// 			return fmt.Errorf("not successful %d", resp.StatusCode)
// 		}
// 	}
//
// 	// TODO(patrik): Check Content-Type
//
// 	decoder := json.NewDecoder(resp.Body)
//
// 	err = decoder.Decode(&dest)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
