package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	spotifyauth "github.com/cjflan/spotify-scrobbling/auth"
	"github.com/cjflan/spotify-scrobbling/spotify"
)

var (
	auth  = spotifyauth.New(spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying))
	ch    = make(chan *spotify.Client)
	state = RandomString(16)
)

func main() {

	if os.Getenv("SPOTIFY_ID") == "" {
		log.Fatal("SPOTIFY_ID not set as enviornment variable")
	}

	if os.Getenv("SPOTIFY_SECRET") == "" {
		log.Fatal("SPOTIFY_SECRET not set as enviornment variable")
	}

	fmt.Println("Hello World")
	http.HandleFunc("/callback", callback)
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

func callback(w http.ResponseWriter, r *http.Request) {
	fmt.Println()
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client

}
