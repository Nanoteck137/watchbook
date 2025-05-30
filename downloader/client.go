package downloader

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"
)

type RLHTTPClient struct {
	Client      *http.Client
	Ratelimiter *rate.Limiter
}

func (c *RLHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()

	err := c.Ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
