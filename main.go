package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cjflan/spotify-scrobbling/controllers"
	database "github.com/cjflan/spotify-scrobbling/db"
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

	if os.Getenv("MYSQL_USER") == "" {
		log.Fatal("MYSQL_USER not set as enviornment variable")
	}
	if os.Getenv("MYSQL_PASSWORD") == "" {
		log.Fatal("MYSQL_PASSWORD not set as enviornment variable")
	}
	if os.Getenv("MYSQL_DATABASE") == "" {
		log.Fatal("MYSQL_DATABASE not set as enviornment variable")
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

	db := database.Connect()

	defer db.Close()

	err := db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch
	for {
		currentlyPlaying, err := client.GetCurrentlyPlaying(context.Background())
		song := currentlyPlaying.Item.Name
		artist := currentlyPlaying.Item.Artists[0].Name
		if err != nil {
			log.Fatal(err)
		}

		if currentlyPlaying.IsPlaying {
			fmt.Printf("You are currently playing: %s - %s", song, artist)
		} else {
			fmt.Println("No song playing")
		}
		time.Sleep(5 * time.Second)
	}
}
