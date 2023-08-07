package scrobbling

import (
	"context"
	"fmt"
	"log"
	"time"
)

type stateBag struct {
	title    string
	artist   string
	duration int
}

var pollingRate int = 5000 // in miliseconds

func (c *Client) Listen() {
	var state stateBag
	var progress int
	for {
		currentlyPlaying, err := c.GetCurrentlyPlaying(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		playing := stateBag{
			title:    currentlyPlaying.Item.Name,
			artist:   currentlyPlaying.Item.Artists[0].Name,
			duration: currentlyPlaying.Item.DurationMs,
		}

		if state == playing {
			progress++
			if progress > playing.duration {

			}
			state = playing
			progress = 0
		}

		if currentlyPlaying.IsPlaying {
			fmt.Printf("%s - %s", playing.title, playing.artist)
		} else {
			fmt.Println("No song playing")
		}
		time.Sleep(time.Duration(pollingRate) * time.Millisecond)
	}
}
