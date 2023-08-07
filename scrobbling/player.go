package scrobbling

import (
	"context"
	"fmt"
)

type CurrentlyPlaying struct {
	Timestamp  int  `json:"timestamp"`
	ProgressMs int  `json:"progress_ms"`
	IsPlaying  bool `json:"is_playing"`
	Item       struct {
		Album struct {
			AlbumType            string   `json:"album_type"`
			TotalTracks          int      `json:"total_tracks"`
			AvailableMarkets     []string `json:"available_markets"`
			ID                   string   `json:"id"`
			Name                 string   `json:"name"`
			ReleaseDate          string   `json:"release_date"`
			ReleaseDatePrecision string   `json:"release_date_precision"`
			Genres               []string `json:"genres"`
			Popularity           int      `json:"popularity"`
		} `json:"album"`
		Artists []struct {
			Genres     []string `json:"genres"`
			Name       string   `json:"name"`
			Popularity int      `json:"popularity"`
		} `json:"artists"`
		DiscNumber  int    `json:"disc_number"`
		DurationMs  int    `json:"duration_ms"`
		Name        string `json:"name"`
		Popularity  int    `json:"popularity"`
		TrackNumber int    `json:"track_number"`
		Type        string `json:"type"`
	} `json:"item"`
}

func (c *Client) GetCurrentlyPlaying(ctx context.Context) (*CurrentlyPlaying, error) {
	requestURL := fmt.Sprintf("%s/me/player/currently-playing", c.baseURL)

	var result CurrentlyPlaying

	err := c.get(ctx, requestURL, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
