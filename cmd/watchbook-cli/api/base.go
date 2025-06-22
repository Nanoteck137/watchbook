package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	addr      string
	authToken string
	apiToken  string
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) SetAuthToken(token string) {
	c.authToken = token
}

func (c *Client) SetApiToken(token string) {
	c.apiToken = token
}

type Options struct {
	QueryParams map[string]string
	Boundary    string
}

func createUrl(addr, path string, query map[string]string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}

	u.Path = path

	params := u.Query()
	for k, v := range query {
		params.Set(k, v)
	}
	u.RawQuery = params.Encode()

	return u.String(), nil
}

type ApiError[E any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Extra   E      `json:"extra,omitempty"`
}

func (err *ApiError[E]) Error() string {
	return err.Message
}

type ApiResponse[D any, E any] struct {
	Success bool         `json:"success"`
	Data    D            `json:"data,omitempty"`
	Error   *ApiError[E] `json:"error,omitempty"`
}

type RequestData struct {
	Url    string
	Method string

	AuthToken string
	ApiToken  string
	Body      any
}

func rawRequest(data *RequestData, contentType string, bodyReader io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(data.Method, data.Url, bodyReader)
	if err != nil {
		return nil, err
	}

	if data.AuthToken != "" {
		req.Header.Add("Authorization", "Bearer "+data.AuthToken)
	}

	if data.ApiToken != "" {
		req.Header.Add("X-Api-Token", data.ApiToken)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Request[D any](data RequestData) (*D, error) {
	var bodyReader io.Reader

	if data.Body != nil {
		buf := bytes.Buffer{}

		err := json.NewEncoder(&buf).Encode(data.Body)
		if err != nil {
			return nil, err
		}

		bodyReader = &buf
	}

	resp, err := rawRequest(&data, "application/json", bodyReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res ApiResponse[D, any]
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, res.Error
	}

	return &res.Data, nil
}

// NOTE(patrik): Copied from multipart.Writer.FormDataContentType
func createFormContentType(b string) string {
	if strings.ContainsAny(b, `()<>@,;:\"/[]?= `) {
		b = `"` + b + `"`
	}
	return "multipart/form-data; boundary=" + b
}

func RequestForm[D any](data RequestData, boundary string, body Reader) (*D, error) {
	ct := createFormContentType(boundary)
	resp, err := rawRequest(&data, ct, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res ApiResponse[D, any]
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, res.Error
	}

	return &res.Data, nil
}

// Simple wrapper for Sprintf
func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

// Copy of io.Reader interface
type Reader interface {
	Read(p []byte) (n int, err error)
}
