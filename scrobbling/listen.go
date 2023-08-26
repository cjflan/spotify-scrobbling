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
	progress int
	isPlaying bool
}

func (c *Client) Listen() *CurrentlyPlaying {
	var state stateBag
	var playing stateBag
	var listened *CurrentlyPlaying

	for { 
		currentlyPlaying, err := c.GetCurrentlyPlaying(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if currentlyPlaying == nil {
			fmt.Println("Nothing Playing")
			continue
		}
		
		playing = stateBag{
			title:    currentlyPlaying.Item.Name,
			artist:   currentlyPlaying.Item.Artists[0].Name,
			duration: currentlyPlaying.Item.DurationMs,
			progress: currentlyPlaying.ProgressMs,
			isPlaying: currentlyPlaying.IsPlaying,
		}

		if playing.isPlaying {
			fmt.Printf("%s - %s\n", playing.title, playing.artist)
			if playing.progress > playing.duration/2 {
				listened = currentlyPlaying
			}
		} else if playing == state {
			fmt.Printf("%s - %s is paused\n", playing.title, playing.artist)
		}
		if listened != nil && state.progress > playing.progress {
			return listened
		}
		state = playing 
		time.Sleep(time.Second * 5)
	}
}

