package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cjflan/spotify-scrobbling/controllers"
)

var (
	auth  = controllers.GetAuth()
	state = controllers.GetState()
	ch    = controllers.GetCh()
)

func main() {

	if os.Getenv("SPOTIFY_ID") == "" {
		log.Fatal("SPOTIFY_ID not set as enviornment variable")
	}

	if os.Getenv("SPOTIFY_SECRET") == "" {
		log.Fatal("SPOTIFY_SECRET not set as enviornment variable")
	}

	http.HandleFunc("/callback", controllers.Callback)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		fmt.Println("Starting Server!")
		err := http.ListenAndServe("localhost:8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch

	currentlyPlaying, err := client.GetCurrentlyPlaying(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if currentlyPlaying.IsPlaying {
		fmt.Println("You are currently playing:", currentlyPlaying.Item.Name)
	} else {
		fmt.Println("No song playing")
	}
}
