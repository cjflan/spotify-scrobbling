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

var pollingRate int = 5000 // miliseconds

func (c *Client) Listen() *CurrentlyPlaying {
	var state stateBag
	var progress int
	var playing stateBag

	for { 
		currentlyPlaying, err := c.GetCurrentlyPlaying(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		
		playing = stateBag{
			title:    currentlyPlaying.Item.Name,
			artist:   currentlyPlaying.Item.Artists[0].Name,
			duration: currentlyPlaying.Item.DurationMs,
		}

		if currentlyPlaying.IsPlaying {
			fmt.Printf("%s - %s\n", currentlyPlaying.Item.Name, currentlyPlaying.Item.Artists[0].Name)
			progress += pollingRate
			if progress > playing.duration/2 {
				return currentlyPlaying
			}
		} else if playing == state {
			fmt.Printf("%s - %s is paused\n", currentlyPlaying.Item.Name, currentlyPlaying.Item.Artists[0].Name)
		}

		state = playing 
		time.Sleep(time.Millisecond * time.Duration(pollingRate))
	}
}

