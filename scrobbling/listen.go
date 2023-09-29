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

func (c *Client) Listen() *CurrentlyPlaying {
	var state stateBag
	var playing stateBag
	var listened *CurrentlyPlaying
	var progress int
	var pollingRate time.Duration = 5 * time.Second
	for {
		currentlyPlaying, err := c.GetCurrentlyPlaying(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		if currentlyPlaying.Item.Name == "" {
			fmt.Println("Nothing Playing")
			time.Sleep(pollingRate)
			continue
		}
		
		playing = stateBag{
			title:    currentlyPlaying.Item.Name,
			artist:   currentlyPlaying.Item.Artists[0].Name,
			duration: currentlyPlaying.Item.DurationMs,
		}

		if currentlyPlaying.IsPlaying {
			fmt.Printf("%s - %s\n", playing.title, playing.artist)
			if currentlyPlaying.ProgressMs > currentlyPlaying.Item.DurationMs/2 {
				listened = currentlyPlaying
			}
		} else if playing == state {
			fmt.Printf("%s - %s is paused\n", playing.title, playing.artist)
		}
		if listened != nil && progress > currentlyPlaying.ProgressMs {
			fmt.Printf("Finished playing %s - %s, writing to database\n", listened.Item.Name, listened.Item.Artists[0].Name)
			return listened
		}

		progress = currentlyPlaying.ProgressMs
		state = playing
		time.Sleep(pollingRate)

	}
}
