package spotify

import (
	"context"
	"fmt"
)

type CurrentlyPlaying struct {
	Timestamp int64 `json:"timestamp"`
	Progress  int   `json:"progress_ms"`
	Playing   bool  `json:"is_playing"`
}

func (c *Client) GetCurrentlyPlaying(ctx context.Context) (*CurrentlyPlaying, error) {
	requestURL := fmt.Sprintf("%s/me/player/currectly-playing", c.baseURL)

	var result CurrentlyPlaying

	err := c.get(ctx, requestURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
