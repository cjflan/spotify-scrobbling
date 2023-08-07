package scrobbling

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type Client struct {
	http      *http.Client
	baseURL   string
	autoRetry bool
	backoff   int
}

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e Error) Error() string {
	return e.Message
}

type ClientOption func(c *Client)

func WithRetry(shouldRetry bool) ClientOption {
	return func(c *Client) {
		c.autoRetry = shouldRetry
	}
}

func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

func New(httpClient *http.Client, opts ...ClientOption) *Client {
	c := &Client{
		http:    httpClient,
		baseURL: "https://api.spotify.com/v1",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) get(ctx context.Context, url string, result interface{}) error {
	for {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}
		resp, err := c.http.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests && c.autoRetry {
			select {
			case <-ctx.Done():
				// If the context is cancelled, return the original error
			case <-time.After(time.Second * time.Duration(c.backoff)):
				continue
			}
		}
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if resp.StatusCode != http.StatusOK {
			return c.decodeError(resp)
		}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		}

		break
	}

	return nil
}

func (c *Client) decodeError(resp *http.Response) error {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBody) == 0 {
		return fmt.Errorf("spotify: HTTP %d: %s (body empty)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	buf := bytes.NewBuffer(responseBody)

	var e struct {
		E Error `json:"error"`
	}
	err = json.NewDecoder(buf).Decode(&e)
	if err != nil {
		return fmt.Errorf("spotify: couldn't decode error: (%d) [%s]", len(responseBody), responseBody)
	}

	if e.E.Message == "" {
		e.E.Message = fmt.Sprintf("spotify: unexpected HTTP %d: %s (empty error)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return e.E
}

func (c *Client) Token() (*oauth2.Token, error) {
	transport, ok := c.http.Transport.(*oauth2.Transport)
	if !ok {
		return nil, errors.New("spotify: client not backed by oauth2 transport")
	}
	t, err := transport.Source.Token()
	if err != nil {
		return nil, err
	}
	return t, nil
}
